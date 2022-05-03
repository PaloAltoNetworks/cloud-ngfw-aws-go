package firewall

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/tag"
)

// V1 list.

type ListInput struct {
	Firewall   string
	AccountId  string
	NextToken  string
	MaxResults int
}

type ListOutput struct {
	Response ListOutputDetails `json:"Response"`
	Status   api.Status        `json:"ResponseStatus"`
}

func (o ListOutput) Failed() *api.Status {
	return o.Status.Failed()
}

type ListOutputDetails struct {
	Firewall  string        `json:"ResourceName"`
	NextToken string        `json:"NextToken"`
	Tags      []tag.Details `json:"Tags"`
}

// V1 create.

type Info struct {
	Firewall  string        `json:"-"`
	AccountId string        `json:"AccountId,omitempty"`
	Tags      []tag.Details `json:"Tags,omitempty"`
}

type UntagInput struct {
	Firewall  string   `json:"-"`
	AccountId string   `json:"AccountId"`
	Tags      []string `json:"TagKeys"`
}
