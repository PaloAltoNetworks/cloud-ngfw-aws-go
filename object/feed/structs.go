package feed

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
)

// V1 list.

type ListInput struct {
	Rulestack   string `json:"-"`
	Candidate   bool   `json:"Candidate,omitempty"`
	Running     bool   `json:"Running,omitempty"`
	Uncommitted bool   `json:"Uncommitted,omitempty"`
	Type        string `json:"Type,omitempty"`
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
	Candidates  []string          `json:"FeedCandidate"`
	Running     []string          `json:"FeedRunning"`
	Uncommitted []ListUncommitted `json:"FeedUncommitted"`
	NextToken   string            `json:"NextToken"`
}

type ListUncommitted struct {
	Name      string `json:"Name"`
	Operation string `json:"Operation"`
}

// V1 create / update.

type Info struct {
	Rulestack    string `json:"-"`
	Name         string `json:"Name,omitempty"`
	Description  string `json:"Description,omitempty"`
	Certificate  string `json:"Certificate,omitempty"`
	Url          string `json:"FeedURL"`
	Type         string `json:"Type"`
	Frequency    string `json:"Frequency"`
	Time         int    `json:"Time"`
	AuditComment string `json:"AuditComment,omitempty"`
	UpdateToken  string `json:"UpdateToken,omitempty"`
}

// V1 read.

type ReadInput struct {
	Rulestack string `json:"-"`
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
	Candidate *Info  `json:"FeedCandidate"`
	Running   *Info  `json:"FeedRunning"`
}
