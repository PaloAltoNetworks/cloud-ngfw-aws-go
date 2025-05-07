package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/stack"

	"github.com/aws/aws-sdk-go/aws"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v5"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
)

const (
	flowUsernamePassword = "USER_PASSWORD_AUTH"
	AuthEndpoint         = "v1/mgmt/tokens/cloudmanager"
	PanoramaEndpoint     = "v1/mgmt/cloudservicetokens/panorama"
	DefaultExpiryTime    = 30 * time.Minute
	MaxBackoffTime       = 10 * time.Minute
)

// externalID based auth over mTLS

type AuthInfo struct {
	TenantID         string
	ExternalID       string
	ExpiryTIme       int
	HttpClient       *http.Client
	SecureHttpClient *http.Client
	Region           string
	RegionURL        string
	AuthURL          string
	RegionV2URL      string
}

// Setup configures the HttpClient param according to the combination of
// locally defined params, environment variables, and the JSON config file.
func (c *Client) SetupUsingCredentials(regionURL string, httpClient *http.Client) error {
	var wg sync.WaitGroup
	c.HttpClient = httpClient
	// Configure the uri prefix.
	c.apiPrefix = regionURL
	go func() {
		for {
			wg.Add(1)
			flow := aws.String(flowUsernamePassword)
			params := map[string]*string{
				"USERNAME": aws.String(c.UserName),
				"PASSWORD": aws.String(c.Password),
			}

			authTry := &cognito.InitiateAuthInput{
				AuthFlow:       flow,
				AuthParameters: params,
				ClientId:       aws.String(c.AppClientID),
			}

			res, err := c.CognitoClient.InitiateAuth(authTry)
			if err != nil {
				log.Printf("error occured in refreshing token:%+v", err)
				wg.Done()
				continue
			}

			c.FirewallAdminMutex.Lock()
			c.FirewallAdminJwt = *res.AuthenticationResult.IdToken
			c.FirewallAdminMutex.Unlock()
			c.RulestackAdminMutex.Lock()
			c.RulestackAdminJwt = *res.AuthenticationResult.IdToken
			c.RulestackAdminMutex.Unlock()
			c.GlobalRulestackAdminMutex.Lock()
			c.GlobalRulestackAdminJwt = *res.AuthenticationResult.IdToken
			c.GlobalRulestackAdminMutex.Unlock()
			wg.Done()
			log.Printf("refreshed token.")
			<-time.After(8 * time.Minute)
		}
	}()
	log.Printf("waiting for token intialization")
	<-time.After(1 * time.Second)
	wg.Wait()
	log.Printf("done")
	return nil
}

