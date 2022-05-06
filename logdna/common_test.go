package logdna

import (
	"fmt"
	"regexp"
	"strings"
)

const tmplPc = `provider "logdna" {
	%s
}`
const tmplRs = `%s %q %q {
%s
}`

var nilLst = []string{}
var nilOpt = map[string]map[string]string{}

var rsDefaults = map[string]map[string]string{
	"alert": {
		"name": `"test"`,
	},
	"view": {
		"apps":       "",
		"categories": "",
		"presetid":   "",
		"hosts":      "",
		"levels":     "",
		"name":       `"test"`,
		"query":      `"test"`,
		"tags":       "",
	},
	"category": {
		"name": `"test"`,
		"type": `"views"`,
	},
}
var chnlDefaults = map[string]map[string]string{
	"email": {
		"emails":          `["test@logdna.com"]`,
		"immediate":       `"false"`,
		"operator":        `"absence"`,
		"terminal":        `"true"`,
		"timezone":        `"Pacific/Samoa"`,
		"triggerinterval": `"15m"`,
		"triggerlimit":    `15`,
	},
	"pagerduty": {
		"immediate":       `"false"`,
		"operator":        `"presence"`,
		"key":             `"Your PagerDuty API key goes here"`,
		"terminal":        `"true"`,
		"triggerinterval": `"15m"`,
		"triggerlimit":    `15`,
	},
	"slack": {
		"immediate":       `"false"`,
		"operator":        `"absence"`,
		"terminal":        `"true"`,
		"triggerinterval": `"30m"`,
		"triggerlimit":    `15`,
		"url":             `"https://hooks.slack.com/services/identifier/secret"`,
	},
	"webhook": {
		"headers": "{\n" +
			"\t\t\thello = \"test3\"\n" +
			"\t\t\ttest = \"test2\"\n" +
			"\t\t}",
		"bodytemplate": `jsonencode({
				fields = {
					description = "{{ matches }} matches found for {{ name }}"
					issuetype = {
						name = "Bug"
					}
					project = {
						key = "test"
					}
					summary = "Alert from {{ name }}"
				}
			})`,
		"immediate":       `"false"`,
		"method":          `"post"`,
		"operator":        `"presence"`,
		"terminal":        `"true"`,
		"triggerinterval": `"15m"`,
		"triggerlimit":    `15`,
		"url":             `"https://yourwebhook/endpoint"`,
	},
}

func cloneDefaults(dfts map[string]string) map[string]string {
	clone := make(map[string]string)
	for k, v := range dfts {
		clone[k] = v
	}
	return clone
}

func fmtTestConfigResource(objTyp, rsName string, pcArgs []string, rsArgs map[string]string, chArgs map[string]map[string]string, dependencies []string) string {
	pc := fmtProviderBlock(pcArgs...)
	rs := fmtResourceBlock(objTyp, rsName, rsArgs, chArgs, dependencies)
	return fmt.Sprintf("%s\n%s", pc, rs)
}

func fmtProviderBlock(args ...string) string {
	opts := []string{serviceKey, ""}
	copy(opts, args)
	sk, ul := opts[0], opts[1]

	pcCfg := fmt.Sprintf(`servicekey = %q`, sk)
	if ul != "" {
		pcCfg = pcCfg + fmt.Sprintf("\n\turl = %q", ul)
	}

	return fmt.Sprintf(tmplPc, pcCfg)
}

func fmtResourceBlock(objTyp, rsName string, rsArgs map[string]string, chArgs map[string]map[string]string, dependencies []string) string {
	var rsCfg strings.Builder
	fmt.Fprint(&rsCfg, fmtBlockArgs(1, rsArgs))

	rgxDgt := regexp.MustCompile(`\d+`)
	for chName, chArgs := range chArgs {
		fmt.Fprintf(&rsCfg, "\t%s_channel {\n", rgxDgt.ReplaceAllString(chName, ""))
		fmt.Fprint(&rsCfg, fmtBlockArgs(2, chArgs))
		fmt.Fprintf(&rsCfg, "\t}\n")
	}

	if len(dependencies) > 0 {
		fmt.Fprintf(&rsCfg, "\tdepends_on = [\"%s\"]\n", strings.Join(dependencies[:], "\",\""))
	}

	rsType := fmt.Sprintf("logdna_%s", objTyp)
	return fmt.Sprintf(tmplRs, "resource", rsType, rsName, rsCfg.String())
}

func fmtBlockArgs(nstLvl int, opts map[string]string) string {
	var numTab strings.Builder
	for i := 0; i < nstLvl; i++ {
		fmt.Fprint(&numTab, "\t")
	}
	tabs := numTab.String()
	var blkCfg strings.Builder
	for arg, val := range opts {
		if val != "" {
			fmt.Fprintf(&blkCfg, "%s%s = %s\n", tabs, arg, val)
		}
	}
	return blkCfg.String()
}
