package country

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/response"
)

// V1 list.

type ListInput struct {
	NextToken  string `json:"NextToken,omitempty"`
	MaxResults int    `json:"MaxResults,omitempty"`
}

type ListOutput struct {
	Response *ListOutputDetails `json:"Response"`
	Status   response.Status    `json:"ResponseStatus"`
}

func (o ListOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type ListOutputDetails struct {
	Countries []Country `json:"CountryCodes"`
	NextToken string    `json:"NextToken"`
}

type Country struct {
	Code        string `json:"Code"`
	Description string `json:"Description"`
}
