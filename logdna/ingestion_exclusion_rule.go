package logdna

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type ingestionExclusionRule struct {
	exclusionRule
	IndexOnly bool `json:"indexonly"`
}

var ingestionExclusionRuleSchema = map[string]*schema.Schema{
	"indexonly": {
		Type:     schema.TypeBool,
		Default:  true,
		Optional: true,
	},
}

func init() {
	for k, v := range exclusionRuleSchema {
		ingestionExclusionRuleSchema[k] = v
	}
}
