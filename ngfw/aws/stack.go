package aws

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/stack"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/tag"
	"net/http"
	"net/url"
	"strconv"
)

const (
	LocalScope  = "Local"
	GlobalScope = "Global"
)

// List returns a list of objects.
func (c *Client) ListRuleStack(ctx context.Context, input stack.ListInput) (stack.ListOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return stack.ListOutput{}, permErr
	}

	c.Log(http.MethodGet, "list rulestacks")

	var uv url.Values
	if input.Scope != "" ||
		input.TagKey != "" ||
		input.TagValue != "" ||
		input.Candidate ||
		input.Running ||
		input.Uncommitted ||
		input.Describe ||
		input.NextToken != "" ||
		input.MaxResults != 0 {
		uv = url.Values{}
		if input.Scope != "" {
			uv.Set("scope", input.Scope)
		}
		if input.TagKey != "" {
			uv.Set("tagkey", input.TagKey)
		}
		if input.TagValue != "" {
			uv.Set("tagvalue", input.TagValue)
		}
		if input.Candidate {
			uv.Set("candidate", "true")
		}
		if input.Running {
			uv.Set("running", "true")
		}
		if input.Describe {
			uv.Set("describe", "true")
		}
		if input.NextToken != "" {
			uv.Set("nexttoken", input.NextToken)
		}
		if input.MaxResults != 0 {
			uv.Set("maxresults", strconv.Itoa(input.MaxResults))
		}
	}

	var ans stack.ListOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks"},
		uv,
		nil,
		&ans,
	)

	return ans, err
}

// Create creates an object.
func (c *Client) CreateRuleStack(ctx context.Context, input stack.Info) error {
	perm, permErr := GetPermission(input.Entry.Scope)
	if permErr != nil {
		return permErr
	}
	c.Log(http.MethodPost, "create rulestack: %s", input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks"},
		nil,
		input,
		nil,
	)

	return err
}

// Read returns information on the given object.
func (c *Client) ReadRuleStack(ctx context.Context, input stack.ReadInput) (stack.ReadOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return stack.ReadOutput{}, permErr
	}

	c.Log(http.MethodGet, "describe rulestack: %s", input.Name)

	var uv url.Values
	if input.Candidate || input.Running {
		uv = url.Values{}
		if input.Candidate {
			uv.Set("candidate", "true")
		}
		if input.Running {
			uv.Set("running", "true")
		}
	}

	var ans stack.ReadOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Name},
		uv,
		nil,
		&ans,
	)

	return ans, err
}

// export returns the rulestack XML.
func (c *Client) ExportRuleStackXML(ctx context.Context, input stack.ReadInput) (stack.ExportRulestackXmlOutput, error) {
	scope := LocalScope
	if input.Scope != "" {
		scope = input.Scope
	}
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return stack.ExportRulestackXmlOutput{}, permErr
	}

	c.Log(http.MethodGet, "export rulestack: %s", input.Name)

	uv := url.Values{"scope": []string{scope}}
	if input.Candidate || input.Running {
		if input.Candidate {
			uv.Set("candidate", "true")
		}
		if input.Running {
			uv.Set("running", "true")
		}
	}

	var ans stack.ExportRulestackXmlOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Name, "export"},
		uv,
		nil,
		&ans,
	)

	return ans, err
}

func B64EncodeGzip(data []byte) (string, error) {
	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	if _, err := gz.Write(data); err != nil {
		return "", err
	}
	if err := gz.Flush(); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// savepanrs saves the panorama rulestack XML in S3 bucket.
func (c *Client) SaveRuleStackXML(ctx context.Context, input stack.SaveRulestackXmlInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	//gzip and b64 encode the xml
	out, err := B64EncodeGzip([]byte(input.RuleStackEntryXml.Xml))
	if err != nil {
		return err
	}
	input.RuleStackEntryXml.Xml = out
	c.Log(http.MethodPost, "save rulestack xml: %s", input.Name)

	_, err = c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Name, "xml"},
		nil,
		input,
		nil,
	)

	return err
}

func (c *Client) CreateSCMRuleStack(ctx context.Context, input stack.CreateSCMRuleStackInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	//gzip and b64 encode the xml
	out, err := B64EncodeGzip([]byte(input.RuleStackEntryXml.Xml))
	if err != nil {
		return err
	}
	input.RuleStackEntryXml.Xml = out
	c.Log(http.MethodPost, "save rulestack xml: %s", input.Name)

	_, err = c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Name, "scm"},
		nil,
		input,
		nil,
	)

	return err
}

// Update updates the given object.
func (c *Client) UpdateRuleStack(ctx context.Context, input stack.Info) error {
	perm, permErr := GetPermission(input.Entry.Scope)
	if permErr != nil {
		return permErr
	}

	name := input.Name
	input.Name = ""

	c.Log(http.MethodPut, "updating rulestack: %s", name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPut,
		[]string{"v1", "config", "rulestacks", name},
		nil,
		input,
		nil,
	)

	return err
}

// Delete removes the given object from the config.
func (c *Client) DeleteRuleStack(ctx context.Context, input stack.SimpleInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodDelete, "delete rulestack: %s", input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", input.Name},
		nil,
		nil,
		nil,
	)

	return err
}

