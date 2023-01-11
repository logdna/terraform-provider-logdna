package logdna

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestChildOrg_ErrorOrgType(t *testing.T) {
	pcArgs := []string{serviceKey, apiHostUrl}
	orgArgs := map[string]string{
		"servicekey": fmt.Sprintf(`"%s"`, serviceKey),
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("child_organization", "new", pcArgs, orgArgs, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: Only enterprise organizations can instantiate a \"logdna_child_organization\" resource"),
			},
		},
	})
}

func TestChildOrg_Basic(t *testing.T) {
	orgArgs := map[string]string{
		"servicekey": fmt.Sprintf(`"%s"`, serviceKey),
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// NOTE: This tests detach childOrg operation
				Config: fmtTestConfigResource("child_organization", "delete", []string{enterpriseServiceKey, apiHostUrl, "enterprise"}, orgArgs, nilOpt, nilLst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("logdna_child_organization.delete", "servicekey", strings.Replace(orgArgs["servicekey"], "\"", "", 2)),
				),
			},
			{
				ResourceName:      "logdna_child_organization.delete",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
