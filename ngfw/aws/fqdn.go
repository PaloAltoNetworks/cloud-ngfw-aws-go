package aws

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/fqdn"
	"context"
	"net/http"
)

// ListFqdn returns a fqdn.List of objects.
func (c *Client) ListFqdn(ctx context.Context, input fqdn.ListInput) (fqdn.ListOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return fqdn.ListOutput{}, permErr
	}

	c.Log(http.MethodGet, "list rulestack %q fqdn fqdn.Lists", input.Rulestack)

	var ans fqdn.ListOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "fqdnlists"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) CreateFqdn(ctx context.Context, input fqdn.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodPost, "create rulestack %q fqdn fqdn.List: %s", input.Rulestack, input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "fqdnlists"},
		nil,
		input,
		nil,
	)

	return err
}

// Read returns information on the given object.
func (c *Client) ReadFqdn(ctx context.Context, input fqdn.ReadInput) (fqdn.ReadOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return fqdn.ReadOutput{}, permErr
	}

	c.Log(http.MethodGet, "describe rulestack %q fqdn fqdn.List: %s", input.Rulestack, input.Name)

	var ans fqdn.ReadOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "fqdnlists", input.Name},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func (c *Client) UpdateFqdn(ctx context.Context, input fqdn.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	name := input.Name
	input.Name = ""

	c.Log(http.MethodPut, "updating rulestack %q fqdn fqdn.List: %s", input.Rulestack, name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPut,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "fqdnlists", name},
		nil,
		input,
		nil,
	)

	return err
}

// Delete removes the given object from the config.
func (c *Client) DeleteFqdn(ctx context.Context, input fqdn.DeleteInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodDelete, "delete rulestack %q fqdn fqdn.List: %s", input.Rulestack, input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "fqdnlists", input.Name},
		nil,
		nil,
		nil,
	)

	return err
}
