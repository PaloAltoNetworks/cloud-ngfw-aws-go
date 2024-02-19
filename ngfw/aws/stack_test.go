package aws

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/stack"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPushRulestackXML(t *testing.T) {
	c := Client{Mock: true, MockedResp: func() ([]byte, error) { return nil, nil }}
	apiClient := api.NewAPIClient(&c, context.TODO(), 1500, "", true)
	input := stack.PushRulestackCMInput{
		Name: "test-dg",
	}

	err := apiClient.PushRuleStackCM(input)
	assert.Nil(t, err, "PushRulestackCM returned %s", err)
}
