package appid

import (
	"context"
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/permissions"
)

// Client is a client for this collection.
type Client struct {
	client api.Client
}

// NewClient returns a new client for this collection.
func NewClient(client api.Client) *Client {
	return &Client{client: client}
}

// List returns a list of objects.
func (c *Client) List(ctx context.Context, input ListInput) (ListOutput, error) {
	c.client.Log(http.MethodGet, "list app-id versions")

	var ans ListOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "appidversions"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// ReadAppId returns information on the given app-id version.
func (c *Client) Read(ctx context.Context, input ReadInput) (ReadOutput, error) {
	c.client.Log(http.MethodGet, "describe app-id version: %s", input.Version)

	var ans ReadOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
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
func (c *Client) ReadApplication(ctx context.Context, version, app string) (ReadApplicationOutput, error) {
	c.client.Log(http.MethodGet, "describe app-id %q application: %s", version, app)

	var ans ReadApplicationOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "appidversions", version, "appids", app},
		nil,
		nil,
		&ans,
	)

	return ans, err
}
