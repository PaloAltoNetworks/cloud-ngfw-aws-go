package url

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
)

// V1 list.

type ListInput struct {
	Rulestack   string `json:"-"`
	Scope       string `json:"-"`
	Candidate   bool   `json:"Candidate,omitempty"`
	Running     bool   `json:"Running,omitempty"`
	Uncommitted bool   `json:"Uncommitted,omitempty"`
	NextToken   string `json:"NextToken,omitempty"`
	MaxResults  int    `json:"MaxResults,omitempty"`
}

type ListOutput struct {
	Response *ListOutputDetails `json:"Response"`
	Status   api.Status         `json:"ResponseStatus"`
}

func (o ListOutput) Failed() *api.Status {
	return o.Status.Failed()
}

type ListOutputDetails struct {
	Rulestack   string            `json:"RuleStackName"`
	Candidates  []string          `json:"CategoriesCandidate"`
	Running     []string          `json:"CategoriesRunning"`
	Uncommitted []ListUncommitted `json:"CategoriesUncommitted"`
	NextToken   string            `json:"NextToken"`
}

type ListUncommitted struct {
	Name      string `json:"Name"`
	Operation string `json:"Operation"`
}

// V1 create / update.

type Info struct {
	Rulestack    string   `json:"-"`
	Scope        string   `json:"-"`
	Name         string   `json:"Name,omitempty"`
	Description  string   `json:"Description,omitempty"`
	UrlList      []string `json:"URLTargets,omitempty"`
	Action       string   `json:"Action,omitempty"`
	AuditComment string   `json:"AuditComment,omitempty"`
	UpdateToken  string   `json:"UpdateToken,omitempty"`
}

// V1 read.

type ReadInput struct {
	Rulestack string `json:"-"`
	Scope     string `json:"-"`
	Name      string `json:"-"`
	Candidate bool   `json:"Candidate,omitempty"`
	Running   bool   `json:"Running,omitempty"`
}

type ReadOutput struct {
	Response *ReadResponse `json:"Response"`
	Status   api.Status    `json:"ResponseStatus"`
}

func (o ReadOutput) Failed() *api.Status {
	return o.Status.Failed()
}

type ReadResponse struct {
	Rulestack string `json:"RuleStackName"`
	Name      string `json:"Name"`
	Candidate *Info  `json:"URLCategoryCandidate"`
	Running   *Info  `json:"URLCategoryRunning"`
}

// V1 delete.

type DeleteInput struct {
	Rulestack string
	Name      string
	Scope     string
}
