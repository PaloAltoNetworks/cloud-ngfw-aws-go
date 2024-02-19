package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/url"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) ListUrlCustomCategory(ctx context.Context, a url.ListInput) (url.ListOutput, error) {
	out, err := c.client.ListUrlCustomCategory(ctx, a)
	if err != nil {
		return url.ListOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ReadUrlCustomCategory(ctx context.Context, a url.ReadInput) (url.ReadOutput, error) {
	out, err := c.client.ReadUrlCustomCategory(ctx, a)
	if err != nil {
		return url.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) CreateUrlCustomCategory(ctx context.Context, f url.Info) error {
	if err := c.client.CreateUrlCustomCategory(ctx, f); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) UpdateUrlCustomCategory(ctx context.Context, f url.Info) error {
	if err := c.client.UpdateUrlCustomCategory(ctx, f); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) DeleteUrlCustomCategory(ctx context.Context, f url.DeleteInput) error {
	input := url.DeleteInput{
		Rulestack: f.Rulestack,
		Scope:     f.Scope,
		Name:      f.Name,
	}
	if err := c.client.DeleteUrlCustomCategory(ctx, input); err != nil {
		return err
	}
	return nil
}
