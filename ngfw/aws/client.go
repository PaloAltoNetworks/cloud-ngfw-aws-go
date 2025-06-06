package aws

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/sts"
	awsngfw "github.com/paloaltonetworks/cloud-ngfw-aws-go/v2"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/response"
)

// Client is the client.
type Client struct {
	CognitoClient   *cognito.CognitoIdentityProvider
	Tenant          string            `json:"tenant"`
	TenantVersion   string            `json:"tenant_version"`
	ExternalID      string            `json:"externalID"`
	Region          string            `json:"region"`
	MPRegion        string            `json:"mp_region"`
	UserName        string            `json:"userName"`
	Password        string            `json:"b64"`
	UserPoolID      string            `json:"userPoolID"`
	AppClientID     string            `json:"appClientID"`
	AppClientSecret string            `json:"appClientSecret"`
	Host            string            `json:"host"`
	MPRegionHost    string            `json:"mp_region_host"`
	V2Host          string            `json:"v2_host"`
	AccessKey       string            `json:"access-key"`
	Profile         string            `json:"profile"`
	SyncMode        bool              `json:"sync_mode"`
	SecretKey       string            `json:"secret-key"`
	Protocol        string            `json:"protocol"`
	Timeout         int               `json:"timeout"`
	ResourceTimeout int               `json:"resource_timeout"`
	Headers         map[string]string `json:"headers"`
	Agent           string            `json:"agent"`
	Origin          string            `json:"-"`

	AuthType string `json:"-"`

	LfaArn       string `json:"lfa-arn"`
	LraArn       string `json:"lra-arn"`
	GraArn       string `json:"gra-arn"`
	AcctAdminArn string `json:"account-admin-arn"`
	Arn          string `json:"arn"`

	AuthFile         string `json:"auth-file"`
	CheckEnvironment bool   `json:"-"`

	SkipVerifyCertificate bool            `json:"skip-verify-certificate"`
	Transport             *http.Transport `json:"-"`

	// Various logging params.
	Logging               uint32   `json:"-"`
	LoggingFromInitialize []string `json:"logging"`

	// Configured by Initialize().
	FirewallAdminJwt               string     `json:"-"`
	FirewallAdminJwtExpTime        time.Time  `json:"-"`
	FirewallSubscriptionKey        string     `json:"-"`
	FirewallAdminMutex             sync.Mutex `json:"-"`
	RulestackAdminJwt              string     `json:"-"`
	RulestackAdminJwtExpTime       time.Time  `json:"-"`
	RulestackSubscriptionKey       string     `json:"-"`
	RulestackAdminMutex            sync.Mutex `json:"-"`
	GlobalRulestackAdminJwt        string     `json:"-"`
	GlobalRulestackAdminJwtExpTime time.Time  `json:"-"`
	GlobalRulestackSubscriptionKey string     `json:"-"`
	GlobalRulestackAdminMutex      sync.Mutex `json:"-"`
	CloudRulestackAdminJwt         string     `json:"-"`
	CloudRulestackAdminJwtExpTime  time.Time  `json:"-"`
	CloudRulestackSubscriptionKey  string     `json:"-"`
	CloudRulestackAdminMutex       sync.Mutex `json:"-"`
	AccountAdminJwt                string     `json:"-"`
	AccountAdminJwtExpTime         time.Time  `json:"-"`
	AccountAdminSubscriptionKey    string     `json:"-"`
	AccountAdminMutex              sync.Mutex `json:"-"`

	// Internal variables.
	apiPrefix   string
	mpApiPrefix string
	v2ApiPrefix string

	// Initialized during Setup().
	HttpClient       *http.Client
	SecureHttpClient *http.Client
	AuthURL          string

	// Variables for testing.
	testData        [][]byte
	testErrors      []error
	testIndex       int
	authFileContent []byte

	// Used for unit tests
	Mock       bool
	MockedResp func() ([]byte, error)
}

// NgfwAuthInput struct
type NgfwAuthInput struct {
	ExternalID       string
	Timeout          int
	HttpClient       *http.Client
	SecureHttpClient *http.Client
	RegionURL        string
	AuthURL          string
}

