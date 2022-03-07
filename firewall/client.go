package firewall

import (
	"context"
	"net/http"
	"strings"

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

// List returns a list of firewalls.
func (c *Client) List(ctx context.Context, input ListInput) (ListOutput, error) {
	if len(input.VpcIds) == 0 {
		c.client.Log(http.MethodGet, "list NGFirewalls in all the VPCs")
	} else {
		c.client.Log(http.MethodGet, "list NGFirewalls in %q VPCs", strings.Join(input.VpcIds, ","))
	}

	var ans ListOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodGet,
		[]string{"v1", "config", "ngfirewalls"},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) Create(ctx context.Context, input Info) (ReadOutput, error) {
	c.client.Log(http.MethodPost, "create firewall %q", input.Name)

	var ans ReadOutput
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

// Read returns information on the given object.
func (c *Client) Read(ctx context.Context, input ReadInput) (ReadOutput, error) {
	name := input.Name
	c.client.Log(http.MethodGet, "describe firewall: %s", name)

	var ans ReadOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "ngfirewalls", name},
		nil,
		input,
		&ans,
	)

	return ans, err
}

// UpdateDescription updates the description of the firewall.
func (c *Client) UpdateDescription(ctx context.Context, input Info) error {
	name := input.Name
	input.Name = ""

	c.client.Log(http.MethodPut, "updating firewall %q description: %s", name, input.Description)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodPut,
		[]string{"v1", "config", "ngfirewalls", name, "description"},
		nil,
		input,
		nil,
	)

	return err
}

// UpdateNGFirewallContentVersion updates the content version of the firewall.
func (c *Client) UpdateNGFirewallContentVersion(ctx context.Context, input Info) error {
	name := input.Name
	input.Name = ""

	c.client.Log(http.MethodPut, "updating firewall %q content version: %s", name, input.AppIdVersion)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodPut,
		[]string{"v1", "config", "ngfirewalls", name, "contentversion"},
		nil,
		input,
		nil,
	)

	return err
}

// UpdateSubnetMappings updates the subnet mappings of the firewall.
func (c *Client) UpdateSubnetMappings(ctx context.Context, input Info) error {
	name := input.Name
	input.Name = ""

	c.client.Log(http.MethodPut, "updating firewall %q associate subnet mappings: %s disassociate subnet mappings: %s", name, input.AssociateSubnetMappings, input.DisassociateSubnetMappings)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodPut,
		[]string{"v1", "config", "ngfirewalls", name, "subnets"},
		nil,
		input,
		nil,
	)

	return err
}

// Delete the given firewall.
func (c *Client) Delete(ctx context.Context, input ReadInput) error {
	name := input.Name

	c.client.Log(http.MethodDelete, "delete firewall: %s", name)

	_, err := c.client.Communicate(
		ctx,
		permissions.Firewall,
		http.MethodDelete,
		[]string{"v1", "config", "ngfirewalls", name},
		nil,
		input,
		nil,
	)

	return err
}
