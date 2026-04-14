package aws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	cloudngfwgosdk "github.com/paloaltonetworks/cloud-ngfw-aws-go/v2"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/firewall"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/response"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/tag"
)

type updateFirewall struct {
	GeneralUpdate bool
	DrsCommit     bool
	RsCommit      bool
}

func getSchemaVersion(ctx context.Context) string {
	schemaVersion := cloudngfwgosdk.SchemaVersionV2
	v := ctx.Value("SchemaVersion")
	if v != nil {
		schemaVersion = v.(string)
	}
	return schemaVersion
}

// List returns a list of firewalls.
func (c *Client) ListFirewall(ctx context.Context, input firewall.ListInput) (firewall.ListOutput, error) {
	if len(input.VpcIds) == 0 {
		c.Log(http.MethodGet, "list NGFirewalls in all the VPCs")
	} else {
		c.Log(http.MethodGet, "list NGFirewalls in %q VPCs", strings.Join(input.VpcIds, ","))
	}

	uv := url.Values{}
	if input.Rulestack != "" {
		uv.Set("rulestackname", input.Rulestack)
	}
	if input.Describe {
		uv.Set("describe", "true")
	}
	if input.Region != "" {
		uv.Set("region", input.Region)
	}
	if input.MaxResults != 0 {
		maxResults := strconv.Itoa(input.MaxResults)
		uv.Set("maxresults", maxResults)
	}
	c.Log(http.MethodGet, "list firewalls, tenant version: %s", c.TenantVersion)
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls"},
		V2Path: []string{"v2", "config", "ngfirewalls"},
	}
	var ans firewall.ListOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodGet,
		path,
		uv,
		nil,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) CreateFirewall(ctx context.Context, input firewall.Info) (firewall.CreateOutput, error) {
	c.Log(http.MethodPost, "create firewall %q", input.Name)

	var ans firewall.CreateOutput
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls"},
		V2Path: []string{"v2", "config", "ngfirewalls"},
	}
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPost,
		path,
		nil,
		input,
		&ans,
	)

	return ans, err
}

func (c *Client) CreateFirewallWithWait(ctx context.Context, input firewall.Info) (firewall.CreateOutput, error) {
	ans, err := c.CreateFirewall(ctx, input)
	if err != nil {
		return ans, err
	}
	err = c.WaitForFirewallStatus(ctx, c, ans.Response.Id, []string{FirewallStatusCreateComplete.String(), FirewallStatusCreateFail.String()})
	if err != nil {
		return ans, err
	}
	return ans, nil
}

func (c *Client) ListTagsForFirewall(ctx context.Context, input firewall.ListTagsInput) (firewall.ListTagsOutput, error) {
	c.Log(http.MethodGet, "list tags for firewall: %s", input.Firewall)

	var uv url.Values
	if input.AccountId != "" || input.NextToken != "" || input.MaxResults != 0 {
		uv = url.Values{}
		if input.AccountId != "" {
			uv.Set("accountid", input.AccountId)
		}
		if input.NextToken != "" {
			uv.Set("nexttoken", input.NextToken)
		}
		if input.MaxResults != 0 {
			uv.Set("maxresults", strconv.Itoa(input.MaxResults))
		}
	}
	uv.Set("v1route", "true")
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.Firewall, "tags"},
	}
	var ans firewall.ListTagsOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodGet,
		path,
		uv,
		nil,
		&ans,
	)

	return ans, err
}

// UpdateDescription updates the description of the firewall.
func (c *Client) UpdateFirewallDescription(ctx context.Context, input firewall.UpdateDescriptionInput) error {
	c.Log(http.MethodPut, "updating firewall description: %s", input.Firewall)
	uv := url.Values{}
	uv.Set("v1route", "true")
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.Firewall, "description"},
	}
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPut,
		path,
		uv,
		input,
		nil,
	)

	return err
}