// this API should be called with c.CloudRulestackAdminMutex taken
func (c *Client) doAuth(ctx context.Context, info AuthInfo) error {
	api.Logger.Debugf("[tenant:%s][region:%s] refreshing token...",
		c.ExternalID, c.Region)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s?externalid=%s", info.AuthURL, AuthEndpoint, c.ExternalID),
		nil,
	)
	api.Logger.Debugf("[tenant:%s][region:%s] req:%+v", c.ExternalID, c.Region, req)
	if err != nil {
		api.Logger.Errorf("[tenant:%s][region:%s] err:%+v", c.ExternalID, c.Region, err)
		return c.updateToken(c.ExternalID, stack.AuthOutputDetails{}, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.Agent)
	resp, err := c.SecureHttpClient.Do(req)
	if err != nil {
		api.Logger.Errorf("[tenant:%s][region:%s] err:%+v", c.ExternalID, c.Region, err)
		return c.updateToken(c.ExternalID, stack.AuthOutputDetails{}, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return c.updateToken(c.ExternalID, stack.AuthOutputDetails{}, fmt.Errorf(string(body)))
	}
	var output stack.AuthOutput
	if err = json.Unmarshal(body, &output); err != nil {
		return c.updateToken(c.ExternalID, output.Response, err)
	}
	if e2 := output.Failed(); e2 != nil {
		return c.updateToken(c.ExternalID, output.Response, err)
	}

	return c.updateToken(c.ExternalID, output.Response, nil)
}

// this API should be called with c.CloudRulestackAdminMutex taken
func (c *Client) updateToken(externalID string, resp stack.AuthOutputDetails, err error) error {
	if err == nil {
		c.CloudRulestackAdminJwtExpTime = time.Now().Add(30 * time.Minute)
		c.CloudRulestackAdminJwt = resp.TokenId
		c.CloudRulestackSubscriptionKey = resp.SubscriptionKey
		if err := c.SetTenantVersion(resp.TokenId); err != nil {
			return err
		}
		api.Logger.Debugf("[tenant:%s][region:%s] jwt:%s", externalID, c.Region, c.CloudRulestackAdminJwt)
	} else {
		c.CloudRulestackAdminJwtExpTime = time.Now().Add(-time.Minute * 30)
		c.CloudRulestackAdminJwt = ""
		c.CloudRulestackSubscriptionKey = ""
		api.Logger.Errorf("[tenant:%s][region:%s] no jwt retrieved, err:%+v", externalID, c.Region, err)
	}
	return err
}

// externalID based auth over mTLS
func (c *Client) SetupUsingCreds(ctx context.Context, info AuthInfo) error {
	c.HttpClient = info.HttpClient
	c.SecureHttpClient = info.SecureHttpClient
	c.apiPrefix = info.RegionURL
	c.v2ApiPrefix = info.RegionV2URL
	c.ExternalID = info.ExternalID
	c.Region = info.Region
	c.AuthURL = info.AuthURL
	return nil
}

// RefreshJwts refreshes all JWTs and stores them for future API calls.
func (c *Client) RefreshCloudRulestackAdminJwt(ctx context.Context) error {
	c.CloudRulestackAdminMutex.Lock()
	defer c.CloudRulestackAdminMutex.Unlock()

	if c.CloudRulestackAdminJwtExpTime.Sub(time.Now()) > 60*time.Second {
		// the jwt is valid for 60 or more seconds. let's not replenish it now
		api.Logger.Debugf("[tenant:%s][region:%s] using existing token:%+v",
			c.ExternalID, c.Region, c.CloudRulestackAdminJwt)
		return nil
	}
	info := AuthInfo{
		ExternalID:       c.ExternalID,
		HttpClient:       c.HttpClient,
		SecureHttpClient: c.SecureHttpClient,
		Region:           c.Region,
		RegionURL:        c.apiPrefix,
		AuthURL:          c.AuthURL,
	}
	c.doAuth(ctx, info)
	return nil
}

func (c *Client) SetTenantVersion(tokenStr string) error {
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		api.Logger.Errorf("[tenant:%s][region:%s] failed to parse token: %v", c.ExternalID, c.Region, err)
		return err
	}
	api.Logger.Debug("setting tenant version")
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("Failed to parse token claims")
		api.Logger.Errorf("[tenant:%s][region:%s] failed to parse token claims: %v", c.ExternalID, c.Region, err)
	}
	if tenantVersion, ok := claims["tenant_version"]; !ok {
		return fmt.Errorf("tenant_version claim not found in token")
		api.Logger.Errorf("[tenant:%s][region:%s] tenant_version claim not found in token", c.ExternalID, c.Region)
	} else {
		c.TenantVersion = tenantVersion.(string)
		api.Logger.Errorf("[tenant:%s][region:%s] tenant version:%s", c.ExternalID, c.Region, c.TenantVersion)
	}
	api.Logger.Debugf("[tenant:%s][region:%s] set tenant version to: %s", c.ExternalID, c.Region, c.TenantVersion)
	if c.TenantVersion == TenantVersionV1 && c.Origin == OriginPA {
		api.Logger.Errorf("[tenant:%s][region:%s] unsupported provider version, please use provider version 2.0.20 or below", c.ExternalID, c.Region)
		return fmt.Errorf("unsupported provider version, please use provider version 2.0.20 or below")
	}
	return nil
}
