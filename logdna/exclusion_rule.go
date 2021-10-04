package logdna

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type exclusionRule struct {
	ID     string   `json:"id,omitempty"`
	Title  string   `json:"title"`
	Active bool     `json:"active"`
	Apps   []string `json:"apps"`
	Hosts  []string `json:"hosts"`
	Query  string   `json:"query"`
}

var exclusionRuleAtLeastOneOfFields = []string{"apps", "hosts", "query"}

var exclusionRuleSchema = map[string]*schema.Schema{
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
	"apps": {
		Type:         schema.TypeList,
		Elem:         &schema.Schema{Type: schema.TypeString},
		MinItems:     1,
		Optional:     true,
		AtLeastOneOf: exclusionRuleAtLeastOneOfFields,
	},
	"hosts": {
		Type:         schema.TypeList,
		Elem:         &schema.Schema{Type: schema.TypeString},
		MinItems:     1,
		Optional:     true,
		AtLeastOneOf: exclusionRuleAtLeastOneOfFields,
	},
	"query": {
		Type:         schema.TypeString,
		Optional:     true,
		AtLeastOneOf: exclusionRuleAtLeastOneOfFields,
	},
}
