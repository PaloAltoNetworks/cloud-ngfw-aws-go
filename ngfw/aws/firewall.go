package aws

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/firewall"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/tag"
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// List returns a list of firewalls.
func (c *Client) ListFirewall(ctx context.Context, input firewall.ListInput) (firewall.ListOutput, error) {
	if len(input.VpcIds) == 0 {
		c.Log(http.MethodGet, "list NGFirewalls in all the VPCs")
	} else {
		c.Log(http.MethodGet, "list NGFirewalls in %q VPCs", strings.Join(input.VpcIds, ","))
	}

	var uv url.Values
	if input.Rulestack != "" || input.Describe == true {
		uv = url.Values{}
	}
	if input.Rulestack != "" {
		uv.Set("rulestackname", input.Rulestack)
	}
	if input.Describe == true {
		uv.Set("describe", "true")
	}

	var ans firewall.ListOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodGet,
		[]string{"v1", "config", "ngfirewalls"},
		uv,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) CreateFirewall(ctx context.Context, input firewall.Info) (firewall.CreateOutput, error) {
	c.Log(http.MethodPost, "create firewall %q", input.Name)

	var ans firewall.CreateOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPost,
		[]string{"v1", "config", "ngfirewalls"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Modify updates the modifiable parts of a NGFW.
//
// This includes:
//   - description
//   - subnet mappings
//   - app id version / automatic upgrade app id version
//   - rulestack
//   - tags
func (c *Client) ModifyFirewall(ctx context.Context, input firewall.Info) error {
	ans, err := c.ReadFirewall(ctx, firewall.ReadInput{Name: input.Name, AccountId: input.AccountId})
	if err != nil {
		return err
	}
	cur := ans.Response.Firewall
	curTags := cur.Tags

	// No idea if this is needed or not, but do it for now just to be safe.
	tin := firewall.ListTagsInput{
		Firewall:   input.Name,
		AccountId:  input.AccountId,
		MaxResults: 100,
	}
	tans, err := c.ListTagsForFirewall(ctx, tin)
	if err != nil {
		return err
	}
	curTags = tans.Response.Tags

	// Update description.
	if input.Description != cur.Description {
		v := firewall.UpdateDescriptionInput{
			Firewall:    input.Name,
			AccountId:   input.AccountId,
			Description: input.Description,
		}
		if err = c.UpdateFirewallDescription(ctx, v); err != nil {
			return err
		}
	}

	// Update subnet mappings.
	assoc := make([]firewall.SubnetMapping, 0, len(input.SubnetMappings))
	disassoc := make([]firewall.SubnetMapping, 0, len(cur.SubnetMappings))
	for _, x := range input.SubnetMappings {
		found := false
		for _, y := range cur.SubnetMappings {
			if x.SubnetId != "" && x.SubnetId == y.SubnetId {
				found = true
			} else if x.AvailabilityZone != "" && x.AvailabilityZone == y.AvailabilityZone {
				found = true
			}
			if found {
				break
			}
		}

		if !found {
			assoc = append(assoc, firewall.SubnetMapping{
				SubnetId:         x.SubnetId,
				AvailabilityZone: x.AvailabilityZone,
			})
		}
	}
	if len(assoc) == 0 {
		assoc = nil
	}
	for _, x := range cur.SubnetMappings {
		found := false
		for _, y := range input.SubnetMappings {
			if x.SubnetId != "" && x.SubnetId == y.SubnetId {
				found = true
			} else if x.AvailabilityZone != "" && x.AvailabilityZone == y.AvailabilityZone {
				found = true
			}
			if found {
				break
			}
		}

		if !found {
			disassoc = append(assoc, firewall.SubnetMapping{
				SubnetId:         x.SubnetId,
				AvailabilityZone: x.AvailabilityZone,
			})
		}
	}
	if len(disassoc) == 0 {
		disassoc = nil
	}
	if assoc != nil || disassoc != nil {
		v := firewall.UpdateSubnetMappingsInput{
			Firewall:                   input.Name,
			AccountId:                  input.AccountId,
			AssociateSubnetMappings:    assoc,
			DisassociateSubnetMappings: disassoc,
		}
		if err = c.UpdateFirewallSubnetMappings(ctx, v); err != nil {
			return err
		}
	}

	// Update content version.
	if input.AppIdVersion != cur.AppIdVersion || input.AutomaticUpgradeAppIdVersion != cur.AutomaticUpgradeAppIdVersion {
		v := firewall.UpdateContentVersionInput{
			Firewall:                     input.Name,
			AccountId:                    input.AccountId,
			AppIdVersion:                 input.AppIdVersion,
			AutomaticUpgradeAppIdVersion: input.AutomaticUpgradeAppIdVersion,
		}
		if err = c.UpdateFirewallContentVersion(ctx, v); err != nil {
			return err
		}
	}

	// Update rulestack.
	if input.Rulestack != cur.Rulestack {
		v := firewall.UpdateRulestackInput{
			Firewall:  input.Name,
			AccountId: input.AccountId,
			Rulestack: input.Rulestack,
		}
		if err = c.UpdateFirewallRulestack(ctx, v); err != nil {
			return err
		}
	}

	// Update tags.
	addTags := make([]tag.Details, 0, len(input.Tags))
	rmTags := make([]string, 0, len(curTags))
	for _, x := range input.Tags {
		ok := false
		for _, y := range curTags {
			if x.Key == y.Key {
				if x.Value == y.Value {
					ok = true
				} else {
					rmTags = append(rmTags, x.Key)
				}
				break
			}
		}
		if !ok {
			addTags = append(addTags, x)
		}
	}
	for _, x := range curTags {
		found := false
		for _, y := range input.Tags {
			if x.Key == y.Key {
				found = true
				break
			}
		}
		if !found {
			rmTags = append(rmTags, x.Key)
		}
	}
	// Due to the 50 tag limit, removing tags must happen before adding tags.
	if len(rmTags) > 0 {
		v := firewall.RemoveTagsInput{
			Firewall:  input.Name,
			AccountId: input.AccountId,
			Tags:      rmTags,
		}
		if err = c.RemoveTagsForFirewall(ctx, v); err != nil {
			return err
		}
	}
	if len(addTags) > 0 {
		v := firewall.AddTagsInput{
			Firewall:  input.Name,
			AccountId: input.AccountId,
			Tags:      addTags,
		}
		if err = c.AddTagsForFirewall(ctx, v); err != nil {
			return err
		}
	}

	// Done.
	return nil
}

// Read returns information on the given object.
func (c *Client) ReadFirewall(ctx context.Context, input firewall.ReadInput) (firewall.ReadOutput, error) {
	name := input.Name
	c.Log(http.MethodGet, "describe firewall: %s", name)

	var ans firewall.ReadOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodGet,
		[]string{"v1", "config", "ngfirewalls", name},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// UpdateDescription updates the description of the firewall.
func (c *Client) UpdateFirewallDescription(ctx context.Context, input firewall.UpdateDescriptionInput) error {
	c.Log(http.MethodPut, "updating firewall description: %s", input.Firewall)

	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPut,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "description"},
		nil,
		input,
		nil,
	)

	return err
}

// UpdateContentVersion updates the content version of the firewall.
func (c *Client) UpdateFirewallContentVersion(ctx context.Context, input firewall.UpdateContentVersionInput) error {
	c.Log(http.MethodPut, "updating firewall content version: %s", input.Firewall)

	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPut,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "contentversion"},
		nil,
		input,
		nil,
	)

	return err
}

