package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceRdsPgPlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsPgPluginCreate,
		ReadContext:   resourceRdsPgPluginRead,
		DeleteContext: resourceRdsPgPluginDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared_preload_libraries": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildCreatePgPluginBody(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"database_name":  utils.ValueIngoreEmpty(d.Get("database_name")),
		"extension_name": utils.ValueIngoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func queryPluginDetail(client *golangsdk.ServiceClient, instanceId, databaseName, name string) (interface{}, error) {
	reqUri := "v3/{project_id}/instances/{instance_id}/extensions?database_name={database_name}"
	reqPath := client.Endpoint + reqUri
	reqPath = strings.ReplaceAll(reqPath, "{project_id}", client.ProjectID)
	reqPath = strings.ReplaceAll(reqPath, "{instance_id}", instanceId)
	reqPath = strings.ReplaceAll(reqPath, "{database_name}", databaseName)

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		reqPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		errCode := parseErrCode(resp)
		if errCode == "DBS.280238" || errCode == "DBS.200823" {
			fmt.Printf("[WARN] error retrieving PG plugin, error code: %s", errCode)
			return nil, golangsdk.ErrDefault404{}
		}
		return nil, fmt.Errorf("error retrieving PG plugin: %s", err)
	}
	bodyBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error marshaling PG plugin: %s", err)
	}

	var bodyJson interface{}
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal PG plugin: %s", err)
	}

	log.Printf("[DEBUG] reading RDS PG plugin response: %#v", bodyJson)

	pluginDetail := utils.PathSearch(fmt.Sprintf("extensions[?name=='%s']|[?created]|[0]", name), bodyJson, nil)
	if pluginDetail == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return pluginDetail, nil
}

func parseErrCode(resp interface{}) string {
	bodyBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[ERROR] error marshaling PG plugin: %s", err)
		return ""
	}

	var bodyJson interface{}
	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		log.Printf("[ERROR] error unmarshal PG plugin: %s", err)
		return ""
	}
	errCode := utils.PathSearch(fmt.Sprintf("errCode"), bodyJson, "")
	return errCode.(string)
}

func resourceRdsPgPluginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS PG plugin client: %s", err)
	}

	// 判断插件是否已创建
	pluginDetail, err := queryPluginDetail(client, d.Get("instance_id").(string), d.Get("database_name").(string), d.Get("name").(string))
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return diag.FromErr(err)
		}
	}

	if pluginDetail != nil {
		return diag.Errorf("error RDS PG plugin already created: %#v", pluginDetail)
	}

	reqUri := "v3/{project_id}/instances/{instance_id}/extensions"
	reqPath := client.Endpoint + reqUri
	reqPath = strings.Replace(reqPath, "{project_id}", client.ProjectID, -1)
	reqPath = strings.Replace(reqPath, "{instance_id}", d.Get("instance_id").(string), -1)

	jsonBody := utils.RemoveNil(buildCreatePgPluginBody(d))
	log.Printf("[DEBUG] create RDS PG plugin: %#v", jsonBody)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         jsonBody,
		OkCodes: []int{
			200,
		},
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err = client.Request("POST", reqPath, &createOpt)
		retryable, err := handleMultiOperationsError(err)
		if retryable {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error creating RDS PG plugin: %s", err)
	}

	id := fmt.Sprintf("%s/%s/%s", d.Get("instance_id").(string), d.Get("database_name").(string), d.Get("name").(string))
	d.SetId(id)

	return resourceRdsPgPluginRead(ctx, d, cfg)
}

func resourceRdsPgPluginRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS PG plugin client: %s", err)
	}

	arr := strings.Split(d.Id(), "/")
	if len(arr) != 3 {
		return diag.Errorf("invalid ID format in ReadContext: %s", d.Id())
	}

	pluginDetail, err := queryPluginDetail(client, arr[0], arr[1], arr[2])
	log.Printf("[DEBUG] RDS PG plugin detail: %#v", pluginDetail)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving PG plugin")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", utils.PathSearch("name", pluginDetail, "")),
		d.Set("database_name", utils.PathSearch("database_name", pluginDetail, "")),
		d.Set("version", utils.PathSearch("version", pluginDetail, "")),
		d.Set("shared_preload_libraries", utils.PathSearch("shared_preload_libraries", pluginDetail, "")),
		d.Set("created", utils.PathSearch("created", pluginDetail, false)),
		d.Set("description", utils.PathSearch("description", pluginDetail, "")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDeletePgPluginBody(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"database_name":  utils.ValueIngoreEmpty(d.Get("database_name")),
		"extension_name": utils.ValueIngoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func resourceRdsPgPluginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error deleting RDS client: %s", err)
	}

	reqUri := "v3/{project_id}/instances/{instance_id}/extensions"
	reqPath := client.Endpoint + reqUri
	reqPath = strings.Replace(reqPath, "{project_id}", client.ProjectID, -1)
	reqPath = strings.Replace(reqPath, "{instance_id}", d.Get("instance_id").(string), -1)

	jsonBody := utils.RemoveNil(buildDeletePgPluginBody(d))
	log.Printf("[DEBUG] delete RDS PG plugin: %#v", jsonBody)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         jsonBody,
		OkCodes: []int{
			200,
		},
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = client.Request("DELETE", reqPath, &deleteOpt)
		retryable, err := handleMultiOperationsError(err)
		if retryable {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error deleting PG plugin: %s", err)
	}

	return nil
}
