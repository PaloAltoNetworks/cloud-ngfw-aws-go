package api

import (
	"context"
	context2 "context"
	"go.uber.org/zap"
	"log"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/appid"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/certificate"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/country"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/feed"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/firewall"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/fqdn"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/logprofile"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/predefinedurl"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/prefix"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/security"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/stack"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/api/url"
)

// vendor specific ngfw clients(AWS, Azure) implement apiClient under ngfw directory
type Client interface {
	ListFeed(ctx context.Context, input feed.ListInput) (feed.ListOutput, error)
	CreateFeed(ctx context.Context, input feed.Info) error
	ReadFeed(ctx context.Context, input feed.ReadInput) (feed.ReadOutput, error)
	UpdateFeed(ctx context.Context, input feed.Info) error
	DeleteFeed(ctx context.Context, input feed.DeleteInput) error

	ListSecurityRule(ctx context.Context, input security.ListInput) (security.ListOutput, error)
	CreateSecurityRule(ctx context.Context, input security.Info) error
	ReadSecurityRule(ctx context.Context, input security.ReadInput) (security.ReadOutput, error)
	UpdateSecurityRule(ctx context.Context, input security.Info) error
	DeleteSecurityRule(ctx context.Context, input security.DeleteInput) error

	ListRuleStack(ctx context.Context, input stack.ListInput) (stack.ListOutput, error)
	CreateRuleStack(ctx context.Context, input stack.Info) error
	ReadRuleStack(ctx context.Context, input stack.ReadInput) (stack.ReadOutput, error)
	ExportRuleStackXML(ctx context.Context, input stack.ReadInput) (stack.ExportRulestackXmlOutput, error)
	SaveRuleStackXML(ctx context.Context, input stack.SaveRulestackXmlInput) error
	UpdateRuleStack(ctx context.Context, input stack.Info) error
	CreateSCMRuleStack(ctx context.Context, input stack.CreateSCMRuleStackInput) error
	DeleteRuleStack(ctx context.Context, input stack.SimpleInput) error
	CommitRuleStack(ctx context.Context, input stack.SimpleInput) error
	PollCommitRuleStack(ctx context.Context, input stack.SimpleInput) (stack.CommitStatus, error)
	CommitStatusRuleStack(ctx context.Context, input stack.SimpleInput) (stack.CommitStatus, error)
	RevertRuleStack(ctx context.Context, input stack.SimpleInput) error
	ValidateRuleStack(ctx context.Context, input stack.SimpleInput) error
	ListTagsRuleStack(ctx context.Context, input stack.ListTagsInput) (stack.ListTagsOutput, error)
	AddTagsRuleStack(ctx context.Context, input stack.AddTagsInput) error
	RemoveTagsRuleStack(ctx context.Context, input stack.RemoveTagsInput) error
	ApplyTagsRuleStack(ctx context.Context, input stack.AddTagsInput) error

	ListAppID(ctx context.Context, input appid.ListInput) (appid.ListOutput, error)
	ReadAppID(ctx context.Context, input appid.ReadInput) (appid.ReadOutput, error)
	ReadApplication(ctx context.Context, version, app string) (appid.ReadApplicationOutput, error)

	ListCertificate(ctx context.Context, input certificate.ListInput) (certificate.ListOutput, error)
	CreateCertificate(ctx context.Context, input certificate.Info) error
	ReadCertificate(ctx context.Context, input certificate.ReadInput) (certificate.ReadOutput, error)
	UpdateCertificate(ctx context.Context, input certificate.Info) error
	DeleteCertificate(ctx context.Context, input certificate.DeleteInput) error

	ListCountry(ctx context.Context, input country.ListInput) (country.ListOutput, error)

	ListFqdn(ctx context.Context, input fqdn.ListInput) (fqdn.ListOutput, error)
	CreateFqdn(ctx context.Context, input fqdn.Info) error
	ReadFqdn(ctx context.Context, input fqdn.ReadInput) (fqdn.ReadOutput, error)
	UpdateFqdn(ctx context.Context, input fqdn.Info) error
	DeleteFqdn(ctx context.Context, input fqdn.DeleteInput) error

	ReadFirewallLogprofile(ctx context.Context, input logprofile.ReadInput) (logprofile.ReadOutput, error)
	UpdateFirewallLogprofile(ctx context.Context, input logprofile.Info) error

	ListUrlPredefinedCategories(ctx context.Context, input predefinedurl.ListInput) (predefinedurl.ListOutput, error)
	ListUrlCategoriesActionOverride(ctx context.Context, input predefinedurl.ListOverridesInput) (predefinedurl.ListOverridesOutput, error)
	DescribeUrlCategoryActionOverride(ctx context.Context, input predefinedurl.GetOverrideInput) (predefinedurl.GetOverrideOutput, error)
	UpdateUrlCategoryActionOverride(ctx context.Context, input predefinedurl.OverrideInput) error

	ListPrefixList(ctx context.Context, input prefix.ListInput) (prefix.ListOutput, error)
	CreatePrefixList(ctx context.Context, input prefix.Info) error
	ReadPrefixList(ctx context.Context, input prefix.ReadInput) (prefix.ReadOutput, error)
	UpdatePrefixList(ctx context.Context, input prefix.Info) error
	DeletePrefixList(ctx context.Context, input prefix.DeleteInput) error

	ListUrlCustomCategory(ctx context.Context, input url.ListInput) (url.ListOutput, error)
	CreateUrlCustomCategory(ctx context.Context, input url.Info) error
	ReadUrlCustomCategory(ctx context.Context, input url.ReadInput) (url.ReadOutput, error)
	UpdateUrlCustomCategory(ctx context.Context, input url.Info) error
	DeleteUrlCustomCategory(ctx context.Context, input url.DeleteInput) error

	ListFirewall(ctx context.Context, input firewall.ListInput) (firewall.ListOutput, error)
	CreateFirewall(ctx context.Context, input firewall.Info) (firewall.CreateOutput, error)
	ModifyFirewall(ctx context.Context, input firewall.Info) error
	ReadFirewall(ctx context.Context, input firewall.ReadInput) (firewall.ReadOutput, error)
	UpdateFirewallDescription(ctx context.Context, input firewall.UpdateDescriptionInput) error
	UpdateFirewallContentVersion(ctx context.Context, input firewall.UpdateContentVersionInput) error
	UpdateFirewallSubnetMappings(ctx context.Context, input firewall.UpdateSubnetMappingsInput) error
	UpdateFirewallRulestack(ctx context.Context, input firewall.UpdateRulestackInput) error
	ListTagsForFirewall(ctx context.Context, input firewall.ListTagsInput) (firewall.ListTagsOutput, error)
	RemoveTagsForFirewall(ctx context.Context, input firewall.RemoveTagsInput) error
	AddTagsForFirewall(ctx context.Context, input firewall.AddTagsInput) error
	DeleteFirewall(ctx context.Context, input firewall.DeleteInput) error
	AssociateGlobalRuleStack(ctx context.Context, input firewall.AssociateInput) (firewall.AssociateOutput, error)
	DisAssociateGlobalRuleStack(ctx context.Context, input firewall.DisAssociateInput) (firewall.DisAssociateOutput, error)
	SetEndpoint(ctx context.Context, input EndPointInput) error
	GetCloudNGFWServiceToken(ctx context.Context, info stack.AuthInput) (stack.AuthOutput, error)
	IsSyncModeEnabled(ctx context.Context) bool
	GetResourceTimeout(ctx context.Context) int
}

type ApiClient struct {
	client    Client
	ctx       context.Context
	maxGortns int
	XSLPath   string
}

type EndPointInput struct {
	ApiEndpoint     string
	ApiAuthEndpoint string
}

var Logger *zap.SugaredLogger

func SetLogger(logger *zap.SugaredLogger) {
	Logger = logger
}

func (c *ApiClient) SetEndpoint(ctx context2.Context, input EndPointInput) error {
	return c.client.SetEndpoint(ctx, input)
}

func (c *ApiClient) IsSyncModeEnabled(ctx context2.Context) bool {
	return c.client.IsSyncModeEnabled(ctx)
}

func (c *ApiClient) GetResourceTimeout(ctx context2.Context) int {
	return c.client.GetResourceTimeout(ctx)
}

// sdk consumers instantiate APIClient using NewAPIClient() and invoke APIs under api directory
func NewAPIClient(client Client, ctx context.Context, maxGortns int, XSLPath string, mock bool) *ApiClient {
	if !mock && Logger == nil {
		log.Fatalf("Initialize logger using SetLogger()")
	}
	return &ApiClient{client: client, ctx: ctx, maxGortns: maxGortns, XSLPath: XSLPath}
}
