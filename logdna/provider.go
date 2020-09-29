package logdna

import (
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type config struct {
	servicekey string
	url        string
	httpClient *http.Client
}

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
			"logdna_view": resourceView(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	servicekey := d.Get("servicekey").(string)
	url := d.Get("url").(string)

	return &config{servicekey: servicekey, url: url, httpClient: &http.Client{Timeout: 15 * time.Second}}, nil
}
