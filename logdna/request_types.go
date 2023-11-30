package logdna

// This separation of concerns between request and response bodies is only due
// to inconsistencies in the API data types returned by the PUT versus the ones
// returned by the GET. In a perfect world, they would use the same types.

import (
	"encoding/json"
	"errors"
	"fmt"

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
	PresetId string           `json:"presetid,omitempty"`
}

type alertRequest struct {
	Name     string           `json:"name,omitempty"`
	Channels []channelRequest `json:"channels,omitempty"`
}

type channelRequest struct {
	BodyTemplate        map[string]interface{} `json:"bodyTemplate,omitempty"`
	Emails              []string               `json:"emails,omitempty"`
	Headers             map[string]string      `json:"headers,omitempty"`
	Immediate           string                 `json:"immediate,omitempty"`
	Integration         string                 `json:"integration,omitempty"`
	Key                 string                 `json:"key,omitempty"`
	Method              string                 `json:"method,omitempty"`
	Operator            string                 `json:"operator,omitempty"`
	Terminal            string                 `json:"terminal,omitempty"`
	TriggerInterval     string                 `json:"triggerinterval,omitempty"`
	TriggerLimit        int                    `json:"triggerlimit,omitempty"`
	AutoResolve         bool                   `json:"autoresolve,omitempty"`
	AutoResolveInterval string                 `json:"autoresolveinterval,omitempty"`
	AutoResolveLimit    int                    `json:"autoresolvelimit,omitempty"`
	Timezone            string                 `json:"timezone,omitempty"`
	URL                 string                 `json:"url,omitempty"`
}

type categoryRequest struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type keyRequest struct {
	Name string `json:"name,omitempty"`
}

type indexRateAlertWebhookRequest struct {
	URL          string                 `json:"url,omitempty"`
	Method       string                 `json:"method,omitempty"`
	Headers      map[string]string      `json:"headers,omitempty"`
	BodyTemplate map[string]interface{} `json:"bodyTemplate,omitempty"`
}

type indexRateAlertChannelRequest struct {
	Email     []string                       `json:"email,omitempty"`
	Pagerduty []string                       `json:"pagerduty,omitempty"`
	Slack     []string                       `json:"slack,omitempty"`
	Webhook   []indexRateAlertWebhookRequest `json:"webhook,omitempty"`
}

type indexRateAlertRequest struct {
	MaxLines       int                          `json:"max_lines,omitempty"`
	MaxZScore      int                          `json:"max_z_score,omitempty"`
	ThresholdAlert string                       `json:"threshold_alert,omitempty"`
	Frequency      string                       `json:"frequency,omitempty"`
	Channels       indexRateAlertChannelRequest `json:"channels,omitempty"`
	Enabled        bool                         `json:"enabled,omitempty"`
}

type memberRequest struct {
	Email  string   `json:"email,omitempty"`
	Role   string   `json:"role,omitempty"`
	Groups []string `json:"groups,omitempty"`
}

type memberPutRequest struct {
	Role   string   `json:"role,omitempty"`
	Groups []string `json:"groups"`
}

type childOrgPutRequest struct {
	Retention int    `json:"retention"`
	Owner     string `json:"owner"`
}

func (view *viewRequest) CreateRequestBody(d *schema.ResourceData) diag.Diagnostics {
	// This function pulls from the schema in preparation to JSON marshal
	var diags diag.Diagnostics

	// Scalars
	view.Name = d.Get("name").(string)
	view.Query = d.Get("query").(string)

	// Simple arrays
	view.Apps = listToStrings(d.Get("apps").([]interface{}))
	view.Category = listToStrings(d.Get("categories").([]interface{}))
	view.Hosts = listToStrings(d.Get("hosts").([]interface{}))
	view.Levels = listToStrings(d.Get("levels").([]interface{}))
	view.Tags = listToStrings(d.Get("tags").([]interface{}))

	view.PresetId = d.Get("presetid").(string)

	// Complex array interfaces
	view.Channels = *aggregateAllChannelsFromSchema(d, &diags)

	return diags
}

func (alert *alertRequest) CreateRequestBody(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	// Scalars
	alert.Name = d.Get("name").(string)

	// Complex array interfaces
	alert.Channels = *aggregateAllChannelsFromSchema(d, &diags)

	return diags
}

func (category *categoryRequest) CreateRequestBody(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	// Scalars
	category.Name = d.Get("name").(string)

	return diags
}

func (key *keyRequest) CreateRequestBody(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	// Scalars
	key.Name = d.Get("name").(string)

	return diags
}

func (doc *indexRateAlertRequest) CreateRequestBody(d *schema.ResourceData) diag.Diagnostics {
	// This function pulls from the schema in preparation to JSON marshal
	var diags diag.Diagnostics

	var channels = d.Get("channels").([]interface{})
	if len(channels) > 1 {
		return diag.FromErr(
			errors.New("Index rate alert resource supports only one channels object"),
		)
	}

	doc.MaxLines = d.Get("max_lines").(int)
	doc.MaxZScore = d.Get("max_z_score").(int)
	doc.Enabled = d.Get("enabled").(bool)
	doc.ThresholdAlert = d.Get("threshold_alert").(string)
	doc.Frequency = d.Get("frequency").(string)

	var indexRateAlertChannel indexRateAlertChannelRequest
	var channel = channels[0].(map[string]interface{})

	indexRateAlertChannel.Email = listToStrings(channel["email"].([]interface{}))
	indexRateAlertChannel.Pagerduty = listToStrings(channel["pagerduty"].([]interface{}))
	indexRateAlertChannel.Slack = listToStrings(channel["slack"].([]interface{}))
	indexRateAlertChannel.Webhook = *aggregateIndexRateAlertWebhookFromSchema(d, &diags)
	doc.Channels = indexRateAlertChannel

	return diags
}

