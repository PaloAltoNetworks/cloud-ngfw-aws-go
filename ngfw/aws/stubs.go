package aws

import (
	"context"
	"fmt"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/stack"
)

func (c *Client) GetCloudNGFWServiceToken(ctx context.Context, info stack.AuthInput) (stack.AuthOutput, error) {
	return stack.AuthOutput{}, fmt.Errorf("not implemented")
}
