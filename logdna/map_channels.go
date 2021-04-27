package logdna

import (
	"strconv"
  "fmt"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func (view *ViewResponse) MapChannelsToSchema() (map[string][]interface{}, diag.Diagnostics) {
	// This function iterations through the channels types and prepares the values
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
        Summary: fmt.Sprintf("The remote view resource contains an unsupported integration: %s", integration),
        Detail: fmt.Sprintf("%s integration ignored since it does not map to the schema", integration),
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

func mapChannelEmail(channel *ChannelResponse) map[string]interface{} {
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

func mapChannelPagerDuty(channel *ChannelResponse) map[string]interface{} {
	c := make(map[string]interface{})

	c["immediate"] = strconv.FormatBool(channel.Immediate)
  c["key"] = channel.Key
	c["operator"] = channel.Operator
	c["terminal"] = strconv.FormatBool(channel.Terminal)
	c["triggerlimit"] = channel.TriggerLimit
	c["triggerinterval"] = channel.TriggerInterval

	return c
}
func mapChannelWebhook(channel *ChannelResponse) map[string]interface{} {
	c := make(map[string]interface{})

	c["bodytemplate"] = channel.BodyTemplate
  c["headers"] = channel.Headers
	c["immediate"] = strconv.FormatBool(channel.Immediate)
  c["method"] = channel.Method
	c["operator"] = channel.Operator
	c["terminal"] = strconv.FormatBool(channel.Terminal)
	c["triggerlimit"] = channel.TriggerLimit
	c["triggerinterval"] = channel.TriggerInterval

	return c
}