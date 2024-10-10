// Generated by PMS #106
package cc

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCcConnectionRoutes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcConnectionRoutesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"cloud_connection_route_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies cloud connection route ID.`,
			},
			"cloud_connection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies cloud connection ID.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies network instance ID of cloud connection route.`,
			},
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies region ID of cloud connection route.`,
			},
			"cloud_connection_routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of cloud connection routes.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cloud connection route ID.`,
						},
						"cloud_connection_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cloud connection ID.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The network instance ID of cloud connection route.`,
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The region ID of cloud connection route.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The project ID of cloud connection route.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the network instance that the next hop of a route points to.`,
						},
						"destination": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The destination address.`,
						},
					},
				},
			},
		},
	}
}

type ConnectionRoutesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newConnectionRoutesDSWrapper(d *schema.ResourceData, meta interface{}) *ConnectionRoutesDSWrapper {
	return &ConnectionRoutesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCcConnectionRoutesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newConnectionRoutesDSWrapper(d, meta)
	lisCloConRouRst, err := wrapper.ListCloudConnectionRoutes()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listCloudConnectionRoutesToSchema(lisCloConRouRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CC GET /v3/{domain_id}/ccaas/cloud-connection-routes
func (w *ConnectionRoutesDSWrapper) ListCloudConnectionRoutes() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cc")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{domain_id}/ccaas/cloud-connection-routes"
	uri = strings.ReplaceAll(uri, "{domain_id}", w.Config.DomainID)
	params := map[string]any{
		"cloud_connection_id": w.Get("cloud_connection_id"),
		"instance_id":         w.Get("instance_id"),
		"region_id":           w.Get("region_id"),
		"id":                  w.Get("cloud_connection_route_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("cloud_connection_routes", "page_info.next_marker", "marker").
		Request().
		Result()
}

func (w *ConnectionRoutesDSWrapper) listCloudConnectionRoutesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("cloud_connection_routes", schemas.SliceToList(body.Get("cloud_connection_routes"),
			func(cloConRou gjson.Result) any {
				return map[string]any{
					"id":                  cloConRou.Get("id").Value(),
					"cloud_connection_id": cloConRou.Get("cloud_connection_id").Value(),
					"instance_id":         cloConRou.Get("instance_id").Value(),
					"region_id":           cloConRou.Get("region_id").Value(),
					"project_id":          cloConRou.Get("project_id").Value(),
					"type":                cloConRou.Get("type").Value(),
					"destination":         cloConRou.Get("destination").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
