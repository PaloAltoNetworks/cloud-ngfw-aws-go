package api

import (
	"context"
	"os"
	"strings"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/firewall"
	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/response"
)

const (
	FW_AMI_VERSION_10_2_7 = "FW_AMI_VERSION_10_2_7"
	FW_AMI_VERSION_11_2_7 = "FW_AMI_VERSION_11_2_7"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) ListFirewall(ctx context.Context, a firewall.ListInput) (firewall.ListOutput, error) {
	out, err := c.client.ListFirewall(ctx, a)
	if err != nil {
		return firewall.ListOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) CreateFirewall(ctx context.Context, input firewall.Info) (firewall.CreateOutput, error) {
	out, err := c.client.CreateFirewall(ctx, input)
	if err != nil {
		return firewall.CreateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) CreateFirewallWithWait(ctx context.Context, input firewall.Info) (firewall.CreateOutput, error) {
	out, err := c.client.CreateFirewallWithWait(ctx, input)
	if err != nil {
		return firewall.CreateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ModifyFirewallV1(ctx context.Context, input firewall.Info) error {
	err := c.client.ModifyFirewallV1(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) ModifyFirewall(ctx context.Context, input firewall.Info) (firewall.UpdateOutput, error) {
	out, err := c.client.ModifyFirewall(ctx, input)
	if err != nil {
		return firewall.UpdateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ModifyFirewallWithWait(ctx context.Context, input firewall.Info) error {
	return c.client.ModifyFirewallWithWait(ctx, input)
}

func (c *ApiClient) ReadFirewall(ctx context.Context, input firewall.ReadInput) (firewall.ReadOutput, error) {
	if c.Mock {
		featureConfigs := make(map[string]firewall.FeatureConfig)
		l := os.Getenv("MOCK_FEATURES")
		s := strings.Split(l, ",")
		for _, v := range s {
			// Populate FeatureConfigs map based on feature type
			switch strings.ToUpper(v) {
			case "USERID":
				featureConfigs["USERID"] = &firewall.UserId{}
			case "LDAP":
				featureConfigs["LDAP"] = &firewall.Ldap{}
			case "DLP":
				// DLP doesn't have a specific FeatureConfig type, but add a placeholder
				featureConfigs["DLP"] = &firewall.Dlp{} // Use a dummy type for now
			}
		}
		fwSwVersion := os.Getenv("MOCK_FW_SW_VERSION")
		return firewall.ReadOutput{
			Response: firewall.ReadResponse{
				Firewall: firewall.Info{
					SoftwareVersion: fwSwVersion,
					FeatureConfigs:  featureConfigs,
				},
			},
			Status: response.Status{},
		}, nil
	}

	out, err := c.client.ReadFirewall(ctx, input)
	if err != nil {
		return firewall.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) AssociateRulestack(ctx context.Context, input firewall.AssociateInput) (firewall.AssociateOutput, error) {
	out, err := c.client.AssociateRulestack(ctx, input)
	if err != nil {
		return firewall.AssociateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) AssociateGlobalRulestack(ctx context.Context, input firewall.AssociateInput) (firewall.AssociateOutput, error) {
	out, err := c.client.AssociateGlobalRuleStack(ctx, input)
	if err != nil {
		return firewall.AssociateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) AssociateRulestackWithWait(ctx context.Context, input firewall.AssociateInput) error {
	return c.client.AssociateRulestackWithWait(ctx, input)
}

func (c *ApiClient) DeleteFirewall(ctx context.Context, input firewall.DeleteInput) error {
	if _, err := c.client.DeleteFirewall(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) DeleteFirewallWithWait(ctx context.Context, input firewall.DeleteInput) error {
	if err := c.client.DeleteFirewallWithWait(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) AssociateGlobalRuleStack(ctx context.Context, input firewall.AssociateInput) (firewall.AssociateOutput, error) {
	if c.Mock {
		return firewall.AssociateOutput{}, nil
	}
	out, err := c.client.AssociateGlobalRuleStack(ctx, input)
	if err != nil {
		return firewall.AssociateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) DisAssociateGlobalRuleStack(ctx context.Context, input firewall.DisAssociateInput) (firewall.DisAssociateOutput, error) {
	out, err := c.client.DisAssociateGlobalRuleStack(ctx, input)
	if err != nil {
		return firewall.DisAssociateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) DisassociateRuleStack(ctx context.Context, input firewall.DisAssociateInput) (firewall.DisAssociateOutput, error) {
	out, err := c.client.DisassociateRuleStack(ctx, input)
	if err != nil {
		return firewall.DisAssociateOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) DisassociateRuleStackWithWait(ctx context.Context, input firewall.DisAssociateInput) error {
	return c.client.DisassociateRuleStackWithWait(ctx, input)
}
