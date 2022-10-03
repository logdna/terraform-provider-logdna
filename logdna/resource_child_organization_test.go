package logdna

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestChildOrg_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// NOTE: This tests detach childOrg operation
				Config: fmtTestConfigResource("enterprise", "delete", []string{serviceKey, apiHostUrl, "enterprise"}, nil, nilOpt, nilLst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("logdna_enterprise_child_org.delete", "account"),
					resource.TestCheckNoResourceAttr("logdna_enterprise_child_org.delete", "enterprise-servicekey"),
				),
			},
			{
				ResourceName:      "logdna_enterprise_child_org.delete",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
