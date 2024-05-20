package logprofile

import "github.com/paloaltonetworks/cloud-ngfw-aws-go/api/response"

// V1 create / update.

type UpdateInput struct {
	Firewall string `json:"-"`
}

type Info struct {
	AccountId                 string           `json:"AccountId,omitempty"`
	Firewall                  string           `json:"FirewallName,omitempty"`
	LogDestinations           []LogDestination `json:"LogDestinationConfigs"`
	CloudWatchMetricNamespace string           `json:"CloudWatchMetricNamespace,omitempty"`
	AdvancedThreatLog         bool             `json:"AdvancedThreatLog,omitempty"`
	CloudWatchMetricsFields   []string         `json:"CloudWatchMetricsFields,omitempty"`
}

type LogDestination struct {
	Destination     string `json:"LogDestination,omitempty"`
	DestinationType string `json:"LogDestinationType,omitempty"`
	LogType         string `json:"LogType,omitempty"`
}

// V1 read.

type ReadInput struct {
	Firewall  string `json:"-"`
	AccountId string `json:"AccountId,omitempty"`
}

type ReadOutput struct {
	Response *Info           `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

func (o ReadOutput) Failed() *response.Status {
	return o.Status.Failed()
}
