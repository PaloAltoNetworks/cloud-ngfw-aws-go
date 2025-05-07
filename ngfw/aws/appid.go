package aws

import (
	"context"
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/appid"
)

// List returns a list of objects.
func (c *Client) ListAppID(ctx context.Context, input appid.ListInput) (appid.ListOutput, error) {
	c.Log(http.MethodGet, "list app-id versions")
	path := Path{
		V1Path: []string{"v1", "config", "appidversions"},
	}
	var ans appid.ListOutput
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodGet,
		path,
		nil,
		input,
		&ans,
	)

	return ans, err
}

// ReadAppId returns information on the given app-id version.
func (c *Client) ReadAppID(ctx context.Context, input appid.ReadInput) (appid.ReadOutput, error) {
	c.Log(http.MethodGet, "describe app-id version: %s", input.Version)
	path := Path{
		V1Path: []string{"v1", "config", "appidversions"},
	}
	var ans appid.ReadOutput
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodGet,
		path,
		nil,
		input,
		&ans,
	)

	return ans, err
}

// ReadApplication returns information on the given application in the specified
// app-id.
func (c *Client) ReadApplication(ctx context.Context, version, app string) (appid.ReadApplicationOutput, error) {
	c.Log(http.MethodGet, "describe app-id %q application: %s", version, app)
	path := Path{
		V1Path: []string{"v1", "config", "appidversions", version, "appids", app},
	}
	var ans appid.ReadApplicationOutput
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodGet,
		path,
		nil,
		nil,
		&ans,
	)

	return ans, err
}
