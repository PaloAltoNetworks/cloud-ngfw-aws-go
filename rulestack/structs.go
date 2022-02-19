package rulestack

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
)

// V1 list.

type ListInput struct {
	Scope       string `json:"Scope,omitempty"`
	Candidate   bool   `json:"Candidate,omitempty"`
	Running     bool   `json:"Running,omitempty"`
	Uncommitted bool   `json:"Uncommitted,omitempty"`
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
	Candidates  []string          `json:"RuleStackCandidate"`
	Running     []string          `json:"RuleStackRunning"`
	Uncommitted []ListUncommitted `json:"RuleStackUncommitted"`
	NextToken   string            `json:"NextToken"`
}

type ListUncommitted struct {
	Name      string `json:"RuleStackName"`
	Operation string `json:"Operation"`
}

// V1 create / update.

type Info struct {
	Name  string  `json:"RuleStackName,omitempty"`
	Entry Details `json:"RuleStackEntry"`
}

type Details struct {
	Description         string        `json:"Description,omitempty"`
	Scope               string        `json:"Scope,omitempty"`
	AccountId           string        `json:"AccountId,omitempty"`
	AccountGroup        string        `json:"AccountGroup,omitempty"`
	MinimumAppIdVersion string        `json:"MinAppIdVersion,omitempty"`
	Profile             ProfileConfig `json:"Profiles"`

	UpdateToken string `json:"UpdateToken,omitempty"`

	Tags []string `json:"Tags,omitempty"`
}

type ProfileConfig struct {
	AntiSpyware   string `json:"AntiSpywareProfile,omitempty"`
	AntiVirus     string `json:"AntiVirusProfile,omitempty"`
	Vulnerability string `json:"VulnerabilityProfile,omitempty"`
	UrlFiltering  string `json:"URLFilteringProfile,omitempty"`
	FileBlocking  string `json:"FileBlockingProfile,omitempty"`

	OutboundTrustCertificate   string `json:"OutboundTrustCertificate,omitempty"`
	OutboundUntrustCertificate string `json:"OutboundUntrustCertificate,omitempty"`
}

// V1 read.

type ReadInput struct {
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
	Name      string   `json:"RuleStackName"`
	State     string   `json:"RuleStackState"`
	Candidate *Details `json:"RuleStackCandidate"`
	Running   *Details `json:"RuleStackRunning"`
}
