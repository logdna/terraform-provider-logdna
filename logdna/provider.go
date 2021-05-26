package logdna

import (
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type providerConfig struct {
	serviceKey string
	baseURL       string
	httpClient *http.Client
}

// Provider initializes the schema with a service key and hooks for our resources
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"servicekey": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://api.logdna.com",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"logdna_alert": resourceAlert(),
			"logdna_view":  resourceView(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	servicekey := d.Get("servicekey").(string)
	url := d.Get("url").(string)

	return &providerConfig{
		serviceKey: servicekey,
		baseURL:       url,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}, nil
}
