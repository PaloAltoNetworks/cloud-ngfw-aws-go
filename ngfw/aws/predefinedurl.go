package aws

import (
	"context"
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/predefinedurl"
)

// List returns a list of objects.
func (c *Client) ListUrlPredefinedCategories(ctx context.Context, input predefinedurl.ListInput) (predefinedurl.ListOutput, error) {
	c.Log(http.MethodGet, "list predefined url categories")
	path := Path{
		V1Path: []string{"v1", "config", "urlcategories"},
	}
	var ans predefinedurl.ListOutput
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodGet,
		path,
		nil,
		input,
		&ans,
	)

	return ans, err
}

// ListOverrides returns URL categories with overrides specified.
func (c *Client) ListUrlCategoriesActionOverride(ctx context.Context, input predefinedurl.ListOverridesInput) (predefinedurl.ListOverridesOutput, error) {
	c.Log(http.MethodGet, "list predefined url category overrides for rulestack %q", input.Rulestack)
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "urlfilteringprofiles", "custom", "urlcategories"},
	}
	var ans predefinedurl.ListOverridesOutput
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodGet,
		path,
		nil,
		input,
		&ans,
	)

	return ans, err
}

// GetOverride returns the URL category override info.
func (c *Client) DescribeUrlCategoryActionOverride(ctx context.Context, input predefinedurl.GetOverrideInput) (predefinedurl.GetOverrideOutput, error) {
	c.Log(http.MethodGet, "get %q predefined url category override: %s", input.Rulestack, input.Name)
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "urlfilteringprofiles", "custom", "urlcategories", input.Name},
	}
	var ans predefinedurl.GetOverrideOutput
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodGet,
		path,
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Override specifies an override for a predefined URL category.
func (c *Client) UpdateUrlCategoryActionOverride(ctx context.Context, input predefinedurl.OverrideInput) error {
	c.Log(http.MethodPut, "override %q predefined url category: %s", input.Rulestack, input.Name)
	path := Path{
		V1Path: []string{"v1", "config", "rulestacks", input.Rulestack, "urlfilteringprofiles", "custom", "urlcategories", input.Name, "action"},
	}
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodPut,
		path,
		nil,
		input,
		nil,
	)

	return err
}
