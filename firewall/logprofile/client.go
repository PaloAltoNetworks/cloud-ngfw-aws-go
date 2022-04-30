package logprofile

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

// Read returns information on the given object.
func (c *Client) Read(ctx context.Context, input ReadInput) (ReadOutput, error) {
	name := input.Firewall
	c.client.Log(http.MethodGet, "describe firewall log profile: %s", name)

	var ans ReadOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodGet,
		[]string{"v1", "config", "ngfirewalls", name, "logprofile"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func (c *Client) Update(ctx context.Context, input Info) error {
	name := input.Firewall
	input.Firewall = ""

	c.client.Log(http.MethodPut, "updating firewall log profile: %s", name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodPut,
		[]string{"v1", "config", "ngfirewalls", name, "logprofile"},
		nil,
		input,
		nil,
	)

	return err
}
