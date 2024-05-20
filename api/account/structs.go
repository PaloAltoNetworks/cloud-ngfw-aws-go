package account

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/response"
)

// V1 create.
type CreateInput struct {
	AccountId string `json:"AccountId,omitempty"`
	Origin    string `json:"Origin,omitempty"`
}

type Info struct {
	TrustedAccount string `json:"ServiceAccountId,omitempty"`
	ExternalId     string `json:"ExternalId,omitempty"`
	SNSTopicArn    string `json:"SNSTopicArn,omitempty"`
	Origin         string `json:"Origin,omitempty"`
}

type CreateOutput struct {
	Response Info            `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

// V1 read.
type ReadInput struct {
	AccountId string `json:"AccountId,omitempty"`
}

type ReadOutput struct {
	Response ReadResponse    `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

func (o ReadOutput) Failed() *response.Status {
	return o.Status.Failed()
}

func (o CreateOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type AccountDetail struct {
	AccountId                 string `json:"AccountId,omitempty"`
	CloudFormationTemplateURL string `json:"CloudFormationTemplateURL,omitempty"`
	OnboardingStatus          string `json:"OnboardingStatus,omitempty"`
	ExternalId                string `json:"ExternalId,omitempty"`
	ServiceAccountId          string `json:"ServiceAccountId,omitempty"`
	SNSTopicArn               string `json:"SNSTopicArn,omitempty"`
}

type ReadResponse struct {
	AccountDetail
	UpdateToken string `json:"UpdateToken,omitempty"`
}

type ListInput struct {
	Describe   bool   `json:"Describe,omitempty"`
	MaxResults int    `json:"MaxResults,omitempty"`
	NextToken  string `json:"NextToken,omitempty"`
}

type ListAccount struct {
	AccountId string
}

type ListResponse struct {
	AccountIds     []string        `json:"AccountIds"`
	AccountDetails []AccountDetail `json:"AccountDetails,omitempty"`
	NextToken      string          `json:"NextToken,omitempty"`
}

type ListOutput struct {
	Response ListResponse    `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

func (o ListOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type DeleteInput struct {
	AccountId string `json:"AccountId"`
}
