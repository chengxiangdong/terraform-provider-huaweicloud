// Generated by PMS #276
package ccm

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceCertificateExport() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCertificateExportRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the certificate ID.`,
			},
			"enc_private_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The encryption certificate private key. This attribute is only meaningful in the state secret certificate.`,
			},
			"entire_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The certificate content and certificate chain.`,
			},
			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The certificate content.`,
			},
			"certificate_chain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The certificate chain.`,
			},
			"private_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private key of the certificate.`,
			},
			"enc_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The encryption certificate content. This attribute is only meaningful in the state secret certificate.`,
			},
		},
	}
}

type CertificateExportDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCertificateExportDSWrapper(d *schema.ResourceData, meta interface{}) *CertificateExportDSWrapper {
	return &CertificateExportDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCertificateExportRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCertificateExportDSWrapper(d, meta)
	exportCertificateRst, err := wrapper.ExportCertificate()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.exportCertificateToSchema(exportCertificateRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API SCM POST /v3/scm/certificates/{certificate_id}/export
func (w *CertificateExportDSWrapper) ExportCertificate() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "scm")
	if err != nil {
		return nil, err
	}

	uri := "/v3/scm/certificates/{certificate_id}/export"
	uri = strings.ReplaceAll(uri, "{certificate_id}", w.Get("certificate_id").(string))
	return httphelper.New(client).
		Method("POST").
		URI(uri).
		Request().
		Result()
}

func (w *CertificateExportDSWrapper) exportCertificateToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("enc_private_key", body.Get("enc_private_key").Value()),
		d.Set("entire_certificate", body.Get("entire_certificate").Value()),
		d.Set("certificate", body.Get("certificate").Value()),
		d.Set("certificate_chain", body.Get("certificate_chain").Value()),
		d.Set("private_key", body.Get("private_key").Value()),
		d.Set("enc_certificate", body.Get("enc_certificate").Value()),
	)
	return mErr.ErrorOrNil()
}
