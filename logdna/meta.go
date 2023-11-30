package logdna

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type OrgType string
type TerraformType string
type TerraformInfo struct {
	name          string
	orgType       OrgType
	terraformType TerraformType
	schema        *schema.Resource
}

const (
	OrgTypeRegular    OrgType = "regular"
	OrgTypeEnterprise OrgType = "enterprise"
)

const (
	TerraformTypeResource   TerraformType = "resource"
	TerraformTypeDataSource TerraformType = "data source"
)

var terraformRegistry []TerraformInfo

func registerTerraform(info TerraformInfo) *TerraformInfo {
	terraformRegistry = append(terraformRegistry, info)
	infoPt := &terraformRegistry[len(terraformRegistry)-1]

	if infoPt.schema.CreateContext != nil {
		infoPt.schema.CreateContext = buildTerraformFunc(infoPt.schema.CreateContext, infoPt)
	}
	if infoPt.schema.ReadContext != nil {
		infoPt.schema.ReadContext = buildTerraformFunc(infoPt.schema.ReadContext, infoPt)
	}
	if infoPt.schema.UpdateContext != nil {
		infoPt.schema.UpdateContext = buildTerraformFunc(infoPt.schema.UpdateContext, infoPt)
	}
	if infoPt.schema.DeleteContext != nil {
		infoPt.schema.DeleteContext = buildTerraformFunc(infoPt.schema.DeleteContext, infoPt)
	}

	return infoPt
}

func filterRegistry(terraformType TerraformType) []TerraformInfo {
	newSlice := []TerraformInfo{}

	for _, info := range terraformRegistry {
		if info.terraformType == terraformType {
			newSlice = append(newSlice, info)
		}
	}

	return newSlice
}

func buildSchemaMap(a []TerraformInfo) map[string]*schema.Resource {
	m := make(map[string]*schema.Resource)

	for _, e := range a {
		m[e.name] = e.schema
	}

	return m
}

func buildTerraformFunc(contextFunc func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics, info *TerraformInfo) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		var diags diag.Diagnostics
		pc := m.(*providerConfig)

		if pc.orgType != info.orgType {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Only %s organizations can instantiate a \"%s\" %s", info.orgType, info.name, info.terraformType),
			})
			return diags
		}

		return contextFunc(ctx, d, m)
	}
}
