package api

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go/service/sts"
)

type Client interface {
	// Logging.
	Log(string, string, ...interface{})

	// API functions.
	Communicate(context.Context, string, string, []string, url.Values, interface{}, Failure, ...*sts.Credentials) ([]byte, error)
}
