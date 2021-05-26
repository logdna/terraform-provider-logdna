package logdna

// This separation of concerns between request and response bodies is only due
// to inconsistencies in the API data types returned by the PUT versus the ones
// returned by the GET. In a perfect world, they would use the same types.

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type viewResponse struct {
	Apps     []string          `json:"apps,omitempty"`
	Category []string          `json:"category,omitempty"`
	Channels []channelResponse `json:"channels,omitempty"`
	Error    string            `json:"error,omitempty"`
	Hosts    []string          `json:"hosts,omitempty"`
	Levels   []string          `json:"levels,omitempty"`
	Name     string            `json:"name,omitempty"`
	Query    string            `json:"query,omitempty"`
	Tags     []string          `json:"tags,omitempty"`
	ViewID   string            `json:"viewID,omitempty"`
}

// channelResponse contains channel data returned from the logdna APIs
// NOTE - Properties with `interface` are due to the APIs returning
// some things as strings (PUT/emails) and other times arrays (GET/emails)
type channelResponse struct {
	AlertID         string            `json:"alertid,omitempty"`
	BodyTemplate    string            `json:"bodyTemplate,omitempty"`
	Emails          interface{}       `json:"emails,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	Immediate       bool              `json:"immediate,omitempty"`
	Integration     string            `json:"integration,omitempty"`
	Key             string            `json:"key,omitempty"`
	Method          string            `json:"method,omitempty"`
	Operator        string            `json:"operator,omitempty"`
	Terminal        bool              `json:"terminal,omitempty"`
	TriggerInterval interface{}       `json:"triggerinterval,omitempty"`
	TriggerLimit    int               `json:"triggerlimit,omitempty"`
	Timezone        string            `json:"timezone,omitempty"`
	URL             string            `json:"url,omitempty"`
}

func (view *viewResponse) MapChannelsToSchema() (map[string][]interface{}, diag.Diagnostics) {
	// This function iterates through the channel types and prepares the values
	// to be set on the schema in the correct keys
	var prepared interface{}
	var diags diag.Diagnostics

	channels := view.Channels
	channelIntegrations := make(map[string][]interface{})

	if channels == nil {
		return channelIntegrations, diags
	}
	for _, c := range channels {
		prepared = nil
		integration := c.Integration
		switch integration {
		case EMAIL:
			prepared = mapChannelEmail(&c)
		case PAGERDUTY:
			prepared = mapChannelPagerDuty(&c)
		case WEBHOOK:
			prepared = mapChannelWebhook(&c)
		default:
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("The remote view resource contains an unsupported integration: %s", integration),
				Detail:   fmt.Sprintf("%s integration ignored since it does not map to the schema", integration),
			})
		}
		if prepared == nil {
			continue
		}
		list := channelIntegrations[integration]
		if list == nil {
			channelIntegrations[integration] = make([]interface{}, 0)
		}
		channelIntegrations[integration] = append(list, prepared)
	}
	return channelIntegrations, diags
}

func mapChannelEmail(channel *channelResponse) map[string]interface{} {
	c := make(map[string]interface{})

	c["emails"] = channel.Emails
	c["immediate"] = strconv.FormatBool(channel.Immediate)
	c["operator"] = channel.Operator
	c["terminal"] = strconv.FormatBool(channel.Terminal)
	c["timezone"] = channel.Timezone
	c["triggerlimit"] = channel.TriggerLimit
	c["triggerinterval"] = channel.TriggerInterval

	return c
}

func mapChannelPagerDuty(channel *channelResponse) map[string]interface{} {
	c := make(map[string]interface{})

	c["immediate"] = strconv.FormatBool(channel.Immediate)
	c["key"] = channel.Key
	c["operator"] = channel.Operator
	c["terminal"] = strconv.FormatBool(channel.Terminal)
	c["triggerlimit"] = channel.TriggerLimit
	c["triggerinterval"] = channel.TriggerInterval

	return c
}
func mapChannelWebhook(channel *channelResponse) map[string]interface{} {
	c := make(map[string]interface{})

	c["bodytemplate"] = channel.BodyTemplate
	c["headers"] = channel.Headers
	c["immediate"] = strconv.FormatBool(channel.Immediate)
	c["method"] = channel.Method
	c["operator"] = channel.Operator
	c["terminal"] = strconv.FormatBool(channel.Terminal)
	c["triggerlimit"] = channel.TriggerLimit
	c["triggerinterval"] = channel.TriggerInterval
	c["url"] = channel.URL

	return c
}
