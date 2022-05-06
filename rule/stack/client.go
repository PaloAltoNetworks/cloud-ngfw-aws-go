package stack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/permissions"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/tag"
)

// Client is a client for this collection.
type Client struct {
	client api.Client
}

// NewClient returns a new client for this collection.
func NewClient(client api.Client) *Client {
	return &Client{client: client}
}

// List returns a list of objects.
func (c *Client) List(ctx context.Context, input ListInput) (ListOutput, error) {
	perm, permErr := selectPermissions(input.Scope)
	if permErr != nil {
		return ListOutput{}, permErr
	}

	c.client.Log(http.MethodGet, "list rulestacks")

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

	var ans ListOutput
	_, err := c.client.Communicate(
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
func (c *Client) Create(ctx context.Context, input Info) error {
	perm, permErr := selectPermissions(input.Entry.Scope)
	if permErr != nil {
		return permErr
	}

	c.client.Log(http.MethodPost, "create rulestack: %s", input.Name)

	_, err := c.client.Communicate(
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
func (c *Client) Read(ctx context.Context, input ReadInput) (ReadOutput, error) {
	perm, permErr := selectPermissions(input.Scope)
	if permErr != nil {
		return ReadOutput{}, permErr
	}

	c.client.Log(http.MethodGet, "describe rulestack: %s", input.Name)

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

	var ans ReadOutput
	_, err := c.client.Communicate(
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

// Update updates the given object.
func (c *Client) Update(ctx context.Context, input Info) error {
	perm, permErr := selectPermissions(input.Entry.Scope)
	if permErr != nil {
		return permErr
	}

	name := input.Name
	input.Name = ""

	c.client.Log(http.MethodPut, "updating rulestack: %s", name)

	_, err := c.client.Communicate(
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
func (c *Client) Delete(ctx context.Context, input SimpleInput) error {
	perm, permErr := selectPermissions(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.client.Log(http.MethodDelete, "delete rulestack: %s", input.Name)

	_, err := c.client.Communicate(
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
func (c *Client) Commit(ctx context.Context, input SimpleInput) error {
	perm, permErr := selectPermissions(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.client.Log(http.MethodPost, "commit rulestack: %s", input.Name)

	_, err := c.client.Communicate(
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

// CommitStatus gets the commit status.
func (c *Client) CommitStatus(ctx context.Context, input SimpleInput) (CommitStatus, error) {
	perm, permErr := selectPermissions(input.Scope)
	if permErr != nil {
		return CommitStatus{}, permErr
	}

	c.client.Log(http.MethodGet, "commit status for rulestack: %s", input.Name)

	var ans CommitStatus
	_, err := c.client.Communicate(
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
func (c *Client) Revert(ctx context.Context, input SimpleInput) error {
	perm, permErr := selectPermissions(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.client.Log(http.MethodPost, "revert rulestack: %s", input.Name)

	_, err := c.client.Communicate(
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
func (c *Client) Validate(ctx context.Context, input SimpleInput) error {
	perm, permErr := selectPermissions(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.client.Log(http.MethodPost, "validate rulestack: %s", input.Name)

	_, err := c.client.Communicate(
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
func (c *Client) ListTags(ctx context.Context, input ListTagsInput) (ListTagsOutput, error) {
	perm, permErr := selectPermissions(input.Scope)
	if permErr != nil {
		return ListTagsOutput{}, permErr
	}

	c.client.Log(http.MethodGet, "list rulestack tags: %s", input.Rulestack)

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

	var ans ListTagsOutput
	_, err := c.client.Communicate(
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
func (c *Client) AddTags(ctx context.Context, input AddTagsInput) error {
	perm, permErr := selectPermissions(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.client.Log(http.MethodPost, "adding tags to the rulestack: %s", input.Rulestack)

	_, err := c.client.Communicate(
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
func (c *Client) RemoveTags(ctx context.Context, input RemoveTagsInput) error {
	perm, permErr := selectPermissions(input.Scope)
	if permErr != nil {
		return permErr
	}

	c.client.Log(http.MethodDelete, "removing tags from rulestack: %s", input.Rulestack)

	_, err := c.client.Communicate(
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
func (c *Client) ApplyTags(ctx context.Context, input AddTagsInput) error {
	lti := ListTagsInput{
		Rulestack:  input.Rulestack,
		Scope:      input.Scope,
		MaxResults: 100,
	}
	lans, err := c.ListTags(ctx, lti)
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
		fi := RemoveTagsInput{
			Rulestack: input.Rulestack,
			Scope:     input.Scope,
			Tags:      toDelete,
		}
		if err = c.RemoveTags(ctx, fi); err != nil {
			return err
		}
	}

	// Tag next.
	if len(toAdd) > 0 {
		fi := AddTagsInput{
			Rulestack: input.Rulestack,
			Scope:     input.Scope,
			Tags:      toAdd,
		}
		if err = c.AddTags(ctx, fi); err != nil {
			return err
		}
	}

	// Done.
	return nil
}

func selectPermissions(v string) (string, error) {
	switch v {
	case "", "Local":
		return permissions.Rulestack, nil
	case "Global":
		return permissions.GlobalRulestack, nil
	}

	return "", fmt.Errorf("Unknown permission: %s", v)
}
