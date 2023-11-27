package logdna

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	alert := alertRequest{}

	if diags = alert.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"POST",
		"/v1/config/presetalert",
		alert,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	createdAlert := alertResponse{}
	err = json.Unmarshal(body, &createdAlert)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] After %s presetalert, the created alert is %+v", req.method, createdAlert)

	d.SetId(createdAlert.PresetID)

	return resourceAlertRead(ctx, d, m)
}

func resourceAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	presetID := d.Id()

	req := newRequestConfig(
		pc,
		"GET",
		fmt.Sprintf("/v1/config/presetalert/%s", presetID),
		nil,
	)

	body, err := req.MakeRequest()

	log.Printf("[DEBUG] GET presetalert raw response body %s\n", body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote presetalert resource",
			Detail:   err.Error(),
		})
		return diags
	}

	alert := alertResponse{}
	err = json.Unmarshal(body, &alert)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote presetalert resource",
			Detail:   err.Error(),
		})
		return diags
	}
	log.Printf("[DEBUG] The GET presetalert structure is as follows: %+v\n", alert)

	// Top level keys can be set directly
	appendError(d.Set("name", alert.Name), &diags)

	// Convert types to maps for setting the schema
	integrations, diags := alert.MapChannelsToSchema()
	log.Printf("[DEBUG] presetalert MapChannelsToSchema result: %+v\n", integrations)

	// Store the responses in the schema - note that this should also NUKE missing
	// integrations since we have done a PUT operation. Thus, remove non-existing things.
	for name, value := range integrations {
		schemaKey := fmt.Sprintf("%s_channel", name)
		appendError(d.Set(schemaKey, value), &diags)
	}

	return diags
}

func resourceAlertUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)
	presetID := d.Id()
	alert := alertRequest{}

	if diags = alert.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"PUT",
		fmt.Sprintf("/v1/config/presetalert/%s", presetID),
		alert,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s %s SUCCESS. Remote resource updated.", req.method, req.apiURL)

	return resourceAlertRead(ctx, d, m)
}

func resourceAlertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	presetID := d.Id()

	req := newRequestConfig(
		pc,
		"DELETE",
		fmt.Sprintf("/v1/config/presetalert/%s", presetID),
		nil,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s presetalert %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertCreate,
		ReadContext:   resourceAlertRead,
		UpdateContext: resourceAlertUpdate,
		DeleteContext: resourceAlertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email_channel": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"emails": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"immediate": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"terminal": {
							Type:     schema.TypeString,
							Required: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"triggerinterval": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"triggerlimit": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 100000 {
									errs = append(errs, fmt.Errorf("%q must be between 1 and 100,000 inclusive, got: %d", key, v))
								}
								return
							},
						},
					},
				},
			},
			"pagerduty_channel": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"immediate": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "presence",
						},
						"terminal": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"triggerinterval": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"triggerlimit": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 100000 {
									errs = append(errs, fmt.Errorf("%q must be between 1 and 100,000 inclusive, got: %d", key, v))
								}
								return
							},
						},
						"autoresolve": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"autoresolveinterval": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"autoresolvelimit": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 100000 {
									errs = append(errs, fmt.Errorf("%q must be between 1 and 100,000 inclusive, got: %d", key, v))
								}
								return
							},
						},
					},
				},
			},
			"slack_channel": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"immediate": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "presence",
						},
						"terminal": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"triggerinterval": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"triggerlimit": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 100000 {
									errs = append(errs, fmt.Errorf("%q must be between 1 and 100,000 inclusive, got: %d", key, v))
								}
								return
							},
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"webhook_channel": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
								log.Println("[DEBUG] Does presetalert 'bodytemplate' value in state appear the same as remote?", shouldSuppress)
								return shouldSuppress
							},
						},
						"headers": {
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"immediate": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "presence",
						},
						"terminal": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"triggerinterval": {
							Type:     schema.TypeString,
							Required: true,
						},
						"triggerlimit": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(int)
								if v < 1 || v > 100000 {
									errs = append(errs, fmt.Errorf("%q must be between 1 and 100,000 inclusive, got: %d", key, v))
								}
								return
							},
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}
