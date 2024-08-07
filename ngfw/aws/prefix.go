package aws

import (
	"context"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/prefix"
	"net/http"
)

// List returns a list of objects.
func (c *Client) ListPrefixList(ctx context.Context, input prefix.ListInput) (prefix.ListOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return prefix.ListOutput{}, permErr
	}

	c.Log(http.MethodGet, "list rulestack %q prefix lists", input.Rulestack)

	var ans prefix.ListOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "prefixlists"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) CreatePrefixList(ctx context.Context, input prefix.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodPost, "create rulestack %q prefix list: %s", input.Rulestack, input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "prefixlists"},
		nil,
		input,
		nil,
	)

	return err
}

// Read returns information on the given object.
func (c *Client) ReadPrefixList(ctx context.Context, input prefix.ReadInput) (prefix.ReadOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return prefix.ReadOutput{}, permErr
	}

	c.Log(http.MethodGet, "describe rulestack %q prefix list: %s", input.Rulestack, input.Name)

	var ans prefix.ReadOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "prefixlists", input.Name},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func (c *Client) UpdatePrefixList(ctx context.Context, input prefix.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	name := input.Name
	input.Name = ""

	c.Log(http.MethodPut, "updating rulestack %q prefix list: %s", input.Rulestack, name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPut,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "prefixlists", name},
		nil,
		input,
		nil,
	)

	return err
}

// Delete removes the given object from the config.
func (c *Client) DeletePrefixList(ctx context.Context, input prefix.DeleteInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodDelete, "delete rulestack %q prefix list: %s", input.Rulestack, input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "prefixlists", input.Name},
		nil,
		nil,
		nil,
	)

	return err
}
