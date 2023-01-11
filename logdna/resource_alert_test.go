package logdna

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var alertDefaults = cloneDefaults(rsDefaults["alert"])

func TestAlert_ErrorProviderUrl(t *testing.T) {
	pcArgs := []string{serviceKey, "https://api.logdna.co"}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("alert", "new", pcArgs, alertDefaults, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"https://api.logdna.co/v1/config/presetalert\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestAlert_ErrorResourceName(t *testing.T) {
	args := cloneDefaults(chnlDefaults["alert_channel"])
	args["name"] = ""

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("alert", "new", globalPcArgs, args, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_ErrorOrgType(t *testing.T) {
	pcArgs := []string{enterpriseServiceKey, apiHostUrl, "enterprise"}
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("alert", "new", pcArgs, alertDefaults, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: Only regular organizations can instantiate a \"logdna_alert\" resource"),
			},
		},
	})
}

func TestAlert_ErrorsChannel(t *testing.T) {
	imArgs := map[string]map[string]string{"email_channel": cloneDefaults(chnlDefaults["email_channel"])}
	imArgs["email_channel"]["immediate"] = `"not a bool"`
	immdte := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, imArgs, nilLst)

	opArgs := map[string]map[string]string{"pagerduty_channel": cloneDefaults(chnlDefaults["pagerduty_channel"])}
	opArgs["pagerduty_channel"]["operator"] = `1000`
	opratr := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, opArgs, nilLst)

	trArgs := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	trArgs["webhook_channel"]["terminal"] = `"invalid"`
	trmnal := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, trArgs, nilLst)

	tiArgs := map[string]map[string]string{"email_channel": cloneDefaults(chnlDefaults["email_channel"])}
	tiArgs["email_channel"]["triggerinterval"] = `18`
	tintvl := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, tiArgs, nilLst)

	tlArgs := map[string]map[string]string{"slack_channel": cloneDefaults(chnlDefaults["slack_channel"])}
	tlArgs["slack_channel"]["triggerlimit"] = `0`
	tlimit := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, tlArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      immdte,
				ExpectError: regexp.MustCompile(`"\\"channels\[0\].immediate\\" must be a boolean"`),
			},
			{
				Config:      opratr,
				ExpectError: regexp.MustCompile(`"\\"channels\[0\].operator\\" must be one of \[presence, absence\]"`),
			},
			{
				Config:      trmnal,
				ExpectError: regexp.MustCompile(`"\\"channels\[0\].terminal\\" must be a boolean"`),
			},
			{
				Config:      tintvl,
				ExpectError: regexp.MustCompile(`"\\"channels\[0\].triggerinterval\\" must be one of \[15m, 30m, 1h, 6h, 12h, 24h\]"`),
			},
			{
				Config:      tlimit,
				ExpectError: regexp.MustCompile(`Error: ".*channel.0.triggerlimit" must be between 1 and 100,000 inclusive`),
			},
		},
	})
}

func TestAlert_ErrorsEmailChannel(t *testing.T) {
	msArgs := map[string]map[string]string{"email_channel": cloneDefaults(chnlDefaults["email_channel"])}
	msArgs["email_channel"]["emails"] = ""
	misngE := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, msArgs, nilLst)

	inArgs := map[string]map[string]string{"email_channel": cloneDefaults(chnlDefaults["email_channel"])}
	inArgs["email_channel"]["emails"] = `"not an array of strings"`
	invldE := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, inArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      misngE,
				ExpectError: regexp.MustCompile("The argument \"emails\" is required, but no definition was found."),
			},
			{
				Config:      invldE,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"emails\": list of string required"),
			},
		},
	})
}

