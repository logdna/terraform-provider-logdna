package logdna

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const archiveConfigID = "archive"

var _ = registerTerraform(TerraformInfo{
	name:          "logdna_archive",
	orgType:       OrgTypeRegular,
	terraformType: TerraformTypeResource,
	schema:        resourceArchiveConfig(),
})

type ibmConfig struct {
	Bucket             string `json:"bucket"`
	Endpoint           string `json:"endpoint"`
	APIKey             string `json:"apikey"`
	ResourceInstanceID string `json:"resourceinstanceid"`
}

type s3Config struct {
	Bucket string `json:"bucket"`
}

type azblobConfig struct {
	AccountName string `json:"accountname"`
	AccountKey  string `json:"accountkey"`
}

type gcsConfig struct {
	Bucket    string `json:"bucket"`
	ProjectID string `json:"projectid"`
}

type dosConfig struct {
	Space     string `json:"space"`
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"accesskey"`
	SecretKey string `json:"secretkey"`
}

type swiftConfig struct {
	AuthURL    string `json:"authurl"`
	Expires    int    `json:"expires,omitempty"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	TenantName string `json:"tenantname"`
}

func generateArchiveConfig(d *schema.ResourceData) (interface{}, error) {
	integration := d.Get("integration").(string)
	configKey := fmt.Sprintf(`%s_config`, integration)
	configRaw := d.Get(configKey).([]interface{})
	if len(configRaw) == 0 {
		err := fmt.Errorf("expected %s_config for integration: %s", integration, integration)
		return nil, err
	}
	config := configRaw[0].(map[string]interface{})

	if integration == "ibm" {
		ibm := ibmConfig{
			Bucket:             config["bucket"].(string),
			Endpoint:           config["endpoint"].(string),
			APIKey:             config["apikey"].(string),
			ResourceInstanceID: config["resourceinstanceid"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			ibmConfig
		}{integration, ibm}, nil
	} else if integration == "s3" {
		s3 := s3Config{
			Bucket: config["bucket"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			s3Config
		}{integration, s3}, nil
	} else if integration == "azblob" {
		azblob := azblobConfig{
			AccountName: config["accountname"].(string),
			AccountKey:  config["accountkey"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			azblobConfig
		}{integration, azblob}, nil
	} else if integration == "gcs" {
		gcs := gcsConfig{
			Bucket:    config["bucket"].(string),
			ProjectID: config["projectid"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			gcsConfig
		}{integration, gcs}, nil
	} else if integration == "dos" {
		dos := dosConfig{
			Space:     config["space"].(string),
			Endpoint:  config["endpoint"].(string),
			AccessKey: config["accesskey"].(string),
			SecretKey: config["secretkey"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			dosConfig
		}{integration, dos}, nil
	} else {
		swift := swiftConfig{
			AuthURL:    config["authurl"].(string),
			Expires:    config["expires"].(int),
			Username:   config["username"].(string),
			Password:   config["password"].(string),
			TenantName: config["tenantname"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			swiftConfig
		}{integration, swift}, nil
	}
}

func setArchiveConfig(cn archiveResponse, d *schema.ResourceData, diags diag.Diagnostics) {
	integration := cn.Integration
	appendError(d.Set("integration", integration), &diags)

	switch integration {
	case "ibm":
		ibmConfig := make(map[string]interface{})
		ibmConfig["bucket"] = cn.Bucket
		ibmConfig["endpoint"] = cn.Endpoint
		ibmConfig["apikey"] = cn.APIKey
		ibmConfig["resourceinstanceid"] = cn.ResourceInstanceID
		appendError(d.Set("ibm_config", []interface{}{ibmConfig}), &diags)
	case "s3":
		s3Config := make(map[string]interface{})
		s3Config["bucket"] = cn.Bucket
		appendError(d.Set("s3_config", []interface{}{s3Config}), &diags)
	case "azblob":
		azblobConfig := make(map[string]interface{})
		azblobConfig["accountname"] = cn.AccountName
		azblobConfig["accountkey"] = cn.AccountKey
		appendError(d.Set("azblob_config", []interface{}{azblobConfig}), &diags)
	case "gcs":
		gcsConfig := make(map[string]interface{})
		gcsConfig["bucket"] = cn.Bucket
		gcsConfig["projectid"] = cn.ProjectID
		appendError(d.Set("gcs_config", []interface{}{gcsConfig}), &diags)
	case "dos":
		dosConfig := make(map[string]interface{})
		dosConfig["space"] = cn.Space
		dosConfig["endpoint"] = cn.Endpoint
		dosConfig["accesskey"] = cn.AccessKey
		dosConfig["secretkey"] = cn.SecretKey
		appendError(d.Set("dos_config", []interface{}{dosConfig}), &diags)
	case "swift":
		swiftConfig := make(map[string]interface{})
		swiftConfig["authurl"] = cn.AuthURL
		swiftConfig["expires"] = cn.Expires
		swiftConfig["username"] = cn.Username
		swiftConfig["password"] = cn.Password
		swiftConfig["tenantname"] = cn.TenantName
		appendError(d.Set("swift_config", []interface{}{swiftConfig}), &diags)
	}
}

func resourceArchiveConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	c, err := generateArchiveConfig(d)

	if err != nil {
		return diag.FromErr(err)
	}

	req := newRequestConfig(
		pc,
		"POST",
		"/v1/config/archiving",
		c,
	)

	body, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	cn := archiveResponse{}
	err = json.Unmarshal(body, &cn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(archiveConfigID)

	return resourceArchiveConfigRead(ctx, d, m)
}

func resourceArchiveConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	pc := m.(*providerConfig)
	req := newRequestConfig(
		pc,
		"GET",
		"/v1/config/archiving",
		nil,
	)

	body, err := req.MakeRequest()

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot read the remote archive resource",
			Detail:   err.Error(),
		})
		return diags
	}

	c := archiveResponse{}
	err = json.Unmarshal(body, &c)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Cannot unmarshal response from the remote archive resource",
			Detail:   err.Error(),
		})
		return diags
	}

	setArchiveConfig(c, d, diags)
	return diags
}

func resourceArchiveConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	c, err := generateArchiveConfig(d)

	if err != nil {
		return diag.FromErr(err)
	}

	req := newRequestConfig(
		pc,
		"PUT",
		"/v1/config/archiving",
		c,
	)

	body, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	cn := archiveResponse{}
	err = json.Unmarshal(body, &cn)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceArchiveConfigRead(ctx, d, m)
}

func resourceArchiveConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	req := newRequestConfig(
		pc,
		"DELETE",
		"/v1/config/archiving",
		nil,
	)

	_, err := req.MakeRequest()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceArchiveConfig() *schema.Resource {
	validIntegrations := []string{"ibm", "s3", "azblob", "gcs", "dos", "swift"}

	return &schema.Resource{
		CreateContext: resourceArchiveConfigCreate,
		ReadContext:   resourceArchiveConfigRead,
		UpdateContext: resourceArchiveConfigUpdate,
		DeleteContext: resourceArchiveConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"integration": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					valid := false
					v := val.(string)
					for _, str := range validIntegrations {
						if v == str {
							valid = true
						}
					}
					if !valid {
						errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validIntegrations, v))
					}
					return
				},
			},
			"ibm_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"apikey": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resourceinstanceid": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"s3_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"azblob_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accountname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"accountkey": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"gcs_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"projectid": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"dos_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint": {
							Type:     schema.TypeString,
							Required: true,
						},
						"space": {
							Type:     schema.TypeString,
							Required: true,
						},
						"accesskey": {
							Type:     schema.TypeString,
							Required: true,
						},
						"secretkey": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"swift_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authurl": {
							Type:     schema.TypeString,
							Required: true,
						},
						"expires": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tenantname": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}
