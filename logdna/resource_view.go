package logdna

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	EMAIL = "email"
	PAGERDUTY = "pagerduty"
	WEBHOOK = "webhook"
)

func buildChannels(emailChannels []interface{}, pagerDutyChannels []interface{}, webhookChannels []interface{}) ([]ChannelRequest, error) {
	var channels []ChannelRequest
	for _, emailChannel := range emailChannels {
		i := emailChannel.(map[string]interface{})

		emails := i["emails"].([]interface{})
		immediate := i["immediate"].(string)
		operator := i["operator"].(string)
		terminal := i["terminal"].(string)
		triggerInterval := i["triggerinterval"].(string)
		triggerLimit := i["triggerlimit"].(int)
		timezone := i["timezone"].(string)

		var emailStrings []string

		for _, email := range emails {
			emailStrings = append(emailStrings, email.(string))
		}

		email := ChannelRequest{
			Emails:          emailStrings,
			Immediate:       immediate,
			Integration:     EMAIL,
			Operator:        operator,
			Terminal:        terminal,
			TriggerInterval: triggerInterval,
			TriggerLimit:    triggerLimit,
			Timezone:        timezone,
		}

		channels = append(channels, email)
	}

	for _, pagerDutyChannel := range pagerDutyChannels {
		i := pagerDutyChannel.(map[string]interface{})

		immediate := i["immediate"].(string)
		key := i["key"].(string)
		operator := i["operator"].(string)
		terminal := i["terminal"].(string)
		triggerInterval := i["triggerinterval"].(string)
		triggerLimit := i["triggerlimit"].(int)

		pagerDuty := ChannelRequest{
			Immediate:       immediate,
			Integration:     "pagerduty",
			Key:             key,
			Operator:        operator,
			Terminal:        terminal,
			TriggerInterval: triggerInterval,
			TriggerLimit:    triggerLimit,
		}

		channels = append(channels, pagerDuty)
	}

	for _, webhookChannel := range webhookChannels {
		i := webhookChannel.(map[string]interface{})

		bodytemplate := i["bodytemplate"].(string)
		headers := i["headers"].(map[string]interface{})
		immediate := i["immediate"].(string)
		method := i["method"].(string)
		operator := i["operator"].(string)
		terminal := i["terminal"].(string)
		triggerInterval := i["triggerinterval"].(string)
		triggerLimit := i["triggerlimit"].(int)
		url := i["url"].(string)

		headersMap := make(map[string]string)

		for k, v := range headers {
			headersMap[k] = v.(string)
		}

		webhook := ChannelRequest{
			Headers:         headersMap,
			Immediate:       immediate,
			Integration:     WEBHOOK,
			Operator:        operator,
			Method:          method,
			TriggerInterval: triggerInterval,
			TriggerLimit:    triggerLimit,
			URL:             url,
			Terminal:        terminal,
		}

		if bodytemplate != "" {
			var bt map[string]interface{}
			err := json.Unmarshal([]byte(bodytemplate), &bt)

			if err != nil {
				return channels, err
			}
			webhook.BodyTemplate = bt
		}
		channels = append(channels, webhook)
	}
	return channels, nil
}

func resourceViewCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	config := m.(*config)

	view := ViewRequest{}

	if diags = view.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	client := Client{
		ServiceKey: config.ServiceKey,
		HTTPClient: config.HTTPClient,
		ApiUrl: fmt.Sprintf("%s/v1/config/view", config.URL),
		Method: "POST",
		Body: view,
	}

	body, err := client.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", client.Method, client.ApiUrl, body)

	if err != nil {
		return diag.FromErr(err)
	}

	createdView := ViewResponse{}
	err = json.Unmarshal(body, &createdView)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] After %s view, the created view is %+v", client.Method, createdView)

	d.SetId(createdView.ViewID)
	return diags
}

func resourceViewRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*config)
	viewID := d.Id()

	c := Client{
		ServiceKey: config.ServiceKey,
		HTTPClient: config.HTTPClient,
		ApiUrl: fmt.Sprintf("%s/v1/config/view/%s", config.URL, viewID),
		Method: "GET",
	}

	body, err := c.MakeRequest()

	log.Printf("[DEBUG] GET view raw response body %s\n", body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary: "Cannot read the remote view resource",
			Detail: err.Error(),
		})
		return diags
	}

	view := ViewResponse{}
	err = json.Unmarshal(body, &view)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary: "Cannot unmarshal response from the remote view resource",
			Detail: err.Error(),
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

	// Convert types to maps for setting the schema
	integrations, diags := view.MapChannelsToSchema()
	log.Printf("[DEBUG] MapChannelsToSchema result: %+v\n", integrations)
	if emailChannels := integrations[EMAIL]; emailChannels != nil {
		appendError(d.Set("email_channel", emailChannels), &diags)
	}
	if pagerdutyChannels := integrations[PAGERDUTY]; pagerdutyChannels != nil {
		appendError(d.Set("pagerduty_channel", pagerdutyChannels), &diags)
	}
	if webhookChannels := integrations[WEBHOOK]; webhookChannels != nil {
		appendError(d.Set("webhook_channel", webhookChannels), &diags)
	}

	return diags
}

func appendError(err error, diags *diag.Diagnostics) *diag.Diagnostics {
	if err != nil {
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary: "There was a problem setting the view schema",
			Detail: err.Error(),
		})
	}
	return diags
}

func resourceViewUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  var diags diag.Diagnostics
	config := m.(*config)
  viewId := d.Id()
	view := ViewRequest{}

	if diags = view.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	client := Client{
		ServiceKey: config.ServiceKey,
		HTTPClient: config.HTTPClient,
		ApiUrl: fmt.Sprintf("%s/v1/config/view/%s", config.URL, viewId),
		Method: "PUT",
		Body: view,
	}

	body, err := client.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", client.Method, client.ApiUrl, body)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s %s SUCCESS. Remote resource updated.", client.Method, client.ApiUrl)

	return diags
}

func resourceViewDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	viewID := d.Id()

	client := Client{
		ServiceKey: config.ServiceKey,
		HTTPClient: config.HTTPClient,
		ApiUrl: fmt.Sprintf("%s/v1/config/view/%s", config.URL, viewID),
		Method: "DELETE",
	}

	body, err := client.MakeRequest()
	log.Printf("[DEBUG] %s %s view %s", client.Method, client.ApiUrl, body)

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
							Default: "false",
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"terminal": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "false",
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
							Default: "false",
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"terminal": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "false",
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
			"webhook_channel": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bodytemplate": {
							Type:     schema.TypeString,
							Optional: true,
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
							Default: "false",
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"terminal": {
							Type:     schema.TypeString,
							Optional: true,
							Default: "false",
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
