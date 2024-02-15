package aws

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/logprofile"
	"context"
	"net/http"
)

// Read returns information on the given object.
func (c *Client) ReadFirewallLogprofile(ctx context.Context, input logprofile.ReadInput) (logprofile.ReadOutput, error) {
	name := input.Firewall
	c.Log(http.MethodGet, "describe firewall log profile: %s", name)

	var ans logprofile.ReadOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodGet,
		[]string{"v1", "config", "ngfirewalls", name, "logprofile"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func (c *Client) UpdateFirewallLogprofile(ctx context.Context, input logprofile.Info) error {
	name := input.Firewall
	input.Firewall = ""

	c.Log(http.MethodPut, "updating firewall log profile: %s", name)

	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPut,
		[]string{"v1", "config", "ngfirewalls", name, "logprofile"},
		nil,
		input,
		nil,
	)

	return err
}