// UpdateSubnetMappings updates the subnet mappings of the firewall.
func (c *Client) UpdateFirewallSubnetMappings(ctx context.Context, input firewall.UpdateSubnetMappingsInput) error {
	c.Log(http.MethodPut, "updating firewall subnet mappings: %s", input.Firewall)
	uv := url.Values{}
	uv.Set("v1route", "true")
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.Firewall, "subnets"},
	}
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPut,
		path,
		uv,
		input,
		nil,
	)
	return err
}

func (c *Client) RemoveTagsForFirewall(ctx context.Context, input firewall.RemoveTagsInput) error {
	c.Log(http.MethodDelete, "removing tags from firewall: %s", input.Firewall)
	uv := url.Values{}
	uv.Set("v1route", "true")
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.Firewall, "tags"},
	}
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodDelete,
		path,
		uv,
		input,
		nil,
	)

	return err
}

// AddTags adds the given tags to the firewall.
func (c *Client) AddTagsForFirewall(ctx context.Context, input firewall.AddTagsInput) error {
	c.Log(http.MethodPost, "adding tags to the firewall: %s", input.Firewall)
	uv := url.Values{}
	uv.Set("v1route", "true")
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.Firewall, "tags"},
	}
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPost,
		path,
		uv,
		input,
		nil,
	)

	return err
}

// UpdateRulestack updates the rulestack for the given firewall.
func (c *Client) UpdateFirewallRulestackV1(ctx context.Context, input firewall.UpdateRulestackInput) error {
	c.Log(http.MethodPost, "updating firewall rulestack: %s", input.Firewall)
	uv := url.Values{}
	uv.Set("v1route", "true")
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.Firewall, "rulestack"},
	}
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPost,
		path,
		uv,
		input,
		nil,
	)

	return err
}

func (c *Client) UpdateFirewallFeatures(ctx context.Context, input firewall.UpdateFeaturesAPIInput) error {
	c.Log(http.MethodPut, "updating firewall features: %+v", input.Features)
	uv := url.Values{}
	uv.Set("v1route", "true")
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.FirewallName, "features"},
	}
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPut,
		path,
		uv,
		input,
		nil,
	)

	return err
}

