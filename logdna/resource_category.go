package logdna

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "strings"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCategoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  var diags diag.Diagnostics
  pc := m.(*providerConfig)  

  // NOTE Type is't a part of a request body
  categoryType := d.Get("type").(string)
  category := categoryRequest{}

  if diags = category.CreateRequestBody(d); diags.HasError() {
    return diags
  }

  req := newRequestConfig(
    pc,
    "POST",
    fmt.Sprintf("/v1/config/categories/%s", categoryType),
    category,
  )

  body, err := req.MakeRequest()
  log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

  if err != nil {
    return diag.FromErr(err)
  }

  createdCategory := categoryResponse{}
  err = json.Unmarshal(body, &createdCategory)

  if err != nil {
    return diag.FromErr(err)
  }

  log.Printf("[DEBUG] After %s categories, the created category is %+v", req.method, createdCategory)

  // NOTE Type is added as a part of category ID to support import of categories
  //      Because type is required field even for read operation
  d.SetId(fmt.Sprintf("%s:%s", createdCategory.Type, createdCategory.Id))

  return resourceCategoryRead(ctx, d, m)
}

func resourceCategoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  var diags diag.Diagnostics
  pc := m.(*providerConfig)

  categoryType, categoryId, err := parseCategoryId(d.Id())

  if err != nil {
    return diag.FromErr(err)
  }

  category := categoryRequest{}

  if diags = category.CreateRequestBody(d); diags.HasError() {
    return diags
  }

  req := newRequestConfig(
    pc,
    "PUT",
    fmt.Sprintf("/v1/config/categories/%s/%s", categoryType, categoryId),
    category,
  )

  body, err := req.MakeRequest()
  log.Printf("[DEBUG] %s %s, payload is: %s", req.method, req.apiURL, body)

  if err != nil {
    return diag.FromErr(err)
  }

  log.Printf("[DEBUG] %s %s SUCCESS. Remote resource updated.", req.method, req.apiURL)

  return resourceCategoryRead(ctx, d, m)
}

func resourceCategoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  var diags diag.Diagnostics

  pc := m.(*providerConfig)

  categoryType, categoryId, err := parseCategoryId(d.Id())

  if err != nil {
    return diag.FromErr(err)
  }

  req := newRequestConfig(
    pc,
    "GET",
    fmt.Sprintf("/v1/config/categories/%s/%s", categoryType, categoryId),
    nil,
  )

  body, err := req.MakeRequest()

  log.Printf("[DEBUG] GET categories raw response body %s\n", body)
  if err != nil {
    diags = append(diags, diag.Diagnostic{
      Severity: diag.Error,
      Summary:  "Cannot read the remote categories resource",
      Detail:   err.Error(),
    })
    return diags
  }

  category := categoryResponse{}
  err = json.Unmarshal(body, &category)
  if err != nil {
    diags = append(diags, diag.Diagnostic{
      Severity: diag.Error,
      Summary:  "Cannot unmarshal response from the remote categories resource",
      Detail:   err.Error(),
    })
    return diags
  }
  log.Printf("[DEBUG] The GET categories structure is as follows: %+v\n", category)

  appendError(d.Set("type", category.Type), &diags)
  appendError(d.Set("name", category.Name), &diags)

  return diags
}

func resourceCategoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  pc := m.(*providerConfig)

  categoryType, categoryId, err := parseCategoryId(d.Id())

  if err != nil {
    return diag.FromErr(err)
  }

  req := newRequestConfig(
    pc,
    "DELETE",
    fmt.Sprintf("/v1/config/categories/%s/%s", categoryType, categoryId),
    nil,
  )

  body, err := req.MakeRequest()
  log.Printf("[DEBUG] %s %s presetalert %s", req.method, req.apiURL, body)

  if err != nil {
    return diag.FromErr(err)
  }
  d.SetId("")
  return nil
}

func parseCategoryId(id string) (string, string, error) {
  parts := strings.SplitN(id, ":", 2)

  if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
    return "", "", fmt.Errorf("Unexpected format of category ID (%s), expected Type:Id", id)
  }

  return parts[0], parts[1], nil
}

func resourceCategory() *schema.Resource {
  return &schema.Resource{
    CreateContext: resourceCategoryCreate,
    UpdateContext: resourceCategoryUpdate,
    ReadContext:   resourceCategoryRead,
    DeleteContext: resourceCategoryDelete,
    Importer: &schema.ResourceImporter{
      State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
        categoryType, categoryId, err := parseCategoryId(d.Id())

        if err != nil {
          return nil, err
        }

        if err := d.Set("type", categoryType); err != nil {
          return nil, err
        }

        d.SetId(fmt.Sprintf("%s:%s", categoryType, categoryId))

        return []*schema.ResourceData{d}, nil
      },
    },
    Schema: map[string]*schema.Schema{
      "name": {
        Type:     schema.TypeString,
        Required: true,
      },
      // NOTE Type is added to the schema but it's not used in a request body
      //      as the type is used just as a part of a url
      "type": {
        Type:     schema.TypeString,
        Optional: true,
        Default:  "views",
      },
    },
  }
}
