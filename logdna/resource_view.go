package logdna

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Channel contains optional and required fields for creating an alert with LogDNA
type Channel struct {
	BodyTemplate    map[string]interface{} `json:"bodyTemplate,omitempty"`
	Emails          []string               `json:"emails,omitempty"`
	Headers         map[string]string      `json:"headers,omitempty"`
	Immediate       string                 `json:"immediate,omitempty"`
	Integration     string                 `json:"integration,omitempty"`
	Key             string                 `json:"key,omitempty"`
	Method          string                 `json:"method,omitempty"`
	Operator        string                 `json:"operator,omitempty"`
	Terminal        string                 `json:"terminal,omitempty"`
	TriggerInterval string                 `json:"triggerinterval,omitempty"`
	TriggerLimit    int                    `json:"triggerlimit,omitempty"`
	Timezone        string                 `json:"timezone,omitempty"`
	URL             string                 `json:"url,omitempty"`
}

func buildChannels(emailChannels []interface{}, pagerDutyChannels []interface{}, webhookChannels []interface{}) ([]Channel, error) {
	var channels []Channel
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

		email := Channel{
			Emails:          emailStrings,
			Immediate:       immediate,
			Integration:     "email",
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

		pagerDuty := Channel{
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

		webhook := Channel{
			Headers:         headersMap,
			Immediate:       immediate,
			Integration:     "webhook",
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
	config := m.(*config)
	client := Client{ServiceKey: config.ServiceKey, HTTPClient: config.HTTPClient}
	name := d.Get("name").(string)
	query := d.Get("query").(string)
	categories := d.Get("categories").([]interface{})
	hosts := d.Get("hosts").([]interface{})
	tags := d.Get("tags").([]interface{})
	apps := d.Get("apps").([]interface{})
	levels := d.Get("levels").([]interface{})
	emailChannels := d.Get("email_channel").([]interface{})
	pagerDutyChannels := d.Get("pagerduty_channel").([]interface{})
	webhookChannels := d.Get("webhook_channel").([]interface{})

	channels, err := buildChannels(emailChannels, pagerDutyChannels, webhookChannels)
	if err != nil {
		return diag.FromErr(errors.New("bodytemplate json configuration is invalid"))
	}
	var categoryStrings []string
	var hostStrings []string
	var tagStrings []string
	var appStrings []string
	var levelStrings []string

	for _, app := range apps {
		appStrings = append(appStrings, app.(string))
	}
	for _, category := range categories {
		categoryStrings = append(categoryStrings, category.(string))
	}
	for _, host := range hosts {
		hostStrings = append(hostStrings, host.(string))
	}
	for _, level := range levels {
		levelStrings = append(levelStrings, level.(string))
	}
	for _, tag := range tags {
		tagStrings = append(tagStrings, tag.(string))
	}

	payload := ViewPayload{Name: name, Query: query, Apps: appStrings, Levels: levelStrings, Hosts: hostStrings, Category: categoryStrings, Tags: tagStrings, Channels: channels}
	resp, err := client.CreateView(config.URL, payload)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp)
	return resourceViewRead(ctx, d, m)
}

func resourceViewRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	viewid := d.Id()
	client := Client{ServiceKey: config.ServiceKey, HTTPClient: config.HTTPClient}
	resp, err := client.GetView(config.URL, viewid)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("query", resp.Query); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("apps", resp.Apps); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("categories", resp.Category); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("hosts", resp.Hosts); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("levels", resp.Levels); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tags", resp.Tags); err != nil {
		return diag.FromErr(err)
	}

	emailChannels := flattenChannelsData(&resp.Channels, "email")
	if emailChannels != nil {
		if err = d.Set("email_channel", emailChannels); err != nil {
			return diag.FromErr(err)
		}
	}

	webhookChannels := flattenChannelsData(&resp.Channels, "webhook")
	if webhookChannels != nil {
		if err = d.Set("webhook_channel", webhookChannels); err != nil {
			return diag.FromErr(err)
		}
	}

	pagerDutyChannels := flattenChannelsData(&resp.Channels, "pagerduty")
	if pagerDutyChannels != nil {
		if err = d.Set("pagerduty_channel", pagerDutyChannels); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceViewUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	client := Client{ServiceKey: config.ServiceKey, HTTPClient: config.HTTPClient}
	viewID := d.Id()
	name := d.Get("name").(string)
	query := d.Get("query").(string)
	categories := d.Get("categories").([]interface{})
	hosts := d.Get("hosts").([]interface{})
	tags := d.Get("tags").([]interface{})
	apps := d.Get("apps").([]interface{})
	levels := d.Get("levels").([]interface{})
	emailChannels := d.Get("email_channel").([]interface{})
	pagerDutyChannels := d.Get("pagerduty_channel").([]interface{})
	webhookChannels := d.Get("webhook_channel").([]interface{})

	channels, err := buildChannels(emailChannels, pagerDutyChannels, webhookChannels)
	if err != nil {
		return diag.FromErr(errors.New("bodytemplate json configuration is invalid"))
	}
	var categoryStrings []string
	var hostStrings []string
	var tagStrings []string
	var appStrings []string
	var levelStrings []string

	for _, app := range apps {
		appStrings = append(appStrings, app.(string))
	}
	for _, category := range categories {
		categoryStrings = append(categoryStrings, category.(string))
	}
	for _, host := range hosts {
		hostStrings = append(hostStrings, host.(string))
	}
	for _, level := range levels {
		levelStrings = append(levelStrings, level.(string))
	}
	for _, tag := range tags {
		tagStrings = append(tagStrings, tag.(string))
	}

	payload := ViewPayload{Name: name, Query: query, Apps: appStrings, Levels: levelStrings, Hosts: hostStrings, Category: categoryStrings, Tags: tagStrings, Channels: channels}
	err = client.UpdateView(config.URL, viewID, payload)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceViewDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	client := Client{ServiceKey: config.ServiceKey, HTTPClient: config.HTTPClient}
	viewID := d.Id()
	err := client.DeleteView(config.URL, viewID)
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
			State: schema.ImportStatePassthrough,
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
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"terminal": {
							Type:     schema.TypeString,
							Optional: true,
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
