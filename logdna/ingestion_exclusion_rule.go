package logdna

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type ingestionExclusionRule struct {
	ID        string   `json:"id,omitempty"`
	Title     string   `json:"title"`
	Active    bool     `json:"active"`
	IndexOnly bool     `json:"indexonly"`
	Apps      []string `json:"apps"`
	Hosts     []string `json:"hosts"`
	Query     string   `json:"query"`
}

var ingestionExclusionRuleAtLeastOneOfFields = []string{"apps", "hosts", "query"}

var ingestionExclusionRuleSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"title": {
		Type:     schema.TypeString,
		Default:  nil,
		Optional: true,
	},
	"active": {
		Type:     schema.TypeBool,
		Default:  false,
		Optional: true,
	},
	"indexonly": {
		Type:     schema.TypeBool,
		Default:  true,
		Optional: true,
	},
	"apps": {
		Type:         schema.TypeList,
		Elem:         &schema.Schema{Type: schema.TypeString},
		MinItems:     1,
		Optional:     true,
		AtLeastOneOf: ingestionExclusionRuleAtLeastOneOfFields,
	},
	"hosts": {
		Type:         schema.TypeList,
		Elem:         &schema.Schema{Type: schema.TypeString},
		MinItems:     1,
		Optional:     true,
		AtLeastOneOf: ingestionExclusionRuleAtLeastOneOfFields,
	},
	"query": {
		Type:         schema.TypeString,
		Optional:     true,
		AtLeastOneOf: ingestionExclusionRuleAtLeastOneOfFields,
	},
}
