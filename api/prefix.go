package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/prefix"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) ListPrefixList(ctx context.Context, a prefix.ListInput) (prefix.ListOutput, error) {
	out, err := c.client.ListPrefixList(ctx, a)
	if err != nil {
		return prefix.ListOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ReadPrefixList(ctx context.Context, a prefix.ReadInput) (prefix.ReadOutput, error) {
	out, err := c.client.ReadPrefixList(ctx, a)
	if err != nil {
		return prefix.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) CreatePrefixList(ctx context.Context, f prefix.Info) error {
	Logger.Infof("prefix:%+v", f)
	if err := c.client.CreatePrefixList(ctx, f); err != nil {
		Logger.Errorf("err:%+v", err)
		return err
	}
	return nil
}

func (c *ApiClient) UpdatePrefixList(ctx context.Context, f prefix.Info) error {
	if err := c.client.UpdatePrefixList(ctx, f); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) DeletePrefixList(ctx context.Context, f prefix.DeleteInput) error {
	input := prefix.DeleteInput{
		Rulestack: f.Rulestack,
		Scope:     f.Scope,
		Name:      f.Name,
	}
	if err := c.client.DeletePrefixList(ctx, input); err != nil {
		return err
	}
	return nil
}
