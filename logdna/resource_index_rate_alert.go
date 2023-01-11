package logdna

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const indexRateAlertConfigID = "config"

var _ = registerTerraform(TerraformInfo{
	name:          "logdna_index_rate_alert",
	orgType:       OrgTypeRegular,
	terraformType: TerraformTypeResource,
	schema:        resourceIndexRateAlert(),
})

/**
 * Create/Update index rate alert resource
 * As API does not allow the POST method, this method calls PUT to be used for both create and update.
 * which allows upsert and create a new index rate alert config record is not exist
 * Only one config per account is allowed
 */
func resourceIndexRateAlertUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	indexRateAlert := indexRateAlertRequest{}

	if diags = indexRateAlert.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"PUT",
		"/v1/config/index-rate",
		indexRateAlert,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	createdIndexRateAlert := indexRateAlertResponse{}
	err = json.Unmarshal(body, &createdIndexRateAlert)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s %s SUCCESS. Remote resource updated.", req.method, req.apiURL)

	d.SetId(indexRateAlertConfigID)

	return resourceIndexRateAlertRead(ctx, d, m)
}

func resourceIndexRateAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	req := newRequestConfig(
		pc,
		"GET",
		"/v1/config/index-rate",
		nil,
	)

	body, err := req.MakeRequest()

	log.Printf("[DEBUG] GET IndexRateAlert raw response body %s\n", body)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote IndexRateAlert resource",
			Detail:   err.Error(),
		})
		return diags
	}

	indexRateAlert := indexRateAlertResponse{}

	err = json.Unmarshal(body, &indexRateAlert)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote indexRateAlert resource",
			Detail:   err.Error(),
		})
		return diags
	}
	log.Printf("[DEBUG] The GET indexRateAlert structure is as follows: %+v\n", indexRateAlert)

	var channels []interface{}

	integrations := make(map[string]interface{})

	integrations["email"] = indexRateAlert.Channels.Email
	integrations["pagerduty"] = indexRateAlert.Channels.Pagerduty
	integrations["slack"] = indexRateAlert.Channels.Slack
	webhooks := mapIndexRateAlertWebhookToSchema(indexRateAlert)

	appendError(d.Set("webhook_channel", webhooks), &diags)

	channels = append(channels, integrations)

	appendError(d.Set("max_lines", indexRateAlert.MaxLines), &diags)
	appendError(d.Set("max_z_score", indexRateAlert.MaxZScore), &diags)
	appendError(d.Set("threshold_alert", indexRateAlert.ThresholdAlert), &diags)
	appendError(d.Set("frequency", indexRateAlert.Frequency), &diags)
	appendError(d.Set("channels", channels), &diags)
	appendError(d.Set("enabled", indexRateAlert.Enabled), &diags)

	d.SetId(indexRateAlertConfigID)

	return diags
}

/**
 * Delete index rate alert resource
 * As API does not allow DELETE method this method calls PUT
 * We considering delete action as just disabling a config
 */
func resourceIndexRateAlertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	resourceIndexRateAlertRead(ctx, d, m)

	indexRateAlert := indexRateAlertRequest{}

	if diags = indexRateAlert.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	indexRateAlert.Enabled = false

	req := newRequestConfig(
		pc,
		"PUT",
		"/v1/config/index-rate",
		indexRateAlert,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s disable IndexRateAlert %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceIndexRateAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIndexRateAlertUpdate,
		UpdateContext: resourceIndexRateAlertUpdate,
		ReadContext:   resourceIndexRateAlertRead,
		DeleteContext: resourceIndexRateAlertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"max_lines": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Max number of lines for alert",
			},
			"max_z_score": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Max Z score before alerting",
			},
			"threshold_alert": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"separate", "both"}, false),
			},
			"frequency": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"hourly", "daily"}, false),
			},
			"channels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"pagerduty": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"slack": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"webhook_channel":{
				Type: schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
				  Schema: map[string]*schema.Schema{
					"url": {
					  Type: schema.TypeString,
					  Required: true,
					},
					"method": {
					  Type: schema.TypeString,
					  Required: true,
					  ValidateFunc: validation.StringInSlice([]string{"GET", "POST","PUT","DELETE"}, false),
					},
					"headers": &schema.Schema{
					  Type: schema.TypeMap,
					  Optional:true,
					  Elem: &schema.Schema{
						Type: schema.TypeString,
					  },
					  Computed: true,
					},
					"bodytemplate": {
					  Type:     schema.TypeString,
					  Optional: true,
					  // This function compares JSON, ignoring whitespace that can occur in a .tf config.
					  // Without this, `terraform apply` will think values are different from remote to state.
					  DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						var jsonOld, jsonNew interface{}
						var err error
						err = json.Unmarshal([]byte(old), &jsonOld)
						if err != nil {
						  return false
						}
						err = json.Unmarshal([]byte(new), &jsonNew)
						if err != nil {
						  return false
						}
						shouldSuppress := reflect.DeepEqual(jsonNew, jsonOld)
						log.Println("[DEBUG] Does view 'bodytemplate' value in state appear the same as remote?", shouldSuppress)
						return shouldSuppress
					  },
					},
				  },
				},
			  },
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}
