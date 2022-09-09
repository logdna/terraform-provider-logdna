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
	Apps      []string          `json:"apps,omitempty"`
	Category  []string          `json:"category,omitempty"`
	Channels  []channelResponse `json:"channels,omitempty"`
	Error     string            `json:"error,omitempty"`
	Hosts     []string          `json:"hosts,omitempty"`
	Levels    []string          `json:"levels,omitempty"`
	Name      string            `json:"name,omitempty"`
	Query     string            `json:"query,omitempty"`
	Tags      []string          `json:"tags,omitempty"`
	PresetIds []string          `json:"presetids,omitempty"`
	ViewID    string            `json:"viewID"`
}

type alertResponse struct {
	Name     string            `json:"name,omitempty"`
	Channels []channelResponse `json:"channels,omitempty"`
	PresetID string            `json:"presetid"`
}

type keyResponse struct {
	KeyID   string `json:"id"`
	Key     string `json:"key"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Created int    `json:"created,omitempty"`
}

type memberResponse struct {
	Email  string   `json:"email"`
	Role   string   `json:"role"`
	Groups []string `json:"groups,omitempty"`
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

type archiveResponse struct {
	Integration        string `json:"integration"`
	Bucket             string `json:"bucket,omitempty"`
	Endpoint           string `json:"endpoint,omitempty"`
	APIKey             string `json:"apikey,omitempty"`
	ResourceInstanceID string `json:"resourceinstanceid,omitempty"`
	AccountName        string `json:"accountname,omitempty"`
	AccountKey         string `json:"accountkey,omitempty"`
	ProjectID          string `json:"projectid,omitempty"`
	Space              string `json:"space,omitempty"`
	AccessKey          string `json:"accesskey,omitempty"`
	SecretKey          string `json:"secretkey,omitempty"`
	AuthURL            string `json:"authurl,omitempty"`
	Expires            int    `json:"expires,omitempty"`
	Username           string `json:"username,omitempty"`
	Password           string `json:"password,omitempty"`
	TenantName         string `json:"tenantname,omitempty"`
}

type categoryResponse struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Id   string `json:"id"`
}

type indexRateAlertChannelResponse struct {
	Email     []string `json:"email,omitempty"`
	Pagerduty []string `json:"pagerduty,omitempty"`
	Slack     []string `json:"slack,omitempty"`
}

type indexRateAlertResponse struct {
	MaxLines       int                           `json:"max_lines,omitempty"`
	MaxZScore      int                           `json:"max_z_score,omitempty"`
	ThresholdAlert string                        `json:"threshold_alert,omitempty"`
	Frequency      string                        `json:"frequency,omitempty"`
	Channels       indexRateAlertChannelResponse `json:"channels,omitempty"`
	Enabled        bool                          `json:"enabled,omitempty"`
}

func (view *viewResponse) MapChannelsToSchema() (map[string][]interface{}, diag.Diagnostics) {
	channels := view.Channels
	channelIntegrations, diags := mapAllChannelsToSchema("view", &channels)
	return channelIntegrations, *diags
}

func (alert *alertResponse) MapChannelsToSchema() (map[string][]interface{}, diag.Diagnostics) {
	channels := alert.Channels
	channelIntegrations, diags := mapAllChannelsToSchema("alert", &channels)
	return channelIntegrations, *diags
}

func mapAllChannelsToSchema(
	resourceName string,
	channels *[]channelResponse,
) (map[string][]interface{}, *diag.Diagnostics) {
	// This function iterates through the channel types and prepares the values
	// to be set on the schema in the correct keys
	var prepared interface{}
	var diags diag.Diagnostics

	channelIntegrations := map[string][]interface{}{
		EMAIL:     make([]interface{}, 0),
		PAGERDUTY: make([]interface{}, 0),
		SLACK:     make([]interface{}, 0),
		WEBHOOK:   make([]interface{}, 0),
	}

	if len(*channels) == 0 {
		return channelIntegrations, &diags
	}
	for _, c := range *channels {
		prepared = nil
		integration := c.Integration

		switch integration {
		case EMAIL:
			prepared = mapChannelEmail(&c)
		case PAGERDUTY:
			prepared = mapChannelPagerDuty(&c)
		case SLACK:
			prepared = mapChannelSlack(&c)
		case WEBHOOK:
			prepared = mapChannelWebhook(&c)
		default:
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("The remote %s resource contains an unsupported integration: %s", resourceName, integration),
				Detail:   fmt.Sprintf("%s integration ignored since it does not map to the schema", integration),
			})
		}
		if prepared == nil {
			continue
		}
		channelIntegrations[integration] = append(
			channelIntegrations[integration],
			prepared,
		)
	}
	return channelIntegrations, &diags
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

func mapChannelSlack(channel *channelResponse) map[string]interface{} {
	c := make(map[string]interface{})

	c["immediate"] = strconv.FormatBool(channel.Immediate)
	c["operator"] = channel.Operator
	c["terminal"] = strconv.FormatBool(channel.Terminal)
	c["triggerlimit"] = channel.TriggerLimit
	c["triggerinterval"] = channel.TriggerInterval
	c["url"] = channel.URL

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

func appendError(err error, diags *diag.Diagnostics) *diag.Diagnostics {
	if err != nil {
		*diags = append(*diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "There was a problem setting the schema",
			Detail:   err.Error(),
		})
	}
	return diags
}
