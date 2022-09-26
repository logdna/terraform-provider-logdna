package logdna

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	diags = pc.CheckOrgType(resourceInfoMap[ResourceTypeMember], diags)
	if diags.HasError() {
		return diags
	}

	member := memberRequest{}

	if diags = member.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"POST",
		"/v1/config/members",
		member,
	)
	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	createdMember := memberResponse{}
	err = json.Unmarshal(body, &createdMember)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] After %s member, the created member is %+v", req.method, createdMember)

	d.SetId(createdMember.Email)

	return resourceMemberRead(ctx, d, m)
}

func resourceMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	diags = pc.CheckOrgType(resourceInfoMap[ResourceTypeMember], diags)
	if diags.HasError() {
		return diags
	}

	memberID := d.Id()

	req := newRequestConfig(
		pc,
		"GET",
		fmt.Sprintf("/v1/config/members/%s", memberID),
		nil,
	)

	body, err := req.MakeRequest()

	log.Printf("[DEBUG] GET member raw response body %s\n", body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote member resource",
			Detail:   err.Error(),
		})
		return diags
	}

	member := memberResponse{}
	err = json.Unmarshal(body, &member)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote member resource",
			Detail:   err.Error(),
		})
		return diags
	}
	log.Printf("[DEBUG] The GET member structure is as follows: %+v\n", member)

	// Top level keys can be set directly
	appendError(d.Set("email", member.Email), &diags)
	appendError(d.Set("role", member.Role), &diags)
	appendError(d.Set("groups", member.Groups), &diags)

	return diags
}

func resourceMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	diags = pc.CheckOrgType(resourceInfoMap[ResourceTypeMember], diags)
	if diags.HasError() {
		return diags
	}

	memberID := d.Id()

	member := memberPutRequest{}
	if diags = member.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"PUT",
		fmt.Sprintf("/v1/config/members/%s", memberID),
		member,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s %s SUCCESS. Remote resource updated.", req.method, req.apiURL)

	return resourceMemberRead(ctx, d, m)
}

func resourceMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	diags = pc.CheckOrgType(resourceInfoMap[ResourceTypeMember], diags)
	if diags.HasError() {
		return diags
	}

	memberID := d.Id()

	req := newRequestConfig(
		pc,
		"DELETE",
		fmt.Sprintf("/v1/config/members/%s", memberID),
		nil,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s key %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberCreate,
		UpdateContext: resourceMemberUpdate,
		ReadContext:   resourceMemberRead,
		DeleteContext: resourceMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"owner", "admin", "member", "readonly"}, false),
			},
			"groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}
