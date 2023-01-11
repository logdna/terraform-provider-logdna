package logdna

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var escapeChar = regexp.MustCompile(`\"|\[|\]|,`)

func TestIndexRateAlert_ErrorProviderUrl(t *testing.T) {
	pcArgs := []string{serviceKey, "https://api.logdna.co"}
	iraArgs := map[string]string{
		"max_lines":       `3`,
		"max_z_score":     `3`,
		"threshold_alert": `"separate"`,
		"frequency":       `"hourly"`,
		"enabled":         `false`,
	}

	chArgs := map[string]map[string]string{
		"channels": {
			"email": `["test@logdna.com", "test2@logdna.com"]`,
		},
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("index_rate_alert", "test_config", pcArgs, iraArgs, chArgs, nilLst),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Put \"https://api.logdna.co/v1/config/index-rate\": dial tcp: lookup api.logdna.co.+?"),
			},
		},
	})
}

func TestIndexRateAlert_ErrorOrgType(t *testing.T) {
	pcArgs := []string{enterpriseServiceKey, apiHostUrl, "enterprise"}
	iraArgs := map[string]string{
		"max_lines":       `3`,
		"max_z_score":     `3`,
		"threshold_alert": `"separate"`,
		"frequency":       `"hourly"`,
		"enabled":         `false`,
	}

	chArgs := map[string]map[string]string{
		"channels": {
			"email":     `["test@logdna.com", "test2@logdna.com"]`,
			"slack":     `["https://hooks.slack.com/KEY"]`,
			"pagerduty": `["ndt3k75rsw520d8t55dv35decdyt3mkcb3r"]`,
		},
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("index_rate_alert", "new", pcArgs, iraArgs, chArgs, nilLst),
				ExpectError: regexp.MustCompile("Error: Only regular organizations can instantiate a \"logdna_index_rate_alert\" resource"),
			},
		},
	})
}

func TestIndexRateAlert_ErrorResourceThresholdAlertInvalid(t *testing.T) {
	iraArgs := map[string]string{
		"max_lines":       `3`,
		"max_z_score":     `3`,
		"threshold_alert": `"invalid"`,
		"frequency":       `"hourly"`,
		"enabled":         `false`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("index_rate_alert", "test_config", globalPcArgs, iraArgs, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: expected threshold_alert to be one of .+?"),
			},
		},
	})
}

func TestIndexRateAlert_ErrorResourceThresholdAlertMissed(t *testing.T) {
	iraArgs := map[string]string{
		"max_lines":   `3`,
		"max_z_score": `3`,
		"frequency":   `"hourly"`,
		"enabled":     `false`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("index_rate_alert", "test_config", globalPcArgs, iraArgs, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("The argument \"threshold_alert\" is required, but no definition was found."),
			},
		},
	})
}

func TestIndexRateAlert_ErrorResourceFrequencyInvalid(t *testing.T) {
	iraArgs := map[string]string{
		"max_lines":       `3`,
		"max_z_score":     `3`,
		"threshold_alert": `"separate"`,
		"frequency":       `"ivalid"`,
		"enabled":         `false`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("index_rate_alert", "test_config", globalPcArgs, iraArgs, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: expected frequency to be one of .+?"),
			},
		},
	})
}

func TestIndexRateAlert_ErrorResourceFrequencyMissed(t *testing.T) {
	iraArgs := map[string]string{
		"max_lines":       `3`,
		"max_z_score":     `3`,
		"threshold_alert": `"separate"`,
		"enabled":         `false`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("index_rate_alert", "test_config", globalPcArgs, iraArgs, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("The argument \"frequency\" is required, but no definition was found."),
			},
		},
	})
}

func TestIndexRateAlert_ErrorResourceEnable(t *testing.T) {
	iraArgs := map[string]string{
		"max_lines":       `3`,
		"max_z_score":     `3`,
		"threshold_alert": `"separate"`,
		"frequency":       `"hourly"`,
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("index_rate_alert", "test_config", globalPcArgs, iraArgs, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("The argument \"enabled\" is required, but no definition was found."),
			},
		},
	})
}

func TestIndexRateAlert_ErrorChannels(t *testing.T) {
	iraArgs := map[string]string{
		"max_lines":       `3`,
		"max_z_score":     `3`,
		"threshold_alert": `"separate"`,
		"frequency":       `"hourly"`,
		"enabled":         `false`,
	}

	chArgs := map[string]map[string]string{
		"channels": {
			"email": `["test@logdna.com"]`,
		},
		"channels1": {
			"email": `["test2@logdna.com"]`,
		},
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("index_rate_alert", "test_config", globalPcArgs, iraArgs, chArgs, nilLst),
				ExpectError: regexp.MustCompile("Index rate alert resource supports only one channels object"),
			},
		},
	})
}

func TestIndexRateAlert_ErrorInvalidTokenWebhookBodyTemplate(t *testing.T) {
	iraArgs := map[string]string{
		"max_lines":       `3`,
		"max_z_score":     `3`,
		"threshold_alert": `"separate"`,
		"frequency":       `"hourly"`,
		"enabled":         `false`,
	}

	chArgs := map[string]map[string]string{
		"channels": {
			"email": `["test@logdna.com"]`,
		},
		"webhook_channel": {
			"url":    `"https://something.com"`,
			"method": `"POST"`,
			"headers": `{
        field2 = "value2"
      }`,
			"bodytemplate": `jsonencode({
        something = "{{maxLines}}"
      })`,
		},
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("index_rate_alert", "test_config", globalPcArgs, iraArgs, chArgs, nilLst),
				ExpectError: regexp.MustCompile("Invalid bodyTemplate: {{maxLines}} is not a valid token"),
			},
		},
	})
}

