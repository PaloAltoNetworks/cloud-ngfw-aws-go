package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/predefinedurl"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) ListUrlPredefinedCategories(ctx context.Context, input predefinedurl.ListInput) (predefinedurl.ListOutput, error) {
	out, err := c.client.ListUrlPredefinedCategories(ctx, input)
	if err != nil {
		return predefinedurl.ListOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ListUrlCategoriesActionOverride(ctx context.Context, input predefinedurl.ListOverridesInput) (predefinedurl.ListOverridesOutput, error) {
	out, err := c.client.ListUrlCategoriesActionOverride(ctx, input)
	if err != nil {
		return predefinedurl.ListOverridesOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) DescribeUrlCategoryActionOverride(ctx context.Context, input predefinedurl.GetOverrideInput) (predefinedurl.GetOverrideOutput, error) {
	out, err := c.client.DescribeUrlCategoryActionOverride(ctx, input)
	if err != nil {
		return predefinedurl.GetOverrideOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) UpdateUrlCategoryActionOverride(ctx context.Context, input predefinedurl.OverrideInput) error {
	if err := c.client.UpdateUrlCategoryActionOverride(ctx, input); err != nil {
		return err
	}
	return nil
}