func (member *memberRequest) CreateRequestBody(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	// Scalars
	member.Email = d.Get("email").(string)
	member.Role = d.Get("role").(string)
	member.Groups = listToStrings(d.Get("groups").([]interface{}))

	return diags
}

func (member *memberPutRequest) CreateRequestBody(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	// Scalars
	member.Role = d.Get("role").(string)
	member.Groups = listToStrings(d.Get("groups").([]interface{}))

	return diags
}

func aggregateIndexRateAlertWebhookFromSchema(
	d *schema.ResourceData,
	diags *diag.Diagnostics,
) *[]indexRateAlertWebhookRequest {

	allWebhookEntries := make([]indexRateAlertWebhookRequest, 0)

	allWebhookEntries = append(
		allWebhookEntries,
		*iterateIndexRateAlertWebhookType(
			d.Get("webhook_channel").([]interface{}),
			diags,
		)...,
	)

	return &allWebhookEntries
}

func (childOrg *childOrgPutRequest) CreateRequestBody(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	// Scalars
	childOrg.Retention = d.Get("retention").(int)
	childOrg.Owner = d.Get("owner").(string)

	return diags
}

func aggregateAllChannelsFromSchema(
	d *schema.ResourceData,
	diags *diag.Diagnostics,
) *[]channelRequest {
	allChannelEntries := make([]channelRequest, 0)

	allChannelEntries = append(
		allChannelEntries,
		*iterateIntegrationType(
			d.Get("email_channel").([]interface{}),
			EMAIL,
			diags,
		)...,
	)

	allChannelEntries = append(
		allChannelEntries,
		*iterateIntegrationType(
			d.Get("pagerduty_channel").([]interface{}),
			PAGERDUTY,
			diags,
		)...,
	)

	allChannelEntries = append(
		allChannelEntries,
		*iterateIntegrationType(
			d.Get("slack_channel").([]interface{}),
			SLACK,
			diags,
		)...,
	)

	allChannelEntries = append(
		allChannelEntries,
		*iterateIntegrationType(
			d.Get("webhook_channel").([]interface{}),
			WEBHOOK,
			diags,
		)...,
	)

	return &allChannelEntries
}

func iterateIntegrationType(
	listEntries []interface{},
	integration string,
	diags *diag.Diagnostics,
) *[]channelRequest {
	var prepared interface{}
	channelRequests := []channelRequest{}

	if len(listEntries) == 0 {
		return &channelRequests
	}

	for _, entry := range listEntries {
		e := entry.(map[string]interface{})
		prepared = nil
		switch integration {
		case EMAIL:
			prepared = emailChannelRequest(e)
		case PAGERDUTY:
			prepared = pagerDutyChannelRequest(e)
		case SLACK:
			prepared = slackChannelRequest(e)
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

func iterateIndexRateAlertWebhookType(
	listEntries []interface{},
	diags *diag.Diagnostics,
) *[]indexRateAlertWebhookRequest {
	webhookRequests := []indexRateAlertWebhookRequest{}

	for _, entry := range listEntries {
		e := entry.(map[string]interface{})
		headersMap := make(map[string]string)

		for k, v := range e["headers"].(map[string]interface{}) {
			headersMap[k] = v.(string)
		}

		var c interface{}
		var bt map[string]interface{}

		if bodyTemplate := e["bodytemplate"].(string); bodyTemplate != "" {
			// See if the JSON is valid, but don't use the value or it will double encode
			err := json.Unmarshal([]byte(bodyTemplate), &bt)

			if err != nil {
				*diags = append(*diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "bodytemplate is not a valid JSON string",
					Detail:   err.Error(),
				})
			}
		}

		c = indexRateAlertWebhookRequest{
			Headers:      headersMap,
			Method:       e["method"].(string),
			URL:          e["url"].(string),
			BodyTemplate: bt,
		}

		webhookRequests = append(webhookRequests, c.(indexRateAlertWebhookRequest))
	}

	return &webhookRequests
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
		Immediate:           s["immediate"].(string),
		Integration:         PAGERDUTY,
		Key:                 s["key"].(string),
		Operator:            s["operator"].(string),
		Terminal:            s["terminal"].(string),
		TriggerInterval:     s["triggerinterval"].(string),
		TriggerLimit:        s["triggerlimit"].(int),
		AutoResolve:         s["autoresolve"].(bool),
		AutoResolveInterval: s["autoresolveinterval"].(string),
		AutoResolveLimit:    s["autoresolvelimit"].(int),
	}

	return c
}

func slackChannelRequest(s map[string]interface{}) channelRequest {
	c := channelRequest{
		Immediate:       s["immediate"].(string),
		Integration:     SLACK,
		Operator:        s["operator"].(string),
		Terminal:        s["terminal"].(string),
		TriggerInterval: s["triggerinterval"].(string),
		TriggerLimit:    s["triggerlimit"].(int),
		URL:             s["url"].(string),
	}

	return c
}

func webHookChannelRequest(
	s map[string]interface{},
	diags *diag.Diagnostics,
) channelRequest {
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

	if bodyTemplate := s["bodytemplate"].(string); bodyTemplate != "" {
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

func listToStrings(list []interface{}) []string {
	strs := make([]string, 0, len(list))
	for _, elem := range list {
		val, ok := elem.(string)
		if ok && val != "" {
			strs = append(strs, val)
		}
	}

	return strs
}
