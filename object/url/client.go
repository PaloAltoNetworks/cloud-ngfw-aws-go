package url

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
	perm, permErr := permissions.Choose(input.Scope)
	if permErr != nil {
		return ListOutput{}, permErr
	}

	c.client.Log(http.MethodGet, "list rulestack %q custom url categories", input.Rulestack)

	var ans ListOutput
	_, err := c.client.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "urlcustomcategories"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) Create(ctx context.Context, input Info) error {
	perm, permErr := permissions.Choose(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.client.Log(http.MethodPost, "create rulestack %q custom url category: %s", input.Rulestack, input.Name)

	_, err := c.client.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "urlcustomcategories"},
		nil,
		input,
		nil,
	)

	return err
}

// Read returns information on the given object.
func (c *Client) Read(ctx context.Context, input ReadInput) (ReadOutput, error) {
	perm, permErr := permissions.Choose(input.Scope)
	if permErr != nil {
		return ReadOutput{}, permErr
	}

	c.client.Log(http.MethodGet, "describe rulestack %q custom url category: %s", input.Rulestack, input.Name)

	var ans ReadOutput
	_, err := c.client.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "urlcustomcategories", input.Name},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func (c *Client) Update(ctx context.Context, input Info) error {
	perm, permErr := permissions.Choose(input.Scope)
	if permErr != nil {
		return permErr
	}

	name := input.Name
	input.Name = ""

	c.client.Log(http.MethodPut, "updating rulestack %q custom url category: %s", input.Rulestack, name)

	_, err := c.client.Communicate(
		ctx,
		perm,
		http.MethodPut,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "urlcustomcategories", name},
		nil,
		input,
		nil,
	)

	return err
}

// Delete removes the given object from the config.
func (c *Client) Delete(ctx context.Context, input DeleteInput) error {
	perm, permErr := permissions.Choose(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.client.Log(http.MethodDelete, "delete rulestack %q custom url category: %s", input.Rulestack, input.Name)

	_, err := c.client.Communicate(
		ctx,
		perm,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "urlcustomcategories", input.Name},
		nil,
		nil,
		nil,
	)

	return err
}
