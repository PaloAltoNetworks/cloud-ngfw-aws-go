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
	c.client.Log(http.MethodGet, "list predefined url categories")

	var ans ListOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "urlcategories"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// ListOverrides returns URL categories with overrides specified.
func (c *Client) ListOverrides(ctx context.Context, input ListOverridesInput) (ListOverridesOutput, error) {
	c.client.Log(http.MethodGet, "list predefined url category overrides for rulestack %q", input.Rulestack)

	var ans ListOverridesOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "urlfilteringprofiles", "custom", "urlcategories"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// GetOverride returns the URL category override info.
func (c *Client) GetOverride(ctx context.Context, input GetOverrideInput) (GetOverrideOutput, error) {
	c.client.Log(http.MethodGet, "get %q predefined url category override: %s", input.Rulestack, input.Name)

	var ans GetOverrideOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "urlfilteringprofiles", "custom", "urlcategories", input.Name},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Override specifies an override for a predefined URL category.
func (c *Client) Override(ctx context.Context, input OverrideInput) error {
	c.client.Log(http.MethodPut, "override %q predefined url category: %s", input.Rulestack, input.Name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodPut,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "urlfilteringprofiles", "custom", "urlcategories", input.Name, "action"},
		nil,
		input,
		nil,
	)

	return err
}
