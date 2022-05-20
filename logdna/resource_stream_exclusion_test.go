package logdna

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestStreamExclusion_expectInvalidURLError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStreamExclusion(`
					title = "test-title"
					query = "query-foo AND query-bar"			
				`, "http://api.logdna.co"),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"http://api.logdna.co/v1/config/stream/exclusions\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestStreamExclusion_expectInvalidError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStreamExclusion(`
					title = "test-title"
				`, ""),
				ExpectError: regexp.MustCompile("one of `apps,hosts,query` must be specified"),
			},
			{
				Config: testStreamExclusion(`
					title = "test-title"
					apps = []
				`, apiHostUrl),
				ExpectError: regexp.MustCompile("requires 1 item minimum, but config has only 0 declared"),
			},
		},
	})
}

func TestStreamExclusion_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStreamExclusion(`
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
					query = "query-foo AND query-bar"
				`, apiHostUrl),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("stream_exclusion", "new"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "title", "test-title"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "active", "true"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "apps.0", "app-1"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "apps.1", "app-2"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "hosts.0", "host-1"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "hosts.1", "host-2"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "query", "query-foo AND query-bar"),
				),
			},
			{
				Config: testStreamExclusion(`
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
					query = "query-foo AND query-bar"
				`, apiHostUrl),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("stream_exclusion", "new"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "title", "test-title-update"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "active", "false"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "query", "query-foo AND query-bar"),
				),
			},
			{
				Config: testStreamExclusion(`
					title = "test-title-update"
					active = false
					apps = [
						"app-1",
						"app-2"
					]
					query = "query-foo AND query-bar"
				`, apiHostUrl),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("stream_exclusion", "new"),
					resource.TestCheckResourceAttr("logdna_stream_exclusion.new", "hosts.#", "0"),
				),
			},
			{
				ResourceName:      "logdna_stream_exclusion.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testStreamExclusion(fields string, url string) string {
	uc := ""
	if url != "" {
		uc = fmt.Sprintf(`url = "%s"`, url)
	}

	return fmt.Sprintf(`
		provider "logdna" {
			servicekey = "%s"
			%s
		}
		resource "logdna_stream_exclusion" "new" {
			%s
		}
	`, serviceKey, uc, fields)
}
