package security

import (
	"context"
	"net/http"
	"strconv"

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

	stack, rlist := input.Rulestack, input.RuleList
	c.client.Log(http.MethodGet, "list %s %q security rules", rlist, stack)

	var ans ListOutput
	_, err := c.client.Communicate(
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
func (c *Client) Create(ctx context.Context, input Info) error {
	perm, permErr := permissions.Choose(input.Scope)
	if permErr != nil {
		return permErr
	}

	stack, rlist := input.Rulestack, input.RuleList
	input.Rulestack, input.RuleList = "", ""

	c.client.Log(http.MethodPost, "create %s security rule in %q: %s", rlist, stack, input.Entry.Name)

	_, err := c.client.Communicate(
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
func (c *Client) Read(ctx context.Context, input ReadInput) (ReadOutput, error) {
	perm, permErr := permissions.Choose(input.Scope)
	if permErr != nil {
		return ReadOutput{}, permErr
	}

	stack, rlist, priority := input.Rulestack, input.RuleList, input.Priority

	c.client.Log(http.MethodGet, "describe %s security rule in %q: %d", rlist, stack, priority)

	var ans ReadOutput
	_, err := c.client.Communicate(
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
func (c *Client) Update(ctx context.Context, input Info) error {
	perm, permErr := permissions.Choose(input.Scope)
	if permErr != nil {
		return permErr
	}

	stack, rlist, priority := input.Rulestack, input.RuleList, input.Priority
	input.Rulestack, input.RuleList, input.Priority = "", "", 0

	c.client.Log(http.MethodPut, "updating %s security rule in %q: priority %d", rlist, stack, priority)

	_, err := c.client.Communicate(
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
func (c *Client) Delete(ctx context.Context, input DeleteInput) error {
	perm, permErr := permissions.Choose(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.client.Log(http.MethodDelete, "delete %s security rule in %q: priority %d", input.RuleList, input.Rulestack, input.Priority)

	_, err := c.client.Communicate(
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