// Log logs an API action.
func (c *Client) Log(method, msg string, i ...interface{}) {
	switch method {
	case http.MethodGet:
		if c.Logging&awsngfw.LogGet != awsngfw.LogGet {
			return
		}
	case http.MethodPatch:
		if c.Logging&awsngfw.LogPatch != awsngfw.LogPatch && c.Logging&awsngfw.LogAction != awsngfw.LogAction {
			return
		}
	case http.MethodPost:
		if c.Logging&awsngfw.LogPost != awsngfw.LogPost && c.Logging&awsngfw.LogAction != awsngfw.LogAction {
			return
		}
	case http.MethodPut:
		if c.Logging&awsngfw.LogPut != awsngfw.LogPut && c.Logging&awsngfw.LogAction != awsngfw.LogAction {
			return
		}
	case http.MethodDelete:
		if c.Logging&awsngfw.LogDelete != awsngfw.LogDelete && c.Logging&awsngfw.LogAction != awsngfw.LogAction {
			return
		}
	default:
		return
	}

	log.Printf("(%s) %s", strings.ToLower(method), fmt.Sprintf(msg, i...))
}

// Setup configures the HttpClient param according to the combination of
// locally defined params, environment variables, and the JSON config file.
func (c *Client) Setup() error {
	var err error
	var tout time.Duration

	// Load up the JSON config file.
	json_client := &Client{}
	if c.AuthFile != "" {
		var b []byte
		if len(c.testData) == 0 {
			b, err = ioutil.ReadFile(c.AuthFile)
		} else {
			b, err = c.authFileContent, nil
		}

		if err != nil {
			return err
		}

		if err = json.Unmarshal(b, &json_client); err != nil {
			return err
		}
	}

	// Host.
	if c.Host == "" {
		if val := os.Getenv("CLOUDNGFWAWS_HOST"); c.CheckEnvironment && val != "" {
			c.Host = val
		} else if json_client.Host != "" {
			c.Host = json_client.Host
		}
	}
	if c.Host == "" {
		c.Host = "api.us-east-1.aws.cloudngfw.paloaltonetworks.com"
	}

	if c.V2Host == "" {
		if val := os.Getenv("CLOUDNGFWAWS_V2_HOST"); c.CheckEnvironment && val != "" {
			c.V2Host = val
		} else if json_client.V2Host != "" {
			c.V2Host = json_client.V2Host
		}
	}
	if c.V2Host == "" {
		c.V2Host = "api.us-east-1.aws.cloudngfw.paloaltonetworks.com"
	}

	// Host.
	if c.MPRegionHost == "" {
		if val := os.Getenv("CLOUDNGFWAWS_MP_REGION_HOST"); c.CheckEnvironment && val != "" {
			c.MPRegionHost = val
		} else if json_client.MPRegionHost != "" {
			c.MPRegionHost = json_client.MPRegionHost
		}
	}
	if c.MPRegionHost == "" {
		c.MPRegionHost = DefaultMPRegionHost
	}

	// Access key.
	if c.AccessKey == "" {
		if val := os.Getenv("CLOUDNGFWAWS_ACCESS_KEY"); c.CheckEnvironment && val != "" {
			c.AccessKey = val
		} else if json_client.AccessKey != "" {
			c.AccessKey = json_client.AccessKey
		}
	}

	// Secret key.
	if c.SecretKey == "" {
		if val := os.Getenv("CLOUDNGFWAWS_SECRET_KEY"); c.CheckEnvironment && val != "" {
			c.SecretKey = val
		} else if json_client.SecretKey != "" {
			c.SecretKey = json_client.SecretKey
		}
	}

	// Profile.
	if c.Profile == "" {
		if val := os.Getenv("CLOUDNGFWAWS_PROFILE"); c.CheckEnvironment && val != "" {
			c.Profile = val
		} else if json_client.Profile != "" {
			c.Profile = json_client.Profile
		}
	}

	// SyncMode.
	if c.SyncMode == false {
		if val := os.Getenv("CLOUDNGFWAWS_SYNC_MODE"); c.CheckEnvironment && strings.ToLower(val) == "true" {
			c.SyncMode = true
		} else if json_client.SyncMode == true {
			c.SyncMode = true
		}
	}

	// LFA ARN.
	if c.LfaArn == "" {
		if val := os.Getenv("CLOUDNGFWAWS_LFA_ARN"); c.CheckEnvironment && val != "" {
			c.LfaArn = val
		} else if json_client.LfaArn != "" {
			c.LfaArn = json_client.LfaArn
		}
	}

	// LRA ARN.
	if c.LraArn == "" {
		if val := os.Getenv("CLOUDNGFWAWS_LRA_ARN"); c.CheckEnvironment && val != "" {
			c.LraArn = val
		} else if json_client.LraArn != "" {
			c.LraArn = json_client.LraArn
		}
	}

	// GRA ARN.
	if c.GraArn == "" {
		if val := os.Getenv("CLOUDNGFWAWS_GRA_ARN"); c.CheckEnvironment && val != "" {
			c.GraArn = val
		} else if json_client.GraArn != "" {
			c.GraArn = json_client.GraArn
		}
	}

	// Account Admin ARN.
	if c.AcctAdminArn == "" {
		if val := os.Getenv("CLOUDNGFWAWS_ACCT_ADMIN_ARN"); c.CheckEnvironment && val != "" {
			c.AcctAdminArn = val
		} else if json_client.AcctAdminArn != "" {
			c.AcctAdminArn = json_client.AcctAdminArn
		}
	}

	// ARN.
	if c.Arn == "" {
		if val := os.Getenv("CLOUDNGFWAWS_ARN"); c.CheckEnvironment && val != "" {
			c.Arn = val
		} else if json_client.Arn != "" {
			c.Arn = json_client.Arn
		}
	}

	// Region.
	if c.Region == "" {
		if val := os.Getenv("CLOUDNGFWAWS_REGION"); c.CheckEnvironment && val != "" {
			c.Region = val
		} else if json_client.Region != "" {
			c.Region = json_client.Region
		} else {
			return fmt.Errorf("No region was specified")
		}
	}

	// MP Region.
	if c.MPRegion == "" {
		if val := os.Getenv("CLOUDNGFWAWS_MP_REGION"); c.CheckEnvironment && val != "" {
			c.MPRegion = val
		} else if json_client.MPRegion != "" {
			c.MPRegion = json_client.MPRegion
		}
	}

	if c.MPRegion == "" {
		c.MPRegion = DefaultMPRegion
	}

	// Protocol.
	if c.Protocol == "" {
		if val := os.Getenv("CLOUDNGFWAWS_PROTOCOL"); c.CheckEnvironment && val != "" {
			c.Protocol = val
		} else if json_client.Protocol != "" {
			c.Protocol = json_client.Protocol
		} else {
			c.Protocol = "https"
		}
	}
	if c.Protocol != "http" && c.Protocol != "https" {
		return fmt.Errorf("Invalid protocol %q; expected 'https' or 'http'", c.Protocol)
	}

	// Timeout.
	if c.Timeout == 0 {
		if val := os.Getenv("CLOUDNGFWAWS_TIMEOUT"); c.CheckEnvironment && val != "" {
			if ival, err := strconv.Atoi(val); err != nil {
				return fmt.Errorf("Failed to parse timeout env var as int: %s", err)
			} else {
				c.Timeout = ival
			}
		} else if json_client.Timeout != 0 {
			c.Timeout = json_client.Timeout
		} else {
			c.Timeout = 30
		}
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("Timeout for %q must be a positive int", c.Host)
	}
	tout = time.Duration(time.Duration(c.Timeout) * time.Second)

	// Resource Timeout.
	if c.ResourceTimeout == 0 {
		if val := os.Getenv("CLOUDNGFWAWS_RESOURCE_TIMEOUT"); c.CheckEnvironment && val != "" {
			if ival, err := strconv.Atoi(val); err != nil {
				return fmt.Errorf("Failed to parse resource timeout env var as int: %s", err)
			} else {
				c.ResourceTimeout = ival
			}
		} else if json_client.ResourceTimeout != 0 {
			c.ResourceTimeout = json_client.ResourceTimeout
		} else {
			c.ResourceTimeout = 7200
		}
	}
	if c.ResourceTimeout <= 0 {
		return fmt.Errorf("Resource Timeout must be a positive int")
	}

	// Headers.
	if len(c.Headers) == 0 {
		if val := os.Getenv("CLOUDNGFWAWS_HEADERS"); c.CheckEnvironment && val != "" {
			if err := json.Unmarshal([]byte(val), &c.Headers); err != nil {
				return err
			}
		}
		if len(c.Headers) == 0 && len(json_client.Headers) > 0 {
			c.Headers = make(map[string]string)
			for k, v := range json_client.Headers {
				c.Headers[k] = v
			}
		}
	}

	// Verify cert.
	if !c.SkipVerifyCertificate {
		if val := os.Getenv("CLOUDNGFWAWS_SKIP_VERIFY_CERTIFICATE"); c.CheckEnvironment && val != "" {
			if vcb, err := strconv.ParseBool(val); err != nil {
				return err
			} else if vcb {
				c.SkipVerifyCertificate = vcb
			}
		}
		if !c.SkipVerifyCertificate && json_client.SkipVerifyCertificate {
			c.SkipVerifyCertificate = json_client.SkipVerifyCertificate
		}
	}

	// Logging.
	if c.Logging == 0 {
		var ll []string
		if val := os.Getenv("CLOUDNGFWAWS_LOGGING"); c.CheckEnvironment && val != "" {
			ll = strings.Split(val, ",")
		} else {
			ll = json_client.LoggingFromInitialize
		}
		if len(ll) > 0 {
			var lv uint32
			for _, x := range ll {
				switch x {
				case "quiet":
					lv |= awsngfw.LogQuiet
				case "login":
					lv |= awsngfw.LogLogin
				case "get":
					lv |= awsngfw.LogGet
				case "patch":
					lv |= awsngfw.LogPatch
				case "post":
					lv |= awsngfw.LogPost
				case "put":
					lv |= awsngfw.LogPut
				case "delete":
					lv |= awsngfw.LogDelete
				case "action":
					lv |= awsngfw.LogPatch | awsngfw.LogPost | awsngfw.LogPut | awsngfw.LogDelete
				case "path":
					lv |= awsngfw.LogPath
				case "send":
					lv |= awsngfw.LogSend
				case "receive":
					lv |= awsngfw.LogReceive
				default:
					return fmt.Errorf("Unknown logging requested: %s", x)
				}
			}
			c.Logging = lv
		} else {
			c.Logging = awsngfw.LogLogin | awsngfw.LogGet | awsngfw.LogAction
		}
	}

	// Setup the https client.
	if c.Transport == nil {
		c.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.SkipVerifyCertificate,
			},
		}
	}
	c.HttpClient = &http.Client{
		Transport: c.Transport,
		Timeout:   tout,
	}

	// Configure the uri prefix.
	c.apiPrefix = fmt.Sprintf("%s://%s", c.Protocol, c.Host)
	c.mpApiPrefix = fmt.Sprintf("%s://%s", c.Protocol, c.MPRegionHost)
	c.v2ApiPrefix = fmt.Sprintf("%s://%s", c.Protocol, c.V2Host)
	return nil
}

