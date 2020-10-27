package logdna

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type clientViewPayload struct {
	Payload viewPayload
	URL     string
}

type channel struct {
	BodyTemplate    map[string]string `json:"bodyTemplate,omitempty"`
	Emails          []string          `json:"emails,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	Immediate       string            `json:"immediate,omitempty"`
	Integration     string            `json:"integration,omitempty"`
	Key             string            `json:"key,omitempty"`
	Method          string            `json:"method,omitempty"`
	Operator        string            `json:"operator,omitempty"`
	Terminal        string            `json:"terminal,omitempty"`
	Triggerinterval string            `json:"triggerinterval,omitempty"`
	Triggerlimit    int               `json:"triggerlimit,omitempty"`
	Timezone        string            `json:"timezone,omitempty"`
	URL             string            `json:"url,omitempty"`
}

func buildChannels(emailchannels []interface{}, pagerdutychannels []interface{}, webhookchannels []interface{}) []channel {
	var myChannels []channel
	for _, currChannel := range emailchannels {
		i := currChannel.(map[string]interface{})

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

		channel := channel{
			Emails:          emailStrings,
			Immediate:       immediate,
			Integration:     "email",
			Operator:        operator,
			Terminal:        terminal,
			Triggerinterval: triggerInterval,
			Triggerlimit:    triggerLimit,
			Timezone:        timezone,
		}

		myChannels = append(myChannels, channel)
	}

	for _, currChannel := range pagerdutychannels {
		i := currChannel.(map[string]interface{})

		immediate := i["immediate"].(string)
		key := i["key"].(string)
		operator := i["operator"].(string)
		terminal := i["terminal"].(string)
		triggerInterval := i["triggerinterval"].(string)
		triggerLimit := i["triggerlimit"].(int)

		channel := channel{
			Immediate:       immediate,
			Integration:     "pagerduty",
			Key:             key,
			Operator:        operator,
			Terminal:        terminal,
			Triggerinterval: triggerInterval,
			Triggerlimit:    triggerLimit,
		}

		myChannels = append(myChannels, channel)
	}

	for _, currChannel := range webhookchannels {
		i := currChannel.(map[string]interface{})

		bodytemplate := i["bodytemplate"].(map[string]interface{})
		headers := i["headers"].(map[string]interface{})
		immediate := i["immediate"].(string)
		method := i["method"].(string)
		operator := i["operator"].(string)
		terminal := i["terminal"].(string)
		triggerInterval := i["triggerinterval"].(string)
		triggerLimit := i["triggerlimit"].(int)
		url := i["url"].(string)

		headersMap := make(map[string]string)
		templateMap := make(map[string]string)

		for k, v := range headers {
			headersMap[k] = v.(string)
		}

		for k, v := range bodytemplate {
			templateMap[k] = v.(string)
		}

		channel := channel{
			BodyTemplate:    templateMap,
			Headers:         headersMap,
			Immediate:       immediate,
			Integration:     "webhook",
			Operator:        operator,
			Method:          method,
			Triggerinterval: triggerInterval,
			Triggerlimit:    triggerLimit,
			URL:             url,
			Terminal:        terminal,
		}

		myChannels = append(myChannels, channel)
	}
	return myChannels
}

func resourceViewCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	client := Client{servicekey: config.servicekey, httpClient: config.httpClient}
	name := d.Get("name").(string)
	query := d.Get("query").(string)
	categories := d.Get("categories").([]interface{})
	hosts := d.Get("hosts").([]interface{})
	tags := d.Get("tags").([]interface{})
	apps := d.Get("apps").([]interface{})
	levels := d.Get("levels").([]interface{})
	emailchannels := d.Get("email_channel").([]interface{})
	pagerdutychannels := d.Get("pagerduty_channel").([]interface{})
	webhookchannels := d.Get("webhook_channel").([]interface{})

	channels := buildChannels(emailchannels, pagerdutychannels, webhookchannels)
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

	payload := viewPayload{Name: name, Query: query, Apps: appStrings, Levels: levelStrings, Hosts: hostStrings, Category: categoryStrings, Tags: tagStrings, Channels: channels}
	resp, err := client.CreateView(config.url, payload)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp)
	return nil
}

func resourceViewRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceViewUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	client := Client{servicekey: config.servicekey, httpClient: config.httpClient}
	viewid := d.Id()
	name := d.Get("name").(string)
	query := d.Get("query").(string)
	categories := d.Get("categories").([]interface{})
	hosts := d.Get("hosts").([]interface{})
	tags := d.Get("tags").([]interface{})
	apps := d.Get("apps").([]interface{})
	levels := d.Get("levels").([]interface{})
	emailchannels := d.Get("email_channel").([]interface{})
	pagerdutychannels := d.Get("pagerduty_channel").([]interface{})
	webhookchannels := d.Get("webhook_channel").([]interface{})

	channels := buildChannels(emailchannels, pagerdutychannels, webhookchannels)

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

	payload := viewPayload{Name: name, Query: query, Apps: appStrings, Levels: levelStrings, Hosts: hostStrings, Category: categoryStrings, Tags: tagStrings, Channels: channels}
	err := client.UpdateView(config.url, viewid, payload)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceViewDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	client := Client{servicekey: config.servicekey, httpClient: config.httpClient}
	viewid := d.Id()
	err := client.DeleteView(config.url, viewid)
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
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
