package logdna

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const ctgies = `["DEMOCATEGORY1", "DemoCategory2"]`

var viewDefaults = cloneDefaults(rsDefaults["view"])

func TestView_ErrorProviderUrl(t *testing.T) {
	pcArgs := []string{serviceKey, "https://api.logdna.co"}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("view", "new", pcArgs, viewDefaults, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"https://api.logdna.co/v1/config/view\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestView_ErrorOrgType(t *testing.T) {
	pcArgs := []string{enterpriseServiceKey, apiHostUrl, "enterprise"}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("view", "new", pcArgs, viewDefaults, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: Only regular organizations can instantiate a \"logdna_view\" resource"),
			},
		},
	})
}

func TestView_ErrorsResourceFields(t *testing.T) {
	nme := cloneDefaults(rsDefaults["view"])
	nme["name"] = ""
	nmeCfg := fmtTestConfigResource("view", "new", globalPcArgs, nme, nilOpt, nilLst)

	app := cloneDefaults(rsDefaults["view"])
	app["apps"] = `"invalid apps value"`
	appCfg := fmtTestConfigResource("view", "new", globalPcArgs, app, nilOpt, nilLst)

	ctg := cloneDefaults(rsDefaults["view"])
	ctg["categories"] = `"invalid categories value"`
	ctgCfg := fmtTestConfigResource("view", "new", globalPcArgs, ctg, nilOpt, nilLst)

	hst := cloneDefaults(rsDefaults["view"])
	hst["hosts"] = `"invalid hosts value"`
	hstCfg := fmtTestConfigResource("view", "new", globalPcArgs, hst, nilOpt, nilLst)

	lvl := cloneDefaults(rsDefaults["view"])
	lvl["levels"] = `"invalid levels value"`
	lvlCfg := fmtTestConfigResource("view", "new", globalPcArgs, lvl, nilOpt, nilLst)

	tgs := cloneDefaults(rsDefaults["view"])
	tgs["tags"] = `"invalid tags value"`
	tgsCfg := fmtTestConfigResource("view", "new", globalPcArgs, tgs, nilOpt, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      nmeCfg,
				ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
			},
			{
				Config:      appCfg,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"apps\": list of string required."),
			},
			{
				Config:      ctgCfg,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"categories\": list of string required."),
			},
			{
				Config:      hstCfg,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"hosts\": list of string required."),
			},
			{
				Config:      lvlCfg,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"levels\": list of string required."),
			},
			{
				Config:      tgsCfg,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"tags\": list of string required."),
			},
		},
	})
}

func TestView_ErrorsChannel(t *testing.T) {
	imArgs := map[string]map[string]string{"email_channel": cloneDefaults(chnlDefaults["email_channel"])}
	imArgs["email_channel"]["immediate"] = `"not a bool"`
	immdte := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, imArgs, nilLst)

	opArgs := map[string]map[string]string{"pagerduty_channel": cloneDefaults(chnlDefaults["pagerduty_channel"])}
	opArgs["pagerduty_channel"]["operator"] = `1000`
	opratr := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, opArgs, nilLst)

	trArgs := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	trArgs["webhook_channel"]["terminal"] = `"invalid"`
	trmnal := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, trArgs, nilLst)

	tiArgs := map[string]map[string]string{"email_channel": cloneDefaults(chnlDefaults["email_channel"])}
	tiArgs["email_channel"]["triggerinterval"] = `18`
	tintvl := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, tiArgs, nilLst)

	tlArgs := map[string]map[string]string{"slack_channel": cloneDefaults(chnlDefaults["slack_channel"])}
	tlArgs["slack_channel"]["triggerlimit"] = `0`
	tlimit := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, tlArgs, nilLst)

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

func TestView_ErrorsEmailChannel(t *testing.T) {
	msArgs := map[string]map[string]string{"email_channel": cloneDefaults(chnlDefaults["email_channel"])}
	msArgs["email_channel"]["emails"] = ""
	misngE := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, msArgs, nilLst)

	inArgs := map[string]map[string]string{"email_channel": cloneDefaults(chnlDefaults["email_channel"])}
	inArgs["email_channel"]["emails"] = `"not an array of strings"`
	invldE := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, inArgs, nilLst)

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

