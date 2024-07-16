package aws

import (
	"context"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/security"
	"net/http"
	"strconv"
)

// List returns a list of objects.
func (c *Client) ListSecurityRule(ctx context.Context, input security.ListInput) (security.ListOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return security.ListOutput{}, permErr
	}

	stack, rlist := input.Rulestack, input.RuleList
	c.Log(http.MethodGet, "list %s %q security rules", rlist, stack)

	var ans security.ListOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", stack, "rulelists", rlist},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) CreateSecurityRule(ctx context.Context, input security.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	stack, rlist := input.Rulestack, input.RuleList
	input.Rulestack, input.RuleList = "", ""

	c.Log(http.MethodPost, "create %s security rule in %q: %s", rlist, stack, input.Entry.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", stack, "rulelists", rlist},
		nil,
		input,
		nil,
	)

	return err
}

// Read returns information on the given object.
func (c *Client) ReadSecurityRule(ctx context.Context, input security.ReadInput) (security.ReadOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return security.ReadOutput{}, permErr
	}

	stack, rlist, priority := input.Rulestack, input.RuleList, input.Priority

	c.Log(http.MethodGet, "describe %s security rule in %q: %d", rlist, stack, priority)

	var ans security.ReadOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", stack, "rulelists", rlist, "priorities", strconv.Itoa(priority)},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Update updates the given object.
func (c *Client) UpdateSecurityRule(ctx context.Context, input security.Info) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	stack, rlist, priority := input.Rulestack, input.RuleList, input.Priority
	input.Rulestack, input.RuleList, input.Priority = "", "", 0

	c.Log(http.MethodPut, "updating %s security rule in %q: priority %d", rlist, stack, priority)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPut,
		[]string{"v1", "config", "rulestacks", stack, "rulelists", rlist, "priorities", strconv.Itoa(priority)},
		nil,
		input,
		nil,
	)

	return err
}

// Delete removes the given object from the config.
func (c *Client) DeleteSecurityRule(ctx context.Context, input security.DeleteInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodDelete, "delete %s security rule in %q: priority %d", input.RuleList, input.Rulestack, input.Priority)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "rulelists", input.RuleList, "priorities", strconv.Itoa(input.Priority)},
		nil,
		nil,
		nil,
	)

	return err
}
