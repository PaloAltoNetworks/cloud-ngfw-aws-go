package url

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
)

// v1 list url predefined categories.

type ListInput struct {
	NextToken  string `json:"NextToken,omitempty"`
	MaxResults int    `json:"MaxResults,omitempty"`
}

type ListOutput struct {
	Response ListResponse `json:"Response"`
	Status   api.Status   `json:"ResponseStatus"`
}

func (o ListOutput) Failed() *api.Status {
	return o.Status.Failed()
}

type ListResponse struct {
	NextToken  string     `json:"NextToken"`
	Categories []Category `json:"CategoriesRunning"`
}

type Category struct {
	Name   string `json:"Name"`
	Action string `json:"Action"`
}

// v1 list overrides.

type ListOverridesInput struct {
	Rulestack   string `json:"-"`
	NextToken   string `json:"NextToken,omitempty"`
	MaxResults  int    `json:"MaxResults,omitempty"`
	Candidate   bool   `json:"Candidate,omitempty"`
	Running     bool   `json:"Running,omitempty"`
	Uncommitted bool   `json:"Uncommitted,omitempty"`
}

type ListOverridesOutput struct {
	Response ListOverridesOutputResponse `json:"Response"`
	Status   api.Status                  `json:"ResponseStatus"`
}

func (o ListOverridesOutput) Failed() *api.Status {
	return o.Status.Failed()
}

type ListOverridesOutputResponse struct {
	Rulestack   string                `json:"RuleStackName"`
	NextToken   string                `json:"NextToken"`
	Running     []string              `json:"Running"`
	Candidate   []string              `json:"Candidate"`
	Uncommitted []UncommittedOverride `json:"CategoriesUncommitted"`
}

type UncommittedOverride struct {
	Name      string `json:"Name"`
	Operation string `json:"Operation"`
}

// v1 get override.

type GetOverrideInput struct {
	Rulestack string `json:"-"`
	Name      string `json:"-"`
	Candidate bool   `json:"Candidate,omitempty"`
	Running   bool   `json:"Running,omitempty"`
}

type GetOverrideOutput struct {
	Response GetOverrideOutputResponse `json:"Response"`
	Status   api.Status                `json:"ResponseStatus"`
}

func (o GetOverrideOutput) Failed() *api.Status {
	return o.Status.Failed()
}

type GetOverrideOutputResponse struct {
	Rulestack string          `json:"RuleStackName"`
	Name      string          `json:"Name"`
	Candidate OverrideDetails `json:"URLCategoryCandidate"`
	Running   OverrideDetails `json:"URLCategoryRunning"`
}

type OverrideDetails struct {
	Action       string `json:"Action"`
	AuditComment string `json:"AuditComment"`
	UpdateToken  string `json:"UpdateToken"`
}

// v1 override.

type OverrideInput struct {
	Rulestack    string `json:"-"`
	Name         string `json:"-"`
	UpdateToken  string `json:"UpdateToken,omitempty"`
	Action       string `json:"Action"`
	AuditComment string `json:"AuditComment,omitempty"`
}
