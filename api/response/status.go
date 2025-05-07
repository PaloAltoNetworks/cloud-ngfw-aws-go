package response

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
	return strings.Contains(s.Reason, " does not exist")
}

func (s Status) InvalidRequest() bool {
	return strings.Contains(s.Reason, "invalid request")
}

func (s Status) Error() string {
	return fmt.Sprintf("Error(%d): %s", s.Code, s.Reason)
}

func (s Status) TokenConflict() bool {
	return strings.Contains(s.Reason, "please provide latest token")
}

func (s Status) LrsDoesNotExist() bool {
	return strings.Contains(s.Reason, "rulestack does not exist")
}

func (s Status) LrsAssociateDoesNotExist() bool {
	return strings.Contains(s.Reason, "rulestack association does not exist")
}

func (s Status) LrsResourceEntryMissing() bool {
	return strings.Contains(s.Reason, "could not find an associated rulestack")
}

func IsResponseWithError(body []byte) error {
	var ans unknownApi
	if err := json.Unmarshal(body, &ans); err == nil {
		if ans.Message != "" {
			return &Status{
				Code:   -1,
				Reason: fmt.Sprintf("error: %s", ans.Message),
			}
		}
	}

	return nil
}

func NewUnknownPathError(v []string) *Status {
	return &Status{
		Code:   -1,
		Reason: fmt.Sprintf("error: /%s", strings.Join(v, "/")),
	}
}

type unknownApi struct {
	Message string `json:"message"`
}
