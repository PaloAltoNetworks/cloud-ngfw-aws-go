package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/country"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) ListCountry(ctx context.Context, a country.ListInput) (country.ListOutput, error) {
	out, err := c.client.ListCountry(ctx, a)
	if err != nil {
		return country.ListOutput{}, err
	}
	return out, nil
}
