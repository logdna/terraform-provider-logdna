package logdna

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestView_expectInvalidURLError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewInvalidURL(),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"http://api.logdna.co/v1/config/view\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestView_expectInvalidJSONError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigMultipleChannelsInvalidJSON(),
				ExpectError: regexp.MustCompile("Error: bodytemplate is not a valid JSON string"),
			},
		},
	})
}

func TestView_expectTriggerIntervalError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigTriggerIntervalError(),
				ExpectError: regexp.MustCompile(`"message":"\\"channels\[0\]\.triggerinterval\\" must be one of \[15m, 30m, 1h, 6h, 12h, 24h\]"`),
			},
		},
	})
}

func TestView_expectImmediateError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigImmediateError(),
				ExpectError: regexp.MustCompile(`"message":"\\"channels\[0\]\.immediate\\" must be a boolean"`),
			},
		},
	})
}

func TestView_expectURLError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigURLError(),
				ExpectError: regexp.MustCompile(`"message":"\\"channels\[0\]\.url\\" must be a valid uri"`),
			},
		},
	})
}

func TestView_expectMethodError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigMethodError(),
				ExpectError: regexp.MustCompile(`"message":"\\"channels\[0\].method\\" must be one of \[post, put, patch, get, delete\]"`),
			},
		},
	})
}

func TestView_expectServiceKeyError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigServiceKeyError(),
				ExpectError: regexp.MustCompile("The argument \"servicekey\" is required, but no definition was found."),
			},
		},
	})
}

func TestView_expectNameError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigNameError(),
				ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
			},
		},
	})
}

func TestView_expectAppsError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigAppsError(),
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"apps\": list of string required."),
			},
		},
	})
}
func TestView_expectCategoriesError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigCategoriesError(),
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"categories\": list of string required."),
			},
		},
	})
}

func TestView_expectHostsError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigHostsError(),
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"hosts\": list of string required."),
			},
		},
	})
}

func TestView_expectLevelsError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigLevelsError(),
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"levels\": list of string required."),
			},
		},
	})
}

func TestView_expectTagsError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigTagsError(),
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"tags\": list of string required."),
			},
		},
	})
}

func TestView_expectEmailTriggerLimitError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigEmailTriggerLimitError(),
				ExpectError: regexp.MustCompile("Error: \"email_channel.0.triggerlimit\" must be between 1 and 100,000 inclusive, got: 0"),
			},
		},
	})
}

func TestView_expectPagerDutyTriggerLimitError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigPagerDutyTriggerLimitError(),
				ExpectError: regexp.MustCompile("Error: \"pagerduty_channel.0.triggerlimit\" must be between 1 and 100,000 inclusive, got: 0"),
			},
		},
	})
}

func TestView_expectWebhookTriggerLimitError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigWebhookTriggerLimitError(),
				ExpectError: regexp.MustCompile("Error: \"webhook_channel.0.triggerlimit\" must be between 1 and 100,000 inclusive, got: 0"),
			},
		},
	})
}

func TestView_expectMissingEmails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigMissingEmails(),
				ExpectError: regexp.MustCompile("The argument \"emails\" is required, but no definition was found."),
			},
		},
	})
}

func TestView_expectMissingKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigMissingKey(),
				ExpectError: regexp.MustCompile("The argument \"key\" is required, but no definition was found."),
			},
		},
	})
}

func TestView_expectMissingURL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testViewConfigMissingURL(),
				ExpectError: regexp.MustCompile("The argument \"url\" is required, but no definition was found."),
			},
		},
	})
}
func TestViewBasic(t *testing.T) {
	name := "test"
	query := "test"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testViewConfigBasic(name, query),
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", name),
					resource.TestCheckResourceAttr("logdna_view.new", "query", query),
				),
			},
		},
	})
}

