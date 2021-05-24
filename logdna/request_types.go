package logdna

// This separation of concerns between request and response bodies is only due
// to inconsistencies in the API data types returned by the PUT versus the ones
// returned by the GET. In a perfect world, they would use the same types.

import (
	"fmt"
  "encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type viewRequest struct {
	Apps     []string         `json:"apps,omitempty"`
	Category []string         `json:"category,omitempty"`
	Channels []channelRequest `json:"channels,omitempty"`
	Hosts    []string         `json:"hosts,omitempty"`
	Levels   []string         `json:"levels,omitempty"`
	Name     string           `json:"name,omitempty"`
	Query    string           `json:"query,omitempty"`
	Tags     []string         `json:"tags,omitempty"`
}

type channelRequest struct {
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

func (view *viewRequest) CreateRequestBody(d *schema.ResourceData) diag.Diagnostics {
	// This function pulls from the schema in preparation to JSON marshal
	var diags diag.Diagnostics

	// Scalars
	view.Name = d.Get("name").(string)
	view.Query = d.Get("query").(string)

	// Simple arrays
	for _, app := range d.Get("apps").([]interface{}) {
		view.Apps = append(view.Apps, app.(string))
	}
	for _, category := range d.Get("categories").([]interface{}) {
		view.Category = append(view.Category, category.(string))
	}
	for _, host := range d.Get("hosts").([]interface{}) {
		view.Hosts = append(view.Hosts, host.(string))
	}
	for _, level := range d.Get("levels").([]interface{}) {
		view.Levels = append(view.Levels, level.(string))
	}
	for _, tag := range d.Get("tags").([]interface{}) {
		view.Tags = append(view.Tags, tag.(string))
	}

	// Complex array interfaces

	view.Channels = append(
		view.Channels,
		*mapChannelsFromSchema(
			d.Get("email_channel").([]interface{}),
			EMAIL,
			&diags,
		)...,
	)

	view.Channels = append(
		view.Channels,
		*mapChannelsFromSchema(
			d.Get("pagerduty_channel").([]interface{}),
			PAGERDUTY,
			&diags,
		)...,
	)

	view.Channels = append(
		view.Channels,
		*mapChannelsFromSchema(
			d.Get("webhook_channel").([]interface{}),
			WEBHOOK,
			&diags,
		)...,
	)
	return diags
}

func mapChannelsFromSchema(listEntries []interface{}, integration string, diags *diag.Diagnostics) *[]channelRequest {
	var prepared interface{}
	channelRequests := []channelRequest{}

	if listEntries == nil {
		return nil
	}
	for _, entry := range listEntries {
		e := entry.(map[string]interface{})
		prepared = nil
		switch integration {
		case EMAIL:
			prepared = emailChannelRequest(e)
		case PAGERDUTY:
			prepared = pagerDutyChannelRequest(e)
		case WEBHOOK:
			prepared = webHookChannelRequest(e, diags)
		default:
			*diags = append(*diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Cannot format integration channel for outbound request",
				Detail:   fmt.Sprintf("Unrecognized integration: %s", integration),
			})
		}
		if prepared == nil {
			continue
		}
		channelRequests = append(channelRequests, prepared.(channelRequest))
	}
	return &channelRequests
}

func emailChannelRequest(s map[string]interface{}) channelRequest {
	var emails []string
	for _, email := range s["emails"].([]interface{}) {
		emails = append(emails, email.(string))
	}

	c := channelRequest{
		Emails:          emails,
		Immediate:       s["immediate"].(string),
		Integration:     EMAIL,
		Operator:        s["operator"].(string),
		Terminal:        s["terminal"].(string),
		TriggerInterval: s["triggerinterval"].(string),
		TriggerLimit:    s["triggerlimit"].(int),
		Timezone:        s["timezone"].(string),
	}

	return c
}

func pagerDutyChannelRequest(s map[string]interface{}) channelRequest {
	c := channelRequest{
		Immediate:       s["immediate"].(string),
		Integration:     PAGERDUTY,
		Key:             s["key"].(string),
		Operator:        s["operator"].(string),
		Terminal:        s["terminal"].(string),
		TriggerInterval: s["triggerinterval"].(string),
		TriggerLimit:    s["triggerlimit"].(int),
	}

	return c
}

func webHookChannelRequest(s map[string]interface{}, diags *diag.Diagnostics) channelRequest {
	headersMap := make(map[string]string)

	for k, v := range s["headers"].(map[string]interface{}) {
		headersMap[k] = v.(string)
	}

	c := channelRequest{
		Headers:         headersMap,
		Immediate:       s["immediate"].(string),
		Integration:     WEBHOOK,
		Operator:        s["operator"].(string),
		Method:          s["method"].(string),
		TriggerInterval: s["triggerinterval"].(string),
		TriggerLimit:    s["triggerlimit"].(int),
		URL:             s["url"].(string),
		Terminal:        s["terminal"].(string),
	}

  if bodyTemplate := s["bodytemplate"].(string) ;bodyTemplate != "" {
    var bt map[string]interface{}
    // See if the JSON is valid, but don't use the value or it will double encode
    err := json.Unmarshal([]byte(bodyTemplate), &bt)

    if err == nil {
      c.BodyTemplate = bt
    } else {
      *diags = append(*diags, diag.Diagnostic{
        Severity: diag.Error,
        Summary:  "bodytemplate is not a valid JSON string",
        Detail:   err.Error(),
      })
    }
  }

	return c
}
