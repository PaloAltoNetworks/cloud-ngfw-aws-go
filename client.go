package awsngfw

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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/permissions"
)

// Client is the client.
type Client struct {
	Host      string            `json:"host"`
	AccessKey string            `json:"access-key"`
	SecretKey string            `json:"secret-key"`
	Region    string            `json:"region"`
	Protocol  string            `json:"protocol"`
	Timeout   int               `json:"timeout"`
	Headers   map[string]string `json:"headers"`
	Agent     string            `json:"agent"`

	LfaArn string `json:"lfa-arn"`
	LraArn string `json:"lra-arn"`
	Arn    string `json:"arn"`

	AuthFile         string `json:"auth-file"`
	CheckEnvironment bool   `json:"-"`

	SkipVerifyCertificate bool            `json:"skip-verify-certificate"`
	Transport             *http.Transport `json:"-"`

	// Various logging params.
	Logging               uint32   `json:"-"`
	LoggingFromInitialize []string `json:"logging"`

	// Configured by Initialize().
	FirewallJwt  string `json:"-"`
	RulestackJwt string `json:"-"`

	// Internal variables.
	apiPrefix string

	// Initialized during Setup().
	HttpClient *http.Client

	// Variables for testing.
	testData        [][]byte
	testErrors      []error
	testIndex       int
	authFileContent []byte
}

