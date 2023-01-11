package logdna

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _ = registerTerraform(TerraformInfo{
	name:          "logdna_enterprise_key",
	orgType:       OrgTypeEnterprise,
	terraformType: TerraformTypeResource,
	schema:        resourceEnterpriseKey(),
})

func resourceEnterpriseKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pc := m.(*providerConfig)

	key := keyRequest{}

	if diags = key.CreateRequestBody(d); diags.HasError() {
		return diags
	}

	req := newRequestConfig(
		pc,
		"POST",
		"/v1/enterprise/keys",
		nil,
	)

	body, err := req.MakeRequest()
	log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

	if err != nil {
		return diag.FromErr(err)
	}

	createdKey := keyResponse{}
	err = json.Unmarshal(body, &createdKey)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] After %s key, the created key is %+v", req.method, createdKey)

	d.SetId(createdKey.KeyID)

	return resourceEnterpriseKeyRead(ctx, d, m)
}

func resourceEnterpriseKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	keyID := d.Id()

	req := newRequestConfig(
		pc,
		"GET",
		fmt.Sprintf("/v1/enterprise/keys/%s", keyID),
		nil,
	)

	body, err := req.MakeRequest()

	log.Printf("[DEBUG] GET key raw response body %s\n", body)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote key resource",
			Detail:   err.Error(),
		})
		return diags
	}

	key := keyResponse{}
	err = json.Unmarshal(body, &key)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote key resource",
			Detail:   err.Error(),
		})
		return diags
	}
	log.Printf("[DEBUG] The GET key structure is as follows: %+v\n", key)

	// Top level keys can be set directly
	appendError(d.Set("name", key.Name), &diags)
	appendError(d.Set("id", key.KeyID), &diags)
	appendError(d.Set("key", key.Key), &diags)
	appendError(d.Set("created", key.Created), &diags)

	return diags
}

func resourceEnterpriseKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	keyID := d.Id()

	req := newRequestConfig(
		pc,
		"DELETE",
		fmt.Sprintf("/v1/enterprise/keys/%s", keyID),
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

func resourceEnterpriseKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnterpriseKeyCreate,
		ReadContext:   resourceEnterpriseKeyRead,
		DeleteContext: resourceEnterpriseKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"key": {
				Type:      schema.TypeString,
				ForceNew:  true,
				Sensitive: true,
				Computed:  true,
			},
			"created": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}