func (c *Client) ModifyFirewallV1(ctx context.Context, input firewall.Info) error {
	ans, err := c.ReadFirewall(ctx, firewall.ReadInput{Name: input.Name, AccountId: input.AccountId})
	if err != nil {
		return err
	}
	cur := ans.Response.Firewall
	curTags := cur.Tags

	// No idea if this is needed or not, but do it for now just to be safe.
	tin := firewall.ListTagsInput{
		Firewall:   input.Name,
		AccountId:  input.AccountId,
		MaxResults: 100,
	}
	tans, err := c.ListTagsForFirewall(ctx, tin)
	if err != nil {
		return err
	}
	curTags = tans.Response.Tags

	// Update description.
	if input.Description != cur.Description {
		v := firewall.UpdateDescriptionInput{
			Firewall:    input.Name,
			AccountId:   input.AccountId,
			Description: input.Description,
		}
		if err = c.UpdateFirewallDescription(ctx, v); err != nil {
			return err
		}
	}

	log.Printf("current subnet mappings: %+v", cur.SubnetMappings)
	log.Printf("input subnet mappings: %+v", input.SubnetMappings)
	// Update subnet mappings.
	assoc := make([]firewall.SubnetMapping, 0, len(input.SubnetMappings))
	disassoc := make([]firewall.SubnetMapping, 0, len(cur.SubnetMappings))
	for _, x := range input.SubnetMappings {
		found := false
		for _, y := range cur.SubnetMappings {
			if x.SubnetId != "" && x.SubnetId == y.SubnetId {
				found = true
			} else if x.AvailabilityZone != "" && x.AvailabilityZone == y.AvailabilityZone {
				found = true
			}
			if found {
				break
			}
		}

		if !found {
			assoc = append(assoc, firewall.SubnetMapping{
				SubnetId:         x.SubnetId,
				AvailabilityZone: x.AvailabilityZone,
			})
		}
	}
	if len(assoc) == 0 {
		assoc = nil
	}
	for _, x := range cur.SubnetMappings {
		found := false
		for _, y := range input.SubnetMappings {
			if x.SubnetId != "" && x.SubnetId == y.SubnetId {
				found = true
			} else if x.AvailabilityZone != "" && x.AvailabilityZone == y.AvailabilityZone {
				found = true
			}
			if found {
				break
			}
		}

		if !found {
			disassoc = append(assoc, firewall.SubnetMapping{
				SubnetId:         x.SubnetId,
				AvailabilityZone: x.AvailabilityZone,
			})
		}
	}
	if len(disassoc) == 0 {
		disassoc = nil
	}
	if assoc != nil || disassoc != nil {
		v := firewall.UpdateSubnetMappingsInput{
			Firewall:                   input.Name,
			AccountId:                  input.AccountId,
			AssociateSubnetMappings:    assoc,
			DisassociateSubnetMappings: disassoc,
		}
		if err = c.UpdateFirewallSubnetMappings(ctx, v); err != nil {
			return err
		}
	}

	// Update rulestack.
	if input.Rulestack != cur.Rulestack {
		v := firewall.UpdateRulestackInput{
			Firewall:  input.Name,
			AccountId: input.AccountId,
			Rulestack: input.Rulestack,
		}
		if err = c.UpdateFirewallRulestackV1(ctx, v); err != nil {
			return err
		}
	}

	// Update tags.
	addTags := make([]tag.Details, 0, len(input.Tags))
	rmTags := make([]string, 0, len(curTags))
	for _, x := range input.Tags {
		ok := false
		for _, y := range curTags {
			if x.Key == y.Key {
				if x.Value == y.Value {
					ok = true
				} else {
					rmTags = append(rmTags, x.Key)
				}
				break
			}
		}
		if !ok {
			addTags = append(addTags, x)
		}
	}
	for _, x := range curTags {
		found := false
		for _, y := range input.Tags {
			if x.Key == y.Key {
				found = true
				break
			}
		}
		if !found {
			rmTags = append(rmTags, x.Key)
		}
	}
	// Due to the 50 tag limit, removing tags must happen before adding tags.
	if len(rmTags) > 0 {
		v := firewall.RemoveTagsInput{
			Firewall:  input.Name,
			AccountId: input.AccountId,
			Tags:      rmTags,
		}
		if err = c.RemoveTagsForFirewall(ctx, v); err != nil {
			return err
		}
	}
	if len(addTags) > 0 {
		v := firewall.AddTagsInput{
			Firewall:  input.Name,
			AccountId: input.AccountId,
			Tags:      addTags,
		}
		if err = c.AddTagsForFirewall(ctx, v); err != nil {
			return err
		}
	}
	if c.endpointsUpdate(input.SecurityZones, ans.Response.Firewall.SecurityZones) || c.featureUpdate(input, ans.Response.Firewall) {
		c.Log(http.MethodPatch, "Firewall update required for endpoints or features")
		i := firewall.UpdateFeaturesAPIInput{
			FirewallName: input.Name,
			AccountId:    input.AccountId,
			UpdateToken:  &input.UpdateToken,
		}
		if input.EgressNAT != nil {
			i.Features.EgressNat = input.EgressNAT
		}
		if input.EgressNAT != nil {
			i.Features.UserId = input.UserID
		}
		if input.EgressNAT != nil {
			i.Features.PrivateAccess = input.PrivateAccess
		}
		if input.SecurityZones != nil {
			i.Features.SecurityZones = input.SecurityZones
		}
		if err := c.UpdateFirewallFeatures(ctx, i); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ModifyFirewall(ctx context.Context, input firewall.Info) (firewall.UpdateOutput, error) {
	c.Log(http.MethodPut, "updating firewall: %s", input.Id)
	//var ans firewall.CreateOutput
	path := Path{
		V2Path: []string{"v2", "config", "ngfirewalls", input.Id},
	}
	output := &firewall.UpdateOutput{}
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPatch,
		path,
		nil,
		input,
		output,
	)
	if err != nil {
		return *output, err
	}
	return *output, nil
}

// Helper function to sort the list to make order irrelevant
func sortUserIDCustomSubnetFilterList(list []firewall.UserIDCustomSubnetFilter) {
	sort.Slice(list, func(i, j int) bool {
		// Define sorting criteria
		return list[i].Name < list[j].Name
	})
}

// Deep diff function
func deepDiff(list1, list2 []firewall.UserIDCustomSubnetFilter) bool {
	// Sort lists to ensure order doesn't matter
	sortUserIDCustomSubnetFilterList(list1)
	sortUserIDCustomSubnetFilterList(list2)

	// Check if lengths are different
	if len(list1) != len(list2) {
		return true
	}

	// Compare each element
	for i := range list1 {
		if !reflect.DeepEqual(list1[i], list2[i]) {
			return true
		}
	}

	return false
}

func (c *Client) endpointsUpdate(inputEps, respEps []firewall.EndpointConfig) bool {
	inputIdMap := EpIdMap(inputEps)
	respIdMap := EpIdMap(respEps)
	for id, ep := range inputIdMap {
		if _, ok := respIdMap[id]; !ok {
			continue
		}
		respEp := respIdMap[id]
		if ep.Prefixes == nil && respEp.Prefixes == nil {
			continue
		}
		if ep.EgressNATEnabled != respEp.EgressNATEnabled {
			c.Log(http.MethodPatch, "Endpoint %s has different EgressNAT enabled status", ep.EndpointId)
			return true
		}
		if ep.Prefixes != nil && respEp.Prefixes == nil {
			c.Log(http.MethodPatch, "Endpoint %s has different prefixes", ep.EndpointId)
			return true
		}
		if ep.Prefixes == nil && respEp.Prefixes != nil {
			c.Log(http.MethodPatch, "Endpoint %s has different prefixes", ep.EndpointId)
			return true
		}
		sort.Strings(ep.Prefixes.PrivatePrefix.Cidrs)
		sort.Strings(respEp.Prefixes.PrivatePrefix.Cidrs)
		if !slices.Equal(ep.Prefixes.PrivatePrefix.Cidrs, respEp.Prefixes.PrivatePrefix.Cidrs) {
			c.Log(http.MethodPatch, "Endpoint %s has different private prefixes", ep.EndpointId)
			return true
		}
	}
	return false
}

func (c *Client) featureUpdate(input, info firewall.Info) bool {
	if input.EgressNAT == nil && info.EgressNAT == nil {
		c.Log(http.MethodPatch, "No change in EgressNAT configuration")
		return false
	}
	if input.UserID == nil && info.UserID == nil {
		return false
	}
	if input.EgressNAT != nil && info.EgressNAT == nil {
		c.Log(http.MethodPatch, "Firewall EgressNAT configuration changed from enabled to disabled")
		return true
	}
	if input.EgressNAT == nil && info.EgressNAT != nil {
		c.Log(http.MethodPatch, "Firewall EgressNAT configuration changed from disabled to enabled")
		return true
	}
	if input.EgressNAT.Enabled != info.EgressNAT.Enabled {
		c.Log(http.MethodPatch, "Firewall EgressNAT configuration changed from %t to %t", info.EgressNAT.Enabled, input.EgressNAT.Enabled)
		return true
	}
	if input.EgressNAT.Settings != nil && info.EgressNAT.Settings == nil {
		c.Log(http.MethodPatch, "Firewall EgressNAT configuration changed from no settings to %+v", input.EgressNAT.Settings)
		return true
	}
	if input.EgressNAT.Settings == nil && info.EgressNAT.Settings != nil {
		c.Log(http.MethodPatch, "Firewall EgressNAT configuration changed from %+v to no settings", info.EgressNAT.Settings)
		return true
	}
	if input.EgressNAT.Settings == nil && info.EgressNAT.Settings == nil {
		c.Log(http.MethodPatch, "No change in Firewall EgressNAT configuration")
		return false
	}
	if input.EgressNAT.Settings.IPPoolType != info.EgressNAT.Settings.IPPoolType {
		c.Log(http.MethodPatch, "Firewall EgressNAT IP pool type changed from %s to %s", info.EgressNAT.Settings.IPPoolType, input.EgressNAT.Settings.IPPoolType)
		return true
	}
	if input.UserID.Enabled != info.UserID.Enabled {
		return true
	}
	if deepDiff(input.UserID.CustomIncludeExcludeNetwork, info.UserID.CustomIncludeExcludeNetwork) {
		return true
	}
	return false
}

func (c *Client) ModifyFirewallWithWait(ctx context.Context, input firewall.Info) error {
	timeStamp := time.Now().UTC().Unix()
	genUpdate := false
	drsCommit := false
	log.Printf("===============ModifyFirewallWithWait================")

	_, err := c.retryOnTokenConflict(ctx, func() (interface{}, error) {
		return nil, c.UpdateFirewallRulestack(ctx, input)
	})
	if err != nil {
		return err
	}

	readInput := firewall.ReadInput{
		FirewallId: input.Id,
	}
	readRuleStack, _ := c.ReadFirewall(ctx, readInput)

	if input.Rulestack == "" && readRuleStack.Response.Firewall.Rulestack != "" {
		disassociateInput := firewall.DisAssociateInput{
			Firewall:    input.Name,
			AccountId:   input.AccountId,
			UpdateToken: input.UpdateToken,
			FirewallId:  input.Id,
		}
		_, err := c.retryOnTokenConflict(ctx, func() (interface{}, error) {
			return nil, c.DisassociateRuleStackWithWait(ctx, disassociateInput)
		})
		if err != nil {
			return err
		}
	}

	readOutput, _ := c.ReadFirewall(ctx, readInput)

	result, err := c.retryOnTokenConflict(ctx, func() (interface{}, error) {
		return c.ReadAndModifyFirewall(ctx, input)
	})
	if err != nil {
		return err
	}
	ans := result.(firewall.UpdateOutput)
	c.Log(http.MethodPatch, "Firewall %s updated with deployment update token %s", input.Id, ans.Response.DeploymentUpdateToken)
	if readOutput.Response.Firewall.DeploymentUpdateToken != ans.Response.DeploymentUpdateToken {
		c.Log(http.MethodPatch, "Firewall update required due to deployment update token mismatch")
		genUpdate = true
	}
	if c.endpointsUpdate(input.Endpoints, readOutput.Response.Firewall.Endpoints) || c.featureUpdate(input, readOutput.Response.Firewall) {
		c.Log(http.MethodPatch, "Firewall update required for endpoints or features")
		drsCommit = true
	}

	updateFirewall := updateFirewall{
		GeneralUpdate: genUpdate,
		DrsCommit:     drsCommit,
	}

	// TODO: Build firewall update struct based on modified properties.
	log.Printf("===============TOKEN POST CHECK================")
	if updateFirewall.GeneralUpdate {
		err := c.WaitForFirewallStatus(ctx, c, input.Id, []string{FirewallStatusUpdateComplete.String(), FirewallStatusUpdateFail.String()})
		if err != nil {
			return err
		}
	}
	if updateFirewall.DrsCommit {
		err := c.WaitForDRSCommit(ctx, c, input.Id, timeStamp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) ReadAndModifyFirewall(ctx context.Context, input firewall.Info) (firewall.UpdateOutput, error) {
	readInput := firewall.ReadInput{
		FirewallId: input.Id,
	}
	readOutput, err := c.ReadFirewall(ctx, readInput)
	if err != nil {
		return firewall.UpdateOutput{}, err
	}

	updatedInput := input
	updatedInput.UpdateToken = readOutput.Response.Firewall.UpdateToken
	updatedInput.DeploymentUpdateToken = readOutput.Response.Firewall.DeploymentUpdateToken

	ans, err := c.ModifyFirewall(ctx, updatedInput)
	if err != nil {
		return firewall.UpdateOutput{}, err
	}

	return ans, nil
}

// Read returns information on the given object.
func (c *Client) ReadFirewall(ctx context.Context, input firewall.ReadInput) (firewall.ReadOutput, error) {
	name := input.Name
	schemaVersion := getSchemaVersion(ctx)
	fwId := input.FirewallId
	uv := url.Values{}
	if schemaVersion == cloudngfwgosdk.SchemaVersionV1 {
		fwId = name
		uv = url.Values{
			"v1route":   []string{"true"},
			"accountid": []string{input.AccountId},
		}
	} else if schemaVersion == cloudngfwgosdk.SchemaVersionV2 && input.FeatureConfig {
		uv = url.Values{
			"featureconfig": []string{"true"},
		}
	}
	c.Log(http.MethodGet, "describe firewall: %s", name)
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", name},
		V2Path: []string{"v2", "config", "ngfirewalls", fwId},
	}
	var ans firewall.ReadOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodGet,
		path,
		uv,
		input,
		&ans,
	)

	return ans, err
}

// Delete the given firewall.
func (c *Client) DeleteFirewall(ctx context.Context, input firewall.DeleteInput) (firewall.DeleteOutput, error) {
	name := input.Name
	c.Log(http.MethodDelete, "delete firewall: %s", input.Name)
	schemaVersion := getSchemaVersion(ctx)
	fwId := input.FirewallId
	uv := url.Values{}
	if schemaVersion == cloudngfwgosdk.SchemaVersionV1 {
		fwId = name
		uv = url.Values{
			"v1route":   []string{"true"},
			"accountid": []string{input.AccountId},
		}
	}
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.Name},
		V2Path: []string{"v2", "config", "ngfirewalls", fwId},
	}
	var ans firewall.DeleteOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodDelete,
		path,
		uv,
		input,
		&ans,
	)
	return ans, err
}

