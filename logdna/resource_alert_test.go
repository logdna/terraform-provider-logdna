package logdna

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAlert_expectServiceKeyError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigServiceKeyError(),
				ExpectError: regexp.MustCompile("The argument \"servicekey\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_expectNameError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertNameError(),
				ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_expectTriggerLimitError(t *testing.T) {
	name := "test"
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigTriggerLimitError(name),
				ExpectError: regexp.MustCompile("Error: \"email_channel.0.triggerlimit\" must be between 1 and 100,000 inclusive, got: 0"),
			},
		},
	})
}

func TestAlert_expectMissingEmails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigMissingEmails(),
				ExpectError: regexp.MustCompile("The argument \"emails\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_expectMissingKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigMissingKey(),
				ExpectError: regexp.MustCompile("The argument \"key\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_expectMissingURL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigMissingURL(),
				ExpectError: regexp.MustCompile("The argument \"url\" is required, but no definition was found."),
			},
		},
	})
}
func TestAlertBasic(t *testing.T) {
	name := "test"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAlertConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.%", "7"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "0"),
				),
			},
		},
	})
}

func TestAlertBulkChannels(t *testing.T) {
	name := "test"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAlertConfigBulkChannels(name),
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.immediate", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.key", "w4g0ushdfalskdj"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.operator", ""),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "0"),
				),
			},
		},
	})
}

func testAlertConfigServiceKeyError() string {
	return `provider "logdna" {
	  }

	resource "logdna_alert" "new" {
		name = "test"
	}`
}

func testAlertNameError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_alert" "new" {
	}`, servicekey)
}

func testAlertConfigTriggerLimitError(name string) string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	resource "logdna_alert" "new" {
		name = "%s"
		email_channel {
			emails          = ["test@logdna.com"]
			immediate       = "false"
			operator        = "presence"
			triggerlimit    = 0
			triggerinterval = "15m"
			terminal        = "true"
			timezone        = "Pacific/Samoa"
		}
	}`, servicekey, name)
}

func testAlertConfigMissingEmails() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_alert" "new" {
		name     = "test"
		email_channel {
			immediate       = "false"
			operator        = "absence"
			terminal        = "true"
			triggerinterval = "15m"
			triggerlimit    = 15
			timezone        = "Pacific/Samoa"
		}
	}`, servicekey)
}

func testAlertConfigMissingKey() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_alert" "new" {
		name     = "test"
		pagerduty_channel {
			immediate       = "false"
			terminal        = "true"
			triggerinterval = "15m"
			triggerlimit    = 15
		}
	}`, servicekey)
}

func testAlertConfigMissingURL() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_alert" "new" {
		name     = "test"
		webhook_channel {
			triggerlimit    = 15
		}
	}`, servicekey)
}

func testAlertConfigMissingTriggerLimit() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_alert" "new" {
		name     = "test"
		webhook_channel {
			url = "https://yourwebhook/endpoint"
		}
	}`, servicekey)
}

func testAlertConfigBasic(name string) string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	resource "logdna_alert" "new" {
		name = "%s"
		email_channel {
			emails          = ["test@logdna.com"]
			immediate       = "false"
			operator        = "presence"
			triggerlimit    = 15
			triggerinterval = "15m"
			terminal        = "true"
			timezone        = "Pacific/Samoa"
		}
	}`, servicekey, name)
}

func testAlertConfigBulkChannels(name string) string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	resource "logdna_alert" "new" {
		name = "%s"
		email_channel {
			emails          = ["test@logdna.com"]
			immediate       = "false"
			operator        = "absence"
			terminal        = "true"
			timezone        = "Pacific/Samoa"
			triggerinterval = "15m"
			triggerlimit    = 15
		}
		pagerduty_channel {
			immediate       = "true"
			key             = "your pagerduty key goes here"
			terminal        = "true"
			triggerinterval = "15m"
			triggerlimit    = 15
		}
	}`, servicekey, name)
}

func testAlertExists(n string) resource.TestCheckFunc {
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