func TestView_ErrorsPagerDutyChannel(t *testing.T) {
	chArgs := map[string]map[string]string{"pagerduty_channel": cloneDefaults(chnlDefaults["pagerduty_channel"])}
	chArgs["pagerduty_channel"]["key"] = ""

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, chArgs, nilLst),
				ExpectError: regexp.MustCompile("The argument \"key\" is required, but no definition was found."),
			},
		},
	})
}

func TestView_ErrorsSlackChannel(t *testing.T) {
	ulInvd := map[string]map[string]string{"slack_channel": cloneDefaults(chnlDefaults["slack_channel"])}
	ulInvd["slack_channel"]["url"] = `"this is not a valid url"`
	ulCfgE := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, ulInvd, nilLst)

	ulMsng := map[string]map[string]string{"slack_channel": cloneDefaults(chnlDefaults["slack_channel"])}
	ulMsng["slack_channel"]["url"] = ""
	ulCfgM := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, ulMsng, nilLst)

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

func TestView_ErrorsWebhookChannel(t *testing.T) {
	btArgs := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	btArgs["webhook_channel"]["bodytemplate"] = `"{\"test\": }"`
	btCfgE := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, btArgs, nilLst)

	hdArgs := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	hdArgs["webhook_channel"]["headers"] = `["headers", "invalid", "array"]`
	hdCfgE := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, hdArgs, nilLst)

	mdArgs := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	mdArgs["webhook_channel"]["method"] = `"false"`
	mdCfgE := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, mdArgs, nilLst)

	ulInvd := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	ulInvd["webhook_channel"]["url"] = `"this is not a valid url"`
	ulCfgE := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, ulInvd, nilLst)

	ulMsng := map[string]map[string]string{"webhook_channel": cloneDefaults(chnlDefaults["webhook_channel"])}
	ulMsng["webhook_channel"]["url"] = ""
	ulCfgM := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, ulMsng, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      btCfgE,
				ExpectError: regexp.MustCompile("Error: bodytemplate is not a valid JSON string"),
			},
			{
				Config:      hdCfgE,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"headers\": map of string required"),
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

func TestView_Basic(t *testing.T) {
	iniCfg := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, nilOpt, nilLst)

	rsArgs := cloneDefaults(rsDefaults["view"])
	rsArgs["name"] = `"test2"`
	rsArgs["query"] = `"test2"`
	updCfg := fmtTestConfigResource("view", "new", globalPcArgs, rsArgs, nilOpt, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: iniCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("view", "new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
				),
			},
			{
				Config: updCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("view", "new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test2"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test2"),
				),
			},
			{
				ResourceName:      "logdna_view.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestView_BulkChannels(t *testing.T) {
	emArgs := map[string]map[string]string{
		"email_channel":  cloneDefaults(chnlDefaults["email_channel"]),
		"email1_channel": cloneDefaults(chnlDefaults["email_channel"]),
	}
	emsCfg := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, emArgs, nilLst)

	pdArgs := map[string]map[string]string{
		"pagerduty_channel":  cloneDefaults(chnlDefaults["pagerduty_channel"]),
		"pagerduty1_channel": cloneDefaults(chnlDefaults["pagerduty_channel"]),
	}
	pdsCfg := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, pdArgs, nilLst)

	slArgs := map[string]map[string]string{
		"slack_channel":  cloneDefaults(chnlDefaults["slack_channel"]),
		"slack1_channel": cloneDefaults(chnlDefaults["slack_channel"]),
	}
	slsCfg := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, slArgs, nilLst)

	wbArgs := map[string]map[string]string{
		"webhook_channel":  cloneDefaults(chnlDefaults["webhook_channel"]),
		"webhook1_channel": cloneDefaults(chnlDefaults["webhook_channel"]),
	}
	wbsCfg := fmtTestConfigResource("view", "new", globalPcArgs, viewDefaults, wbArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: emsCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("view", "new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.%", "7"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.%", "7"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: pdsCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("view", "new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.1.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: slsCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("view", "new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.1.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: wbsCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("view", "new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.%", "9"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.1.%", "9"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.#", "0"),
				),
			},
		},
	})
}