func (c *Client) DeleteFirewallWithWait(ctx context.Context, input firewall.DeleteInput) error {
	c.Log(http.MethodDelete, "delete firewall: %s", input.Name)
	_, err := c.DeleteFirewall(ctx, input)
	if err != nil {
		return err
	}
	err = c.WaitForFirewallStatus(ctx, c, input.FirewallId, []string{FirewallStatusDeleteComplete.String(), FirewallStatusDeleteFail.String()})
	if err != nil {
		return err
	}
	return nil
}

// AssociateRulestack updates the local rulestack for the given firewall.
func (c *Client) AssociateRulestack(ctx context.Context, input firewall.AssociateInput) (firewall.AssociateOutput, error) {
	c.Log(http.MethodPost, "associating firewall rulestack: %s", input.Firewall)
	name := input.Firewall
	schemaVersion := getSchemaVersion(ctx)
	fwId := input.FirewallId
	uv := url.Values{}
	if schemaVersion == cloudngfwgosdk.SchemaVersionV1 {
		fwId = name
		uv = url.Values{
			"v1route": []string{"true"},
		}
	}
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.Firewall, "rulestack"},
		V2Path: []string{"v2", "config", "ngfirewalls", fwId, "rulestack"},
	}
	var ans firewall.AssociateOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPost,
		path,
		uv,
		input,
		&ans,
	)

	return ans, err
}