func TestIndexRateAlert_Basic(t *testing.T) {
	iraArgs := map[string]string{
		"max_lines":       `3`,
		"max_z_score":     `3`,
		"threshold_alert": `"separate"`,
		"frequency":       `"hourly"`,
		"enabled":         `false`,
	}

	chArgs := map[string]map[string]string{
		"channels": {
			"email":     `["test@logdna.com", "test2@logdna.com"]`,
			"slack":     `["https://hooks.slack.com/KEY"]`,
			"pagerduty": `["ndt3k75rsw520d8t55dv35decdyt3mkcb3r"]`,
		},
		"webhook_channel": {
		  "url" : `"https://something.com"`,
		  "method": `"POST"`,
		  "headers" : `{
			field2 = "value2"
		  }`,
		  "bodytemplate" :`jsonencode({
			something = "something"
		  })`,
		},
	}

	iraUpdArgs := map[string]string{
		"max_lines":       `5`,
		"max_z_score":     `5`,
		"threshold_alert": `"separate"`,
		"frequency":       `"hourly"`,
		"enabled":         `true`,
	}

	chUpdArgs := map[string]map[string]string{
		"channels": {
			"email":     `["test_updated@logdna.com"]`,
			"slack":     `["https://hooks.slack.com/UPDATED_KEY", "https://hooks.slack.com/KEY_2"]`,
			"pagerduty": `["new3k75rsw520d8t55dv35decdyt3mkcnew"]`,
		},
		"webhook_channel": {
			"url" : `"https://something.com"`,
			"method": `"PUT"`,
			"headers" : `{
			  field2 = "value2"
			}`,
			"bodytemplate" :`jsonencode({
			  something = "!something"
			})`,
		  },
	}

	createdEmails := strings.Split(
		escapeChar.ReplaceAllString(chArgs["channels"]["email"], ""),
		" ",
	)

	createdSlack := strings.Split(
		escapeChar.ReplaceAllString(chArgs["channels"]["slack"], ""),
		" ",
	)

	createdPagerduty := strings.Split(
		escapeChar.ReplaceAllString(chArgs["channels"]["pagerduty"], ""),
		" ",
	)

	updatedEmails := strings.Split(
		escapeChar.ReplaceAllString(chUpdArgs["channels"]["email"], ""),
		" ",
	)

	updatedSlack := strings.Split(
		escapeChar.ReplaceAllString(chUpdArgs["channels"]["slack"], ""),
		" ",
	)

	updatedPagerduty := strings.Split(
		escapeChar.ReplaceAllString(chUpdArgs["channels"]["pagerduty"], ""),
		" ",
	)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// NOTE It tests a index rate alert create operation
				Config: fmtTestConfigResource("index_rate_alert", "test_config", globalPcArgs, iraArgs, chArgs, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("index_rate_alert", "test_config"),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"max_lines",
						escapeChar.ReplaceAllString(iraArgs["max_lines"], ""),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"max_z_score",
						escapeChar.ReplaceAllString(iraArgs["max_z_score"], ""),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"threshold_alert",
						escapeChar.ReplaceAllString(iraArgs["threshold_alert"], ""),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"frequency",
						escapeChar.ReplaceAllString(iraArgs["frequency"], ""),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.email.#",
						strconv.Itoa(len(createdEmails)),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.email.0",
						createdEmails[0],
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.email.1",
						createdEmails[1],
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.slack.#",
						strconv.Itoa(len(createdSlack)),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.slack.0",
						createdSlack[0],
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.pagerduty.#",
						strconv.Itoa(len(createdPagerduty)),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.pagerduty.0",
						createdPagerduty[0],
					),
				),
			},
			{
				// NOTE It tests a index rate alert config update operation
				Config: fmtTestConfigResource("index_rate_alert", "test_config", globalPcArgs, iraUpdArgs, chUpdArgs, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("index_rate_alert", "test_config"),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"max_lines",
						escapeChar.ReplaceAllString(iraUpdArgs["max_lines"], ""),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"max_z_score",
						escapeChar.ReplaceAllString(iraUpdArgs["max_z_score"], ""),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"threshold_alert",
						escapeChar.ReplaceAllString(iraUpdArgs["threshold_alert"], ""),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"frequency",
						escapeChar.ReplaceAllString(iraUpdArgs["frequency"], ""),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.email.#",
						strconv.Itoa(len(updatedEmails)),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.email.0",
						updatedEmails[0],
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.slack.#",
						strconv.Itoa(len(updatedSlack)),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.slack.0",
						updatedSlack[0],
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.slack.1",
						updatedSlack[1],
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.pagerduty.#",
						strconv.Itoa(len(updatedPagerduty)),
					),
					resource.TestCheckResourceAttr(
						"logdna_index_rate_alert.test_config",
						"channels.0.pagerduty.0",
						updatedPagerduty[0],
					),
				),
			},
			{
				ResourceName:      "logdna_index_rate_alert.test_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
