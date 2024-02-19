package prefix

import "github.com/paloaltonetworks/cloud-ngfw-aws-go/api/response"

// V1 list.

type ListInput struct {
	Rulestack   string `json:"-"`
	Scope       string `json:"-"`
	Candidate   bool   `json:"Candidate"`
	Running     bool   `json:"Running,omitempty"`
	Uncommitted bool   `json:"Uncommitted,omitempty"`
	MaxResults  int    `json:"MaxResults,omitempty"`
	NextToken   string `json:"NextToken,omitempty"`
}

type ListOutput struct {
	Response *ListOutputDetails `json:"Response"`
	Status   response.Status    `json:"ResponseStatus"`
}

func (o ListOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type ListOutputDetails struct {
	Rulestack   string            `json:"RuleStackName"`
	Candidates  []string          `json:"PrefixListCandidate"`
	Running     []string          `json:"PrefixListRunning"`
	Uncommitted []ListUncommitted `json:"PrefixListUncommitted"`
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
	PrefixList   []string `json:"PrefixList,omitempty"`
	AuditComment string   `json:"AuditComment,omitempty"`
	UpdateToken  string   `json:"UpdateToken,omitempty"`
	DeleteFlag   bool     `json:"-"`
}

// V1 read.

type ReadInput struct {
	Rulestack string `json:"-"`
	Scope     string `json:"-"`
	Name      string `json:"-"`
	Candidate bool   `json:"Candidate,omitempty"`
	Running   bool   `json:"Running,omitempty"`
	NextToken string `json:"NextToken,omitempty"`
	MaxResult string `json:"MaxResult,omitempty"`
}

type ReadOutput struct {
	Response *ReadResponse   `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

func (o ReadOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type ReadResponse struct {
	Rulestack string `json:"RuleStackName"`
	Name      string `json:"Name"`
	NextToken string `json:"NextToken"`
	Candidate *Info  `json:"PrefixListCandidate"`
	Running   *Info  `json:"PrefixListRunning"`
}

// V1 delete.

type DeleteInput struct {
	Rulestack string
	Scope     string
	Name      string
}