func (c *Client) AssociateRulestackWithWait(ctx context.Context, input firewall.AssociateInput) error {
	timeStamp := time.Now().Unix()
	_, err := c.AssociateRulestack(ctx, input)
	if err != nil {
		return err
	}

	// Wait for LRS commit to complete.
	err = c.WaitForLRSCommit(ctx, c, input.FirewallId, timeStamp)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateFirewallRulestack(ctx context.Context, input firewall.Info) error {
	readInput := firewall.ReadInput{
		FirewallId: input.Id,
		Name:       input.Name,
	}
	readOutput, err := c.ReadFirewall(ctx, readInput)
	if err != nil {
		return err
	}

	if input.Rulestack != readOutput.Response.Firewall.Rulestack && input.Rulestack != "" {
		associateInput := firewall.AssociateInput{
			Firewall:    input.Name,
			Rulestack:   input.Rulestack,
			AccountId:   input.AccountId,
			UpdateToken: input.UpdateToken,
			FirewallId:  input.Id,
		}
		err := c.AssociateRulestackWithWait(ctx, associateInput)
		if err != nil {
			return err
		}
	}

	return nil
}

// Disassociate local Firewall to Global rulestack
func (c *Client) DisassociateRuleStack(ctx context.Context, input firewall.DisAssociateInput) (firewall.DisAssociateOutput, error) {
	c.Log(http.MethodDelete, "disassociating firewall to local rulestack: %s", input.Firewall)
	var uv url.Values
	path := Path{
		V2Path: []string{"v2", "config", "ngfirewalls", input.FirewallId, "rulestack"},
	}
	var ans firewall.DisAssociateOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodDelete,
		path,
		uv,
		input,
		&ans,
	)

	return ans, err
}

