package firewall

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/permissions"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/tag"
)

// Client is a client for this collection.
type Client struct {
	client api.Client
}

// NewClient returns a new client for this collection.
func NewClient(client api.Client) *Client {
	return &Client{client: client}
}

// List returns a list of firewalls.
func (c *Client) List(ctx context.Context, input ListInput) (ListOutput, error) {
	if len(input.VpcIds) == 0 {
		c.client.Log(http.MethodGet, "list NGFirewalls in all the VPCs")
	} else {
		c.client.Log(http.MethodGet, "list NGFirewalls in %q VPCs", strings.Join(input.VpcIds, ","))
	}

	var uv url.Values
	if input.Rulestack != "" {
		uv = url.Values{}
		uv.Set("rulestackname", input.Rulestack)
	}

	var ans ListOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodGet,
		[]string{"v1", "config", "ngfirewalls"},
		uv,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) Create(ctx context.Context, input Info) (CreateOutput, error) {
	c.client.Log(http.MethodPost, "create firewall %q multi vpc:%+v", input.Name, input.MultiVpc)

	var ans CreateOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
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
func (c *Client) Modify(ctx context.Context, input Info) error {
	ans, err := c.Read(ctx, ReadInput{Name: input.Name, AccountId: input.AccountId})
	if err != nil {
		return err
	}
	cur := ans.Response.Firewall
	curTags := cur.Tags

	// No idea if this is needed or not, but do it for now just to be safe.
	tin := ListTagsInput{
		Firewall:   input.Name,
		AccountId:  input.AccountId,
		MaxResults: 100,
	}
	tans, err := c.ListTags(ctx, tin)
	if err != nil {
		return err
	}
	curTags = tans.Response.Tags

	// Update description.
	if input.Description != cur.Description {
		v := UpdateDescriptionInput{
			Firewall:    input.Name,
			AccountId:   input.AccountId,
			Description: input.Description,
		}
		if err = c.UpdateDescription(ctx, v); err != nil {
			return err
		}
	}

	// Update subnet mappings.
	assoc := make([]SubnetMapping, 0, len(input.SubnetMappings))
	disassoc := make([]SubnetMapping, 0, len(cur.SubnetMappings))
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
			assoc = append(assoc, SubnetMapping{
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
			disassoc = append(assoc, SubnetMapping{
				SubnetId:         x.SubnetId,
				AvailabilityZone: x.AvailabilityZone,
			})
		}
	}
	if len(disassoc) == 0 {
		disassoc = nil
	}
	if assoc != nil || disassoc != nil {
		v := UpdateSubnetMappingsInput{
			Firewall:                   input.Name,
			AccountId:                  input.AccountId,
			AssociateSubnetMappings:    assoc,
			DisassociateSubnetMappings: disassoc,
			MultiVpc:                   cur.MultiVpc,
		}
		if err = c.UpdateSubnetMappings(ctx, v); err != nil {
			return err
		}
	}

	// update MultiVpcEnable
	if input.MultiVpc != cur.MultiVpc {
		v := UpdateSubnetMappingsInput{
			Firewall:  input.Name,
			AccountId: input.AccountId,
			MultiVpc:  input.MultiVpc,
		}
		if err = c.UpdateSubnetMappings(ctx, v); err != nil {
			return err
		}
	}

	// Update content version.
	if input.AppIdVersion != cur.AppIdVersion || input.AutomaticUpgradeAppIdVersion != cur.AutomaticUpgradeAppIdVersion {
		v := UpdateContentVersionInput{
			Firewall:                     input.Name,
			AccountId:                    input.AccountId,
			AppIdVersion:                 input.AppIdVersion,
			AutomaticUpgradeAppIdVersion: input.AutomaticUpgradeAppIdVersion,
		}
		if err = c.UpdateContentVersion(ctx, v); err != nil {
			return err
		}
	}

	// Update rulestack.
	if input.Rulestack != cur.Rulestack {
		v := UpdateRulestackInput{
			Firewall:  input.Name,
			AccountId: input.AccountId,
			Rulestack: input.Rulestack,
		}
		if err = c.UpdateRulestack(ctx, v); err != nil {
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
		v := RemoveTagsInput{
			Firewall:  input.Name,
			AccountId: input.AccountId,
			Tags:      rmTags,
		}
		if err = c.RemoveTags(ctx, v); err != nil {
			return err
		}
	}
	if len(addTags) > 0 {
		v := AddTagsInput{
			Firewall:  input.Name,
			AccountId: input.AccountId,
			Tags:      addTags,
		}
		if err = c.AddTags(ctx, v); err != nil {
			return err
		}
	}

	// Done.
	return nil
}

// Read returns information on the given object.
func (c *Client) Read(ctx context.Context, input ReadInput) (ReadOutput, error) {
	name := input.Name
	c.client.Log(http.MethodGet, "describe firewall: %s", name)

	var ans ReadOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodGet,
		[]string{"v1", "config", "ngfirewalls", name},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// UpdateDescription updates the description of the firewall.
func (c *Client) UpdateDescription(ctx context.Context, input UpdateDescriptionInput) error {
	c.client.Log(http.MethodPut, "updating firewall description: %s", input.Firewall)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodPut,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "description"},
		nil,
		input,
		nil,
	)

	return err
}

// UpdateContentVersion updates the content version of the firewall.
func (c *Client) UpdateContentVersion(ctx context.Context, input UpdateContentVersionInput) error {
	c.client.Log(http.MethodPut, "updating firewall content version: %s", input.Firewall)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodPut,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "contentversion"},
		nil,
		input,
		nil,
	)

	return err
}

