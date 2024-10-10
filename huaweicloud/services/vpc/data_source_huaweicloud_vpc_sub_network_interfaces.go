// Generated by PMS #235
package vpc

import (
	"context"

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

func DataSourceVpcSubNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcSubNetworkInterfacesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"interface_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the supplementary network interface.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the VPC to which the supplementary network interface belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the subnet to which the supplementary network interface belongs.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the elastic network interface`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the private IPv4 address of the supplementary network interface.`,
			},
			"mac_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the MAC address of the supplementary network interface.`,
			},
			"description": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the description of the supplementary network interface.`,
			},
			"sub_network_interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of supplementary network interfaces.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of supplementary network interface.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the VPC to which the supplementary network interface belongs.`,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the subnet to which the supplementary network interface belongs.`,
						},
						"parent_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the elastic network interface to which the supplementary network interface belongs.`,
						},
						"parent_device_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the parent device.`,
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of the security groups IDs to which the supplementary network interface belongs.`,
						},
						"vlan_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The vlan ID of the supplementary network interface.`,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The private IPv4 address of the supplementary network interface.`,
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The MAC address of the supplementary network interface.`,
						},
						"ipv6_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IPv6 address of the supplementary network interface.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the supplementary network interface.`,
						},
						"security_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the IPv6 address is it enabled of the supplementary network interface.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the project to which the supplementary network interface belongs.`,
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The tags of a supplementary network interface.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the supplementary network interface is created.`,
						},
					},
				},
			},
		},
	}
}

type SubNetworkInterfacesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newSubNetworkInterfacesDSWrapper(d *schema.ResourceData, meta interface{}) *SubNetworkInterfacesDSWrapper {
	return &SubNetworkInterfacesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceVpcSubNetworkInterfacesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newSubNetworkInterfacesDSWrapper(d, meta)
	lisSubNetIntRst, err := wrapper.ListSubNetworkInterfaces()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listSubNetworkInterfacesToSchema(lisSubNetIntRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API VPC GET /v3/{project_id}/vpc/sub-network-interfaces
func (w *SubNetworkInterfacesDSWrapper) ListSubNetworkInterfaces() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "vpc")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/vpc/sub-network-interfaces"
	params := map[string]any{
		"id":                 w.PrimToArray("interface_id"),
		"virsubnet_id":       w.PrimToArray("subnet_id"),
		"private_ip_address": w.PrimToArray("ip_address"),
		"mac_address":        w.PrimToArray("mac_address"),
		"vpc_id":             w.PrimToArray("vpc_id"),
		"description":        w.ListToArray("description"),
		"parent_id":          w.PrimToArray("parent_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("sub_network_interfaces", "page_info.next_marker", "marker").
		Request().
		Result()
}

func (w *SubNetworkInterfacesDSWrapper) listSubNetworkInterfacesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("sub_network_interfaces", schemas.SliceToList(body.Get("sub_network_interfaces"),
			func(subNetInt gjson.Result) any {
				return map[string]any{
					"id":               subNetInt.Get("id").Value(),
					"vpc_id":           subNetInt.Get("vpc_id").Value(),
					"subnet_id":        subNetInt.Get("virsubnet_id").Value(),
					"parent_id":        subNetInt.Get("parent_id").Value(),
					"parent_device_id": subNetInt.Get("parent_device_id").Value(),
					"security_groups":  schemas.SliceToStrList(subNetInt.Get("security_groups")),
					"vlan_id":          subNetInt.Get("vlan_id").Value(),
					"ip_address":       subNetInt.Get("private_ip_address").Value(),
					"mac_address":      subNetInt.Get("mac_address").Value(),
					"ipv6_ip_address":  subNetInt.Get("ipv6_ip_address").Value(),
					"description":      subNetInt.Get("description").Value(),
					"security_enabled": subNetInt.Get("security_enabled").Value(),
					"project_id":       subNetInt.Get("project_id").Value(),
					"tags":             schemas.SliceToStrList(subNetInt.Get("tags")),
					"created_at":       subNetInt.Get("created_at").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
