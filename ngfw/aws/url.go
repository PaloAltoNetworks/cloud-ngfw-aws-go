package aws

import (
	"context"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/url"
	"net/http"
)

// List returns a list of objects.
func (c *Client) ListUrlCustomCategory(ctx context.Context, input url.ListInput) (url.ListOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return url.ListOutput{}, permErr
	}

	c.Log(http.MethodGet, "list rulestack %q custom url categories", input.Rulestack)

	var ans url.ListOutput
	_, err := c.Communicate(
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
func (c *Client) CreateUrlCustomCategory(ctx context.Context, input url.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodPost, "create rulestack %q custom url category: %s", input.Rulestack, input.Name)

	_, err := c.Communicate(
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
func (c *Client) ReadUrlCustomCategory(ctx context.Context, input url.ReadInput) (url.ReadOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return url.ReadOutput{}, permErr
	}

	c.Log(http.MethodGet, "describe rulestack %q custom url category: %s", input.Rulestack, input.Name)

	var ans url.ReadOutput
	_, err := c.Communicate(
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
func (c *Client) UpdateUrlCustomCategory(ctx context.Context, input url.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	name := input.Name
	input.Name = ""

	c.Log(http.MethodPut, "updating rulestack %q custom url category: %s", input.Rulestack, name)

	_, err := c.Communicate(
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
func (c *Client) DeleteUrlCustomCategory(ctx context.Context, input url.DeleteInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodDelete, "delete rulestack %q custom url category: %s", input.Rulestack, input.Name)

	_, err := c.Communicate(
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
