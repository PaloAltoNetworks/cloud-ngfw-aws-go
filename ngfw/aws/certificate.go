package aws

import (
	"context"
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/certificate"
)

// ListCertificate returns a certificate.List of objects.
func (c *Client) ListCertificate(ctx context.Context, input certificate.ListInput) (certificate.ListOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return certificate.ListOutput{}, permErr
	}
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "certificates"},
	}
	c.Log(http.MethodGet, "certificate.List rulestack %q certificate objects", input.Rulestack)

	var ans certificate.ListOutput
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
func (c *Client) CreateCertificate(ctx context.Context, input certificate.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "certificates"},
	}
	c.Log(http.MethodPost, "create rulestack %q certificate object: %s", input.Rulestack, input.Name)

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
func (c *Client) ReadCertificate(ctx context.Context, input certificate.ReadInput) (certificate.ReadOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return certificate.ReadOutput{}, permErr
	}
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "certificates", input.Name},
	}
	c.Log(http.MethodGet, "describe rulestack %q certificate object: %s", input.Rulestack, input.Name)

	var ans certificate.ReadOutput
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
func (c *Client) UpdateCertificate(ctx context.Context, input certificate.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "certificates", input.Name},
	}
	name := input.Name
	input.Name = ""

	c.Log(http.MethodPut, "updating rulestack %q certificate object: %s", input.Rulestack, name)

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
func (c *Client) DeleteCertificate(ctx context.Context, input certificate.DeleteInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "certificates", input.Name},
	}
	c.Log(http.MethodDelete, "delete rulestack %q certificate object: %s", input.Rulestack, input.Name)

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
