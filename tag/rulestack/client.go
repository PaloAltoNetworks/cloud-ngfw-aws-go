package rulestack

import (
	"context"
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

	var ans ListOutput
	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodGet,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "tags"},
		uv,
		nil,
		&ans,
	)

	return ans, err
}

// Tag applies the given tags to the resource.
func (c *Client) Tag(ctx context.Context, input Info) error {
	c.client.Log(http.MethodPost, "adding tags to rulestack: %s", input.Rulestack)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodPost,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "tags"},
		nil,
		input,
		nil,
	)

	return err
}

// Untag removes the given tags from the resource.
func (c *Client) Untag(ctx context.Context, input UntagInput) error {
	c.client.Log(http.MethodDelete, "removing tags from rulestack: %s", input.Rulestack)

	_, err := c.client.Communicate(
		ctx,
		permissions.Rulestack,
		http.MethodDelete,
		[]string{"v1", "config", "rulestacks", input.Rulestack, "tags"},
		nil,
		input,
		nil,
	)

	return err
}

// Apply applies the given tags, performing a create and delete if necessary.
func (c *Client) Apply(ctx context.Context, input Info) error {
	lans, err := c.List(ctx, ListInput{Rulestack: input.Rulestack, MaxResults: 1000})
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
		if err = c.Untag(ctx, UntagInput{Rulestack: input.Rulestack, Tags: toDelete}); err != nil {
			return err
		}
	}

	// Tag next.
	if len(toAdd) > 0 {
		if err = c.Tag(ctx, Info{Rulestack: input.Rulestack, Tags: toAdd}); err != nil {
			return err
		}
	}

	return nil
}
