package aws

import (
	"context"
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/url"
)

// List returns a list of objects.
func (c *Client) ListUrlCustomCategory(ctx context.Context, input url.ListInput) (url.ListOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return url.ListOutput{}, permErr
	}

	c.Log(http.MethodGet, "list rulestack %q custom url categories", input.Rulestack)

	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "urlcustomcategories"},
	}
	var ans url.ListOutput
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
func (c *Client) CreateUrlCustomCategory(ctx context.Context, input url.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodPost, "create rulestack %q custom url category: %s", input.Rulestack, input.Name)

	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "urlcustomcategories"},
	}
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
func (c *Client) ReadUrlCustomCategory(ctx context.Context, input url.ReadInput) (url.ReadOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return url.ReadOutput{}, permErr
	}
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "urlcustomcategories", input.Name},
	}
	c.Log(http.MethodGet, "describe rulestack %q custom url category: %s", input.Rulestack, input.Name)

	var ans url.ReadOutput
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
func (c *Client) UpdateUrlCustomCategory(ctx context.Context, input url.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	name := input.Name
	input.Name = ""

	c.Log(http.MethodPut, "updating rulestack %q custom url category: %s", input.Rulestack, name)
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "urlcustomcategories", name},
	}

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
func (c *Client) DeleteUrlCustomCategory(ctx context.Context, input url.DeleteInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodDelete, "delete rulestack %q custom url category: %s", input.Rulestack, input.Name)
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "urlcustomcategories", input.Name},
	}

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
