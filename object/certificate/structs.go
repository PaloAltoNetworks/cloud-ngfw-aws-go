package certificate

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
)

// V1 list.

type ListInput struct {
	Rulestack   string `json:"-"`
	Candidate   bool   `json:"Candidate,omitempty"`
	Running     bool   `json:"Running,omitempty"`
	Uncommitted bool   `json:"Uncommitted,omitempty"`
	MaxResults  int    `json:"MaxResults,omitempty"`
	NextToken   string `json:"NextToken,omitempty"`
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
	Candidates  []string          `json:"CertificateObjectCandidate"`
	Running     []string          `json:"CertificateObjectRunning"`
	Uncommitted []ListUncommitted `json:"CertificateObjectUncommitted"`
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
	SignerArn    string `json:"CertificateSignerArn,omitempty"`
	SelfSigned   bool   `json:"CertificateSelfSigned,omitempty"`
	AuditComment string `json:"AuditComment,omitempty"`
	UpdateToken  string `json:"UpdateToken,omitempty"`
}

// V1 read.

type ReadInput struct {
	Rulestack string `json:"-"`
	Name      string `json:"-"`
	Candidate bool   `json:"Candidate,omitempty"`
	Running   bool   `json:"Running,omitempty"`
	NextToken string `json:"NextToken,omitempty"`
	MaxResult string `json:"MaxResult,omitempty"`
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
	NextToken string `json:"NextToken"`
	Candidate *Info  `json:"CertificateObjectCandidate"`
	Running   *Info  `json:"CertificateObjectRunning"`
}
