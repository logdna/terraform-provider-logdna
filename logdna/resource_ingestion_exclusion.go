package logdna

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const baseIngestionExclusionUrl = "/v1/config/ingestion/exclusions"

func resourceIngestionExclusionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	ex := ingestionExclusionRule{
		exclusionRule: exclusionRule{
			Title:  d.Get("title").(string),
			Active: d.Get("active").(bool),
			Apps:   listToStrings(d.Get("apps").([]interface{})),
			Hosts:  listToStrings(d.Get("hosts").([]interface{})),
			Query:  d.Get("query").(string),
		},
		IndexOnly: d.Get("indexonly").(bool),
	}

	req := newRequestConfig(
		pc,
		"POST",
		baseIngestionExclusionUrl,
		ex,
	)

	body, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	exn := ingestionExclusionRule{}
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

func resourceIngestionExclusionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	req := newRequestConfig(
		pc,
		"GET",
		fmt.Sprintf("%s/%s", baseIngestionExclusionUrl, d.Id()),
		nil,
	)

	body, err := req.MakeRequest()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote ingestion exclusion resource",
			Detail:   err.Error(),
		})
		return diags
	}

	ex := ingestionExclusionRule{}
	err = json.Unmarshal(body, &ex)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote ingestion exclusion resource",
			Detail:   err.Error(),
		})
		return diags
	}

	appendError(d.Set("title", ex.Title), &diags)
	appendError(d.Set("active", ex.Active), &diags)
	appendError(d.Set("indexonly", ex.IndexOnly), &diags)
	appendError(d.Set("apps", ex.Apps), &diags)
	appendError(d.Set("hosts", ex.Hosts), &diags)
	appendError(d.Set("query", ex.Query), &diags)

	return diags
}

func resourceIngestionExclusionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	ex := ingestionExclusionRule{
		exclusionRule: exclusionRule{
			Title:  d.Get("title").(string),
			Active: d.Get("active").(bool),
			Apps:   listToStrings(d.Get("apps").([]interface{})),
			Hosts:  listToStrings(d.Get("hosts").([]interface{})),
			Query:  d.Get("query").(string),
		},
		IndexOnly: d.Get("indexonly").(bool),
	}

	req := newRequestConfig(
		pc,
		"PATCH",
		fmt.Sprintf("%s/%s", baseIngestionExclusionUrl, d.Id()),
		ex,
	)

	_, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIngestionExclusionRead(ctx, d, m)
}

func resourceIngestionExclusionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	req := newRequestConfig(
		pc,
		"DELETE",
		fmt.Sprintf("%s/%s", baseIngestionExclusionUrl, d.Id()),
		nil,
	)

	_, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceIngestionExclusion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIngestionExclusionCreate,
		ReadContext:   resourceIngestionExclusionRead,
		UpdateContext: resourceIngestionExclusionUpdate,
		DeleteContext: resourceIngestionExclusionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: ingestionExclusionRuleSchema,
	}
}
