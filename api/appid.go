package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/appid"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) ReadAppID(ctx context.Context, a appid.ReadInput) (appid.ReadOutput, error) {
	out, err := c.client.ReadAppID(ctx, a)
	if err != nil {
		return appid.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ListAppID(ctx context.Context, a appid.ListInput) (appid.ListOutput, error) {
	out, err := c.client.ListAppID(ctx, a)
	if err != nil {
		return appid.ListOutput{}, err
	}
	return out, nil
}
