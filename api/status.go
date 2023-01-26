package api

import (
	"encoding/json"
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

func IsErrorMessage(path []string, body []byte, statusCode int) error {
	var ans unknownApi
	if err := json.Unmarshal(body, &ans); err == nil {
		fmt.Sprintf("Error: HTTP %d: %s. path:%s",
			statusCode, ans.Message, path)
	}

	return nil
}

func NewUnknownPathError(v []string) *Status {
	return &Status{
		Code:   -1,
		Reason: fmt.Sprintf("Unknown path: /%s", strings.Join(v, "/")),
	}
}

type unknownApi struct {
	Message string `json:"message"`
}
