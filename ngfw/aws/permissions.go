package aws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	awsngfw "github.com/paloaltonetworks/cloud-ngfw-aws-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	PermissionFirewall        = "firewall"
	PermissionRulestack       = "rulestack"
	PermissionGlobalRulestack = "global rulestack"
	PermissionAccount         = "account admin"
	PermissionAccountAdminJWT = "account admin JWT"

	AuthTypeIAMRole    = "AuthTypeIAMRole"
	AuthTypeCognito    = "AuthTypeCognito"
	AuthTypeExternalID = "AuthTypeExternalID"
)

// Choose returns the correct JWT style for the given scope.
func GetPermission(v string) (string, error) {
	switch v {
	case "", LocalScope:
		return PermissionRulestack, nil
	case GlobalScope:
		return PermissionGlobalRulestack, nil
	}

	return "", fmt.Errorf("Unknown permission: %s", v)
}

// RefreshJwts refreshes all JWTs and stores them for future API calls.
func (c *Client) RefreshFirewallAdminJwt(ctx context.Context) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if c.Logging&awsngfw.LogLogin == awsngfw.LogLogin {
		log.Printf("(login) refreshing JWTs...")
	}
	c.FirewallAdminMutex.Lock()
	defer c.FirewallAdminMutex.Unlock()
	if c.FirewallAdminJwtExpTime.Sub(time.Now()) > 10*time.Second {
		// the jwt is valid for 10 or more seconds. let's not replenish it now
		return nil
	}

	jwtReq := getJwt{
		Expires: 120,
		/*
			KeyInfo: &jwtKeyInfo{
				Region: c.Region,
				Tenant: "XY",
			},
		*/
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
		Profile: *aws.String(c.Profile),
	})
	if err != nil {
		return err
	}

	svc := sts.New(sess)
	// Get firewall JWT.
	var rarn *string
	if c.LfaArn != "" {
		rarn = aws.String(c.LfaArn)
	} else if c.Arn != "" {
		rarn = aws.String(c.Arn)
	} else {
		log.Printf("err: No LFA role is assigned")
		return err
	}

	if c.Logging&awsngfw.LogLogin == awsngfw.LogLogin {
		log.Printf("(login) refreshing firewall JWT...")
	}
	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         rarn,
		RoleSessionName: aws.String("sdk_session"),
	})
	if err != nil {
		return err
	}

	var ans authResponse
	_, err = c.Communicate(
		ctx, "", http.MethodGet, []string{"v1", "mgmt", "tokens", "cloudfirewalladmin"}, nil, jwtReq, &ans, result.Credentials,
	)
	if err != nil {
		log.Printf("err:%+v", err)
		return err
	}

	tNow := time.Now()
	c.FirewallAdminJwtExpTime = tNow.Add(time.Duration(ans.Resp.ExpiryTime) * time.Minute)
	c.FirewallAdminJwt = ans.Resp.Jwt
	c.FirewallSubscriptionKey = ans.Resp.SubscriptionKey
	return nil
}

// RefreshJwts refreshes all JWTs and stores them for future API calls.
func (c *Client) RefreshRulestackAdminJwt(ctx context.Context) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if c.Logging&awsngfw.LogLogin == awsngfw.LogLogin {
		log.Printf("(login) refreshing RulestackAdmin JWT...")
	}

	c.RulestackAdminMutex.Lock()
	defer c.RulestackAdminMutex.Unlock()
	if time.Until(c.RulestackAdminJwtExpTime) > 10*time.Second {
		// the jwt is valid for 10 or more seconds. let's not replenish it now
		log.Printf("exptime:%+v now:%+v", c.RulestackAdminJwtExpTime, time.Now())
		return nil
	}

	jwtReq := getJwt{
		Expires: 120,
		/*
			KeyInfo: &jwtKeyInfo{
				Region: c.Region,
				Tenant: "XY",
			},
		*/
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
		Profile: c.Profile,
	})

	if err != nil {
		log.Printf("err:%+v", err)
		return err
	}

	svc := sts.New(sess)
	// Get rulestack JWT.
	var rarn *string
	if c.LraArn != "" {
		rarn = aws.String(c.LraArn)
	} else if c.Arn != "" {
		rarn = aws.String(c.Arn)
	} else {
		log.Printf("err: No LRA role is assigned")
		return err
	}
	if c.Logging&awsngfw.LogLogin == awsngfw.LogLogin {
		log.Printf("(login) refreshing RulestackAdmin JWT...")
	}
	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         rarn,
		RoleSessionName: aws.String("sdk_session"),
	})
	if err != nil {
		log.Printf("err:%+v", err)
		return err
	}

	var ans authResponse
	_, err = c.RequestJwt(
		ctx, http.MethodGet, []string{"v1", "mgmt", "tokens", "cloudrulestackadmin"}, nil, jwtReq, &ans, result.Credentials,
	)
	if err != nil {
		log.Printf("err:%+v", err)
		return err
	}

	tNow := time.Now()
	c.RulestackAdminJwtExpTime = tNow.Add(time.Duration(ans.Resp.ExpiryTime) * time.Minute)
	c.RulestackAdminJwt = ans.Resp.Jwt
	c.RulestackSubscriptionKey = ans.Resp.SubscriptionKey
	return nil
}

