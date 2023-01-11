package logdna

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Constants for identifying channel names easily
const (
	EMAIL     = "email"
	PAGERDUTY = "pagerduty"
	SLACK     = "slack"
	WEBHOOK   = "webhook"
)

/*
  This resource needs to initialize before Terraform initializes so we can correctly populate the Provider schema.
  We can't use the init() function because Terraform initializes before that.
*/
var _ = registerTerraform(TerraformInfo{
	name:          "logdna_view",
	orgType:       OrgTypeRegular,
	terraformType: TerraformTypeResource,
	schema:        resourceView(),
})

func resourceViewCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	view := viewRequest{}

	if diags = view.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"POST",
		"/v1/config/view",
		view,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	createdView := viewResponse{}
	err = json.Unmarshal(body, &createdView)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] After %s view, the created view is %+v", req.method, createdView)

	d.SetId(createdView.ViewID)

	return resourceViewRead(ctx, d, m)
}

func resourceViewRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	viewID := d.Id()

	req := newRequestConfig(
		pc,
		"GET",
		fmt.Sprintf("/v1/config/view/%s", viewID),
		nil,
	)

	body, err := req.MakeRequest()

	log.Printf("[DEBUG] GET view raw response body %s\n", body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote view resource",
			Detail:   err.Error(),
		})
		return diags
	}

	view := viewResponse{}
	err = json.Unmarshal(body, &view)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote view resource",
			Detail:   err.Error(),
		})
		return diags
	}
	log.Printf("[DEBUG] The GET view structure is as follows: %+v\n", view)

	// Top level keys can be set directly
	appendError(d.Set("name", view.Name), &diags)
	appendError(d.Set("query", view.Query), &diags)
	appendError(d.Set("categories", view.Category), &diags)
	appendError(d.Set("hosts", view.Hosts), &diags)
	appendError(d.Set("tags", view.Tags), &diags)
	appendError(d.Set("apps", view.Apps), &diags)
	appendError(d.Set("levels", view.Levels), &diags)
	// NOTE There is always one element in the PresetIds slice
	appendError(d.Set("presetid", strings.Join(view.PresetIds, "")), &diags)

	// NOTE API does DB denormalization and extend a view record in DB
	//      with a alert channels which break a schema validation here.
	//      We don't need the channels field in case when a presetid exists
	if len(d.Get("presetid").(string)) > 0 {
		return diags
	}

	// Convert types to maps for setting the schema
	integrations, diags := view.MapChannelsToSchema()
	log.Printf("[DEBUG] view MapChannelsToSchema result: %+v\n", integrations)

	// Store the channel responses in the schema - note that this should also NUKE missing
	// integrations since we have done a PUT operation. Thus, remove non-existing things.
	for name, value := range integrations {
		schemaKey := fmt.Sprintf("%s_channel", name)
		appendError(d.Set(schemaKey, value), &diags)
	}

	return diags
}

func resourceViewUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)
	viewID := d.Id()
	view := viewRequest{}

	if diags = view.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"PUT",
		fmt.Sprintf("/v1/config/view/%s", viewID),
		view,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s %s SUCCESS. Remote resource updated.", req.method, req.apiURL)

	return resourceViewRead(ctx, d, m)
}

func resourceViewDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	viewID := d.Id()

	req := newRequestConfig(
		pc,
		"DELETE",
		fmt.Sprintf("/v1/config/view/%s", viewID),
		nil,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s view %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceView() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceViewCreate,
		ReadContext:   resourceViewRead,
		UpdateContext: resourceViewUpdate,
		DeleteContext: resourceViewDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"apps": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"categories": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					shouldSuppress := false
					lowerCaseOld := strings.ToLower(old)
					lowerCaseNew := strings.ToLower(new)
					if lowerCaseOld == lowerCaseNew {
						shouldSuppress = true
					}
					log.Println("[DEBUG] Do view category names appear the same (case-insensitive) between state and remote?", shouldSuppress)
					return shouldSuppress
				},
			},
			"hosts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"levels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"presetid": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"email_channel",
					"pagerduty_channel",
					"slack_channel",
					"webhook_channel",
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
							Default:  "presence",
						},
						"terminal": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "false",
						},
						"timezone": {
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
						"triggerinterval": {
							Type:     schema.TypeString,
							Optional: true,
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
								log.Println("[DEBUG] Does view 'bodytemplate' value in state appear the same as remote?", shouldSuppress)
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
		},
	}
}
