package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/account"
)

func (c *ApiClient) CreateAccount(ctx context.Context, input account.CreateInput) (account.CreateOutput, error) {
	out, err := c.client.CreateAccount(ctx, input)
	if err != nil {
		return account.CreateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ReadAccount(ctx context.Context, input account.ReadInput) (account.ReadOutput, error) {
	out, err := c.client.ReadAccount(ctx, input)
	if err != nil {
		return account.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ListAccounts(ctx context.Context, a account.ListInput) (account.ListOutput, error) {
	out, err := c.client.ListAccounts(ctx, a)
	if err != nil {
		return account.ListOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) DeleteAccount(ctx context.Context, input account.DeleteInput) error {
	err := c.client.DeleteAccount(ctx, input)
	if err != nil {
		return err
	}
	return nil
}
