package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/firewall"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/tag"
)

func ConvertToUTCEpoch(timestamp string) (int64, error) {
	tsList := strings.Split(timestamp, " ")
	zoneInfo := tsList[1]
	layout := fmt.Sprintf("2006-01-02T15:04:05 %s", zoneInfo)
	t, err := time.Parse(layout, timestamp)
	if err != nil {
		return 0, err
	}
	unixTimeStamp := t.UTC().Unix()
	return unixTimeStamp, nil
}

func Contains(inputStr string, strList []string) bool {
	for _, str := range strList {
		if inputStr == str {
			return true
		}
	}
	return false
}

func WaitForOperation(ctx context.Context, op func(ctx context.Context) (bool, error)) error {
	var err error
	for i := 0; i < 120; i++ {
		shouldRetry, err := op(ctx)
		if err == nil {
			return nil
		}
		if !shouldRetry {
			return err
		}
		time.Sleep(30 * time.Second)
	}
	return fmt.Errorf("operation timed out: %v", err)
}

func SliceToMap(s []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}

func TagMap(tags []tag.Details) map[string]string {
	m := make(map[string]string)
	for _, tag := range tags {
		m[tag.Key] = tag.Value
	}
	return m
}

func CompareMapString(a, b map[string]string) bool {
	if a == nil && b != nil {
		return false
	}
	if a != nil && b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if _, ok := b[k]; !ok {
			return false
		}
		if b[k] != v {
			return false
		}
	}
	return true
}

func EpSubnetMap(eps []firewall.EndpointConfig) map[string]firewall.EndpointConfig {
	res := make(map[string]firewall.EndpointConfig)
	for _, ep := range eps {
		if ep.SubnetId == "" {
			continue
		}
		res[ep.SubnetId] = ep
	}
	return res
}

func EpIdMap(eps []firewall.EndpointConfig) map[string]firewall.EndpointConfig {
	res := make(map[string]firewall.EndpointConfig)
	for _, ep := range eps {
		if ep.EndpointId == "" {
			continue
		}
		res[ep.EndpointId] = ep
	}
	return res
}

func CompareEpSubnetMap(a, b map[string]firewall.EndpointConfig) bool {
	if a == nil && b != nil {
		return false
	}
	if a != nil && b == nil {
		return false
	}
	if a == nil && b == nil {
		return true
	}
	for k := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			return false
		}
	}
	return true
}
