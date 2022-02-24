package fqdn

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
	c.client.Log(http.MethodGet, "list rulestack %q pqdn lists", input.Rulestack)

	var ans ListOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "fqdnlists"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) Create(ctx context.Context, input Info) error {
	c.client.Log(http.MethodPost, "create rulestack %q fqdn list: %s", input.Rulestack, input.Name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "fqdnlists"},
		nil,
		input,
		nil,
	)

	return err
}

// Read returns information on the given object.
func (c *Client) Read(ctx context.Context, input ReadInput) (ReadOutput, error) {
	c.client.Log(http.MethodGet, "describe rulestack %q fqdn list: %s", input.Rulestack, input.Name)

	var ans ReadOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "fqdnlists", input.Name},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func (c *Client) Update(ctx context.Context, input Info) error {
	name := input.Name
	input.Name = ""

	c.client.Log(http.MethodPut, "updating rulestack %q fqdn list: %s", input.Rulestack, name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodPut,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "fqdnlists", name},
		nil,
		input,
		nil,
	)

	return err
}

// Delete removes the given object from the config.
func (c *Client) Delete(ctx context.Context, stack, name string) error {
	c.client.Log(http.MethodDelete, "delete rulestack %q fqdn list: %s", stack, name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", stack, "fqdnlists", name},
		nil,
		nil,
		nil,
	)

	return err
}
