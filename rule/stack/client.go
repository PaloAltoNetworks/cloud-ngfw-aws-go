package stack

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
	c.client.Log(http.MethodGet, "list rulestacks")

	var ans ListOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) Create(ctx context.Context, input Info) error {
	c.client.Log(http.MethodPost, "create rulestack: %s", input.Name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks"},
		nil,
		input,
		nil,
	)

	return err
}

// Read returns information on the given object.
func (c *Client) Read(ctx context.Context, input ReadInput) (ReadOutput, error) {
	name := input.Name

	c.client.Log(http.MethodGet, "describe rulestack: %s", name)

	var ans ReadOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", name},
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

	c.client.Log(http.MethodPut, "updating rulestack: %s", name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodPut,
		[]string{"v1", "config", "rulestacks", name},
		nil,
		input,
		nil,
	)

	return err
}

// Delete removes the given object from the config.
func (c *Client) Delete(ctx context.Context, name string) error {
	c.client.Log(http.MethodDelete, "delete rulestack: %s", name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", name},
		nil,
		nil,
		nil,
	)

	return err
}

// Commit commits the rulestack configuration.
func (c *Client) Commit(ctx context.Context, name string) error {
	c.client.Log(http.MethodPost, "commit rulestack: %s", name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", name, "commit"},
		nil,
		nil,
		nil,
	)

	return err
}

// CommitStatus gets the commit status.
func (c *Client) CommitStatus(ctx context.Context, name string) (CommitStatus, error) {
	c.client.Log(http.MethodGet, "commit status for rulestack: %s", name)

	var ans CommitStatus
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", name, "commit"},
		nil,
		nil,
		&ans,
	)

	return ans, err
}

// Revert reverts to the last good config.
func (c *Client) Revert(ctx context.Context, name string) error {
	c.client.Log(http.MethodPost, "revert rulestack: %s", name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", name, "revert"},
		nil,
		nil,
		nil,
	)

	return err
}

// Validate validates the rulestack config.
func (c *Client) Validate(ctx context.Context, name string) error {
	c.client.Log(http.MethodPost, "validate rulestack: %s", name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", name, "validate"},
		nil,
		nil,
		nil,
	)

	return err
}
