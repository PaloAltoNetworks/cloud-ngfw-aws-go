package cloudngfw

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
)

type jwtKeyInfo struct {
	Region string `json:"Region"`
	Tenant string `json:"Tenant"`
}

type getJwt struct {
	Expires int         `json:"ExpiryTime"`
	KeyInfo *jwtKeyInfo `json:"KeyInfo"`
}

type authResponse struct {
	Resp   authData   `json:"Response"`
	Status api.Status `json:"ResponseStatus"`
}

type authData struct {
	Jwt            string  `json:"TokenId"`
	SubscriptionId string  `json:"SubscriptionKey"`
	ExpiryTime     float64 `json:"ExpiryTime"`
	Enabled        bool    `json:"Enabled"`
}

func (a authResponse) Ok() bool {
	return a.Status.Ok()
}

func (a authResponse) Error() string {
	return a.Status.Error()
}