// Commit commits the rulestack configuration.
func (c *Client) CommitRuleStack(ctx context.Context, input stack.SimpleInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodPost, "commit rulestack: %s", input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Name, "commit"},
		nil,
		nil,
		nil,
	)

	return err
}

// PollCommit does the necessary looping to wait for a commit to complete.
func (c *Client) PollCommitRuleStack(ctx context.Context, input stack.SimpleInput) (stack.CommitStatus, error) {
	c.Log(http.MethodGet, "begin commit polling: %s", input.Name)
	defer c.Log(http.MethodGet, "end commit polling: %s", input.Name)

	ans, err := c.CommitStatusRuleStack(ctx, input)
	if err != nil {
		return ans, err
	}

	switch ans.Response.CommitStatus {
	case api.RsCommitStatusPending:
		return ans, nil
	case api.RsCommitStatusSuccess:
		return ans, nil
	default:
		return ans, fmt.Errorf(ans.CommitErrors())
	}
}

// CommitStatus gets the commit status.
func (c *Client) CommitStatusRuleStack(ctx context.Context, input stack.SimpleInput) (stack.CommitStatus, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return stack.CommitStatus{}, permErr
	}

	c.Log(http.MethodGet, "commit status for rulestack: %s", input.Name)

	var ans stack.CommitStatus
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Name, "commit"},
		nil,
		nil,
		&ans,
	)

	return ans, err
}

// Revert reverts to the last good config.
func (c *Client) RevertRuleStack(ctx context.Context, input stack.SimpleInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodPost, "revert rulestack: %s", input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Name, "revert"},
		nil,
		nil,
		nil,
	)

	return err
}

// Validate validates the rulestack config.
func (c *Client) ValidateRuleStack(ctx context.Context, input stack.SimpleInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodPost, "validate rulestack: %s", input.Name)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Name, "validate"},
		nil,
		nil,
		nil,
	)

	return err
}

// ListTags returns the list of tags for this rulestack.
func (c *Client) ListTagsRuleStack(ctx context.Context, input stack.ListTagsInput) (stack.ListTagsOutput, error) {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return stack.ListTagsOutput{}, permErr
	}

	c.Log(http.MethodGet, "list rulestack tags: %s", input.Rulestack)

	var uv url.Values
	if input.NextToken != "" || input.MaxResults != 0 {
		uv = url.Values{}
		if input.NextToken != "" {
			uv.Set("nexttoken", input.NextToken)
		}
		if input.MaxResults != 0 {
			uv.Set("maxresults", strconv.Itoa(input.MaxResults))
		}
	}

	var ans stack.ListTagsOutput
	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "tags"},
		uv,
		nil,
		&ans,
	)

	return ans, err
}

// AddTags adds tags to the specified rulestack.
func (c *Client) AddTagsRuleStack(ctx context.Context, input stack.AddTagsInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodPost, "adding tags to the rulestack: %s", input.Rulestack)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "tags"},
		nil,
		input,
		nil,
	)

	return err
}

// RemoveTags removes the given tags from the resource.
func (c *Client) RemoveTagsRuleStack(ctx context.Context, input stack.RemoveTagsInput) error {
	perm, permErr := GetPermission(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.Log(http.MethodDelete, "removing tags from rulestack: %s", input.Rulestack)

	_, err := c.Communicate(
		ctx,
		perm,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "tags"},
		nil,
		input,
		nil,
	)

	return err
}

// ApplyTags applies the tags to this rulestack, performing a create and delete if needed.
func (c *Client) ApplyTagsRuleStack(ctx context.Context, input stack.AddTagsInput) error {
	lti := stack.ListTagsInput{
		Rulestack:  input.Rulestack,
		Scope:      input.Scope,
		MaxResults: 100,
	}
	lans, err := c.ListTagsRuleStack(ctx, lti)
	if err != nil {
		return err
	}

	toAdd := make([]tag.Details, 0, len(input.Tags))
	toDelete := make([]string, 0, len(lans.Response.Tags))

	// Find tags to add in.
	for _, x := range input.Tags {
		ok := false
		for _, y := range lans.Response.Tags {
			if x.Key == y.Key {
				if x.Value == y.Value {
					ok = true
				} else {
					toDelete = append(toDelete, x.Key)
				}
				break
			}
		}

		if !ok {
			toAdd = append(toAdd, x)
		}
	}

	// Find current tags that shouldn't exist.
	for _, x := range lans.Response.Tags {
		found := false
		for _, y := range input.Tags {
			if x.Key == y.Key {
				found = true
				break
			}
		}

		if !found {
			toDelete = append(toDelete, x.Key)
		}
	}

	// Delete first.
	if len(toDelete) > 0 {
		fi := stack.RemoveTagsInput{
			Rulestack: input.Rulestack,
			Scope:     input.Scope,
			Tags:      toDelete,
		}
		if err = c.RemoveTagsRuleStack(ctx, fi); err != nil {
			return err
		}
	}

	// Tag next.
	if len(toAdd) > 0 {
		fi := stack.AddTagsInput{
			Rulestack: input.Rulestack,
			Scope:     input.Scope,
			Tags:      toAdd,
		}
		if err = c.AddTagsRuleStack(ctx, fi); err != nil {
			return err
		}
	}

	// Done.
	return nil
}
