package logdna

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var s3Bucket = os.Getenv("S3_BUCKET")
var gcsBucket = os.Getenv("GCS_BUCKET")
var gcsProjectid = os.Getenv("GCS_PROJECTID")

func testArchivePreCheck(t *testing.T) {
	if s3Bucket == "" {
		t.Fatal("'S3_BUCKET' environment variable must be set for acceptance tests")
	}
	if gcsBucket == "" {
		t.Fatal("'GCS_BUCKET' environment variable must be set for acceptance tests")
	}
	if gcsProjectid == "" {
		t.Fatal("'GCS_PROJECTID' environment variable must be set for acceptance tests")
	}
}

func TestArchiveConfig_checkEnvServiceKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testArchivePreCheck(t) },
	})
}

func TestArchiveConfig_expectInvalidURLError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testArchivePreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testArchiveConfig(fmt.Sprintf(`
					integration = "s3"
					s3_config {
						bucket = "%s"
					}
				`, s3Bucket), "http://api.logdna.co"),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"http://api.logdna.co/v1/config/archiving\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestArchiveConfig_expectInvalidIntegrationError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testArchiveConfig(`
					integration = "invalid"
				`, ""),
				ExpectError: regexp.MustCompile(`"integration" must be one of \[ibm s3 azblob gcs dos swift\]`),
			},
		},
	})
}

func TestArchiveConfig_expectMissingFieldError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testArchiveConfig(`
					integration = "s3"
					s3_config {
					}
				`, ""),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestArchiveConfig_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testArchivePreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testArchiveConfig(fmt.Sprintf(`
					integration = "s3"
					s3_config {
						bucket = "%s"
					}
				`, s3Bucket), ""),
				Check: resource.ComposeTestCheckFunc(
					testArchiveConfigExists("logdna_archive.archive"),
					resource.TestCheckResourceAttr("logdna_archive.archive", "integration", "s3"),
					resource.TestCheckResourceAttr("logdna_archive.archive", "s3_config.0.bucket", s3Bucket),
				),
			},
			{
				Config: testArchiveConfig(fmt.Sprintf(`
					integration = "gcs"
					gcs_config {
						bucket = "%s"
						projectid = "%s"
					}
				`, gcsBucket, gcsProjectid), ""),
				Check: resource.ComposeTestCheckFunc(
					testArchiveConfigExists("logdna_archive.archive"),
					resource.TestCheckResourceAttr("logdna_archive.archive", "integration", "gcs"),
					resource.TestCheckResourceAttr("logdna_archive.archive", "gcs_config.0.bucket", gcsBucket),
					resource.TestCheckResourceAttr("logdna_archive.archive", "gcs_config.0.projectid", gcsProjectid),
				),
			},
			{
				ResourceName:      "logdna_archive.archive",
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

func testArchiveConfig(fields string, url string) string {
	uc := ""
	if url != "" {
		uc = fmt.Sprintf(`url = "%s"`, url)
	}

	return fmt.Sprintf(`
		provider "logdna" {
			servicekey = "%s"
			%s
		}
		
		resource "logdna_archive" "archive" {
			%s
		}
	`, serviceKey, uc, fields)
}
