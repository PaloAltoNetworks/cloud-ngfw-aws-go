package aws

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/response"
)

type jwtKeyInfo struct {
	Region string `json:"Region"`
	Tenant string `json:"Tenant"`
}

type getJwt struct {
	Expires int         `json:"ExpiryTime"`
	KeyInfo *jwtKeyInfo `json:"KeyInfo,omitempty"`
}

type authResponse struct {
	Resp   authData        `json:"Response"`
	Status response.Status `json:"ResponseStatus"`
}

type authData struct {
	Jwt             string  `json:"TokenId"`
	SubscriptionKey string  `json:"SubscriptionKey"`
	ExpiryTime      float64 `json:"ExpiryTime"`
	Enabled         bool    `json:"Enabled"`
}

func (o authResponse) Failed() *response.Status {
	return o.Status.Failed()
}
