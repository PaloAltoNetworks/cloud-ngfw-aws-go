package api

import (
	"context"
	"log"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/stack"
)

const (
	RulestackNotExists      = "Rulestack does not exist"
	PANORAMA_RULESTACK_NAME = "panorama-rulestack" // FIXME: use dg as rulestack name after integration with proxy
	LocalScope              = "Local"
	GlobalScope             = "Global"
	RsCommitStatusPending   = "Pending"
	RsCommitStatusSuccess   = "Success"
	RsCommitStatusFailed    = "Failed"
	FwStatusCommitting      = "Committing"
	FwStatusFailure         = "Failed"
	FwStatusSuccess         = "Success"
	FwStatusValidating      = "Validating"
	RuleStackTypeSCM        = "scm"
	RuleStackTypePanorama   = "panorama"
)

func (c *ApiClient) ReadRuleStack(ctx context.Context, input stack.ReadInput) (stack.ReadOutput, error) {
	out, err := c.client.ReadRuleStack(ctx, input)
	if err != nil {
		return stack.ReadOutput{}, err
	}

	log.Printf(
		"read rulestack %s",
		input.Name)

	return out, nil
}

func (c *ApiClient) ValidateRuleStack(ctx context.Context, input stack.SimpleInput) error {
	err := c.client.ValidateRuleStack(ctx, input)
	if err != nil {
		return err
	}

	log.Printf(
		"validate rulestack %s",
		input.Name)

	return nil
}

func (c *ApiClient) ExportRuleStackXML(ctx context.Context, input stack.ReadInput) (stack.ExportRulestackXmlOutput, error) {
	out, err := c.client.ExportRuleStackXML(ctx, input)
	if err != nil {
		return stack.ExportRulestackXmlOutput{}, err
	}

	log.Printf(
		"export rulestack:%s",
		input.Name)

	return out, nil
}

func (c *ApiClient) SaveRuleStackXML(ctx context.Context, input stack.SaveRulestackXmlInput) error {
	if err := c.client.SaveRuleStackXML(ctx, input); err != nil {
		return err
	}

	log.Printf(
		"export rulestack:%s",
		input.Name)

	return nil
}

func (c *ApiClient) CreateRuleStack(ctx context.Context, input stack.Info) error {
	log.Printf(
		"create rulestack %s",
		input.Name)
	if err := c.client.CreateRuleStack(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) UpdateRuleStack(ctx context.Context, input stack.Info) error {
	log.Printf(
		"create rulestack %s",
		input.Name)
	if err := c.client.UpdateRuleStack(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) DeleteRuleStack(ctx context.Context, input stack.SimpleInput) error {
	log.Printf(
		"commit rulestack %s %s",
		input.Name,
		input.Scope)
	if err := c.client.DeleteRuleStack(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) CommitRuleStack(ctx context.Context, input stack.SimpleInput) error {
	log.Printf(
		"commit rulestack %s %s",
		input.Name,
		input.Scope)
	if err := c.client.CommitRuleStack(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) CommitStatusRuleStack(ctx context.Context, input stack.SimpleInput) (stack.CommitStatus, error) {
	Logger.Debugf(
		"commit status rulestack %s %s",
		input.Name,
		input.Scope)
	status, err := c.client.CommitStatusRuleStack(ctx, input)
	if err != nil {
		return stack.CommitStatus{}, err
	}
	return status, nil
}

func (c *ApiClient) PollCommitRulestack(ctx context.Context, input stack.SimpleInput) (stack.CommitStatus, error) {
	Logger.Debugf(
		"commit rulestack %s %s",
		input.Name,
		input.Scope)
	status, err := c.client.PollCommitRuleStack(ctx, input)
	if err != nil {
		return stack.CommitStatus{}, err
	}
	return status, nil
}

func (c *ApiClient) ListTagsRuleStack(ctx context.Context, input stack.ListTagsInput) (stack.ListTagsOutput, error) {
	status, err := c.client.ListTagsRuleStack(ctx, input)
	if err != nil {
		return stack.ListTagsOutput{}, err
	}
	return status, nil
}

func (c *ApiClient) AddTagsRuleStack(ctx context.Context, input stack.AddTagsInput) error {
	if err := c.client.AddTagsRuleStack(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) RemoveTagsRuleStack(ctx context.Context, input stack.RemoveTagsInput) error {
	if err := c.client.RemoveTagsRuleStack(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) ApplyTagsRuleStack(ctx context.Context, input stack.AddTagsInput) error {
	if err := c.client.ApplyTagsRuleStack(ctx, input); err != nil {
		return err
	}
	return nil
}
