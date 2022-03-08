package security

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/tag"
)

// V1 list.

type ListInput struct {
	Rulestack   string `json:"-"`
	RuleList    string `json:"-"`
	NextToken   string `json:"NextToken,omitempty"`
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
	Rulestack   string               `json:"RuleStackName"`
	RuleList    string               `json:"RuleListName"`
	Candidates  []ListEntryCandidate `json:"RuleEntryCandidate,omitempty"`
	Running     []ListEntryCandidate `json:"RuleEntryRunning,omitempty"`
	Uncommitted []ListEntryCandidate `json:"RuleEntryUncommitted,omitempty"`
	NextToken   string               `json:"NextToken"`
}

type ListEntryCandidate struct {
	Name     string `json:"RuleName"`
	Priority int    `json:"Priority"`
	// Only uncommitted has this field.
	Operation string `json:"Operation,omitempty"`
}

// V1 create / update.

type Info struct {
	Rulestack string  `json:"RuleStackName,omitempty"`
	RuleList  string  `json:"RuleListName,omitempty"`
	Priority  int     `json:"Priority,omitempty"`
	Entry     Details `json:"RuleEntry"`
}

type Details struct {
	Name               string             `json:"RuleName,omitempty"`
	Description        string             `json:"Description,omitempty"`
	Enabled            bool               `json:"Enabled,omitempty"`
	Source             SourceDetails      `json:"Source"`
	NegateSource       bool               `json:"NegateSource,omitempty"`
	Destination        DestinationDetails `json:"Destination"`
	NegateDestination  bool               `json:"NegateDestination,omitempty"`
	Applications       []string           `json:"Applications"`
	Category           CategoryDetails    `json:"Category"`
	Protocol           string             `json:"Protocol,omitempty"`
	AuditComment       string             `json:"AuditComment,omitempty"`
	Action             string             `json:"Action,omitempty"`
	Logging            bool               `json:"Logging,omitempty"`
	DecryptionRuleType string             `json:"DecryptionRuleType,omitempty"`
	Tags               []tag.Details      `json:"Tags,omitempty"`
	UpdateToken        string             `json:"UpdateToken,omitempty"`
}

type SourceDetails struct {
	Cidrs       []string `json:"Cidrs,omitempty"`
	Countries   []string `json:"Countries,omitempty"`
	Feeds       []string `json:"Feeds,omitempty"`
	PrefixLists []string `json:"PrefixLists,omitempty"`
}

type DestinationDetails struct {
	Cidrs       []string `json:"Cidrs,omitempty"`
	Countries   []string `json:"Countries,omitempty"`
	Feeds       []string `json:"Feeds,omitempty"`
	PrefixLists []string `json:"PrefixLists,omitempty"`
	FqdnLists   []string `json:"FqdnLists,omitempty"`
}

type CategoryDetails struct {
	UrlCategoryNames []string `json:"URLCategoryNames,omitempty"`
	Feeds            []string `json:"Feeds,omitempty"`
}

// V1 read.

type ReadInput struct {
	Rulestack string `json:"-"`
	RuleList  string `json:"-"`
	Priority  int    `json:"-"`
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
	Rulestack string   `json:"RuleStackName"`
	RuleList  string   `json:"RuleListName"`
	Priority  int      `json:"Priority"`
	Running   *Details `json:"RuleEntryRunning"`
	Candidate *Details `json:"RuleEntryCandidate"`
}
