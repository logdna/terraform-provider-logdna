package logdna

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const streamConfigID = "stream"

var _ = registerTerraform(TerraformInfo{
	name:          "logdna_stream_config",
	orgType:       OrgTypeRegular,
	terraformType: TerraformTypeResource,
	schema:        resourceStreamConfig(),
})

type streamConfig struct {
	Status   string   `json:"status,omitempty"`
	Brokers  []string `json:"brokers"`
	Topic    string   `json:"topic"`
	User     string   `json:"user"`
	Password string   `json:"password"`
}

func resourceStreamConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	c := streamConfig{
		Brokers:  listToStrings(d.Get("brokers").([]interface{})),
		Topic:    d.Get("topic").(string),
		User:     d.Get("user").(string),
		Password: d.Get("password").(string),
	}

	req := newRequestConfig(
		pc,
		"POST",
		"/v1/config/stream",
		c,
	)

	body, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	cn := streamConfig{}
	err = json.Unmarshal(body, &cn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(streamConfigID)
	appendError(d.Set("brokers", cn.Brokers), &diags)
	appendError(d.Set("topic", cn.Topic), &diags)
	appendError(d.Set("user", cn.User), &diags)
	appendError(d.Set("status", cn.Status), &diags)

	return diags
}

func resourceStreamConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	req := newRequestConfig(
		pc,
		"GET",
		"/v1/config/stream",
		nil,
	)

	body, err := req.MakeRequest()

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote stream config resource",
			Detail:   err.Error(),
		})
		return diags
	}

	c := streamConfig{}
	err = json.Unmarshal(body, &c)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote stream config resource",
			Detail:   err.Error(),
		})
		return diags
	}

	appendError(d.Set("brokers", c.Brokers), &diags)
	appendError(d.Set("topic", c.Topic), &diags)
	appendError(d.Set("user", c.User), &diags)
	appendError(d.Set("status", c.Status), &diags)

	return diags
}

func resourceStreamConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	c := streamConfig{
		Brokers:  listToStrings(d.Get("brokers").([]interface{})),
		Topic:    d.Get("topic").(string),
		User:     d.Get("user").(string),
		Password: d.Get("password").(string),
	}

	req := newRequestConfig(
		pc,
		"PUT",
		"/v1/config/stream",
		c,
	)

	_, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStreamConfigRead(ctx, d, m)
}

func resourceStreamConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	req := newRequestConfig(
		pc,
		"DELETE",
		"/v1/config/stream",
		nil,
	)

	_, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceStreamConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStreamConfigCreate,
		ReadContext:   resourceStreamConfigRead,
		UpdateContext: resourceStreamConfigUpdate,
		DeleteContext: resourceStreamConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"brokers": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
				Required: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}
