package logdna

import (
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type config struct {
	ServiceKey string
	URL        string
	HTTPClient *http.Client
}

// Provider sets the schema to a servicekey and url and adds logdna_view as a resource
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

	return &config{ServiceKey: servicekey, URL: url, HTTPClient: &http.Client{Timeout: 15 * time.Second}}, nil
}
