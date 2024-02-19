package aws

import (
	"context"
	"fmt"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/stack"
)

func (c *Client) RefreshCloudRulestackAdminJwt(ctx context.Context) error {
	return nil
}

func (c *Client) GetCloudNGFWServiceToken(ctx context.Context, info stack.AuthInput) (stack.AuthOutput, error) {
	return stack.AuthOutput{}, fmt.Errorf("not implemented")
}