func (c *Client) DisassociateRuleStackWithWait(ctx context.Context, input firewall.DisAssociateInput) error {
	// timeStamp := time.Now().Unix()
	_, err := c.DisassociateRuleStack(ctx, input)
	if err != nil {
		return err
	}

	// err = c.WaitForLRSCommit(ctx, c, input.FirewallId, timeStamp)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// Associate Firewall to Global rulestack
func (c *Client) AssociateGlobalRuleStack(ctx context.Context, input firewall.AssociateInput) (firewall.AssociateOutput, error) {
	c.Log(http.MethodPut, "associating firewall to global rulestack: %s", input.Firewall)
	c.Log(http.MethodPost, "associating firewall rulestack: %s", input.Firewall)
	name := input.Firewall
	schemaVersion := getSchemaVersion(ctx)
	fwId := input.FirewallId
	uv := url.Values{}
	if schemaVersion == cloudngfwgosdk.SchemaVersionV1 {
		fwId = name
		uv = url.Values{
			"v1route": []string{"true"},
		}
	}
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.Firewall, "globalrulestack"},
		V2Path: []string{"v2", "config", "ngfirewalls", fwId, "rulestack"},
	}
	var ans firewall.AssociateOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodPost,
		path,
		uv,
		input,
		&ans,
	)

	return ans, err
}

