package logdna

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
  This resource needs to initialize before Terraform initializes so we can correctly populate the Provider schema.
  We can't use the init() function because Terraform initializes before that.
*/
var _ = registerTerraform(TerraformInfo{
	name:          "logdna_child_organization",
	orgType:       OrgTypeEnterprise,
	terraformType: TerraformTypeResource,
	schema:        resourceChildOrg(),
})

func resourceChildOrgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)

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

	createdChildOrg := childOrgResponse{}
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
	childOrgID := d.Id()

	req := newRequestConfig(
		pc,
		"GET",
		fmt.Sprintf("/v1/enterprise/account/%s", childOrgID),
		nil,
	)

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

	childOrg := childOrgResponse{}
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
	appendError(d.Set("retention", childOrg.Retention), &diags)
	appendError(d.Set("retention_tiers", childOrg.RetentionTiers), &diags)
	appendError(d.Set("owner", childOrg.Owner), &diags)

	return diags
}

func resourceChildOrgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)
	childOrgID := d.Id()

	childOrg := childOrgPutRequest{}
	if diags = childOrg.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"PUT",
		fmt.Sprintf("/v1/enterprise/account/%s", childOrgID),
		childOrg,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s %s SUCCESS. Remote resource updated.", req.method, req.apiURL)

	return resourceChildOrgRead(ctx, d, m)
}

func resourceChildOrgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	childOrgID := d.Id()

	req := newRequestConfig(
		pc,
		"DELETE",
		fmt.Sprintf("/v1/enterprise/account/%s", childOrgID),
		nil,
	)

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
		UpdateContext: resourceChildOrgUpdate,
		DeleteContext: resourceChildOrgDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"retention": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"retention_tiers": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
			},
			"servicekey": {
				Type:      schema.TypeString,
				ForceNew:  true,
				Required:  true,
				Sensitive: true,
				DiffSuppressFunc: func(_, _, _ string, _ *schema.ResourceData) bool {
					return false
				},
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
		},
	}
}
