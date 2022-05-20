package logdna

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestKey_ErrorResourceTypeUndefined(t *testing.T) {
	args := map[string]string{}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("key", "new", []string{serviceKey, apiHostUrl}, args, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("The argument \"type\" is required, but no definition was found."),
			},
		},
	})
}

func TestKey_ErrorResourceTypeInvalid(t *testing.T) {
	args := map[string]string{
		"type": `"incorrect"`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("key", "new", []string{serviceKey, apiHostUrl}, args, nilOpt, nilLst),
				ExpectError: regexp.MustCompile(`Error: expected type to be one of \[ingestion service\], got incorrect`),
			},
		},
	})
}

func TestKey_Basic(t *testing.T) {
	serviceArgs := map[string]string{
		"type": `"service"`,
		"name": `"my first name"`,
	}

	serviceUpdateArgs := map[string]string{
		"type": `"service"`,
		"name": `"my new name"`,
	}

	ingestionArgs := map[string]string{
		"type": `"ingestion"`,
		"name": `"my first name"`,
	}

	ingestionUpdateArgs := map[string]string{
		"type": `"ingestion"`,
		"name": `"my new name"`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// NOTE It tests a service key create operation
				Config: fmtTestConfigResource("key", "new-service-key", []string{serviceKey, apiHostUrl}, serviceArgs, nilOpt, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("key", "new-service-key"),
					resource.TestCheckResourceAttr("logdna_key.new-service-key", "type", strings.Replace(serviceArgs["type"], "\"", "", 2)),
					resource.TestCheckResourceAttr("logdna_key.new-service-key", "name", strings.Replace(serviceArgs["name"], "\"", "", 2)),
					resource.TestCheckResourceAttrSet("logdna_key.new-service-key", "id"),
					resource.TestCheckResourceAttrSet("logdna_key.new-service-key", "key"),
					resource.TestCheckResourceAttrSet("logdna_key.new-service-key", "created"),
				),
			},
			{
				// NOTE It tests a service key update operation
				Config: fmtTestConfigResource("key", "new-service-key", []string{serviceKey, apiHostUrl}, serviceUpdateArgs, nilOpt, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("key", "new-service-key"),
					resource.TestCheckResourceAttr("logdna_key.new-service-key", "type", strings.Replace(serviceUpdateArgs["type"], "\"", "", 2)),
					resource.TestCheckResourceAttr("logdna_key.new-service-key", "name", strings.Replace(serviceUpdateArgs["name"], "\"", "", 2)),
					resource.TestCheckResourceAttrSet("logdna_key.new-service-key", "id"),
					resource.TestCheckResourceAttrSet("logdna_key.new-service-key", "key"),
					resource.TestCheckResourceAttrSet("logdna_key.new-service-key", "created"),
				),
			},
			{
				ResourceName:      "logdna_key.new-service-key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// NOTE It tests an ingestion key create operation
				Config: fmtTestConfigResource("key", "new-ingestion-key", []string{serviceKey, apiHostUrl}, ingestionArgs, nilOpt, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("key", "new-ingestion-key"),
					resource.TestCheckResourceAttr("logdna_key.new-ingestion-key", "type", strings.Replace(ingestionArgs["type"], "\"", "", 2)),
					resource.TestCheckResourceAttr("logdna_key.new-ingestion-key", "name", strings.Replace(ingestionArgs["name"], "\"", "", 2)),
					resource.TestCheckResourceAttrSet("logdna_key.new-ingestion-key", "id"),
					resource.TestCheckResourceAttrSet("logdna_key.new-ingestion-key", "key"),
					resource.TestCheckResourceAttrSet("logdna_key.new-ingestion-key", "created"),
				),
			},
			{
				// NOTE It tests an ingestion key update operation
				Config: fmtTestConfigResource("key", "new-ingestion-key", []string{serviceKey, apiHostUrl}, ingestionUpdateArgs, nilOpt, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("key", "new-ingestion-key"),
					resource.TestCheckResourceAttr("logdna_key.new-ingestion-key", "type", strings.Replace(ingestionUpdateArgs["type"], "\"", "", 2)),
					resource.TestCheckResourceAttr("logdna_key.new-ingestion-key", "name", strings.Replace(ingestionUpdateArgs["name"], "\"", "", 2)),
					resource.TestCheckResourceAttrSet("logdna_key.new-ingestion-key", "id"),
					resource.TestCheckResourceAttrSet("logdna_key.new-ingestion-key", "key"),
					resource.TestCheckResourceAttrSet("logdna_key.new-ingestion-key", "created"),
				),
			},
			{
				ResourceName:      "logdna_key.new-ingestion-key",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
