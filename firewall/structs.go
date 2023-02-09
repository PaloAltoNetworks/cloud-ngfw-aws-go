package firewall

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/tag"
)

// V1 list.

type ListInput struct {
	Rulestack  string   `json:"-"`
	MaxResults int      `json:"MaxResults,omitempty"`
	NextToken  string   `json:"NextToken,omitempty"`
	VpcIds     []string `json:"VpcIds,omitempty"`
}

type ListOutput struct {
	Response ListOutputDetails `json:"Response"`
	Status   api.Status        `json:"ResponseStatus"`
}

func (o ListOutput) Failed() *api.Status {
	return o.Status.Failed()
}

type ListOutputDetails struct {
	Firewalls []ListFirewall `json:"Firewalls"`
	NextToken string         `json:"NextToken"`
}

type ListFirewall struct {
	Name      string `json:"FirewallName"`
	AccountId string `json:"AccountId"`
}

// V1 create.

type Info struct {
	Name                         string          `json:"FirewallName,omitempty"`
	AccountId                    string          `json:"AccountId,omitempty"`
	VpcId                        string          `json:"VpcId,omitempty"`
	AppIdVersion                 string          `json:"AppIdVersion,omitempty"`
	Description                  string          `json:"Description,omitempty"`
	Rulestack                    string          `json:"RuleStackName,omitempty"`
	GlobalRulestack              string          `json:"GlobalRuleStackName,omitempty"`
	MultiVpc                     bool            `json:"MultiVpcEnable"`
	EndpointMode                 string          `json:"EndpointMode,omitempty"`
	EndpointServiceName          string          `json:"EndpointServiceName,omitempty"`
	AutomaticUpgradeAppIdVersion bool            `json:"AutomaticUpgradeAppIdVersion,omitempty"`
	SubnetMappings               []SubnetMapping `json:"SubnetMappings,omitempty"`
	Tags                         []tag.Details   `json:"Tags,omitempty"`
	UpdateToken                  string          `json:"UpdateToken,omitempty"`
}

type SubnetMapping struct {
	SubnetId           string `json:"SubnetId,omitempty"`
	AvailabilityZone   string `json:"AvailabilityZone,omitempty"`
	AvailabilityZoneId string `json:"AvailabilityZoneId,omitempty"`
}

type CreateOutput struct {
	Response Info       `json:"Response"`
	Status   api.Status `json:"ResponseStatus"`
}

func (o CreateOutput) Failed() *api.Status {
	return o.Status.Failed()
}

// V1 update description.

type UpdateDescriptionInput struct {
	Firewall    string `json:"-"`
	AccountId   string `json:"AccountId,omitempty"`
	Description string `json:"Description,omitempty"`
	UpdateToken string `json:"UpdateToken,omitempty"`
}

// V1 update subnet mappings.

type UpdateSubnetMappingsInput struct {
	Firewall                   string          `json:"-"`
	AccountId                  string          `json:"AccountId,omitempty"`
	AssociateSubnetMappings    []SubnetMapping `json:"AssociateSubnetMappings,omitempty"`
	DisassociateSubnetMappings []SubnetMapping `json:"DisassociateSubnetMappings,omitempty"`
	MultiVpc                   bool            `json:"MultiVpcEnable"`
	UpdateToken                string          `json:"UpdateToken,omitempty"`
}

// V1 update content version.

type UpdateContentVersionInput struct {
	Firewall                     string `json:"-"`
	AccountId                    string `json:"AccountId,omitempty"`
	AppIdVersion                 string `json:"AppIdVersion,omitempty"`
	AutomaticUpgradeAppIdVersion bool   `json:"AutomaticUpgradeAppIdVersion,omitempty"`
	UpdateToken                  string `json:"UpdateToken,omitempty"`
}

// V1 update rulestack.

type UpdateRulestackInput struct {
	Firewall    string `json:"-"`
	AccountId   string `json:"AccountId"`
	Rulestack   string `json:"RuleStackName"`
	UpdateToken string `json:"UpdateToken,omitempty"`
}

// V1 remove tags.

type RemoveTagsInput struct {
	Firewall  string   `json:"-"`
	AccountId string   `json:"AccountId"`
	Tags      []string `json:"TagKeys"`
}

// V1 add tags.

type AddTagsInput struct {
	Firewall  string        `json:"-"`
	AccountId string        `json:"AccountId"`
	Tags      []tag.Details `json:"Tags"`
}

// V1 read.

type ReadInput struct {
	Name      string `json:"-"`
	AccountId string `json:"AccountId,omitempty"`
}

type ReadOutput struct {
	Response ReadResponse `json:"Response"`
	Status   api.Status   `json:"ResponseStatus"`
}

func (o ReadOutput) Failed() *api.Status {
	return o.Status.Failed()
}

type ReadResponse struct {
	Firewall Info           `json:"Firewall"`
	Status   FirewallStatus `json:"Status,omitempty"`
}

type FirewallStatus struct {
	FirewallStatus  string       `json:"FirewallStatus,omitempty"`
	FailureReason   string       `json:"FailureReason,omitempty"`
	RulestackStatus string       `json:"RuleStackStatus,omitempty"`
	Attachments     []Attachment `json:"Attachments,omitempty"`
}

type Attachment struct {
	EndpointId     string `json:"EndpointId,omitempty"`
	Status         string `json:"Status,omitempty"`
	RejectedReason string `json:"RejectedReason,omitempty"`
	SubnetId       string `json:"SubnetId,omitempty"`
}

// V1 delete.

type DeleteInput struct {
	Name      string `json:"-"`
	AccountId string `json:"AccountId,omitempty"`
}

// V1 list tags.

type ListTagsInput struct {
	Firewall   string
	AccountId  string
	NextToken  string
	MaxResults int
}

type ListTagsOutput struct {
	Response ListTagsOutputDetails `json:"Response"`
	Status   api.Status            `json:"ResponseStatus"`
}

func (o ListTagsOutput) Failed() *api.Status {
	return o.Status.Failed()
}

type ListTagsOutputDetails struct {
	Firewall  string        `json:"ResourceName"`
	NextToken string        `json:"NextToken"`
	Tags      []tag.Details `json:"Tags"`
}
