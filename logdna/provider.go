package logdna

import (
	"net/http"
	"time"

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
		DataSourcesMap: buildSchemaMap(filterRegistry(TerraformTypeDataSource)),
		ResourcesMap:   buildSchemaMap(filterRegistry(TerraformTypeResource)),
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