// UpdateSubnetMappings updates the subnet mappings of the firewall.
func (c *Client) UpdateSubnetMappings(ctx context.Context, input UpdateSubnetMappingsInput) error {
	c.client.Log(http.MethodPut, "updating firewall subnet mappings: %s", input.Firewall)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodPut,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "subnets"},
		nil,
		input,
		nil,
	)

	return err
}

// UpdateRulestack updates the rulestack for the given firewall.
func (c *Client) UpdateRulestack(ctx context.Context, input UpdateRulestackInput) error {
	c.client.Log(http.MethodPost, "updating firewall rulestack: %s", input.Firewall)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodPost,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "rulestack"},
		nil,
		input,
		nil,
	)

	return err
}

// ListTags gets the tags for the given Firewall.
func (c *Client) ListTags(ctx context.Context, input ListTagsInput) (ListTagsOutput, error) {
	c.client.Log(http.MethodGet, "list tags for firewall: %s", input.Firewall)

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

	var ans ListTagsOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodGet,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "tags"},
		uv,
		nil,
		&ans,
	)

	return ans, err
}

// RemoveTags removes the given tags from the firewall.
func (c *Client) RemoveTags(ctx context.Context, input RemoveTagsInput) error {
	c.client.Log(http.MethodDelete, "removing tags from firewall: %s", input.Firewall)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodDelete,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "tags"},
		nil,
		input,
		nil,
	)

	return err
}

// AddTags adds the given tags to the firewall.
func (c *Client) AddTags(ctx context.Context, input AddTagsInput) error {
	c.client.Log(http.MethodPost, "adding tags to the firewall: %s", input.Firewall)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodPost,
		[]string{"v1", "config", "ngfirewalls", input.Firewall, "tags"},
		nil,
		input,
		nil,
	)

	return err
}

// Delete the given firewall.
func (c *Client) Delete(ctx context.Context, input DeleteInput) error {
	c.client.Log(http.MethodDelete, "delete firewall: %s", input.Name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodDelete,
		[]string{"v1", "config", "ngfirewalls", input.Name},
		nil,
		input,
		nil,
	)

	return err
}
