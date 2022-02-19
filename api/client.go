package api

import (
	"net/url"

	"github.com/aws/aws-sdk-go/service/sts"
)

type Client interface {
	// Logging.
	Log(string, string, ...interface{})

	// API functions.
	Communicate(string, string, []string, url.Values, interface{}, Failure, ...*sts.Credentials) ([]byte, error)
}
