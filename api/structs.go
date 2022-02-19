package api

import (
	"fmt"
)

type Status struct {
	Code   int    `json:"ErrorCode"`
	Reason string `json:"Reason"`
}

func (s Status) Ok() bool {
	return s.Code == 0
}

func (s Status) Error() string {
	return fmt.Sprintf("(%d): %s", s.Code, s.Reason)
}
