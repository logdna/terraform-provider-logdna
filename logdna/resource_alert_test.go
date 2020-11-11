package logdna

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAlert_expectInvalidURLError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertInvalidURL(),
				ExpectError: regexp.MustCompile("Error: Error with alert: Post \"http://api.logdna.co/v1/config/alert\": dial tcp: lookup api.logdna.co on 127.0.0.11:53: no such host"),
			},
		},
	})
}
func TestAlert_expectInvalidJSONError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigMultipleChannelsInvalidJSON(),
				ExpectError: regexp.MustCompile("bodytemplate json configuration is invalid"),
			},
		},
	})
}

func TestAlert_expectTriggerIntervalError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigTriggerIntervalError(),
				ExpectError: regexp.MustCompile(`\"channels\[0]\.triggerinterval\" must be one of \[15m, 30m, 1h, 6h, 12h, 24h]`),
			},
		},
	})
}

func TestAlert_expectImmediateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigImmediateError(),
				ExpectError: regexp.MustCompile(`\"channels\[0\].immediate\" must be \[false\]. \"channels\[0]\.immediate\" must be a boolean`),
			},
		},
	})
}

func TestAlert_expectURLError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigURLError(),
				ExpectError: regexp.MustCompile(`\"channels\[0\].url\" must be a valid uri`),
			},
		},
	})
}

func TestAlert_expectMethodError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigMethodError(),
				ExpectError: regexp.MustCompile(`\"channels\[0\].method\" must be one of \[post, put, patch, get, delete\]`),
			},
		},
	})
}

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

func TestAlert_expectEmailTriggerLimitError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigEmailTriggerLimitError(),
				ExpectError: regexp.MustCompile("Error: \"email_channel.0.triggerlimit\" must be between 1 and 100,000 inclusive, got: 0"),
			},
		},
	})
}

func TestAlert_expectPagerDutyTriggerLimitError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigPagerDutyTriggerLimitError(),
				ExpectError: regexp.MustCompile("Error: \"pagerduty_channel.0.triggerlimit\" must be between 1 and 100,000 inclusive, got: 0"),
			},
		},
	})
}

func TestAlert_expectWebhookTriggerLimitError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAlertConfigWebhookTriggerLimitError(),
				ExpectError: regexp.MustCompile("Error: \"webhook_channel.0.triggerlimit\" must be between 1 and 100,000 inclusive, got: 0"),
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

func TestAlertBasicUpdate(t *testing.T) {
	name := "test"
	name2 := "test2"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAlertConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", name),
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
			{
				Config: testAlertConfigBasic(name2),
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", name2),
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

func TestAlertJSONUpdateError(t *testing.T) {
	name := "test"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAlertConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
				),
			},
			{
				Config:      testAlertConfigMultipleChannelsInvalidJSON(),
				ExpectError: regexp.MustCompile("bodytemplate json configuration is invalid"),
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
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.key", "Your PagerDuty API key goes here"),
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

func testAlertInvalidURL() string {
	return fmt.Sprintf(`provider "logdna" {
		url = "http://api.logdna.co"
		servicekey = "%s"
	  }

	  resource "logdna_alert" "new" {
		name = "test"
	  }`, servicekey)
}

func testAlertConfigMultipleChannelsInvalidJSON() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	  resource "logdna_alert" "new" {
		name = "test"
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
		  immediate       = "false"
		  key             = "Your PagerDuty API key goes here"
		  terminal        = "true"
		  triggerinterval = "15m"
		  triggerlimit    = 15
		}

		webhook_channel {
		  headers = {
			hello = "test3"
			test  = "test2"
		  }
		  bodytemplate = "{\"test\": }"
		  immediate       = "false"
		  method          = "post"
		  url             = "https://yourwebhook/endpoint"
		  terminal        = "true"
		  triggerinterval = "15m"
		  triggerlimit    = 15
		}
	  }`, servicekey)
}

func testAlertConfigTriggerIntervalError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
      }

      resource "logdna_alert" "new" {
        name = "test"
        email_channel {
          emails          = ["test@logdna.com"]
          immediate       = "false"
          operator        = "absence"
          terminal        = "true"
          timezone        = "Pacific/Samoa"
          triggerinterval = "17m"
          triggerlimit    = 15
        }
      }`, servicekey)
}

func testAlertConfigInvalidPagerDutyTriggerLimitError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

      resource "logdna_alert" "new" {
        name = "test"
		pagerduty_channel {
			immediate       = "false"
			key             = "Your PagerDuty API key goes here"
			terminal        = "true"
			triggerinterval = "15m"
			triggerlimit    = 0
		}
      }`, servicekey)
}

func testAlertConfigImmediateError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
      }

      resource "logdna_alert" "new" {
        name = "test"
        email_channel {
          emails          = ["test@logdna.com"]
          immediate       = "no"
          operator        = "absence"
          terminal        = "true"
          timezone        = "Pacific/Samoa"
          triggerinterval = "15m"
          triggerlimit    = 15
        }
      }`, servicekey)
}

func testAlertConfigURLError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	  resource "logdna_alert" "new" {
		name = "test"
		webhook_channel {
		  headers = {
			hello = "test3"
			test  = "test2"
		  }
		  immediate       = "false"
		  method          = "post"
		  url             = "this is not a valid url"
		  terminal        = "true"
		  triggerinterval = "15m"
		  triggerlimit    = 15
		}
	  }`, servicekey)
}

func testAlertConfigMethodError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	  resource "logdna_alert" "new" {
		name = "test"
		webhook_channel {
		  headers = {
			hello = "test3"
			test  = "test2"
		  }
		  immediate       = "false"
		  method          = "false"
		  url             = "http://yourwebhook/test"
		  terminal        = "true"
		  triggerinterval = "15m"
		  triggerlimit    = 15
		}
	  }`, servicekey)
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

func testAlertConfigEmailTriggerLimitError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	resource "logdna_alert" "new" {
		name = "test"
		email_channel {
			emails          = ["test@logdna.com"]
			immediate       = "false"
			operator        = "presence"
			triggerlimit    = 0
			triggerinterval = "15m"
			terminal        = "true"
			timezone        = "Pacific/Samoa"
		}
	}`, servicekey)
}

func testAlertConfigPagerDutyTriggerLimitError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	resource "logdna_alert" "new" {
		name = "test"
		pagerduty_channel {
			immediate       = "false"
			key             = "Your PagerDuty API key goes here"
			terminal        = "true"
			triggerinterval = "15m"
			triggerlimit    = 0
		}
	}`, servicekey)
}

func testAlertConfigWebhookTriggerLimitError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	resource "logdna_alert" "new" {
		name = "test"
		webhook_channel {
			headers = {
			  hello = "test3"
			  test  = "test2"
			}
			immediate       = "false"
			method          = "post"
			url             = "https://yourwebhook/endpoint"
			terminal        = "true"
			triggerinterval = "15m"
			triggerlimit    = 0
		}
	}`, servicekey)
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
			key             = "Your PagerDuty API key goes here"
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
