package logdna

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const ds = `
data "logdna_alert" "remote" {
	presetid = logdna_alert.test.id
}
`

func TestDataAlert_ErrorOrgType(t *testing.T) {
	pcArgs := []string{enterpriseServiceKey, apiHostUrl, "enterprise"}
	alertConfig := fmtTestConfigResource("alert", "test", pcArgs, alertDefaults, nilOpt, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf("%s\n%s", alertConfig, ds),
				ExpectError: regexp.MustCompile("Error: Only regular organizations can instantiate a \"logdna_alert\" resource"),
			},
		},
	})
}

func TestDataAlert_BulkChannels(t *testing.T) {
	emArgs := map[string]map[string]string{
		"email_channel":  cloneDefaults(chnlDefaults["email_channel"]),
		"email1_channel": cloneDefaults(chnlDefaults["email_channel"]),
	}
	emsCfg := fmtTestConfigResource("alert", "test", globalPcArgs, alertDefaults, emArgs, nilLst)

	pdArgs := map[string]map[string]string{
		"pagerduty_channel":  cloneDefaults(chnlDefaults["pagerduty_channel"]),
		"pagerduty1_channel": cloneDefaults(chnlDefaults["pagerduty_channel"]),
	}
	pdsCfg := fmtTestConfigResource("alert", "test", globalPcArgs, alertDefaults, pdArgs, nilLst)

	slArgs := map[string]map[string]string{
		"slack_channel":  cloneDefaults(chnlDefaults["slack_channel"]),
		"slack1_channel": cloneDefaults(chnlDefaults["slack_channel"]),
	}
	slsCfg := fmtTestConfigResource("alert", "test", globalPcArgs, alertDefaults, slArgs, nilLst)

	wbArgs := map[string]map[string]string{
		"webhook_channel":  cloneDefaults(chnlDefaults["webhook_channel"]),
		"webhook1_channel": cloneDefaults(chnlDefaults["webhook_channel"]),
	}
	wbsCfg := fmtTestConfigResource("alert", "test", globalPcArgs, alertDefaults, wbArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf("%s\n%s", emsCfg, ds),
				Check: resource.ComposeTestCheckFunc(
					testDataSourceAlertExists("data.logdna_alert.remote"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.#", "2"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.0.%", "7"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.1.%", "7"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.#", "0"),
				),
			},
			{
				Config: fmt.Sprintf("%s\n%s", pdsCfg, ds),
				Check: resource.ComposeTestCheckFunc(
					testDataSourceAlertExists("data.logdna_alert.remote"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "name", "test"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.#", "2"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.1.%", "6"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.#", "0"),
				),
			},
			{
				Config: fmt.Sprintf("%s\n%s", slsCfg, ds),
				Check: resource.ComposeTestCheckFunc(
					testDataSourceAlertExists("data.logdna_alert.remote"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "name", "test"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.#", "2"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.0.%", "6"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.1.%", "6"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.#", "0"),
				),
			},
			{
				Config: fmt.Sprintf("%s\n%s", wbsCfg, ds),
				Check: resource.ComposeTestCheckFunc(
					testDataSourceAlertExists("data.logdna_alert.remote"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "name", "test"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.#", "2"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.%", "9"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.1.%", "9"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.#", "0"),
				),
			},
		},
	})
}

func TestDataSourceAlert_MultipleChannels(t *testing.T) {
	chArgs := map[string]map[string]string{
		"email_channel":     cloneDefaults(chnlDefaults["email_channel"]),
		"pagerduty_channel": cloneDefaults(chnlDefaults["pagerduty_channel"]),
		"slack_channel":     cloneDefaults(chnlDefaults["slack_channel"]),
		"webhook_channel":   cloneDefaults(chnlDefaults["webhook_channel"]),
	}
	fmtCfg := fmt.Sprintf("%s\n%s", fmtTestConfigResource("alert", "test", globalPcArgs, alertDefaults, chArgs, nilLst), ds)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmtCfg,
				Check: resource.ComposeTestCheckFunc(
					testDataSourceAlertExists("data.logdna_alert.remote"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "name", "test"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.#", "1"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.#", "1"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.0.key", "Your PagerDuty API key goes here"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "pagerduty_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.#", "1"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.0.%", "6"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.0.triggerinterval", "30m"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "slack_channel.0.url", "https://hooks.slack.com/services/identifier/secret"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.#", "1"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.%", "9"),
					// The JSON will have newlines per our API which uses JSON.stringify(obj, null, 2) as the value
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.bodytemplate", "{\n  \"fields\": {\n    \"description\": \"{{ matches }} matches found for {{ name }}\",\n    \"issuetype\": {\n      \"name\": \"Bug\"\n    },\n    \"project\": {\n      \"key\": \"test\"\n    },\n    \"summary\": \"Alert from {{ name }}\"\n  }\n}"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.headers.%", "2"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.headers.hello", "test3"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.headers.test", "test2"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.method", "post"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("data.logdna_alert.remote", "webhook_channel.0.url", "https://yourwebhook/endpoint"),
				),
			},
		},
	})
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
