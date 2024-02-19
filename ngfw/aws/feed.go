package aws

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/feed"
	"context"
	"net/http"
)

// List returns a list of objects.
func (c *Client) ListFeed(ctx context.Context, input feed.ListInput) (feed.ListOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return feed.ListOutput{}, permErr
	}

	c.Log(http.MethodGet, "list rulestack %q intelligent feeds", input.Rulestack)

	var ans feed.ListOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "feeds"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) CreateFeed(ctx context.Context, input feed.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodPost, "create rulestack %q intelligent feed: %s", input.Rulestack, input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "feeds"},
		nil,
		input,
		nil,
	)

	return err
}

// Read returns information on the given object.
func (c *Client) ReadFeed(ctx context.Context, input feed.ReadInput) (feed.ReadOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return feed.ReadOutput{}, permErr
	}

	c.Log(http.MethodGet, "describe rulestack %q intelligent feed: %s", input.Rulestack, input.Name)

	var ans feed.ReadOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "feeds", input.Name},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func (c *Client) UpdateFeed(ctx context.Context, input feed.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	name := input.Name
	input.Name = ""

	c.Log(http.MethodPut, "updating rulestack %q intelligent feed: %s", input.Rulestack, name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPut,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "feeds", name},
		nil,
		input,
		nil,
	)

	return err
}

// Delete removes the given object from the config.
func (c *Client) DeleteFeed(ctx context.Context, input feed.DeleteInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodDelete, "delete rulestack %q intelligent feed: %s", input.Rulestack, input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "feeds", input.Name},
		nil,
		nil,
		nil,
	)

	return err
}
