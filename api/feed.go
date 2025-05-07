package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/feed"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) ReadFeed(ctx context.Context, f feed.ReadInput) (feed.ReadOutput, error) {
	out, err := c.client.ReadFeed(ctx, f)
	if err != nil {
		return feed.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) CreateFeed(ctx context.Context, f feed.Info) error {
	if err := c.client.CreateFeed(ctx, f); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) ListFeed(ctx context.Context, f feed.ListInput) (feed.ListOutput, error) {
	out, err := c.client.ListFeed(ctx, f)
	if err != nil {
		return feed.ListOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) UpdateFeed(ctx context.Context, f feed.Info) error {
	if err := c.client.UpdateFeed(ctx, f); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) DeleteFeed(ctx context.Context, f feed.DeleteInput) error {
	input := feed.DeleteInput{
		Rulestack: f.Rulestack,
		Scope:     f.Scope,
		Name:      f.Name,
	}
	if err := c.client.DeleteFeed(ctx, input); err != nil {
		return err
	}
	return nil
}
