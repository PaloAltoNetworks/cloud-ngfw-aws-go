package api

import (
	"context"

	"github.com/paloaltonetworks/cloud-ngfw-aws-go/v2/api/certificate"
)

/* Cloud vendor agnostic interface APIs to program NGFW
 */
func (c *ApiClient) CreateCertificate(ctx context.Context, input certificate.Info) error {
	if err := c.client.CreateCertificate(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) ListCertificate(ctx context.Context, a certificate.ListInput) (certificate.ListOutput, error) {
	out, err := c.client.ListCertificate(ctx, a)
	if err != nil {
		return certificate.ListOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) ReadCertificate(ctx context.Context, input certificate.ReadInput) (certificate.ReadOutput, error) {
	out, err := c.client.ReadCertificate(ctx, input)
	if err != nil {
		return certificate.ReadOutput{}, err
	}
	return out, nil
}

func (c *ApiClient) UpdateCertificate(ctx context.Context, input certificate.Info) error {
	if err := c.client.UpdateCertificate(ctx, input); err != nil {
		return err
	}
	return nil
}

func (c *ApiClient) DeleteCertificate(ctx context.Context, cert certificate.DeleteInput) error {
	input := certificate.DeleteInput{
		Rulestack: cert.Rulestack,
		Scope:     cert.Scope,
		Name:      cert.Name,
	}
	if err := c.client.DeleteCertificate(ctx, input); err != nil {
		return err
	}
	return nil
}
