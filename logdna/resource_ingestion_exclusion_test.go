package logdna

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestIngestionExclusion_expectInvalidURLError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIngestionExclusion(`
					title = "test-title"
					query = "foo"
				`, "http://api.logdna.co"),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"http://api.logdna.co/v1/config/ingestion/exclusions\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestIngestionExclusion_expectInvalidError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIngestionExclusion(`
					title = "test-title"
				`, apiHostUrl),
				ExpectError: regexp.MustCompile("one of `apps,hosts,query` must be specified"),
			},
			{
				Config: testIngestionExclusion(`
					title = "test-title"
					apps = []
				`, apiHostUrl),
				ExpectError: regexp.MustCompile("requires 1 item minimum, but config has only 0 declared"),
			},
		},
	})
}

func TestIngestionExclusion_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testIngestionExclusion(`
					title = "test-title"
					active = true
					apps = [
						"app-1",
						"app-2"
					]
					hosts = [
						"host-1",
						"host-2"
					]
					query = "foo bar"
				`, apiHostUrl),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("ingestion_exclusion", "new"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "title", "test-title"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "active", "true"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "apps.0", "app-1"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "apps.1", "app-2"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "hosts.0", "host-1"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "hosts.1", "host-2"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "query", "foo bar"),
				),
			},
			{
				Config: testIngestionExclusion(`
					title = "test-title-update"
					active = false
					apps = [
						"app-1",
						"app-2"
					]
					hosts = [
						"host-1",
						"host-2"
					]
					query = "foo bar"
				`, apiHostUrl),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("ingestion_exclusion", "new"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "title", "test-title-update"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "active", "false"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "query", "foo bar"),
				),
			},
			{
				Config: testIngestionExclusion(`
					title = "test-title-update"
					active = false
					apps = [
						"app-1",
						"app-2"
					]
					query = "foo bar"
				`, apiHostUrl),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("ingestion_exclusion", "new"),
					resource.TestCheckResourceAttr("logdna_ingestion_exclusion.new", "hosts.#", "0"),
				),
			},
			{
				ResourceName:      "logdna_ingestion_exclusion.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testIngestionExclusion(fields string, url string) string {
	uc := ""
	if url != "" {
		uc = fmt.Sprintf(`url = "%s"`, url)
	}

	return fmt.Sprintf(`
		provider "logdna" {
			servicekey = "%s"
			%s
		}
		resource "logdna_ingestion_exclusion" "new" {
			%s
		}
	`, serviceKey, uc, fields)
}