// Log logs an API action.
func (c *Client) Log(method, msg string, i ...interface{}) {
	switch method {
	case http.MethodGet:
		if c.Logging&LogGet != LogGet {
			return
		}
	case http.MethodPost:
		if c.Logging&LogPost != LogPost {
			return
		}
	case http.MethodPut:
		if c.Logging&LogPut != LogPut {
			return
		}
	case http.MethodDelete:
		if c.Logging&LogDelete != LogDelete {
			return
		}
	default:
		return
	}

	log.Printf("(%s) %s", method, fmt.Sprintf(msg, i...))
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
		if val := os.Getenv("CLOUD_NGFW_HOST"); c.CheckEnvironment && val != "" {
			c.Host = val
		} else if json_client.Host != "" {
			c.Host = json_client.Host
		}
	}
	if c.Host == "" {
		c.Host = "api.us-east-1.aws.cloudngfw.com"
	}

	// Region.
	if c.Region == "" {
		if val := os.Getenv("CLOUD_NGFW_REGION"); c.CheckEnvironment && val != "" {
			c.Region = val
		} else if json_client.Region != "" {
			c.Region = json_client.Region
		} else {
			return fmt.Errorf("No region was specified")
		}
	}

	// Protocol.
	if c.Protocol == "" {
		if val := os.Getenv("CLOUD_NGFW_PROTOCOL"); c.CheckEnvironment && val != "" {
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
		if val := os.Getenv("CLOUD_NGFW_TIMEOUT"); c.CheckEnvironment && val != "" {
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

	// Headers.
	if len(c.Headers) == 0 {
		if val := os.Getenv("CLOUD_NGFW_HEADERS"); c.CheckEnvironment && val != "" {
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

	// LFA ARN.
	if c.LfaArn == "" {
		if val := os.Getenv("CLOUD_NGFW_LFA_ARN"); c.CheckEnvironment && val != "" {
			c.LfaArn = val
		} else if json_client.LfaArn != "" {
			c.LfaArn = json_client.LfaArn
		}
	}

	// LRA ARN.
	if c.LraArn == "" {
		if val := os.Getenv("CLOUD_NGFW_LRA_ARN"); c.CheckEnvironment && val != "" {
			c.LraArn = val
		} else if json_client.LraArn != "" {
			c.LraArn = json_client.LraArn
		}
	}

	// ARN.
	if c.Arn == "" {
		if val := os.Getenv("CLOUD_NGFW_ARN"); c.CheckEnvironment && val != "" {
			c.Arn = val
		} else if json_client.Arn != "" {
			c.Arn = json_client.Arn
		}
	}

	// Verify cert.
	if !c.SkipVerifyCertificate {
		if val := os.Getenv("CLOUD_NGFW_SKIP_VERIFY_CERTIFICATE"); c.CheckEnvironment && val != "" {
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
		if val := os.Getenv("CLOUD_NGFW_LOGGING"); c.CheckEnvironment && val != "" {
			ll = strings.Split(val, ",")
		} else {
			ll = json_client.LoggingFromInitialize
		}
		if len(ll) > 0 {
			var lv uint32
			for _, x := range ll {
				switch x {
				case "quiet":
					lv |= LogQuiet
				case "login":
					lv |= LogLogin
				case "get":
					lv |= LogGet
				case "post":
					lv |= LogPost
				case "put":
					lv |= LogPut
				case "delete":
					lv |= LogDelete
				case "path":
					lv |= LogPath
				case "send":
					lv |= LogSend
				case "receive":
					lv |= LogReceive
				default:
					return fmt.Errorf("Unknown logging requested: %s", x)
				}
			}
			c.Logging = lv
		} else {
			c.Logging = LogLogin | LogGet | LogPost | LogPut | LogDelete
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
	//c.apiPrefix = fmt.Sprintf("%s://api.%s.aws.awsngfw.com", c.Protocol, c.Region)

	return nil
}

// RefreshJwts refreshes all JWTs and stores them for future API calls.
func (c *Client) RefreshJwts(ctx context.Context) error {
	if c.Logging&LogLogin == LogLogin {
		log.Printf("(login) refreshing JWTs...")
	}

	jwtReq := getJwt{
		Expires: 120,
		KeyInfo: &jwtKeyInfo{
			Region: c.Region,
			Tenant: "XY",
		},
	}

	var creds *credentials.Credentials
	if c.AccessKey != "" || c.SecretKey != "" {
		creds = credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: creds,
			Region:      aws.String(c.Region),
		},
	})

	if err != nil {
		return err
	}

	svc := sts.New(sess)
	results := make(chan error)

	go func() {
		// Get firewall JWT.
		var rarn *string
		if c.LfaArn != "" {
			rarn = aws.String(c.LfaArn)
		} else if c.Arn != "" {
			rarn = aws.String(c.Arn)
		} else {
			results <- nil
			return
		}

		if c.Logging&LogLogin == LogLogin {
			log.Printf("(login) refreshing firewall JWT...")
		}
		result, err := svc.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         rarn,
			RoleSessionName: aws.String("sdk_session"),
		})
		if err != nil {
			results <- err
			return
		}

		var ans authResponse
		_, err = c.Communicate(
			ctx, "", http.MethodGet, []string{"v1", "mgmt", "tokens", "cloudfirewalladmin"}, nil, jwtReq, &ans, result.Credentials,
		)
		if err != nil {
			results <- err
			return
		}

		c.FirewallJwt = ans.Resp.Jwt
		results <- nil
	}()

	go func() {
		// Get rulestack JWT.
		var rarn *string
		if c.LraArn != "" {
			rarn = aws.String(c.LraArn)
		} else if c.Arn != "" {
			rarn = aws.String(c.Arn)
		} else {
			results <- nil
			return
		}

		if c.Logging&LogLogin == LogLogin {
			log.Printf("(login) refreshing rulestack JWT...")
		}
		result, err := svc.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         rarn,
			RoleSessionName: aws.String("sdk_session"),
		})
		if err != nil {
			results <- err
			return
		}

		var ans authResponse
		_, err = c.Communicate(
			ctx, "", http.MethodGet, []string{"v1", "mgmt", "tokens", "cloudrulestackadmin"}, nil, jwtReq, &ans, result.Credentials,
		)
		if err != nil {
			results <- err
			return
		}

		c.RulestackJwt = ans.Resp.Jwt
		results <- nil
	}()

	e1, e2 := <-results, <-results
	if e1 != nil {
		return e1
	} else if e2 != nil {
		return e2
	} else if c.FirewallJwt == "" && c.RulestackJwt == "" {
		return fmt.Errorf("No ARNs were specified")
	}

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
func (c *Client) Communicate(ctx context.Context, auth, method string, path []string, queryParams url.Values, input interface{}, output api.Failure, creds ...*sts.Credentials) ([]byte, error) {
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
	if c.Logging&LogSend == LogSend {
		log.Printf("sending: %s", data)
	}

	if len(c.testData) > 0 {
		// Testing.
		body = []byte(`{"test"}`)
	} else {
		// Create the API request.
		var qp string
		if len(queryParams) > 0 {
			qp = fmt.Sprintf("?%s", queryParams.Encode())
		}
		if c.Logging&LogPath == LogPath {
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

		// Configure headers.
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", c.Agent)
		switch auth {
		case "":
		case permissions.Firewall:
			req.Header.Set("Authorization", c.FirewallJwt)
		case permissions.Rulestack:
			req.Header.Set("Authorization", c.RulestackJwt)
		default:
			return nil, fmt.Errorf("Unknown auth type: %q", auth)
		}
		// Add in the custom headers.
		for k, v := range c.Headers {
			req.Header.Set(k, v)
		}

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
			_, err = signer.Sign(req, strings.NewReader(string(data)), "execute-api", c.Region, time.Now())
			if err != nil {
				return nil, err
			}
		}

		// Perform the API action.
		resp, err := c.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	// Log the response.
	if c.Logging&LogReceive == LogReceive {
		log.Printf("received: %s", body)
	}

	// Check for unknown path error first.
	if err := api.IsPathUnknownError(path, body); err != nil {
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
		var generic api.Response
		if err = json.Unmarshal(body, &generic); err != nil {
			return body, err
		}

		if e2 := generic.Failed(); e2 != nil {
			return body, e2
		}
	}

	return body, nil
}
