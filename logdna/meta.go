package logdna

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

type OrgType string
type DataSourceType int64
type ResourceType int64

const (
	OrgTypeRegular    OrgType = "regular"
	OrgTypeEnterprise         = "enterprise"
)

const (
	DataSourceTypeAlert DataSourceType = iota
)

const (
	ResourceTypeAlert ResourceType = iota
	ResourceTypeView
	ResourceTypeCategory
	ResourceTypeStreamConfig
	ResourceTypeStreamExclusion
	ResourceTypeIngestionExclusion
	ResourceTypeArchive
	ResourceTypeKey
	ResourceTypeIndexRateAlert
	ResourceTypeMember
	ResourceTypeChildOrg
)

type TerraformInfo struct {
	name    string
	orgType OrgType
	schema  *schema.Resource
}

var dataSourceInfoMap map[DataSourceType]TerraformInfo
var resourceInfoMap map[ResourceType]TerraformInfo

func buildSchemaMap(a []TerraformInfo) map[string]*schema.Resource {
	m := make(map[string]*schema.Resource)

	for _, e := range a {
		m[e.name] = e.schema
	}

	return m
}

func collectInfo(x interface{}) []TerraformInfo {
	a := []TerraformInfo{}

	switch m := x.(type) {
	case map[DataSourceType]TerraformInfo:
		for _, v := range m {
			a = append(a, v)
		}
	case map[ResourceType]TerraformInfo:
		for _, v := range m {
			a = append(a, v)
		}
	}

	return a
}

func collectDataSourceInfo(m map[DataSourceType]TerraformInfo) []TerraformInfo {
	a := []TerraformInfo{}

	for _, v := range m {
		a = append(a, v)
	}

	return a
}
