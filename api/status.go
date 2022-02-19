package api

import (
	"fmt"
	"strings"
)

/*
Status is a container for the status of API calls.

This contains the error code and the reason.
*/
type Status struct {
	Code   int    `json:"ErrorCode"`
	Reason string `json:"Reason"`
}

func (s Status) Failed() *Status {
	if s.Code != 0 {
		return &s
	}

	return nil
}

func (s Status) ObjectNotFound() bool {
	return strings.HasSuffix(s.Reason, " does not exist")
}

func (s Status) Error() string {
	return fmt.Sprintf("Error(%d): %s", s.Code, s.Reason)
}
