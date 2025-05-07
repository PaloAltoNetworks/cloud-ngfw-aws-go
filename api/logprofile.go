package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/logprofile"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) ReadFirewallLogProfile(ctx context.Context, f logprofile.ReadInput) (logprofile.ReadOutput, error) {
	out, err := c.client.ReadFirewallLogprofile(ctx, f)
	if err != nil {
		return logprofile.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) UpdateFirewallLogProfile(ctx context.Context, input logprofile.Info) error {
	if err := c.client.UpdateFirewallLogprofile(ctx, input); err != nil {
		return err
	}
	return nil
}
