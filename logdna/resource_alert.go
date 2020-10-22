package logdna

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenChannelsData(channels *[]getChannelResponse, channelType string) []interface{} {
	if channels != nil {
		cs := make([]interface{}, len(*channels), len(*channels))
		actualLength := 0

		for _, channel := range *channels {
			c := make(map[string]interface{})

			if channelType == "pagerduty" && channel.Integration == "pagerduty" {
				c["immediate"] = strconv.FormatBool(channel.Immediate)
				c["key"] = channel.Key
				c["operator"] = channel.Operator
				c["terminal"] = strconv.FormatBool(channel.Terminal)
				c["triggerinterval"] = channel.TriggerInterval
				c["triggerlimit"] = channel.TriggerLimit
				cs[actualLength] = c
				actualLength++
			}

			if channelType == "webhook" && channel.Integration == "webhook" {
				c["bodytemplate"] = channel.BodyTemplate
				c["headers"] = channel.Headers
				c["immediate"] = strconv.FormatBool(channel.Immediate)
				c["method"] = channel.Method
				c["operator"] = channel.Operator
				c["terminal"] = strconv.FormatBool(channel.Terminal)
				c["triggerinterval"] = channel.TriggerInterval
				c["triggerlimit"] = channel.TriggerLimit
				c["url"] = channel.URL
				cs[actualLength] = c
				actualLength++
			}

			if channelType == "email" && channel.Integration == "email" {
				c["immediate"] = strconv.FormatBool(channel.Immediate)
				c["emails"] = channel.Emails
				c["operator"] = channel.Operator
				c["terminal"] = strconv.FormatBool(channel.Terminal)
				c["triggerinterval"] = channel.TriggerInterval
				c["triggerlimit"] = channel.TriggerLimit
				c["timezone"] = channel.Timezone
				cs[actualLength] = c
				actualLength++
			}
		}

		return cs[:actualLength]
	}

	return nil
}

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
	return resourceAlertRead(ctx, d, m)
}

func resourceAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	presetID := d.Id()
	config := m.(*config)
	client := Client{ServiceKey: config.ServiceKey, HTTPClient: config.HTTPClient}
	resp, err := client.GetAlert(config.URL, presetID)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", resp.Name); err != nil {
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
