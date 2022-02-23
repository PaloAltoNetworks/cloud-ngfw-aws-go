package awsngfw

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type provider struct {
	Value credentials.Value
}

func (p provider) Retrieve() (credentials.Value, error) {
	return p.Value, nil
}

func (p provider) IsExpired() bool {
	return false
}
