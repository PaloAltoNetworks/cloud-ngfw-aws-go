package stack

import (
	"fmt"
	"strings"
	"sync"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/certificate"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/feed"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/fqdn"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/predefinedurl"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/prefix"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/response"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/security"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/tag"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/url"
)

type Rulestack struct {
	Info              *Info
	Feed              map[string]*feed.Info
	Certificate       map[string]*certificate.Info
	CustomURLCategory map[string]*url.Info
	FqdnList          map[string]*fqdn.Info
	PrefixList        map[string]*prefix.Info
	PredefinedURL     map[string]*predefinedurl.OverrideInput
	PreRules          map[int]*security.Info
	PostRules         map[int]*security.Info
	LocalRules        map[int]*security.Info
	State             string
	Commit            bool
	sync.RWMutex
}

// V1 list.

type ListInput struct {
	Scope       string
	TagKey      string
	TagValue    string
	Candidate   bool
	Running     bool
	Uncommitted bool
	Describe    bool
	NextToken   string
	MaxResults  int
}

type ListOutput struct {
	Response *ListOutputDetails `json:"Response"`
	Status   response.Status    `json:"ResponseStatus"`
}

func (o ListOutput) Failed() *response.Status {
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
	LookupXForwardedFor string        `json:"LookupXForwardedFor,omitempty"`
	Profile             ProfileConfig `json:"Profiles"`

	UpdateToken string `json:"UpdateToken,omitempty"`

	Tags []tag.Details `json:"Tags,omitempty"`
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
	Name      string
	Scope     string
	Candidate bool
	Running   bool
}

type ReadOutput struct {
	Response *ReadResponse   `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

func (o ReadOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type ReadResponse struct {
	Name      string   `json:"RuleStackName"`
	State     string   `json:"RuleStackState"`
	Candidate *Details `json:"RuleStackCandidate"`
	Running   *Details `json:"RuleStackRunning"`
}

// V1 commit status.

type CommitStatus struct {
	Response CommitResponse  `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

func (o CommitStatus) Failed() *response.Status {
	return o.Status.Failed()
}

func (c CommitStatus) CommitErrors() string {
	var b strings.Builder
	b.Grow(50 * len(c.Response.CommitMessages))

	b.WriteString(fmt.Sprintf("Commit(%s):", c.Response.CommitStatus))
	for i, x := range c.Response.CommitMessages {
		if i != 0 {
			b.WriteString(" |")
		}
		b.WriteString(" ")
		b.WriteString(x)
	}

	return b.String()
}

type CommitResponse struct {
	Name               string   `json:"RuleStackName"`
	CommitStatus       string   `json:"CommitStatus"`
	ValidationStatus   string   `json:"ValidateStatus"`
	CommitMessages     []string `json:"CommitMessages"`
	ValidationMessages []string `json:"ValidateMessages"`
}

// V1 list tags.

type ListTagsInput struct {
	Rulestack  string
	Scope      string
	NextToken  string
	MaxResults int
}

type ListTagsOutput struct {
	Response ListTagsOutputDetails `json:"Response"`
	Status   response.Status       `json:"ResponseStatus"`
}

func (o ListTagsOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type ListTagsOutputDetails struct {
	Rulestack string        `json:"ResourceName"`
	NextToken string        `json:"NextToken"`
	Tags      []tag.Details `json:"Tags"`
}

// V1 create.

type AddTagsInput struct {
	Rulestack string        `json:"-"`
	Scope     string        `json:"-"`
	Tags      []tag.Details `json:"Tags,omitempty"`
}

type RemoveTagsInput struct {
	Rulestack string   `json:"-"`
	Scope     string   `json:"-"`
	Tags      []string `json:"TagKeys"`
}

// V1 delete.

type SimpleInput struct {
	Name  string
	Scope string
}

// v1 export rulestack xml
type ExportRulestackXmlOutput struct {
	Response string          `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

func (o ExportRulestackXmlOutput) Failed() *response.Status {
	return o.Status.Failed()
}

// v1 save rulestack xml
type SaveRulestackXmlInput struct {
	Name              string          `json:"-"`
	FirewallId        string          `json:"-"`
	Scope             string          `json:"-"`
	RuleStackEntryXml XmlString       `json:"RuleStackEntryXml"`
	ReadOnly          bool            `json:"ReadOnly,omitempty"`
	Firewalls         []FirewallEntry `json:"Firewalls,omitempty"`
	RulesStackType    string          `json:"RulesStackType,omitempty"`
}

type XmlString struct {
	Xml string `json:"Xml"`
}

type SaveRulestackXmlOutput struct {
	Response string          `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

func (o SaveRulestackXmlOutput) Failed() *response.Status {
	return o.Status.Failed()
}

// v1 Create SCM rulestack
type CreateSCMRuleStackInput struct {
	Name              string          `json:"-"`
	Scope             string          `json:"-"`
	RuleStackEntryXml XmlString       `json:"RuleStackEntryXml"`
	Firewalls         []FirewallEntry `json:"Firewalls"`
}

type FirewallEntry struct {
	Firewall   string `json:"Firewall"`
	AccountId  string `json:"AccountId,omitempty"`
	ArmId      string `json:"ArmId,omitempty"`
	FirewallId string `json:"FirewallId,omitempty"`
}

type PushRulestackCMOutput struct {
	Response string          `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

func (o PushRulestackCMOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type AuthInput struct {
	ExternalId string `json:"ExternalId"`
	ExpiryTime int    `json:"ExpiryTime"`
	AuthURL    string `json:"AuthURL"`
}

type AuthOutput struct {
	Response AuthOutputDetails `json:"Response,omitempty"`
	Status   response.Status   `json:"ResponseStatus"`
}

type AuthOutputDetails struct {
	TokenId         string `json:"TokenId"`
	SubscriptionKey string `json:"SubscriptionKey"`
	ExpiryTime      int    `json:"ExpiryTime"`
	Enabled         bool   `json:"Enabled"`
}

func (o AuthOutput) Failed() *response.Status {
	return o.Status.Failed()
}
