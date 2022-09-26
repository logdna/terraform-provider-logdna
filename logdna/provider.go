package logdna

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type providerConfig struct {
	serviceKey string
	orgType    OrgType
	baseURL    string
	httpClient *http.Client
}

// Provider initializes the schema with a service key and hooks for our resources
func Provider() *schema.Provider {
	dataSourceInfoMap = map[DataSourceType]TerraformInfo{
		DataSourceTypeAlert: {name: "logdna_alert", orgType: OrgTypeRegular, schema: dataSourceAlert()},
	}
	resourceInfoMap = map[ResourceType]TerraformInfo{
		ResourceTypeAlert:              {name: "logdna_alert", orgType: OrgTypeRegular, schema: resourceAlert()},
		ResourceTypeView:               {name: "logdna_view", orgType: OrgTypeRegular, schema: resourceView()},
		ResourceTypeCategory:           {name: "logdna_category", orgType: OrgTypeRegular, schema: resourceCategory()},
		ResourceTypeStreamConfig:       {name: "logdna_stream_config", orgType: OrgTypeRegular, schema: resourceStreamConfig()},
		ResourceTypeStreamExclusion:    {name: "logdna_stream_exclusion", orgType: OrgTypeRegular, schema: resourceStreamExclusion()},
		ResourceTypeIngestionExclusion: {name: "logdna_ingestion_exclusion", orgType: OrgTypeRegular, schema: resourceIngestionExclusion()},
		ResourceTypeArchive:            {name: "logdna_archive", orgType: OrgTypeRegular, schema: resourceArchiveConfig()},
		ResourceTypeKey:                {name: "logdna_key", orgType: OrgTypeRegular, schema: resourceKey()},
		ResourceTypeIndexRateAlert:     {name: "logdna_index_rate_alert", orgType: OrgTypeRegular, schema: resourceIndexRateAlert()},
		ResourceTypeMember:             {name: "logdna_member", orgType: OrgTypeRegular, schema: resourceMember()},
		ResourceTypeChildOrg:           {name: "logdna_enterprise_child_org", orgType: OrgTypeEnterprise, schema: resourceChildOrg()},
	}

	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"servicekey": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "regular",
				ValidateFunc: validation.StringInSlice([]string{"regular", "enterprise"}, false),
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://api.logdna.com",
			},
		},
		DataSourcesMap: buildSchemaMap(collectInfo(dataSourceInfoMap)),
		ResourcesMap:   buildSchemaMap(collectInfo(resourceInfoMap)),
		ConfigureFunc:  providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	serviceKey := d.Get("servicekey").(string)
	orgTypeRaw := d.Get("type").(string)
	url := d.Get("url").(string)

	orgType := OrgTypeRegular

	switch orgTypeRaw {
	case "regular":
		orgType = OrgTypeRegular
	case "enterprise":
		orgType = OrgTypeEnterprise
	}

	return &providerConfig{
		serviceKey: serviceKey,
		orgType:    orgType,
		baseURL:    url,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}, nil
}

func (pc *providerConfig) CheckOrgType(info TerraformInfo, diags diag.Diagnostics) diag.Diagnostics {
	if pc.orgType != info.orgType {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Only %s organizations can instantiate a \"%s\"", info.orgType, info.name),
		})
	}

	return diags
}
