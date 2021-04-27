package logdna

import (
	"context"
	"encoding/json"
	"errors"
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
	config := m.(*config)
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

	client := Client{
		ServiceKey: config.ServiceKey,
		HTTPClient: config.HTTPClient,
		ApiUrl: fmt.Sprintf("%s/v1/config/view", config.URL),
		Method: "POST",
		Payload: ViewPayload{
			Name: name,
			Query: query,
			Apps: appStrings,
			Levels: levelStrings,
			Hosts: hostStrings,
			Category: categoryStrings,
			Tags: tagStrings,
			Channels: channels,
		},
	}

	body, err := client.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", client.Method, client.ApiUrl, body)

	if err != nil {
		return diag.FromErr(err)
	}

	createdView := ViewResponsePayload{}
	err = json.Unmarshal(body, &createdView)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] After POST view, the created view is %+v", createdView)

	d.SetId(createdView.ViewID)
	return nil
}

func resourceViewRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	viewID := d.Id()

	c := Client{
		ServiceKey: config.ServiceKey,
		HTTPClient: config.HTTPClient,
		ApiUrl: fmt.Sprintf("%s/v1/config/view/%s", config.URL, viewID),
		Method: "GET",
	}

	body, err := c.MakeRequest()

	log.Printf("[DEBUG] %s %s response body %s", c.Method, c.ApiUrl, body)
	if err != nil {
		return diag.FromErr(err)
	}

	view := ViewResponsePayload{}
	err = json.Unmarshal(body, &view)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] After %s, the view structure is %+v", c.Method, view)

	// Top level keys can be set directly
	if err = d.Set("name", view.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("query", view.Query); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("categories", view.Category); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("hosts", view.Hosts); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("tags", view.Tags); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("apps", view.Apps); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("levels", view.Levels); err != nil {
		return diag.FromErr(err)
	}

	// Convert types to maps for setting the schema
	integrations, diags := view.MapChannelsToSchema()
	log.Printf("[DEBUG] Parsed integrations are %+v", integrations)
	if emailChannels := integrations[EMAIL]; emailChannels != nil {
		if err = d.Set("email_channel", emailChannels); err != nil {
			return diag.FromErr(err)
		}
	}
	if pagerdutyChannels := integrations[PAGERDUTY]; pagerdutyChannels != nil {
		if err = d.Set("pagerduty_channel", pagerdutyChannels); err != nil {
			return diag.FromErr(err)
		}
	}
	if webhookChannels := integrations[WEBHOOK]; webhookChannels != nil {
		if err = d.Set("webhook_channel", webhookChannels); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceViewUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
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

	client := Client{
		ServiceKey: config.ServiceKey,
		HTTPClient: config.HTTPClient,
		ApiUrl: fmt.Sprintf("%s/v1/config/view/%s", config.URL, viewID),
		Method: "PUT",
		Payload: ViewPayload{
			Name: name,
			Query: query,
			Apps: appStrings,
			Levels: levelStrings,
			Hosts: hostStrings,
			Category: categoryStrings,
			Tags: tagStrings,
			Channels: channels,
		},
	}

	body, err := client.MakeRequest()
	log.Printf("[DEBUG] %s %s result: %s", client.Method, client.ApiUrl, body)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
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