func TestViewBasicUpdate(t *testing.T) {
	name := "test"
	query := "test"
	name2 := "test2"
	query2 := "test2"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testViewConfigBasic(name, query),
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", name),
					resource.TestCheckResourceAttr("logdna_view.new", "query", query),
				),
			},
			{
				Config: testViewConfigBasic(name2, query2),
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", name2),
					resource.TestCheckResourceAttr("logdna_view.new", "query", query2),
				),
			},
		},
	})
}

func TestViewJSONUpdateError(t *testing.T) {
	name := "test"
	query := "test"
	app1 := "app1"
	app2 := "app2"
	levels1 := "fatal"
	levels2 := "critical"
	host1 := "host1"
	host2 := "host2"
	category1 := "DEMOCATEGORY1"
	category2 := "DemoCategory2"
	tags1 := "host1"
	tags2 := "host2"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testViewConfigMultipleChannels(name, query, app1, app2, levels1, levels2, host1, host2, category1, category2, tags1, tags2),
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
				),
			},
			{
				Config:      testViewConfigMultipleChannelsInvalidJSON(),
				ExpectError: regexp.MustCompile("Error: bodytemplate is not a valid JSON string"),
			},
		},
	})
}

func TestViewBulkEmails(t *testing.T) {
	name := "test"
	query := "test"
	app1 := "app1"
	app2 := "app2"
	levels1 := "fatal"
	levels2 := "critical"
	host1 := "host1"
	host2 := "host2"
	category1 := "DEMOCATEGORY1"
	category2 := "DemoCategory2"
	tags1 := "host1"
	tags2 := "host2"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testViewConfigBulkEmails(name, query, app1, app2, levels1, levels2, host1, host2, category1, category2, tags1, tags2),
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", name),
					resource.TestCheckResourceAttr("logdna_view.new", "query", query),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.0", app1),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.1", app2),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.0", "DemoCategory1"), // This value on the server is mixed case
					resource.TestCheckResourceAttr("logdna_view.new", "categories.1", category2),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.%", "7"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", host1),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.1", host2),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", host1),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.0", levels1),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.1", levels2),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.0", tags1),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.1", tags2),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channels.#", "0"),
				),
			},
		},
	})
}

func TestViewBulkEmailsUpdate(t *testing.T) {
	name := "test"
	query := "test"
	app1 := "app1"
	app2 := "app2"
	levels1 := "fatal"
	levels2 := "critical"
	host1 := "host1"
	host2 := "host2"
	category1 := "DEMOCATEGORY1"
	category2 := "DemoCategory2"
	tags1 := "host1"
	tags2 := "host2"

	name2 := "test2"
	query2 := "query2"
	app3 := "app3"
	app4 := "app4"
	levels3 := "error"
	levels4 := "warning"
	host3 := "host3"
	host4 := "host4"
	tags3 := "tags3"
	tags4 := "tags4"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testViewConfigBulkEmails(name, query, app1, app2, levels1, levels2, host1, host2, category1, category2, tags1, tags2),
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", name),
					resource.TestCheckResourceAttr("logdna_view.new", "query", query),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.0", app1),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.1", app2),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.0", "DemoCategory1"), // This value on the server is mixed case
					resource.TestCheckResourceAttr("logdna_view.new", "categories.1", category2),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.%", "7"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", host1),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.1", host2),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", host1),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.0", levels1),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.1", levels2),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.0", host1),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.1", host2),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channels.#", "0"),
				),
			},
			{
				Config: testViewConfigBulkEmails(name2, query2, app3, app4, levels3, levels4, host3, host4, category1, category2, tags3, tags4),
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", name2),
					resource.TestCheckResourceAttr("logdna_view.new", "query", query2),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.0", app3),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.1", app4),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.0", "DemoCategory1"), // This value on the server is mixed case
					resource.TestCheckResourceAttr("logdna_view.new", "categories.1", category2),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.%", "7"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", host3),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.1", host4),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", host3),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.0", levels3),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.1", levels4),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.0", tags3),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.1", tags4),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channels.#", "0"),
				),
			},
		},
	})
}