// Disassociate Firewall to Global rulestack
func (c *Client) DisAssociateGlobalRuleStack(ctx context.Context, input firewall.DisAssociateInput) (firewall.DisAssociateOutput, error) {
	c.Log(http.MethodDelete, "disassociating firewall to global rulestack: %s", input.Firewall)
	var uv url.Values
	path := Path{
		V1Path: []string{"v1", "config", "ngfirewalls", input.Firewall, "globalrulestack"},
		V2Path: []string{"v2", "config", "ngfirewalls", input.FirewallId, "rulestack"},
	}
	var ans firewall.DisAssociateOutput
	_, err := c.Communicate(
		ctx,
		PermissionFirewall,
		http.MethodDelete,
		path,
		uv,
		input,
		&ans,
	)

	return ans, err
}

func verifyCommitStatus(commitInfo *firewall.RuleStackCommitData, commitStatus string, timestamp int64, name string) (bool, error) {
	commitProcessFinished := false
	if commitStatus == CommitStateFailed.String() {
		messages := commitInfo.CommitMessages
		return commitProcessFinished, fmt.Errorf("commit failed, %v", messages)
	}
	if commitStatus != CommitStateSuccess.String() {
		log.Printf("rulestack %s commit is in progress, status: %s", name, commitStatus)
		return commitProcessFinished, nil
	}

	commitTimestamp := commitInfo.CommitTS
	if len(commitTimestamp) == 0 {
		log.Printf("rulestack commit timestamp is empty: %s", name)
		return commitProcessFinished, nil
	}
	epochTimestamp, err := ConvertToUTCEpoch(commitTimestamp)
	if err != nil {
		return commitProcessFinished, fmt.Errorf("failed to convert to epoch timestamp, err: %s, timestamp: %d", err, epochTimestamp)
	}
	log.Printf("current epoch: %d", timestamp)
	log.Printf("commit epoch: %d", epochTimestamp)
	if epochTimestamp <= timestamp {
		log.Printf("rulestack commit timestamp validation failed: %s", name)
		return commitProcessFinished, nil
	}
	return true, nil
}