// Path holds the V1 and V2 API paths for a given resource.
type Path struct {
	V1Path []string
	V2Path []string
}

func setV2Path(c *Client, path []string, req *http.Request, queryParams url.Values) error {
	var qp string
	if queryParams == nil {
		queryParams = url.Values{}
	}
	if !queryParams.Has("region") {
		queryParams.Set("region", c.Region)
	}
	if len(queryParams) > 0 {
		qp = fmt.Sprintf("?%s", queryParams.Encode())
	}
	u, err := url.Parse(fmt.Sprintf("%s/%s%s", c.v2ApiPrefix, strings.Join(path, "/"), qp))
	if err != nil {
		return err
	}
	req.URL = u
	req.Host = u.Host
	return nil
}

/*
Communicate sends information to the API.

Param auth should be one of the permissions constants or an empty string,
which means not to add any JWTs to the API call.

Param method should be one of http.Method constants.

Param path should be a slice of path parts that will be joined together with
the base apiPrefix to create the final API endpoint.

Param queryParams are the query params that should be appended to the API URL.

Param input is an interface that can be passed in to json.Marshal() to send to
the API.

Param output is a pointer to a struct that will be filled with json.Unmarshal().

Param creds is only used internally for refreshing the JWTs and can otherwise
be ignored.

This function returns the content of the body from the API call and any errors
that may have been present.  If this function got all the way to invoking the
API and getting a response, then the error passed back will be a `api.Status`
if an error was detected.
*/
func (c *Client) Communicate(ctx context.Context, auth, method string, path Path, queryParams url.Values, input interface{}, output response.Failure, creds ...*sts.Credentials) (s []byte, e error) {
	// check if mocking is enabled(for unit test purposes)
	if c.Mock {
		log.Printf("mocking response.")
		return c.MockedResp()
	}

	// Sanity check the input.
	if len(creds) > 1 {
		return nil, fmt.Errorf("[tenant:%s][region:%s] Only one credentials is allowed",
			c.ExternalID, c.Region)
	}

	var err error
	var body []byte
	var data []byte

	// Convert input into JSON.
	if input != nil {
		data, err = json.Marshal(input)
		if err != nil {
			return nil, err
		}
	}

	if c.Logging&awsngfw.LogSend == awsngfw.LogSend {
		log.Printf("sending: %s", data)
	}

	// Create the API request.
	var qp string
	if len(queryParams) > 0 {
		qp = fmt.Sprintf("?%s", queryParams.Encode())
	}
	apiPrefix := c.apiPrefix
	if auth == PermissionAccountAdminJWT || auth == PermissionAccount {
		apiPrefix = c.mpApiPrefix
	}
	if c.Logging&awsngfw.LogPath == awsngfw.LogPath {
		log.Printf("path: %s/%s%s", apiPrefix, strings.Join(path.V1Path, "/"), qp)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("%s/%s%s", apiPrefix, strings.Join(path.V1Path, "/"), qp),
		strings.NewReader(string(data)),
	)
	if err != nil {
		return nil, err
	}

	// Add in the custom headers.
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	// Configure standard headers.
	permErr := "[tenant:%s][region:%s]This connection does not have the required JWT:%s err:%s"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.Agent)
	switch auth {
	case "", PermissionAccountAdminJWT:
	case PermissionFirewall:
		switch c.AuthType {
		case AuthTypeIAMRole:
			err := c.RefreshFirewallAdminJwt(ctx)
			if err != nil {
				return nil, fmt.Errorf(permErr, c.ExternalID, c.Region, PermissionFirewall, err)
			}
			req.Header.Set("Authorization", c.FirewallAdminJwt)
			req.Header.Set("x-api-key", c.FirewallSubscriptionKey)
		case AuthTypeExternalID:
			err := c.RefreshCloudRulestackAdminJwt(ctx)
			if err != nil {
				return nil, fmt.Errorf(permErr, c.ExternalID, c.Region, PermissionFirewall, err)
			}
			req.Header.Set("Authorization", c.CloudRulestackAdminJwt)
			req.Header.Set("x-api-key", c.CloudRulestackSubscriptionKey)
		default:
			req.Header.Set("Authorization", c.FirewallAdminJwt)
			req.Header.Set("x-api-key", c.FirewallSubscriptionKey)
		}
	case PermissionRulestack:
		switch c.AuthType {
		case AuthTypeIAMRole:
			err := c.RefreshRulestackAdminJwt(ctx)
			if err != nil {
				return nil, fmt.Errorf(permErr, c.ExternalID, c.Region, PermissionRulestack, err)
			}
			req.Header.Set("Authorization", c.RulestackAdminJwt)
			req.Header.Set("x-api-key", c.RulestackSubscriptionKey)
		case AuthTypeExternalID:
			err := c.RefreshCloudRulestackAdminJwt(ctx)
			if err != nil {
				return nil, fmt.Errorf(permErr, c.ExternalID, c.Region, PermissionRulestack, err)
			}
			req.Header.Set("Authorization", c.CloudRulestackAdminJwt)
			req.Header.Set("x-api-key", c.CloudRulestackSubscriptionKey)
		default:
			req.Header.Set("Authorization", c.RulestackAdminJwt)
			req.Header.Set("x-api-key", c.RulestackSubscriptionKey)
		}
	case PermissionGlobalRulestack:
		switch c.AuthType {
		case AuthTypeIAMRole:
			err := c.RefreshGlobalRulestackAdminJwt(ctx)
			if err != nil {
				return nil, fmt.Errorf(permErr, c.ExternalID, c.Region, PermissionGlobalRulestack, err)
			}
			req.Header.Set("Authorization", c.GlobalRulestackAdminJwt)
			req.Header.Set("x-api-key", c.GlobalRulestackSubscriptionKey)
		case AuthTypeExternalID:
			err := c.RefreshCloudRulestackAdminJwt(ctx)
			if err != nil {
				return nil, fmt.Errorf(permErr, c.ExternalID, c.Region, PermissionGlobalRulestack, err)
			}
			req.Header.Set("Authorization", c.CloudRulestackAdminJwt)
			req.Header.Set("x-api-key", c.CloudRulestackSubscriptionKey)
		default:
			req.Header.Set("Authorization", c.GlobalRulestackAdminJwt)
			req.Header.Set("x-api-key", c.GlobalRulestackSubscriptionKey)
		}
	case PermissionAccount:
		switch c.AuthType {
		case AuthTypeIAMRole:
			err := c.RefreshAccountAdminJwt(ctx)
			if err != nil {
				return nil, fmt.Errorf(permErr, c.ExternalID, c.MPRegion, PermissionAccount, err)
			}
			req.Header.Set("Authorization", c.AccountAdminJwt)
			req.Header.Set("x-api-key", c.AccountAdminSubscriptionKey)
		default:
			req.Header.Set("Authorization", c.AccountAdminJwt)
			req.Header.Set("x-api-key", c.AccountAdminSubscriptionKey)
		}
	default:
		return nil, fmt.Errorf("[tenant:%s][region:%s] Unknown permission required: %q",
			c.ExternalID, c.MPRegion, auth)
	}
	api.Logger.Debugf("SDK Path: %v", path)

	if c.TenantVersion == awsngfw.TenantVersionV2 && path.V2Path != nil {
		if err := setV2Path(c, path.V2Path, req, queryParams); err != nil {
			return nil, err
		}
	}
	api.Logger.Debugf("SDK Request URL: %s", req.URL.String())

	// Optional: v4 sign the request.
	if len(creds) == 1 {
		prov := provider{
			Value: credentials.Value{
				AccessKeyID:     *creds[0].AccessKeyId,
				SecretAccessKey: *creds[0].SecretAccessKey,
				SessionToken:    *creds[0].SessionToken,
			},
		}
		signer := v4.NewSigner(credentials.NewCredentials(prov))
		region := c.Region
		if auth == PermissionAccountAdminJWT || auth == PermissionAccount {
			region = c.MPRegion
		}
		_, err = signer.Sign(req, strings.NewReader(string(data)), "execute-api", region, time.Now())
		if err != nil {
			return nil, err
		}
	}

	// Perform the API action.
	if len(c.testData) > 0 {
		body = []byte(`{"test"}`)
	} else {
		resp, err := c.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
	}

	if err != nil {
		return nil, err
	}

	// Log the response.
	if c.Logging&awsngfw.LogReceive == awsngfw.LogReceive {
		log.Printf("received: %s", body)
	}

	// Check for unknown path error first.
	if err := response.IsResponseWithError(body); err != nil {
		return body, err
	}

	// Check for errors and unmarshal the response into the given interface.
	if output != nil {
		c.Log(method, "Unmarshaling response into: %T", output)
		if err = json.Unmarshal(body, output); err != nil {
			c.Log(method, "Error unmarshaling response: %v", err)
			return body, err
		}

		if e2 := output.Failed(); e2 != nil {
			c.Log(method, "Error unmarshaling response output: %v", e2)
			return body, e2
		}
	} else {
		c.Log(method, "generic response")
		var generic Response
		if err = json.Unmarshal(body, &generic); err != nil {
			return body, err
		}

		if e2 := generic.Failed(); e2 != nil {
			return body, e2
		}
	}

	return body, nil
}

