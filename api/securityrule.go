package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/security"
)

func (c *ApiClient) CreateSecurityRule(ctx context.Context, f security.Info) error {
	if err := c.client.CreateSecurityRule(ctx, f); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) ReadSecurityRule(ctx context.Context, f security.ReadInput) (security.ReadOutput, error) {
	out, err := c.client.ReadSecurityRule(ctx, f)
	if err != nil {
		return security.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ListSecurityRule(ctx context.Context, f security.ListInput) (security.ListOutput, error) {
	out, err := c.client.ListSecurityRule(ctx, f)
	if err != nil {
		return security.ListOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) UpdateSecurityRule(ctx context.Context, f security.Info) error {
	//FIXME - pass context from app itself
	// ctx := context.Background()
	if err := c.client.UpdateSecurityRule(ctx, f); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) DeleteSecurityRule(ctx context.Context, f security.DeleteInput) error {
	//FIXME - pass context from app itself
	// ctx := context.Background()
	input := security.DeleteInput{
		Rulestack: f.Rulestack,
		RuleList:  f.RuleList,
		Scope:     f.Scope,
		Priority:  f.Priority,
	}
	if err := c.client.DeleteSecurityRule(ctx, input); err != nil {
		return err
	}
	return nil
}
