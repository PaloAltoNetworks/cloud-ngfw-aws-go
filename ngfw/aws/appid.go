package aws

import (
	"context"
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/appid"
)

// List returns a list of objects.
func (c *Client) ListAppID(ctx context.Context, input appid.ListInput) (appid.ListOutput, error) {
	c.Log(http.MethodGet, "list app-id versions")

	var ans appid.ListOutput
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodGet,
		[]string{"v1", "config", "appidversions"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// ReadAppId returns information on the given app-id version.
func (c *Client) ReadAppID(ctx context.Context, input appid.ReadInput) (appid.ReadOutput, error) {
	c.Log(http.MethodGet, "describe app-id version: %s", input.Version)

	var ans appid.ReadOutput
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodGet,
		[]string{"v1", "config", "appidversions", input.Version},
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

	var ans appid.ReadApplicationOutput
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodGet,
		[]string{"v1", "config", "appidversions", version, "appids", app},
		nil,
		nil,
		&ans,
	)

	return ans, err
}
