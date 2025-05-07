package aws

import (
	"context"
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/feed"
)

// List returns a list of objects.
func (c *Client) ListFeed(ctx context.Context, input feed.ListInput) (feed.ListOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return feed.ListOutput{}, permErr
	}
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "feeds"},
	}
	c.Log(http.MethodGet, "list rulestack %q intelligent feeds", input.Rulestack)

	var ans feed.ListOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		path,
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
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "feeds"},
	}
	c.Log(http.MethodPost, "create rulestack %q intelligent feed: %s", input.Rulestack, input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		path,
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
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "feeds", input.Name},
	}
	c.Log(http.MethodGet, "describe rulestack %q intelligent feed: %s", input.Rulestack, input.Name)

	var ans feed.ReadOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		path,
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
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "feeds", name},
	}
	c.Log(http.MethodPut, "updating rulestack %q intelligent feed: %s", input.Rulestack, name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPut,
		path,
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
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "feeds", input.Name},
	}
	c.Log(http.MethodDelete, "delete rulestack %q intelligent feed: %s", input.Rulestack, input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodDelete,
		path,
		nil,
		nil,
		nil,
	)

	return err
}
