package appid

import (
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/response"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/tag"
)

// V1 list.

type ListInput struct {
	NextToken  string `json:"NextTovi ken,omitempty"`
	MaxResults int    `json:"MaxResults,omitempty"`
}

type ListOutput struct {
	Response ListOutputDetails `json:"Response"`
	Status   response.Status   `json:"ResponseStatus"`
}

func (o ListOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type ListOutputDetails struct {
	Versions  []string `json:"AppIdVersions"`
	NextToken string   `json:"NextToken"`
}

// V1 read app-id.

type ReadInput struct {
	Version    string `json:"-"`
	MaxResults int    `json:"MaxResults,omitempty"`
	NextToken  string `json:"NextToken,omitempty"`
}

type ReadOutput struct {
	Response ReadOutputDetails `json:"Response"`
	Status   response.Status   `json:"ResponseStatus"`
}

func (o ReadOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type ReadOutputDetails struct {
	Version      string   `json:"AppIdVersion"`
	Applications []string `json:"Applications"`
	NextToken    string   `json:"NextToken"`
}

// V1 read app-id application.

type ReadApplicationOutput struct {
	Response ApplicationOutputDetails `json:"Response"`
	Status   response.Status          `json:"ResponseStatus"`
}

func (o ReadApplicationOutput) Failed() *response.Status {
	return o.Status.Failed()
}

type ApplicationOutputDetails struct {
	Name    string             `json:"Name"`
	Details ApplicationDetails `json:"AppIdEntry"`
}

type ApplicationDetails struct {
	Description            string                     `json:"Description"`
	Properties             ApplicationProperties      `json:"Properties"`
	Characteristics        ApplicationCharacteristics `json:"Characteristics"`
	Options                ApplicationOptions         `json:"Options"`
	StandardPorts          []string                   `json:"StandardPorts"`
	AdditionalInfo         []tag.Details              `json:"AdditionalInformations"`
	DependsOn              []string                   `json:"DependsOn"`
	ImplicitlyUses         []string                   `json:"ImplicitlyUses"`
	PreviouslyIdentifiedAs string                     `json:"PreviouslyIdentifiedAs"`
}

type ApplicationProperties struct {
	Category    string `json:"Category"`
	Subcategory string `json:"Subcategory"`
	Technology  string `json:"Technology"`
	ParentApp   string `json:"ParentApp"`
	Risk        int    `json:"Risk"`
}

type ApplicationCharacteristics struct {
	Evasive                 bool   `json:"Evasive"`
	UsedByMalware           bool   `json:"UsedByMalware"`
	ExcessiveBandwidth      bool   `json:"ExcessiveBandwidth"`
	CapableFileTransfer     bool   `json:"CapableFileTransfer"`
	HasKnownVulnerability   bool   `json:"HasKnownVulnerability"`
	TunnelOtherApplications bool   `json:"TunnelOtherApplications"`
	ProneToMisuse           bool   `json:"ProneToMisuse"`
	WidelyUsed              bool   `json:"WidelyUsed"`
	Saas                    bool   `json:"SaaS"`
	DenyAction              string `json:"DenyAction"`
}

type ApplicationOptions struct {
	TcpTimeout   int  `json:"TCPTimeout"`
	Alg          bool `json:"ALG"`
	AppIdEnabled bool `json:"AppIdEnabled"`
}