func TestView_MultipleChannels(t *testing.T) {
	chArgs := map[string]map[string]string{
		"email_channel":     cloneDefaults(chnlDefaults["email_channel"]),
		"pagerduty_channel": cloneDefaults(chnlDefaults["pagerduty_channel"]),
		"slack_channel":     cloneDefaults(chnlDefaults["slack_channel"]),
		"webhook_channel":   cloneDefaults(chnlDefaults["webhook_channel"]),
	}

	dependencies := []string{"logdna_category.cat_1", "logdna_category.cat_2"}

	cat1Args := map[string]string{
		"name": `"DemoCategory1"`,
		"type": `"views"`,
	}
	cat2Args := map[string]string{
		"name": `"DemoCategory2"`,
		"type": `"views"`,
	}

	rsArgs := cloneDefaults(rsDefaults["view"])
	rsArgs["apps"] = `["app1", "app2"]`
	rsArgs["categories"] = ctgies
	rsArgs["hosts"] = `["host1", "host2"]`
	rsArgs["levels"] = `["fatal", "critical"]`
	rsArgs["tags"] = `["tags1", "tags2"]`
	iniCfg := fmt.Sprintf(
		"%s\n%s\n%s",
		fmtTestConfigResource("view", "new", globalPcArgs, rsArgs, chArgs, dependencies),
		fmtResourceBlock("category", "cat_1", cat1Args, nilOpt, nilLst),
		fmtResourceBlock("category", "cat_2", cat2Args, nilOpt, nilLst),
	)

	rsUptd := cloneDefaults(rsDefaults["view"])
	rsUptd["apps"] = `["app3", "app4"]`
	rsUptd["categories"] = ctgies
	rsUptd["hosts"] = `["host3", "host4"]`
	rsUptd["levels"] = `["error", "warning"]`
	rsUptd["tags"] = `["tags3", "tags4"]`
	rsUptd["name"] = `"test2"`
	rsUptd["query"] = `"query2"`
	updCfg := fmt.Sprintf(
		"%s\n%s\n%s",
		fmtTestConfigResource("view", "new", globalPcArgs, rsUptd, chArgs, dependencies),
		fmtResourceBlock("category", "cat_1", cat1Args, nilOpt, nilLst),
		fmtResourceBlock("category", "cat_2", cat2Args, nilOpt, nilLst),
	)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: iniCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("view", "new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.0", "app1"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.1", "app2"),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.0", "DemoCategory1"), // This value on the server is mixed case
					resource.TestCheckResourceAttr("logdna_view.new", "categories.1", "DemoCategory2"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", "host1"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.1", "host2"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.0", "fatal"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.1", "critical"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.0", "tags1"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.1", "tags2"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.key", "Your PagerDuty API key goes here"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.triggerinterval", "30m"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.url", "https://hooks.slack.com/services/identifier/secret"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.%", "9"),
					// The JSON will have newlines per our API which uses JSON.stringify(obj, null, 2) as the value
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.bodytemplate", "{\n  \"fields\": {\n    \"description\": \"{{ matches }} matches found for {{ name }}\",\n    \"issuetype\": {\n      \"name\": \"Bug\"\n    },\n    \"project\": {\n      \"key\": \"test\"\n    },\n    \"summary\": \"Alert from {{ name }}\"\n  }\n}"),
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
			{
				Config: updCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("view", "new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test2"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "query2"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.0", "app3"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.1", "app4"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", "host3"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.1", "host4"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.0", "error"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.1", "warning"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.0", "tags3"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.1", "tags4"),
				),
			},
			{
				ResourceName:      "logdna_view.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestView_PresetAlert(t *testing.T) {
	chArgs := map[string]map[string]string{
		"email_channel":     cloneDefaults(chnlDefaults["email_channel"]),
		"pagerduty_channel": cloneDefaults(chnlDefaults["pagerduty_channel"]),
		"slack_channel":     cloneDefaults(chnlDefaults["slack_channel"]),
		"webhook_channel":   cloneDefaults(chnlDefaults["webhook_channel"]),
	}

	dependenciesIns := []string{
		"logdna_alert.test_preset_alert_ins",
		"logdna_category.test_category",
	}

	dependenciesUpd := []string{
		"logdna_alert.test_preset_alert_upd",
		"logdna_category.test_category",
	}

	catArgs := map[string]string{
		"name": `"DemoCategory"`,
		"type": `"views"`,
	}
	alertInsArgs := map[string]string{
		"name": `"Test Alert Ins"`,
	}
	alertUpdArgs := map[string]string{
		"name": `"Test Alert Upd"`,
	}

	rsArgs := cloneDefaults(rsDefaults["view"])
	rsArgs["apps"] = `["app1", "app2"]`
	rsArgs["categories"] = `[logdna_category.test_category.name]`
	rsArgs["hosts"] = `["host1", "host2"]`
	rsArgs["levels"] = `["fatal", "critical"]`
	rsArgs["tags"] = `["tags1", "tags2"]`
	rsArgs["presetid"] = `logdna_alert.test_preset_alert_ins.id`
	iniCfg := fmt.Sprintf(
		"%s\n%s\n%s",
		fmtResourceBlock("category", "test_category", catArgs, nilOpt, nilLst),
		fmtResourceBlock("alert", "test_preset_alert_ins", alertInsArgs, chArgs, nilLst),
		fmtTestConfigResource("view", "test_view", globalPcArgs, rsArgs, nilOpt, dependenciesIns),
	)

	rsUptd := cloneDefaults(rsDefaults["view"])
	rsUptd["apps"] = `["app3", "app4"]`
	rsUptd["categories"] = `[logdna_category.test_category.name]`
	rsUptd["hosts"] = `["host3", "host4"]`
	rsUptd["levels"] = `["error", "warning"]`
	rsUptd["tags"] = `["tags3", "tags4"]`
	rsUptd["name"] = `"test2"`
	rsUptd["query"] = `"query2"`
	rsUptd["presetid"] = `logdna_alert.test_preset_alert_upd.id`
	updCfg := fmt.Sprintf(
		"%s\n%s\n%s",
		fmtResourceBlock("category", "test_category", catArgs, nilOpt, nilLst),
		fmtResourceBlock("alert", "test_preset_alert_upd", alertUpdArgs, chArgs, nilLst),
		fmtTestConfigResource("view", "test_view", globalPcArgs, rsUptd, nilOpt, dependenciesUpd),
	)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: iniCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("view", "test_view"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "apps.0", "app1"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "apps.1", "app2"),
					resource.TestCheckResourceAttrPair(
						"logdna_alert.test_preset_alert_ins",
						"id",
						"logdna_view.test_view",
						"presetid",
					),
					resource.TestCheckResourceAttr("logdna_view.test_view", "categories.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "categories.0", "DemoCategory"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "hosts.0", "host1"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "hosts.1", "host2"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "levels.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "levels.0", "fatal"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "levels.1", "critical"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "tags.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "tags.0", "tags1"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "tags.1", "tags2"),
				),
			},
			{
				Config: updCfg,
				Check: resource.ComposeTestCheckFunc(
					testResourceExists("view", "test_view"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "name", "test2"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "query", "query2"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "apps.0", "app3"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "apps.1", "app4"),
					resource.TestCheckResourceAttrPair(
						"logdna_alert.test_preset_alert_upd",
						"id",
						"logdna_view.test_view",
						"presetid",
					),
					resource.TestCheckResourceAttr("logdna_view.test_view", "categories.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "categories.0", "DemoCategory"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "hosts.0", "host3"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "hosts.1", "host4"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "levels.0", "error"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "levels.1", "warning"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "tags.0", "tags3"),
					resource.TestCheckResourceAttr("logdna_view.test_view", "tags.1", "tags4"),
				),
			},
			{
				ResourceName:      "logdna_view.test_view",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestView_ErrorsConflictPresetId(t *testing.T) {
	chArgs := map[string]map[string]string{
		"email_channel":     cloneDefaults(chnlDefaults["email_channel"]),
		"pagerduty_channel": cloneDefaults(chnlDefaults["pagerduty_channel"]),
		"slack_channel":     cloneDefaults(chnlDefaults["slack_channel"]),
		"webhook_channel":   cloneDefaults(chnlDefaults["webhook_channel"]),
	}

	rsArgs := cloneDefaults(rsDefaults["view"])
	rsArgs["apps"] = `["app1", "app2"]`
	rsArgs["hosts"] = `["host1", "host2"]`
	rsArgs["levels"] = `["fatal", "critical"]`
	rsArgs["tags"] = `["tags1", "tags2"]`
	rsArgs["presetid"] = `"1q2w3e4r5t"`

	incCfg := fmtTestConfigResource("view", "test_view", globalPcArgs, rsArgs, chArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      incCfg,
				ExpectError: regexp.MustCompile("Error: Conflicting configuration arguments"),
			},
		},
	})
}