func (c *Client) RequestJwt(ctx context.Context, method string, path []string, queryParams url.Values, input interface{}, output response.Failure, creds ...*sts.Credentials) ([]byte, error) {
	// Sanity check the input.
	if len(creds) > 1 {
		return nil, fmt.Errorf("Only one credentials is allowed")
	}

	var err error
	var body []byte
	var data []byte

	// Convert input into JSON.
	if input != nil {
		data, err = json.Marshal(input)
		if err != nil {
			return nil, err
		}
	}
	if c.Logging&awsngfw.LogSend == awsngfw.LogSend {
		log.Printf("sending: %s", data)
	}

	// Create the API request.
	var qp string
	if len(queryParams) > 0 {
		qp = fmt.Sprintf("?%s", queryParams.Encode())
	}
	if c.Logging&awsngfw.LogPath == awsngfw.LogPath {
		log.Printf("path: %s/%s%s", c.apiPrefix, strings.Join(path, "/"), qp)
	}
	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("%s/%s%s", c.apiPrefix, strings.Join(path, "/"), qp),
		strings.NewReader(string(data)),
	)
	if err != nil {
		return nil, err
	}

	// Add in the custom headers.
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	// Configure standard headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.Agent)

	// v4 sign the request.
	if len(creds) == 1 {
		prov := provider{
			Value: credentials.Value{
				AccessKeyID:     *creds[0].AccessKeyId,
				SecretAccessKey: *creds[0].SecretAccessKey,
				SessionToken:    *creds[0].SessionToken,
			},
		}
		signer := v4.NewSigner(credentials.NewCredentials(prov))
		_, err = signer.Sign(req, strings.NewReader(string(data)), "execute-api", c.Region, time.Now())
		if err != nil {
			return nil, err
		}
	}

	var resp *http.Response
	// Perform the API action.
	if len(c.testData) > 0 {
		body = []byte(`{"test"}`)
	} else {
		resp, err = c.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
	}

	if err != nil {
		return nil, err
	}

	// Log the response.
	if c.Logging&awsngfw.LogReceive == awsngfw.LogReceive {
		log.Printf("received: %s", body)
	}

	// Check for unknown path error first.
	if err := response.IsResponseWithError(body); err != nil {
		return body, err
	}

	// Check for errors and unmarshal the response into the given interface.
	if output != nil {
		if err = json.Unmarshal(body, output); err != nil {
			return body, err
		}

		if e2 := output.Failed(); e2 != nil {
			return body, e2
		}
	} else {
		var generic Response
		if err = json.Unmarshal(body, &generic); err != nil {
			return body, err
		}

		if e2 := generic.Failed(); e2 != nil {
			return body, e2
		}
	}

	return body, nil
}

func (c *Client) SetEndpoint(ctx context.Context, input api.EndPointInput) error {
	c.apiPrefix = input.ApiEndpoint
	c.AuthURL = input.ApiAuthEndpoint
	return nil
}

func (c *Client) IsSyncModeEnabled(ctx context.Context) bool {
	return c.SyncMode
}

func (c *Client) GetMPRegion(ctx context.Context) string {
	return c.MPRegion
}

func (c *Client) GetRegion(ctx context.Context) string {
	return c.Region
}

func (c *Client) GetProfile(ctx context.Context) string {
	return c.Profile
}

func (c *Client) GetResourceTimeout(ctx context.Context) int {
	return c.ResourceTimeout
}

func (c *Client) GetCloudProvider(ctx context.Context) string {
	return awsngfw.CloudProviderAWS
}
