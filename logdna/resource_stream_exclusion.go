package logdna

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
  This resource needs to initialize before Terraform initializes so we can correctly populate the Provider schema.
  We can't use the init() function because Terraform initializes before that.
*/
var _ = registerTerraform(TerraformInfo{
	name:          "logdna_stream_exclusion",
	orgType:       OrgTypeRegular,
	terraformType: TerraformTypeResource,
	schema:        resourceStreamExclusion(),
})

func resourceStreamExclusionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	ex := exclusionRule{
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

	exn := exclusionRule{}
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

	ex := exclusionRule{}
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
	ex := exclusionRule{
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

		Schema: exclusionRuleSchema,
	}
}
