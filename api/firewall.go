package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/firewall"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) ListFirewall(ctx context.Context, a firewall.ListInput) (firewall.ListOutput, error) {
	out, err := c.client.ListFirewall(ctx, a)
	if err != nil {
		return firewall.ListOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) CreateFirewall(ctx context.Context, input firewall.Info) (firewall.CreateOutput, error) {
	out, err := c.client.CreateFirewall(ctx, input)
	if err != nil {
		return firewall.CreateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ModifyFirewall(ctx context.Context, input firewall.Info) error {
	if err := c.client.ModifyFirewall(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) ReadFirewall(ctx context.Context, input firewall.ReadInput) (firewall.ReadOutput, error) {
	out, err := c.client.ReadFirewall(ctx, input)
	if err != nil {
		return firewall.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) UpdateFirewallDescription(ctx context.Context, input firewall.UpdateDescriptionInput) error {
	if err := c.client.UpdateFirewallDescription(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) UpdateFirewallContentVersion(ctx context.Context, input firewall.UpdateContentVersionInput) error {
	if err := c.client.UpdateFirewallContentVersion(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) UpdateFirewallSubnetMappings(ctx context.Context, input firewall.UpdateSubnetMappingsInput) error {
	if err := c.client.UpdateFirewallSubnetMappings(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) UpdateFirewallRulestack(ctx context.Context, input firewall.UpdateRulestackInput) error {
	if err := c.client.UpdateFirewallRulestack(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) ListTagsForFirewall(ctx context.Context, input firewall.ListTagsInput) (firewall.ListTagsOutput, error) {
	out, err := c.client.ListTagsForFirewall(ctx, input)
	if err != nil {
		return firewall.ListTagsOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) RemoveTagsForFirewall(ctx context.Context, input firewall.RemoveTagsInput) error {
	if err := c.client.RemoveTagsForFirewall(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) AddTagsForFirewall(ctx context.Context, input firewall.AddTagsInput) error {
	if err := c.client.AddTagsForFirewall(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) DeleteFirewall(ctx context.Context, input firewall.DeleteInput) error {
	if err := c.client.DeleteFirewall(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) AssociateGlobalRuleStack(ctx context.Context, input firewall.AssociateInput) (firewall.AssociateOutput, error) {
	out, err := c.client.AssociateGlobalRuleStack(ctx, input)
	if err != nil {
		return firewall.AssociateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) DisAssociateGlobalRuleStack(ctx context.Context, input firewall.DisAssociateInput) (firewall.DisAssociateOutput, error) {
	out, err := c.client.DisAssociateGlobalRuleStack(ctx, input)
	if err != nil {
		return firewall.DisAssociateOutput{}, err
	}
	return out, nil
}
