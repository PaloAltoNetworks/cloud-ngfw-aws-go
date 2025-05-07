package aws

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/response"
)

/*
Response is a generic response container.

This is useful if you don't care about the response from the API, as long as
there wasn't any errors.
*/
type Response struct {
	Status response.Status `json:"ResponseStatus"`
}

func (o Response) Failed() *response.Status {
	return o.Status.Failed()
}
