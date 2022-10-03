package logdna

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceChildOrgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	diags = pc.CheckOrgType(resourceInfoMap[ResourceTypeChildOrg], diags)
	if diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"POST",
		"/v1/enterprise/account",
		nil,
	)
	req.serviceKey = d.Get("servicekey").(string)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	createdChildOrg := childOrgCreateResponse{}
	err = json.Unmarshal(body, &createdChildOrg)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] After %s method, the created child org is %+v", req.method, createdChildOrg)

	d.SetId(createdChildOrg.Account)

	return resourceChildOrgRead(ctx, d, m)
}

func resourceChildOrgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)

	diags = pc.CheckOrgType(resourceInfoMap[ResourceTypeChildOrg], diags)
	if diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"GET",
		"/v1/enterprise/account",
		nil,
	)
	req.serviceKey = d.Get("servicekey").(string)

	body, err := req.MakeRequest()

	log.Printf("[DEBUG] GET child org raw response body %s\n", body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote child org resource",
			Detail:   err.Error(),
		})
		return diags
	}

	childOrg := childOrgGetResponse{}
	err = json.Unmarshal(body, &childOrg)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote childOrg resource",
			Detail:   err.Error(),
		})
		return diags
	}
	log.Printf("[DEBUG] The GET child org structure is as follows: %+v\n", childOrg)

	// Top level keys can be set directly
	// appendError(d.Set("retention", childOrg.Retention), &diags)

	return diags
}

func resourceChildOrgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	diags = pc.CheckOrgType(resourceInfoMap[ResourceTypeChildOrg], diags)
	if diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"DELETE",
		"/v1/enterprise/account",
		nil,
	)
	req.serviceKey = d.Get("serviceKey").(string)
	body, err := req.MakeRequest()
	log.Printf("[DEBUG] DELETE request body : %s", body)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceChildOrg() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChildOrgCreate,
		ReadContext:   resourceChildOrgRead,
		DeleteContext: resourceChildOrgDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"servicekey": {
				Type:      schema.TypeString,
				ForceNew:  true,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}
