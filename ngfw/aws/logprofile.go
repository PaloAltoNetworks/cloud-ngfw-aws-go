package aws

import (
	"context"
	"net/http"
	"net/url"

	cloudngfwgosdk "github.com/paloaltonetworks/cloud-ngfw-aws-go/v2"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/logprofile"
)

// Read returns information on the given object.
func (c *Client) ReadFirewallLogprofile(ctx context.Context, input logprofile.ReadInput) (logprofile.ReadOutput, error) {
	name := input.Firewall
	fwId := input.FirewallId
	uv := url.Values{}
	schemaVersion := ctx.Value("SchemaVersion").(string)
	if schemaVersion == cloudngfwgosdk.SchemaVersionV1 {
		fwId = name
		uv = url.Values{
			"v1route":   []string{"true"},
			"accountid": []string{input.AccountId},
		}
	}
	c.Log(http.MethodGet, "describe firewall log profile: %s", fwId)
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", fwId, "logprofile"},
		V2Path: []string{"v2", "config", "ngfirewalls", fwId, "logprofile"},
	}
	var ans logprofile.ReadOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodGet,
		path,
		uv,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func (c *Client) UpdateFirewallLogprofile(ctx context.Context, input logprofile.Info) error {
	fwId := input.FirewallId
	name := input.Firewall
	uv := url.Values{}
	schemaVersion := ctx.Value("SchemaVersion").(string)
	method := http.MethodPost
	if schemaVersion == cloudngfwgosdk.SchemaVersionV1 {
		fwId = name
		uv = url.Values{
			"v1route": []string{"true"},
		}
		method = http.MethodPut
	}
	c.Log(http.MethodGet, "describe firewall log profile: %s", fwId)
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", fwId, "logprofile"},
		V2Path: []string{"v2", "config", "ngfirewalls", fwId, "logprofile"},
	}
	c.Log(http.MethodPost, "updating firewall log profile: %s", fwId)

	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		method,
		path,
		uv,
		input,
		nil,
	)

	return err
}
