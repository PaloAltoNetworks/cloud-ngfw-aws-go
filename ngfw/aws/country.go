package aws

import (
	"context"
	"net/http"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/country"
)

// List returns a list of objects.
func (c *Client) ListCountry(ctx context.Context, input country.ListInput) (country.ListOutput, error) {
	c.Log(http.MethodGet, "list countries")
	path := Path{
		V1Path: []string{"v1", "config", "countries"},
	}
	var ans country.ListOutput
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
