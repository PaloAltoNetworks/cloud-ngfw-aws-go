package aws

import (
	"context"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/country"
	"net/http"
)

// List returns a list of objects.
func (c *Client) ListCountry(ctx context.Context, input country.ListInput) (country.ListOutput, error) {
	c.Log(http.MethodGet, "list countries")

	var ans country.ListOutput
	_, err := c.Communicate(
		ctx,
		PermissionRulestack,
		http.MethodGet,
		[]string{"v1", "config", "countries"},
		nil,
		input,
		&ans,
	)

	return ans, err
}