// RefreshJwts refreshes all JWTs and stores them for future API calls.
func (c *Client) RefreshGlobalRulestackAdminJwt(ctx context.Context) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if c.Logging&awsngfw.LogLogin == awsngfw.LogLogin {
		log.Printf("(login) refreshing JWTs...")
	}

	c.GlobalRulestackAdminMutex.Lock()
	defer c.GlobalRulestackAdminMutex.Unlock()
	if c.GlobalRulestackAdminJwtExpTime.Sub(time.Now()) > 10*time.Second {
		// the jwt is valid for 10 or more seconds. let's not replenish it now
		return nil
	}

	jwtReq := getJwt{
		Expires: 120,
		/*
			KeyInfo: &jwtKeyInfo{
				Region: c.Region,
				Tenant: "XY",
			},
		*/
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
		Profile: *aws.String(c.Profile),
	})

	if err != nil {
		return err
	}

	svc := sts.New(sess)

	// Get global rulestack JWT.
	var rarn *string
	if c.GraArn != "" {
		rarn = aws.String(c.GraArn)
	} else if c.Arn != "" {
		rarn = aws.String(c.Arn)
	} else {
		log.Printf("err: No GRA role is assigned")
		return err
	}

	if c.Logging&awsngfw.LogLogin == awsngfw.LogLogin {
		log.Printf("(login) refreshing global rulestack JWT...")
	}
	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         rarn,
		RoleSessionName: aws.String("sdk_session"),
	})
	if err != nil {
		return err
	}

	var ans authResponse
	_, err = c.Communicate(
		ctx, "", http.MethodGet, []string{"v1", "mgmt", "tokens", "cloudglobalrulestackadmin"}, nil, jwtReq, &ans, result.Credentials,
	)
	if err != nil {
		return err
	}
	tNow := time.Now()
	c.GlobalRulestackAdminJwtExpTime = tNow.Add(time.Duration(ans.Resp.ExpiryTime) * time.Minute)
	c.GlobalRulestackAdminJwt = ans.Resp.Jwt
	c.GlobalRulestackSubscriptionKey = ans.Resp.SubscriptionKey
	return nil
}

// RefreshJwts refreshes all JWTs and stores them for future API calls.
func (c *Client) RefreshAccountAdminJwt(ctx context.Context) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if c.Logging&awsngfw.LogLogin == awsngfw.LogLogin {
		log.Printf("(login) refreshing JWTs...")
	}
	c.AccountAdminMutex.Lock()
	defer c.AccountAdminMutex.Unlock()
	if c.AccountAdminJwtExpTime.Sub(time.Now()) > 10*time.Second {
		// the jwt is valid for 10 or more seconds. let's not replenish it now
		return nil
	}

	jwtReq := getJwt{
		Expires: 120,
		/*
			KeyInfo: &jwtKeyInfo{
				Region: c.Region,
				Tenant: "XY",
			},
		*/
	}

	var creds *credentials.Credentials
	if c.AccessKey != "" || c.SecretKey != "" {
		creds = credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: creds,
			Region:      aws.String(c.MPRegion),
		},
		Profile: *aws.String(c.Profile),
	})
	if err != nil {
		return err
	}

	svc := sts.New(sess)
	// Get account admin JWT.
	if c.AcctAdminArn == "" {
		return err
	}

	if c.Logging&awsngfw.LogLogin == awsngfw.LogLogin {
		log.Printf("(login) refreshing account admin JWT...")
	}
	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(c.AcctAdminArn),
		RoleSessionName: aws.String("sdk_session"),
	})
	if err != nil {
		return err
	}

	var ans authResponse
	_, err = c.Communicate(
		ctx, PermissionAccountAdminJWT, http.MethodGet, []string{"v1", "mgmt", "tokens", "cloudaccountadmin"}, nil, jwtReq, &ans, result.Credentials,
	)
	if err != nil {
		log.Printf("err:%+v", err)
		return err
	}

	tNow := time.Now()
	c.AccountAdminJwtExpTime = tNow.Add(time.Duration(ans.Resp.ExpiryTime) * time.Minute)
	c.AccountAdminJwt = ans.Resp.Jwt
	c.AccountAdminSubscriptionKey = ans.Resp.SubscriptionKey
	return nil
}
