package logdna

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMember_ErrorRoleEmpty(t *testing.T) {
	args := map[string]string{
		"email": `"user@example.org"`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("member", "new", []string{serviceKey, apiHostUrl}, args, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("The argument \"role\" is required, but no definition was found."),
			},
		},
	})
}

func TestMember_Basic(t *testing.T) {
	memberArgs := map[string]string{
		"email": `"member@example.org"`,
		"role":  `"member"`,
	}

	adminArgs := map[string]string{
		"email": `"admin@example.org"`,
		"role":  `"admin"`,
	}
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmtTestConfigResource("member", "member", []string{serviceKey, apiHostUrl}, memberArgs, nilOpt, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("member", "member"),
					resource.TestCheckResourceAttr("logdna_member.member", "email", strings.Replace(memberArgs["email"], "\"", "", 2)),
					resource.TestCheckResourceAttr("logdna_member.member", "role", strings.Replace(memberArgs["role"], "\"", "", 2)),
				),
			},
			{
				ResourceName:      "logdna_member.member",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmtTestConfigResource("member", "admin", []string{serviceKey, apiHostUrl}, adminArgs, nilOpt, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("member", "admin"),
					resource.TestCheckResourceAttr("logdna_member.admin", "email", strings.Replace(adminArgs["email"], "\"", "", 2)),
					resource.TestCheckResourceAttr("logdna_member.admin", "role", strings.Replace(adminArgs["role"], "\"", "", 2)),
				),
			},
			{
				ResourceName:      "logdna_member.admin",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
