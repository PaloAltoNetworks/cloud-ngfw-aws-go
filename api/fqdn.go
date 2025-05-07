package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/fqdn"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) ListFqdn(ctx context.Context, a fqdn.ListInput) (fqdn.ListOutput, error) {
	out, err := c.client.ListFqdn(ctx, a)
	if err != nil {
		return fqdn.ListOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ReadFqdn(ctx context.Context, f fqdn.ReadInput) (fqdn.ReadOutput, error) {
	out, err := c.client.ReadFqdn(ctx, f)
	if err != nil {
		return fqdn.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) CreateFqdn(ctx context.Context, f fqdn.Info) error {
	if err := c.client.CreateFqdn(ctx, f); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) UpdateFqdn(ctx context.Context, f fqdn.Info) error {
	if err := c.client.UpdateFqdn(ctx, f); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) DeleteFqdn(ctx context.Context, f fqdn.DeleteInput) error {
	if err := c.client.DeleteFqdn(ctx, f); err != nil {
		return err
	}
	return nil
}
