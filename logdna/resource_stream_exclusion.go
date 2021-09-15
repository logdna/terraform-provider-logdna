package logdna

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type streamExclusion struct {
	ID     string   `json:"id,omitempty"`
	Title  string   `json:"title"`
	Active bool     `json:"active"`
	Apps   []string `json:"apps"`
	Hosts  []string `json:"hosts"`
	Query  string   `json:"query"`
}

var streamExclusionAtLeastOneOfFields = []string{"apps", "hosts", "query"}

func resourceStreamExclusionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	ex := streamExclusion{
		Title:  d.Get("title").(string),
		Active: d.Get("active").(bool),
		Apps:   listToStrings(d.Get("apps").([]interface{})),
		Hosts:  listToStrings(d.Get("hosts").([]interface{})),
		Query:  d.Get("query").(string),
	}

	req := newRequestConfig(
		pc,
		"POST",
		"/v1/config/stream/exclusions",
		ex,
	)

	body, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	exn := streamExclusion{}
	err = json.Unmarshal(body, &exn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(exn.ID)
	appendError(d.Set("title", exn.Title), &diags)
	appendError(d.Set("active", exn.Active), &diags)
	appendError(d.Set("apps", exn.Apps), &diags)
	appendError(d.Set("hosts", exn.Hosts), &diags)
	appendError(d.Set("query", exn.Query), &diags)

	return diags
}

func resourceStreamExclusionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	req := newRequestConfig(
		pc,
		"GET",
		fmt.Sprintf("/v1/config/stream/exclusions/%s", d.Id()),
		nil,
	)

	body, err := req.MakeRequest()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote stream exclusion resource",
			Detail:   err.Error(),
		})
		return diags
	}

	ex := streamExclusion{}
	err = json.Unmarshal(body, &ex)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote stream exclusion resource",
			Detail:   err.Error(),
		})
		return diags
	}

	appendError(d.Set("title", ex.Title), &diags)
	appendError(d.Set("active", ex.Active), &diags)
	appendError(d.Set("apps", ex.Apps), &diags)
	appendError(d.Set("hosts", ex.Hosts), &diags)
	appendError(d.Set("query", ex.Query), &diags)

	return diags
}

func resourceStreamExclusionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	ex := streamExclusion{
		Title:  d.Get("title").(string),
		Active: d.Get("active").(bool),
		Apps:   listToStrings(d.Get("apps").([]interface{})),
		Hosts:  listToStrings(d.Get("hosts").([]interface{})),
		Query:  d.Get("query").(string),
	}

	req := newRequestConfig(
		pc,
		"PATCH",
		fmt.Sprintf("/v1/config/stream/exclusions/%s", d.Id()),
		ex,
	)

	_, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStreamExclusionRead(ctx, d, m)
}

func resourceStreamExclusionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	req := newRequestConfig(
		pc,
		"DELETE",
		fmt.Sprintf("/v1/config/stream/exclusions/%s", d.Id()),
		nil,
	)

	_, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceStreamExclusion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStreamExclusionCreate,
		ReadContext:   resourceStreamExclusionRead,
		UpdateContext: resourceStreamExclusionUpdate,
		DeleteContext: resourceStreamExclusionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
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
				AtLeastOneOf: streamExclusionAtLeastOneOfFields,
			},
			"hosts": {
				Type:         schema.TypeList,
				Elem:         &schema.Schema{Type: schema.TypeString},
				MinItems:     1,
				Optional:     true,
				AtLeastOneOf: streamExclusionAtLeastOneOfFields,
			},
			"query": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: streamExclusionAtLeastOneOfFields,
			},
		},
	}
}