func TestViewMultipleChannels(t *testing.T) {
	name := "test"
	query := "test"
	app1 := "app1"
	app2 := "app2"
	levels1 := "fatal"
	levels2 := "critical"
	host1 := "host1"
	host2 := "host2"
	category1 := "DemoCategory1"
	category2 := "DemoCategory2"
	tags1 := "host1"
	tags2 := "host2"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testViewConfigMultipleChannels(name, query, app1, app2, levels1, levels2, host1, host2, category1, category2, tags1, tags2),
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", name),
					resource.TestCheckResourceAttr("logdna_view.new", "query", query),
					resource.TestCheckResourceAttr("logdna_view.new", "%", "11"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.0", app1),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.1", app2),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.0", category1),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.1", category2),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.%", "7"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", host1),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.1", host2),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.0", levels1),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.1", levels2),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.key", "Your PagerDuty API key goes here"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.1", levels2),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.0", tags1),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.1", tags2),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.%", "9"),
					// The JSON will have newlines per our API which uses JSON.stringify(obj, null, 2) as the value
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.bodytemplate", "{\n  \"fields\": {\n    \"description\": \"{{ matches }} matches found for {{ name }}\",\n    \"issuetype\": {\n      \"name\": \"Bug\"\n    },\n    \"project\": {\n      \"key\": \"test\"\n    },\n    \"summary\": \"Alert From {{ name }}\"\n  }\n}"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.headers.%", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.headers.hello", "test3"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.headers.test", "test2"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.method", "post"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.url", "https://yourwebhook/endpoint"),
				),
			},
		},
	})
}

func testViewInvalidURL() string {
	return fmt.Sprintf(`provider "logdna" {
		url = "http://api.logdna.co"
		servicekey = "%s"
	  }

	  resource "logdna_view" "new" {
		name = "test"
		query = "test"
	  }`, serviceKey)
}

func testViewConfigMultipleChannelsInvalidJSON() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	  resource "logdna_view" "new" {
		name = "test"
		query = "test"
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
	  }`, serviceKey)
}

func testViewConfigTriggerIntervalError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
      }

      resource "logdna_view" "new" {
        name = "test"
        query = "test"
        email_channel {
          emails          = ["test@logdna.com"]
          immediate       = "false"
          operator        = "absence"
          terminal        = "true"
          timezone        = "Pacific/Samoa"
          triggerinterval = "17m"
          triggerlimit    = 15
        }
      }`, serviceKey)
}

func testViewConfigImmediateError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
      }

      resource "logdna_view" "new" {
        name = "test"
        query = "test"
        email_channel {
          emails          = ["test@logdna.com"]
          immediate       = "no"
          operator        = "absence"
          terminal        = "true"
          timezone        = "Pacific/Samoa"
          triggerinterval = "15m"
          triggerlimit    = 15
        }
      }`, serviceKey)
}

func testViewConfigURLError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	  resource "logdna_view" "new" {
		name = "test"
		query = "test"
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
	  }`, serviceKey)
}

func testViewConfigMethodError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	  resource "logdna_view" "new" {
		name = "test"
		query = "test"
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
	  }`, serviceKey)
}

func testViewConfigServiceKeyError() string {
	return `provider "logdna" {
	}

	resource "logdna_view" "new" {
		name = "test"
		query = "test"
	}`
}

func testViewConfigNameError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		query = "test"
	}`, serviceKey)
}

func testViewConfigAppsError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name = "test"
		query = "test"
		apps = "test"
	}`, serviceKey)
}

func testViewConfigCategoriesError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name = "test"
		query = "test"
		categories = "test"
	}`, serviceKey)
}

func testViewConfigHostsError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name = "test"
		query = "test"
		hosts = "test"
	}`, serviceKey)
}

func testViewConfigLevelsError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name = "test"
		query = "test"
		levels = "test"
	}`, serviceKey)
}

func testViewConfigTagsError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name = "test"
		query = "test"
		tags = "test"
	}`, serviceKey)
}

func testViewConfigEmailTriggerLimitError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name     = "test"
		query    = "test"
		email_channel {
			emails          = ["test@logdna.com"]
			immediate       = "false"
			operator        = "absence"
			terminal        = "true"
			triggerinterval = "15m"
			triggerlimit    = 0
			timezone        = "Pacific/Samoa"
		}
	}`, serviceKey)
}