func (c *Client) WaitForDRSCommit(ctx context.Context, svc *Client, fid string, timestamp int64) error {
	return WaitForOperation(ctx, func(ctx context.Context) (bool, error) {
		req := firewall.ReadInput{
			FirewallId: fid,
		}
		res, err := svc.ReadFirewall(ctx, req)
		if err != nil {
			return false, err
		}
		commitInfo := res.Response.Status.DeviceRuleStackCommitInfo
		commitStatus := res.Response.Status.DeviceRuleStackCommitStatus
		completed, err := verifyCommitStatus(commitInfo, commitStatus, timestamp, "drs")
		if err != nil {
			return false, err
		}
		if !completed {
			svc.Log("Waiting for DRS commit to be completed..: %s ", fid)
			return true, fmt.Errorf("DRS commit is not yet completed, retrying")
		}
		return false, nil
	})
}

func (c *Client) WaitForLRSCommit(ctx context.Context, svc *Client, fid string, timestamp int64) error {
	return WaitForOperation(ctx, func(ctx context.Context) (bool, error) {
		req := firewall.ReadInput{
			FirewallId: fid,
		}
		res, err := svc.ReadFirewall(ctx, req)
		if err != nil {
			return false, err
		}
		commitInfo := res.Response.Status.RuleStackCommitInfo
		commitStatus := res.Response.Status.RulestackStatus
		completed, err := verifyCommitStatus(commitInfo, commitStatus, timestamp, "lrs")
		if err != nil {
			return false, err
		}
		if !completed {
			log.Printf("Waiting for LRS commit to be completed..")
			svc.Log("Waiting for LRS commit to be completed..: %s ", fid)
			return true, fmt.Errorf("LRS commit is not yet completed, retrying...")
		}
		return false, nil
	})
}

func (c *Client) WaitForFirewallStatus(ctx context.Context, svc *Client, fid string, expStatus []string) error {
	return WaitForOperation(ctx, func(ctx context.Context) (bool, error) {
		req := firewall.ReadInput{
			FirewallId: fid,
		}
		res, err := svc.ReadFirewall(ctx, req)
		if err != nil {
			return false, err
		}
		status := res.Response.Status.FirewallStatus
		if !slices.Contains(expStatus, status) {
			svc.Log("Waiting for firewall status: %s, exp: %s, got: %s", fid, expStatus, res.Response.Status.FirewallStatus)
			return true, fmt.Errorf("firewall status did not match expected status, expected: %v, got: %s", expStatus, status)
		}
		return false, nil
	})
}

func (c *Client) retryOnTokenConflict(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	var err error

	err = WaitForOperation(ctx, func(ctx context.Context) (bool, error) {
		result, err = operation()
		if err != nil {
			if failureResponse, ok := err.(response.Failure); ok {
				if status := failureResponse.Failed(); status != nil && status.TokenConflict() {
					log.Printf("Retrying operation due to token conflict")
					return true, err
				}
			}
			log.Printf("Token conflict not found, returning original error")
			return false, err
		}
		log.Printf("Operation successful, returning result")
		return false, nil
	})
	log.Printf("retry operation count exceeded, returning final error")
	return result, err
}
