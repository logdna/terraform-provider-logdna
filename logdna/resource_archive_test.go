package logdna

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const defaultUrl = "https://api.logdna.com"
const stagingUrl = "https://api.use.stage.logdna.net"

var s3Bucket = os.Getenv("S3_BUCKET")
var gcsBucket = os.Getenv("GCS_BUCKET")
var gcsProjectid = os.Getenv("GCS_PROJECTID")
var stagingKey = os.Getenv("SERVICE_KEY_STAGE")
var baseOpts = map[string]string{
	"url": "",
	"s3":  "",
	"gcs": "",
}

const pcFmt = `
provider "logdna" {
	servicekey = "%s"
	%s
}
`
const rsFmt = `
resource "logdna_archive" "config" {
	integration = "%s"
	%s
}
`
const s3Fmt = `s3_config {
		bucket = "%s"
	}`
const gcsFmt = `gcs_config {
		bucket = "%s"
		projectid = "%s"
	}`

func TestArchiveConfig_expectInvalidURLError(t *testing.T) {
	opts := inheritOpts()
	opts["url"] = "http://api.logdna.co"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testArchiveConfigDyn(false, "s3", opts),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"http://api.logdna.co/v1/config/archiving\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestArchiveConfig_expectInvalidIntegrationError(t *testing.T) {
	opts := inheritOpts()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testArchiveConfigDyn(false, "invalid", opts),
				ExpectError: regexp.MustCompile(""),
			},
		},
	})
}

func TestArchiveConfig_expectMissingFieldError(t *testing.T) {
	opts := inheritOpts()
	opts["s3"] = `s3_config {}`

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testArchiveConfigDyn(false, "s3", opts),
				ExpectError: regexp.MustCompile(`Error: Missing required argument`),
			},
		},
	})
}

func TestArchiveConfig_basic(t *testing.T) {
	opts := inheritOpts()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testArchiveConfigDyn(false, "s3", opts),

				Check: resource.ComposeTestCheckFunc(
					testArchiveConfigExists("logdna_archive.config"),
					resource.TestCheckResourceAttr("logdna_archive.config", "integration", "s3"),
					resource.TestCheckResourceAttr("logdna_archive.config", "s3_config.0.bucket", s3Bucket),
				),
			},
			{
				Config: testArchiveConfigDyn(false, "gcs", opts),
				Check: resource.ComposeTestCheckFunc(
					testArchiveConfigExists("logdna_archive.config"),
					resource.TestCheckResourceAttr("logdna_archive.config", "integration", "gcs"),
					resource.TestCheckResourceAttr("logdna_archive.config", "gcs_config.0.bucket", gcsBucket),
					resource.TestCheckResourceAttr("logdna_archive.config", "gcs_config.0.projectid", gcsProjectid),
				),
			},
			{
				ResourceName:      "logdna_archive.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testArchiveConfigExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set")
		}
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		return nil
	}
}

func testArchiveConfigDyn(staging bool, integration string, opts map[string]string) string {
	sk, ul := "", ""

	if staging {
		sk = stagingKey
		ul = fmt.Sprintf("url = %q", stagingUrl)
	} else {
		sk = serviceKey
		ul = fmt.Sprintf("url = %q", defaultUrl)
	}

	if opts["url"] != "" {
		ul = fmt.Sprintf("url = %q", opts["url"])
	}

	pc := fmt.Sprintf(pcFmt, sk, ul)
	rs, cf := "", ""

	switch integration {
	case "invalid":
		{
			rs = fmt.Sprintf(rsFmt, integration, "")
		}
	case "s3":
		{
			s3 := fmt.Sprintf(s3Fmt, s3Bucket)
			if opts["s3"] != "" {
				s3 = opts["s3"]
			}
			rs = fmt.Sprintf(rsFmt, integration, s3)
		}
	case "gcs":
		{
			gcs := fmt.Sprintf(gcsFmt, gcsBucket, gcsProjectid)
			if opts["gcs"] != "" {
				gcs = opts["gcs"]
			}
			rs = fmt.Sprintf(rsFmt, integration, gcs)
		}
	}

	cf = fmt.Sprintf("%s%s", pc, rs)
	//fmt.Println(cf)
	return cf
}

func inheritOpts() map[string]string {
	opts := make(map[string]string)
	for k, v := range baseOpts {
		opts[k] = v
	}
	return opts
}