func testViewConfigPagerDutyTriggerLimitError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name     = "test"
		query    = "test"
		pagerduty_channel {
			immediate       = "false"
			key             = "Your PagerDuty API key goes here"
			terminal        = "true"
			triggerinterval = "15m"
			triggerlimit    = 0
		}
	}`, serviceKey)
}

func testViewConfigWebhookTriggerLimitError() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name     = "test"
		query    = "test"
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
	}`, serviceKey)
}

func testViewConfigMissingEmails() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name     = "test"
		query    = "test"
		email_channel {
			immediate       = "false"
			operator        = "absence"
			terminal        = "true"
			triggerinterval = "15m"
			triggerlimit    = 15
			timezone        = "Pacific/Samoa"
		}
	}`, serviceKey)
}

func testViewConfigMissingKey() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name     = "test"
		query    = "test"
		pagerduty_channel {
			immediate       = "false"
			terminal        = "true"
			triggerinterval = "15m"
			triggerlimit    = 15
		}
	}`, serviceKey)
}

func testViewConfigMissingURL() string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

	resource "logdna_view" "new" {
		name     = "test"
		query    = "test"
		webhook_channel {
			triggerlimit = 15
		}
	}`, serviceKey)
}

func testViewConfigBasic(name, query string) string {
	return fmt.Sprintf(`provider "logdna" {
			servicekey = "%s"
		}

	resource "logdna_view" "new" {
		name     = "%s"
		query    = "%s"
	}`, serviceKey, name, query)
}

func testViewConfigBulkEmails(name, query, app1, app2, levels1, levels2, host1, host2, category1, category2, tags1, tags2 string) string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	}

  resource "logdna_view" "new" {
	name     = "%s"
	query    = "%s"
	apps     = ["%s", "%s"]
	levels   = ["%s", "%s"]
	hosts    = ["%s", "%s"]
	categories = ["%s", "%s"]
	tags     = ["%s", "%s"]
	email_channel {
	  emails          = ["test@logdna.com"]
	  immediate       = "false"
	  operator        = "absence"
	  terminal        = "true"
	  triggerinterval = "15m"
	  triggerlimit    = 15
	  timezone        = "Pacific/Samoa"
	}
	email_channel {
	  emails          = ["test@logdna.com"]
	  immediate       = "false"
	  operator        = "absence"
	  terminal        = "true"
	  timezone        = "Pacific/Samoa"
	  triggerlimit    = 15
	  triggerinterval = "15m"
	}
  }`, serviceKey, name, query, app1, app2, levels1, levels2, host1, host2, category1, category2, tags1, tags2)
}

func testViewConfigMultipleChannels(name, query, app1, app2, levels1, levels2, host1, host2, category1, category2, tags1, tags2 string) string {
	return fmt.Sprintf(`provider "logdna" {
		servicekey = "%s"
	  }

	  resource "logdna_view" "new" {
		name     = "%s"
		query    = "%s"
		apps     = ["%s", "%s"]
		levels   = ["%s", "%s"]
		hosts    = ["%s", "%s"]
		categories = ["%s", "%s"]
		tags     = ["%s", "%s"]
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
		  bodytemplate = jsonencode({
			fields = {
				description = "{{ matches }} matches found for {{ name }}"
				issuetype = {
					name = "Bug"
				}
				project = {
					key = "test"
				},
				summary = "Alert From {{ name }}"
		   }
		  })
		  immediate       = "false"
		  method          = "post"
		  url             = "https://yourwebhook/endpoint"
		  terminal        = "true"
		  triggerinterval = "15m"
		  triggerlimit    = 15
		}
	  }`, serviceKey, name, query, app1, app2, levels1, levels2, host1, host2, category1, category2, tags1, tags2)
}

func testViewExists(n string) resource.TestCheckFunc {
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
