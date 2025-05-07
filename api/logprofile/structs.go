package logprofile

import "github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/response"

// V1 create / update.

type UpdateInput struct {
	Firewall string `json:"-"`
}

type Info struct {
	LogConfig                 *LogConfig         `json:"LogConfig,omitempty"`
	CloudwatchMetrics         *CloudwatchMetrics `json:"CloudwatchMetrics,omitempty"`
	UpdateToken               string             `json:"UpdateToken,omitempty"`
	Region                    string             `json:"Region,omitempty"`
	AccountId                 string             `json:"AccountId,omitempty"`
	Firewall                  string             `json:"FirewallName,omitempty"`
	FirewallId                string             `json:"FirewallId,omitempty"`
	CloudWatchMetricNamespace string             `json:"CloudWatchMetricNamespace,omitempty"`
	AdvancedThreatLog         bool               `json:"AdvancedThreatLog,omitempty"`
	CloudWatchMetricsFields   []string           `json:"CloudWatchMetricsFields,omitempty"`
}

type LogConfig struct {
	AccountId          string   `json:"AccountId,omitempty"`
	LogDestination     string   `json:"LogDestination,omitempty"`
	LogDestinationType string   `json:"LogDestinationType,omitempty"`
	LogType            []string `json:"LogType,omitempty"`
	RoleType           string   `json:"RoleType,omitempty"`
}

type CloudwatchMetrics struct {
	Namespace string   `json:"CloudWatchMetricNamespace,omitempty"`
	Metrics   []string `json:"CloudWatchMetricsFields,omitempty"`
	AccountId string   `json:"AccountId,omitempty"`
}

// V1 read.

type ReadInput struct {
	Firewall   string `json:"Firewall,omitempty"`
	FirewallId string `json:"FirewallId"`
}

type ReadOutput struct {
	Response *Info           `json:"Response"`
	Status   response.Status `json:"ResponseStatus"`
}

func (o ReadOutput) Failed() *response.Status {
	return o.Status.Failed()
}
