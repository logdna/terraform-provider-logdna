package logdna

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCategory_ErrorProviderUrl(t *testing.T) {
	pcArgs := []string{serviceKey, "https://api.logdna.co"}
	catArgs := map[string]string{
		"name": `"test-category"`,
		"type": `"views"`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("category", "new", pcArgs, catArgs, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"https://api.logdna.co/v1/config/categories/views\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestCategory_ErrorResourceName(t *testing.T) {
	catArgs := map[string]string{
		"type": `"views"`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("category", "new", globalPcArgs, catArgs, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
			},
		},
	})
}

func TestCategory_ErrorResourceType(t *testing.T) {
	catArgs := map[string]string{
		"name": `"test-category"`,
		"type": `"incorrect"`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("category", "new", globalPcArgs, catArgs, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: POST .+?, status 400 NOT OK!"),
			},
		},
	})
}

func TestCategory_Basic(t *testing.T) {
	catInsArgs := map[string]string{
		"name": `"test-category"`,
		"type": `"views"`,
	}

	catUpdArgs := map[string]string{
		"name": `"test-category-updated"`,
		"type": `"views"`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// NOTE It tests a category create operation
				Config: fmtTestConfigResource("category", "new-category", globalPcArgs, catInsArgs, nilOpt, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("category", "new-category"),
					resource.TestCheckResourceAttr("logdna_category.new-category", "name", strings.Replace(catInsArgs["name"], "\"", "", 2)),
				),
			},
			{
				// NOTE It tests a category update operation
				Config: fmtTestConfigResource("category", "new-category", globalPcArgs, catUpdArgs, nilOpt, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("category", "new-category"),
					resource.TestCheckResourceAttr("logdna_category.new-category", "name", strings.Replace(catUpdArgs["name"], "\"", "", 2)),
				),
			},
			{
				ResourceName:      "logdna_category.new-category",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
