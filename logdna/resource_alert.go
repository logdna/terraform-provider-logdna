package logdna

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlertCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	client := Client{ServiceKey: config.ServiceKey, HTTPClient: config.HTTPClient}
	name := d.Get("name").(string)
	emailchannels := d.Get("email_channel").([]interface{})
	pagerdutychannels := d.Get("pagerduty_channel").([]interface{})
	webhookchannels := d.Get("webhook_channel").([]interface{})

	channels, err := buildChannels(emailchannels, pagerdutychannels, webhookchannels)
	if err != nil {
		return diag.FromErr(errors.New("bodytemplate json configuration is invalid"))
	}
	payload := ViewPayload{Name: name, Channels: channels}
	resp, err := client.CreateAlert(config.URL, payload)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp)
	return nil
}

func resourceAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceAlertUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	client := Client{ServiceKey: config.ServiceKey, HTTPClient: config.HTTPClient}
	presetID := d.Id()
	name := d.Get("name").(string)
	emailchannels := d.Get("email_channel").([]interface{})
	pagerdutychannels := d.Get("pagerduty_channel").([]interface{})
	webhookchannels := d.Get("webhook_channel").([]interface{})

	channels, err := buildChannels(emailchannels, pagerdutychannels, webhookchannels)
	if err != nil {
		return diag.FromErr(errors.New("bodytemplate json configuration is invalid"))
	}
	payload := ViewPayload{Name: name, Channels: channels}
	err = client.UpdateAlert(config.URL, presetID, payload)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAlertDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*config)
	client := Client{ServiceKey: config.ServiceKey, HTTPClient: config.HTTPClient}
	presetID := d.Id()
	err := client.DeleteAlert(config.URL, presetID)
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
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"terminal": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
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
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
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
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
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
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
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
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
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
