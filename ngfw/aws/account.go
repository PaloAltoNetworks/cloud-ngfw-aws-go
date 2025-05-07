package aws

import (
	"context"
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/account"
)

// Create creates an object.
func (c *Client) CreateAccount(ctx context.Context, input account.CreateInput) (account.CreateOutput, error) {
	c.Log(http.MethodPost, "create account")
	path := Path{
		V1Path: []string{"v1", "mgmt", "linkaccounts"},
	}
	var ans account.CreateOutput
	_, err := c.Communicate(
		ctx,
		PermissionAccount,
		http.MethodPost,
		path,
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Read returns information on the given object.
func (c *Client) ReadAccount(ctx context.Context, input account.ReadInput) (account.ReadOutput, error) {
	accountId := input.AccountId
	c.Log(http.MethodGet, "describe account: %s", accountId)
	path := Path{
		V1Path: []string{"v1", "mgmt", "linkaccounts", input.AccountId},
	}
	var ans account.ReadOutput
	_, err := c.Communicate(
		ctx,
		PermissionAccount,
		http.MethodGet,
		path,
		nil,
		nil,
		&ans,
	)

	return ans, err
}

// List returns a list of given objects.
func (c *Client) ListAccounts(ctx context.Context, input account.ListInput) (account.ListOutput, error) {
	c.Log(http.MethodGet, "list accounts")
	path := Path{
		V1Path: []string{"v1", "mgmt", "linkaccounts"},
	}
	var ans account.ListOutput
	_, err := c.Communicate(
		ctx,
		PermissionAccount,
		http.MethodGet,
		path,
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Delete the given account.
func (c *Client) DeleteAccount(ctx context.Context, input account.DeleteInput) error {
	c.Log(http.MethodDelete, "delete account: %s", input.AccountId)
	path := Path{
		V1Path: []string{"v1", "mgmt", "linkaccounts", input.AccountId},
	}
	_, err := c.Communicate(
		ctx,
		PermissionAccount,
		http.MethodDelete,
		path,
		nil,
		nil,
		nil,
	)

	return err
}
