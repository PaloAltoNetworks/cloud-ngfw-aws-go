package firewall

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
)

// V1 list.

type ListInput struct {
	MaxResults int      `json:"MaxResults,omitempty"`
	NextToken  string   `json:"NextToken,omitempty"`
	VpcIds     []string `json:"VpcIds,omitempty"`
}

type ListOutput struct {
	Response *ListOutputDetails `json:"Response"`
	Status   api.Status         `json:"ResponseStatus"`
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

// V1 create / update.

type Info struct {
	Name                         string          `json:"FirewallName,omitempty"`
	VpcId                        string          `json:"VpcId,omitempty"`
	AccountId                    string          `json:"AccountId,omitempty"`
	Description                  string          `json:"Description,omitempty"`
	EndpointMode                 string          `json:"EndpointMode,omitempty"`
	SubnetMappings               []SubnetMapping `json:"SubnetMappings,omitempty"`
	AppIdVersion                 string          `json:"AppIdVersion,omitempty"`
	AutomaticUpgradeAppIdVersion bool            `json:"AutomaticUpgradeAppIdVersion,omitempty"`
	RuleStackName                string          `json:"RuleStackName,omitempty"`
	GlobalRuleStackName          string          `json:"GlobalRuleStackName,omitempty"`
	Tags                         []TagDetails    `json:"Tags,omitempty"`

	AssociateSubnetMappings    []SubnetMapping `json:"AssociateSubnetMappings,omitempty"`
	DisassociateSubnetMappings []SubnetMapping `json:"DisassociateSubnetMappings,omitempty"`

	UpdateToken string `json:"UpdateToken,omitempty"`
}

type SubnetMapping struct {
	SubnetId           string `json:"SubnetId,omitempty"`
	AvailabilityZone   string `json:"AvailabilityZone,omitempty"`
	AvailabilityZoneId string `json:"AvailabilityZoneId,omitempty"`
}

type TagDetails struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type CreateOutput struct {
	Response Firewall   `json:"Response"`
	Status   api.Status `json:"ResponseStatus"`
}

func (o CreateOutput) Failed() *api.Status {
	return o.Status.Failed()
}

// V1 read.

type ReadInput struct {
	Name      string `json:"-"`
	AccountId string `json:"AccountId,omitempty"`
}

type ReadOutput struct {
	Response *ReadResponse `json:"Response"`
	Status   api.Status    `json:"ResponseStatus"`
}

func (o ReadOutput) Failed() *api.Status {
	return o.Status.Failed()
}

type ReadResponse struct {
	Firewall Firewall        `json:"Firewall,omitempty"`
	Status   *FirewallStatus `json:"Status,omitempty"`
}

type Firewall struct {
	Name                         string          `json:"FirewallName,omitempty"`
	AccountId                    string          `json:"AccountId,omitempty"`
	SubnetMappings               []SubnetMapping `json:"SubnetMappings,omitempty"`
	VpcId                        string          `json:"VpcId,omitempty"`
	AppIdVersion                 string          `json:"AppIdVersion,omitempty"`
	Description                  string          `json:"Description,omitempty"`
	RuleStackName                string          `json:"RuleStackName,omitempty"`
	GlobalRuleStackName          string          `json:"GlobalRuleStackName,omitempty"`
	EndpointServiceName          string          `json:"EndpointServiceName,omitempty"`
	EndpointMode                 string          `json:"EndpointMode,omitempty"`
	AutomaticUpgradeAppIdVersion bool            `json:"AutomaticUpgradeAppIdVersion,omitempty"`
	Tags                         []TagDetails    `json:"Tags,omitempty"`
	UpdateToken                  string          `json:"UpdateToken,omitempty"`
}

type FirewallStatus struct {
	FirewallStatus  string       `json:"FirewallStatus,omitempty"`
	FailureReason   string       `json:"FailureReason,omitempty"`
	RuleStackStatus string       `json:"RuleStackStatus,omitempty"`
	Attachments     []Attachment `json:"Attachments,omitempty"`
}

type Attachment struct {
	EndpointId     string `json:"EndpointId,omitempty"`
	Status         string `json:"Status,omitempty"`
	RejectedReason string `json:"RejectedReason,omitempty"`
	SubnetId       string `json:"SubnetId,omitempty"`
}
