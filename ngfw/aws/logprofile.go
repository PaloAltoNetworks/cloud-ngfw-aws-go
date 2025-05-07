package aws

import (
	"context"
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/logprofile"
)

// Read returns information on the given object.
func (c *Client) ReadFirewallLogprofile(ctx context.Context, input logprofile.ReadInput) (logprofile.ReadOutput, error) {
	firewallId := input.FirewallId
	name := input.Firewall
	c.Log(http.MethodGet, "describe firewall log profile: %s", firewallId)
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", name, "logprofile"},
		V2Path: []string{"v2", "config", "ngfirewalls", firewallId, "logprofile"},
	}

	var ans logprofile.ReadOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodGet,
		path,
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func (c *Client) UpdateFirewallLogprofile(ctx context.Context, input logprofile.Info) error {
	firewallId := input.FirewallId
	name := input.Firewall
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", name, "logprofile"},
		V2Path: []string{"v2", "config", "ngfirewalls", firewallId, "logprofile"},
	}
	c.Log(http.MethodPost, "updating firewall log profile: %s", firewallId)

	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPost,
		path,
		nil,
		input,
		nil,
	)

	return err
}