// UpdateSubnetMappings updates the subnet mappings of the firewall.
func (c *Client) UpdateFirewallSubnetMappings(ctx context.Context, input firewall.UpdateSubnetMappingsInput) error {
	c.Log(http.MethodPatch, "updating firewall subnet mappings: %s", input.Firewall)

	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPatch,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "subnets"},
		nil,
		input,
		nil,
	)

	return err
}

// UpdateRulestack updates the rulestack for the given firewall.
func (c *Client) UpdateFirewallRulestack(ctx context.Context, input firewall.UpdateRulestackInput) error {
	c.Log(http.MethodPost, "updating firewall rulestack: %s", input.Firewall)

	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPost,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "rulestack"},
		nil,
		input,
		nil,
	)

	return err
}

// ListTags gets the tags for the given Firewall.
func (c *Client) ListTagsForFirewall(ctx context.Context, input firewall.ListTagsInput) (firewall.ListTagsOutput, error) {
	c.Log(http.MethodGet, "list tags for firewall: %s", input.Firewall)

	var uv url.Values
	if input.AccountId != "" || input.NextToken != "" || input.MaxResults != 0 {
		uv = url.Values{}
		if input.AccountId != "" {
			uv.Set("accountid", input.AccountId)
		}
		if input.NextToken != "" {
			uv.Set("nexttoken", input.NextToken)
		}
		if input.MaxResults != 0 {
			uv.Set("maxresults", strconv.Itoa(input.MaxResults))
		}
	}

	var ans firewall.ListTagsOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodGet,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "tags"},
		uv,
		nil,
		&ans,
	)

	return ans, err
}

// RemoveTags removes the given tags from the firewall.
func (c *Client) RemoveTagsForFirewall(ctx context.Context, input firewall.RemoveTagsInput) error {
	c.Log(http.MethodDelete, "removing tags from firewall: %s", input.Firewall)

	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodDelete,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "tags"},
		nil,
		input,
		nil,
	)

	return err
}

// AddTags adds the given tags to the firewall.
func (c *Client) AddTagsForFirewall(ctx context.Context, input firewall.AddTagsInput) error {
	c.Log(http.MethodPost, "adding tags to the firewall: %s", input.Firewall)

	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPost,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "tags"},
		nil,
		input,
		nil,
	)

	return err
}

// Delete the given firewall.
func (c *Client) DeleteFirewall(ctx context.Context, input firewall.DeleteInput) error {
	c.Log(http.MethodDelete, "delete firewall: %s", input.Name)

	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodDelete,
		[]string{"v1", "config", "ngfirewalls", input.Name},
		nil,
		input,
		nil,
	)

	return err
}

// Associate Firewall to Global rulestack
func (c *Client) AssociateGlobalRuleStack(ctx context.Context, input firewall.AssociateInput) (firewall.AssociateOutput, error) {
	c.Log(http.MethodPut, "associating firewall to global rulestack: %s", input.Firewall)

	var ans firewall.AssociateOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPost,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "globalrulestack"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Disassociate Firewall to Global rulestack
func (c *Client) DisAssociateGlobalRuleStack(ctx context.Context, input firewall.DisAssociateInput) (firewall.DisAssociateOutput, error) {
	c.Log(http.MethodDelete, "associating firewall to global rulestack: %s", input.Firewall)

	var ans firewall.DisAssociateOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodDelete,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "globalrulestack"},
		nil,
		input,
		&ans,
	)

	return ans, err
}
