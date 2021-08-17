package logdna

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const pcStr = `
provider "logdna" {
	servicekey = "%s"
}
`
const rsStr = `
resource "logdna_alert" "test" {
	name = "%s"
	%s
}
`
const rsStrMultiple = `
resource "logdna_alert" "test" {
	name = "%s"
	%s
	%s
	%s
}
`
const dsStr = `
data "logdna_alert" "remote" {
	presetid = logdna_alert.test.id
}
`
const email = `
	email_channel {
		emails          = ["test@logdna.com"]
		immediate       = "false"
		operator        = "presence"
		triggerlimit    = 15
		triggerinterval = "15m"
		terminal        = "true"
		timezone        = "Pacific/Samoa"
	}
`
const pagerduty = `
	pagerduty_channel {
		immediate       = "false"
		key             = "Your PagerDuty API key goes here"
		terminal        = "true"
		triggerinterval = "15m"
		triggerlimit    = 15
	}
`
const webhook = `
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
`

func TestDataSourceAlert_basicEmail(t *testing.T) {
	data := "data.logdna_alert.remote"
	rsName := "test"

	args := []string{rsName, email}
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlertConfig(args...),
				Check: resource.ComposeTestCheckFunc(
					testDataSourceAlertExists(data),
					resource.TestCheckResourceAttr(data, "name", rsName),
					resource.TestCheckResourceAttr(data, "email_channel.#", "1"),
					resource.TestCheckResourceAttr(data, "email_channel.0.%", "7"),
					resource.TestCheckResourceAttr(data, "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr(data, "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr(data, "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr(data, "email_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr(data, "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr(data, "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr(data, "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr(data, "email_channel.0.triggerlimit", "15"),
				),
			},
		},
	})
}

func TestDataSourceAlert_basicPagerDuty(t *testing.T) {
	data := "data.logdna_alert.remote"
	rsName := "test"

	args := []string{rsName, pagerduty}
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlertConfig(args...),
				Check: resource.ComposeTestCheckFunc(
					testDataSourceAlertExists(data),
					resource.TestCheckResourceAttr(data, "name", rsName),
					resource.TestCheckResourceAttr(data, "pagerduty_channel.#", "1"),
					resource.TestCheckResourceAttr(data, "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr(data, "pagerduty_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr(data, "pagerduty_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr(data, "pagerduty_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr(data, "pagerduty_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr(data, "pagerduty_channel.0.triggerlimit", "15"),
				),
			},
		},
	})
}

func TestDataSourceAlert_basicWebhook(t *testing.T) {
	data := "data.logdna_alert.remote"
	rsName := "test"

	args := []string{rsName, webhook}
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlertConfig(args...),
				Check: resource.ComposeTestCheckFunc(
					testDataSourceAlertExists(data),
					resource.TestCheckResourceAttr(data, "name", rsName),
					resource.TestCheckResourceAttr(data, "webhook_channel.#", "1"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.%", "9"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.headers.%", "2"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.headers.hello", "test3"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.headers.test", "test2"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.bodytemplate", "{\n  \"fields\": {\n    \"description\": \"{{ matches }} matches found for {{ name }}\",\n    \"issuetype\": {\n      \"name\": \"Bug\"\n    },\n    \"project\": {\n      \"key\": \"test\"\n    },\n    \"summary\": \"Alert From {{ name }}\"\n  }\n}"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.method", "post"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.url", "https://yourwebhook/endpoint"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr(data, "webhook_channel.0.triggerlimit", "15"),
				),
			},
		},
	})
}

func TestDataSourceAlert_multipleChannels(t *testing.T) {
	data := "data.logdna_alert.remote"
	rsName := "test"

	args := []string{rsName, "multiple"}
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlertConfig(args...),
				Check: resource.ComposeTestCheckFunc(
					testDataSourceAlertExists(data),
					resource.TestCheckResourceAttr(data, "name", rsName),
				),
			},
		},
	})
}

func testDataSourceAlertConfig(args ...string) string {
	name := args[0]
	integration := args[1]
	isMultiple := false
	if integration == "multiple" {
		isMultiple = true
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, pcStr, serviceKey)

	if isMultiple {
		fmt.Fprintf(&sb, rsStrMultiple, name, email, pagerduty, webhook)
	} else {
		fmt.Fprintf(&sb, rsStr, name, integration)
	}

	sb.WriteString(dsStr)
	return sb.String()
}

func testDataSourceAlertExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		return nil
	}
}