func TestAlert_ErrorsPagerDutyChannel(t *testing.T) {
	chArgs := map[string]map[string]string{"pagerduty_channel": cloneDefaults(chnlDefaults["pagerduty_channel"])}
	chArgs["pagerduty_channel"]["key"] = ""

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, chArgs, nilLst),
				ExpectError: regexp.MustCompile("The argument \"key\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_ErrorsSlackChannel(t *testing.T) {
	ulInvd := map[string]map[string]string{"slack_channel": cloneDefaults(chnlDefaults["slack_channel"])}
	ulInvd["slack_channel"]["url"] = `"this is not a valid url"`
	ulCfgE := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, ulInvd, nilLst)

	ulMsng := map[string]map[string]string{"slack_channel": cloneDefaults(chnlDefaults["slack_channel"])}
	ulMsng["slack_channel"]["url"] = ""
	ulCfgM := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, ulMsng, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      ulCfgE,
				ExpectError: regexp.MustCompile(`"message":"\\"channels\[0\]\.url\\" must be a valid uri"`),
			},
			{
				Config:      ulCfgM,
				ExpectError: regexp.MustCompile("The argument \"url\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_ErrorsWebhookChannel(t *testing.T) {
	btArgs := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	btArgs["webhook_channel"]["bodytemplate"] = `"{\"test\": }"`
	btCfgE := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, btArgs, nilLst)

	mdArgs := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	mdArgs["webhook_channel"]["method"] = `"false"`
	mdCfgE := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, mdArgs, nilLst)

	ulInvd := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	ulInvd["webhook_channel"]["url"] = `"this is not a valid url"`
	ulCfgE := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, ulInvd, nilLst)

	ulMsng := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	ulMsng["webhook_channel"]["url"] = ""
	ulCfgM := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, ulMsng, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      btCfgE,
				ExpectError: regexp.MustCompile("Error: bodytemplate is not a valid JSON string"),
			},
			{
				Config:      mdCfgE,
				ExpectError: regexp.MustCompile(`"message":"\\"channels\[0\].method\\" must be one of \[post, put, patch, get, delete\]"`),
			},
			{
				Config:      ulCfgE,
				ExpectError: regexp.MustCompile(`"message":"\\"channels\[0\]\.url\\" must be a valid uri"`),
			},
			{
				Config:      ulCfgM,
				ExpectError: regexp.MustCompile("The argument \"url\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_Basic(t *testing.T) {
	chArgs := map[string]map[string]string{"email_channel": cloneDefaults(chnlDefaults["email_channel"])}
	iniCfg := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, chArgs, nilLst)

	rsArgs := cloneDefaults(rsDefaults["alert"])
	rsArgs["name"] = `"test2"`
	updCfg := fmtTestConfigResource("alert", "new", globalPcArgs, rsArgs, chArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: iniCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("alert", "new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.%", "7"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: updCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("alert", "new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test2"),
				),
			},
			{
				ResourceName:      "logdna_alert.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAlert_BulkChannels(t *testing.T) {
	emArgs := map[string]map[string]string{
		"email_channel":  cloneDefaults(chnlDefaults["email_channel"]),
		"email1_channel": cloneDefaults(chnlDefaults["email_channel"]),
	}
	emsCfg := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, emArgs, nilLst)

	pdArgs := map[string]map[string]string{
		"pagerduty_channel":  cloneDefaults(chnlDefaults["pagerduty_channel"]),
		"pagerduty1_channel": cloneDefaults(chnlDefaults["pagerduty_channel"]),
	}
	pdsCfg := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, pdArgs, nilLst)

	slArgs := map[string]map[string]string{
		"slack_channel":  cloneDefaults(chnlDefaults["slack_channel"]),
		"slack1_channel": cloneDefaults(chnlDefaults["slack_channel"]),
	}
	slsCfg := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, slArgs, nilLst)

	wbArgs := map[string]map[string]string{
		"webhook_channel":  cloneDefaults(chnlDefaults["webhook_channel"]),
		"webhook1_channel": cloneDefaults(chnlDefaults["webhook_channel"]),
	}
	wbsCfg := fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, wbArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: emsCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("alert", "new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.%", "7"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.1.%", "7"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: pdsCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("alert", "new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.1.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: slsCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("alert", "new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.1.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: wbsCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("alert", "new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.%", "9"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.1.%", "9"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "0"),
				),
			},
		},
	})
}

func TestAlert_MultipleChannels(t *testing.T) {
	chArgs := map[string]map[string]string{
		"email_channel":     cloneDefaults(chnlDefaults["email_channel"]),
		"pagerduty_channel": cloneDefaults(chnlDefaults["pagerduty_channel"]),
		"slack_channel":     cloneDefaults(chnlDefaults["slack_channel"]),
		"webhook_channel":   cloneDefaults(chnlDefaults["webhook_channel"]),
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmtTestConfigResource("alert", "new", globalPcArgs, alertDefaults, chArgs, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("alert", "new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.key", "Your PagerDuty API key goes here"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.triggerinterval", "30m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.url", "https://hooks.slack.com/services/identifier/secret"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.%", "9"),
					// The JSON will have newlines per our API which uses JSON.stringify(obj, null, 2) as the value
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.bodytemplate", "{\n  \"fields\": {\n    \"description\": \"{{ matches }} matches found for {{ name }}\",\n    \"issuetype\": {\n      \"name\": \"Bug\"\n    },\n    \"project\": {\n      \"key\": \"test\"\n    },\n    \"summary\": \"Alert from {{ name }}\"\n  }\n}"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.headers.%", "2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.headers.hello", "test3"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.headers.test", "test2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.method", "post"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.url", "https://yourwebhook/endpoint"),
				),
			},
			{
				ResourceName:      "logdna_alert.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
