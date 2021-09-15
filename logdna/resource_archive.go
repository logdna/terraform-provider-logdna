package logdna

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const archiveConfigID = "archive"

type ibmConfig struct {
	Bucket             string `json:"bucket"`
	Endpoint           string `json:"endpoint"`
	ApiKey             string `json:"apikey"`
	ResourceInstanceId string `json:"resourceinstanceid"`
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
	ProjectId string `json:"projectid"`
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
	UserName   string `json:"username"`
	Password   string `json:"password"`
	TenantName string `json:"tenantname"`
}

func resourceArchiveConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	c := archiveRequestBody(d)

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
	fmt.Printf("struct (create): %v\n", cn)

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
	fmt.Printf("struct (read): %v\n", c)

	mapStructToSchema(c.Integration, c, d, diags)
	return diags
}

func resourceArchiveConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	pc := m.(*providerConfig)
	c := archiveRequestBody(d)

	req := newRequestConfig(
		pc,
		"PUT",
		"/v1/config/archiving",
		c,
	)

	_, err := req.MakeRequest()
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

func archiveRequestBody(d *schema.ResourceData) interface{} {
	rsint := d.Get("integration").(string)
	rsfld := fmt.Sprintf("%s_config", rsint)
	rscfg := d.Get(rsfld).([]interface{})[0].(map[string]interface{})
	fmt.Printf("integration: %s, field: %s, details: %v\n", rsint, rsfld, rscfg)

	if rsint == "ibm" {
		ibm := ibmConfig{
			Bucket:             rscfg["bucket"].(string),
			Endpoint:           rscfg["endpoint"].(string),
			ApiKey:             rscfg["apikey"].(string),
			ResourceInstanceId: rscfg["resourceinstanceid"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			ibmConfig
		}{rsint, ibm}
	} else if rsint == "s3" {
		s3 := s3Config{
			Bucket: rscfg["bucket"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			s3Config
		}{rsint, s3}
	} else if rsint == "azblob" {
		azblob := azblobConfig{
			AccountName: rscfg["accountname"].(string),
			AccountKey:  rscfg["accountkey"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			azblobConfig
		}{rsint, azblob}
	} else if rsint == "gcs" {
		gcs := gcsConfig{
			Bucket:    rscfg["bucket"].(string),
			ProjectId: rscfg["projectid"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			gcsConfig
		}{rsint, gcs}
	} else if rsint == "dos" {
		dos := dosConfig{
			Endpoint:  rscfg["endpoint"].(string),
			Space:     rscfg["space"].(string),
			AccessKey: rscfg["accesskey"].(string),
			SecretKey: rscfg["secretkey"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			dosConfig
		}{rsint, dos}
	} else { // final case has to be [integration:swift] based on its ValidateFunc
		swift := swiftConfig{
			AuthURL:    rscfg["authurl"].(string),
			Expires:    rscfg["expires"].(int),
			UserName:   rscfg["username"].(string),
			Password:   rscfg["password"].(string),
			TenantName: rscfg["tenantname"].(string),
		}
		return struct {
			Integration string `json:"integration"`
			swiftConfig
		}{rsint, swift}
	}
}

func mapStructToSchema(intgn string, strct archiveResponse, d *schema.ResourceData, diags diag.Diagnostics) {
	appendError(d.Set("integration", intgn), &diags)

	switch intgn {
	case "ibm":
		{
			IbmConfig := make(map[string]interface{})
			IbmConfig["bucket"] = strct.Bucket
			IbmConfig["endpoint"] = strct.Endpoint
			IbmConfig["apikey"] = strct.ApiKey
			IbmConfig["resourceinstanceid"] = strct.ResourceInstanceId
			appendError(d.Set("ibm_config", []interface{}{IbmConfig}), &diags)
		}
	case "s3":
		{
			S3Config := make(map[string]interface{})
			S3Config["bucket"] = strct.Bucket
			appendError(d.Set("s3_config", []interface{}{S3Config}), &diags)
		}
	case "azblob":
		{
			AzblobConfig := make(map[string]interface{})
			AzblobConfig["accountname"] = strct.AccountName
			AzblobConfig["accountkey"] = strct.AccountKey
			appendError(d.Set("azblob_config", []interface{}{AzblobConfig}), &diags)
		}
	case "gcs":
		{
			GcsConfig := make(map[string]interface{})
			GcsConfig["bucket"] = strct.Bucket
			GcsConfig["projectid"] = strct.ProjectId
			appendError(d.Set("gcs_config", []interface{}{GcsConfig}), &diags)
		}
	case "dos":
		{
			DosConfig := make(map[string]interface{})
			DosConfig["endpoint"] = strct.Endpoint
			DosConfig["space"] = strct.Space
			DosConfig["accesskey"] = strct.AccessKey
			DosConfig["secretkey"] = strct.SecretKey
			appendError(d.Set("dos_config", []interface{}{DosConfig}), &diags)
		}
	case "swift":
		{
			SwiftConfig := make(map[string]interface{})
			SwiftConfig["authurl"] = strct.AuthURL
			SwiftConfig["expires"] = strct.Expires
			SwiftConfig["username"] = strct.UserName
			SwiftConfig["password"] = strct.Password
			SwiftConfig["tenantname"] = strct.TenantName
			appendError(d.Set("swift_config", []interface{}{SwiftConfig}), &diags)
		}
	}
}
