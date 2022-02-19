package api

import (
	"github.com/aws/aws-sdk-go/service/sts"
)

type Client interface {
	// Logging.
	Log(string, string, ...interface{})

	// API functions.
	Communicate(string, string, []string, interface{}, Oker, ...*sts.Credentials) ([]byte, error)
}
