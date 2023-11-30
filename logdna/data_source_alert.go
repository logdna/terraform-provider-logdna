package logdna

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var intSchema = &schema.Schema{
	Type:     schema.TypeInt,
	Computed: true,
}
var strSchema = &schema.Schema{
	Type:     schema.TypeString,
	Computed: true,
}
var boolSchema = &schema.Schema{
	Type:     schema.TypeBool,
	Computed: true,
}
var alertProps = map[string]*schema.Schema{
	"immediate":       strSchema,
	"operator":        strSchema,
	"terminal":        strSchema,
	"triggerinterval": strSchema,
	"triggerlimit":    intSchema,
}

var _ = registerTerraform(TerraformInfo{
	name:          "logdna_alert",
	orgType:       OrgTypeRegular,
	terraformType: TerraformTypeDataSource,
	schema:        dataSourceAlert(),
})

func dataSourceAlertRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	id := d.Get("presetid").(string)

	req := newRequestConfig(
		pc,
		"GET",
		fmt.Sprintf("/v1/config/presetalert/%s", id),
		nil,
	)

	body, err := req.MakeRequest()

	log.Printf("[DEBUG] GET presetalert raw response body %s\n", body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote presetalert resource",
			Detail:   err.Error(),
		})
		return diags
	}

	alert := alertResponse{}
	err = json.Unmarshal(body, &alert)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote presetalert resource",
			Detail:   err.Error(),
		})
		return diags
	}
	log.Printf("GET presetalert structure is as follows: %+v\n", alert)

	appendError(d.Set("name", alert.Name), &diags)

	ints, diags := alert.MapChannelsToSchema()
	log.Printf("[DEBUG] presetalert MapChannelsToSchema result: %+v\n", ints)

	for name, value := range ints {
		if len(value) == 0 {
			continue
		}

		key := fmt.Sprintf("%s_channel", name)
		appendError(d.Set(key, value), &diags)
	}

	d.SetId(id)
	return diags
}

func dataSourceAlert() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlertRead,
		Schema: map[string]*schema.Schema{
			"presetid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": strSchema,
			"email_channel": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: getAlertSchema("email"),
				},
				Computed: true,
			},
			"pagerduty_channel": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: getAlertSchema("pagerduty"),
				},
				Computed: true,
			},
			"slack_channel": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: getAlertSchema("slack"),
				},
				Computed: true,
			},
			"webhook_channel": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: getAlertSchema("webhook"),
				},
				Computed: true,
			},
		},
	}
}

func getAlertSchema(chnl string) map[string]*schema.Schema {
	schma := map[string]*schema.Schema{}
	for key, value := range alertProps {
		schma[key] = value
	}

	switch chnl {
	case "email":
		schma["timezone"] = strSchema
		schma["emails"] = &schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
		}
	case "slack":
		schma["url"] = strSchema
	case "pagerduty":
		schma["key"] = strSchema
		schma["autoresolve"] = boolSchema
		schma["autoresolvelimit"] = intSchema
		schma["autoresolveinterval"] = strSchema
	case "webhook":
		schma["bodytemplate"] = strSchema
		schma["method"] = strSchema
		schma["url"] = strSchema
		schma["headers"] = &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
		}
	}

	return schma
}
