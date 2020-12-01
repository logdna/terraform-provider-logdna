# logdna_alert Resource

Manages [LogDNA Preset Alerts](https://docs.logdna.com/docs/alerts). Preset Alerts are alerts that you can define separately from a specific View. Preset Alerts can be created standalone and then attached (or detached) to any View, as opposed to View-specific Alerts, which are created specifically for a certain View.

To get started, all you need to do is to specify a `name` and configuration for one of our currently supported alerts: email, webhook, or PagerDuty.

## Example - Basic Preset Alert

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_alert" "my_alert" {
  name = "My Preset Alert via Terraform"
  email_channel {
    emails          = ["test@logdna.com"]
    immediate       = false
    operator        = "presence"
    triggerlimit    = 15
    triggerinterval = "15m"
    terminal        = true
    timezone        = "Pacific/Samoa"
  }
}


```
## Example - Multi-channel Preset Alert

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_alert" "my_alert" {
  name = "Terraform Multi-channel Preset Alert"
  email_channel {
    emails          = ["test@logdna.com"]
    immediate       = false
    operator        = "absence"
    terminal        = true
    timezone        = "Pacific/Samoa"
    triggerinterval = "15m"
    triggerlimit    = 15
  }

  pagerduty_channel {
    immediate       = true
    key             = "Your PagerDuty API key goes here"
    terminal        = true
    triggerinterval = "15m"
    triggerlimit    = 15
  }

  webhook_channel {
    bodytemplate = jsonencode({
      message = "Alerts from {{name}}"
    })
    headers = {
      "Authentication" = "auth_header_value"
      "HeaderTwo"      = "ValueTwo"
    }
    immediate       = false
    method          = "post"
    terminal        = true
    triggerinterval = "15m"
    triggerlimit    = 15
    url             = "https://yourwebhook/endpoint"
 }
}
```

## Argument Reference

The following arguments are supported:

- `name`: (Required) The name this Preset Alert will be given, type _string_

### email_channel

`email_channel` supports the following arguments:

- `emails`: _(Required)_ An array of email addresses (each email is of type _string_) to notify in the Alert
- `immediate`: _(Optional)_ Whether the Alert will trigger immediately after the trigger limit is reached. _Note: Immediate can only be set to `true` for presence alerts_. Valid options are `true` and `false` for presence Alerts and `false` for absence Alerts, type _bool_ (**Default: false**)
- `operator`: _(Optional)_ Whether the Alert will trigger on the presence or absence of logs. Valid options are `presence` and `absence`, type _string_ (**Default: "presence"**)
- `terminal`: _(Optional)_ Whether the Alert will trigger after the `triggerinterval` if the Alert condition is met (e.g., send an Alert after 30s). Valid options are `true` and `false` for presence Alerts and `true` for absence Alerts, type _bool_ (**Default: true**)
- `timezone`: _(Optional)_ Which time zone the log timestamps will be formatted in. Timezones are represented as [database time zones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones), type _string_
- `triggerinterval`: _(Optional)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. For absence Alerts, valid options are: `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. Type _string_ (**Defaults: "30" for presence; "15m" for absence**)
- `triggerlimit`: _(Required)_ Number of lines before the Alert is triggered. (eg. Setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`), type _integer_

### pagerduty_channel

`pagerduty_channel` supports the following arguments:

- `immediate`: _(Optional)_ Whether the Alert will trigger immediately after the trigger limit is reached. _Note: Immediate can only be set to `true` for presence alerts_. Valid options are `true` and `false` for presence Alerts and `false` for absence Alerts, type _bool_ (**Default: false**)
- `key`: _(Required)_ PagerDuty service key, type _string_
- `operator`: _(Optional)_ Whether the Alert will trigger on the presence or absence of logs. Valid options are `presence` and `absence`, type _string_ (**Default: "presence"**)
- `terminal`: _(Optional)_ Whether the Alert will trigger after the `triggerinterval` if the Alert condition is met (e.g., send an Alert after 30s). Valid options are `true` and `false` for presence Alerts and `true` for absence Alerts, type _bool_ (**Default: true**)
- `triggerinterval`: _(Optional)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. For absence Alerts, valid options are: `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. Type _string_ (**Defaults: "30" for presence; "15m" for absence**)
- `triggerlimit`: _(Required)_ Number of lines before the Alert is triggered. (eg. Setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`), type _integer_

### webhook_channel

`webhook_channel` supports the following arguments:

- `bodytemplate`: _(Optional)_ JSON formatted string for the body of the webhook. We recommend using [`jsonencode()`](https://www.terraform.io/docs/configuration/functions/jsonencode.html) to easily convert a Terraform map into a JSON string. Type _string_
- `headers`: _(Optional)_ Key-value pair for webhook request headers and header values, type Map of _strings_
- `immediate`: _(Optional)_ Whether the Alert will trigger immediately after the trigger limit is reached. _Note: Immediate can only be set to `true` for presence alerts_. Valid options are `true` and `false` for presence Alerts and `false` for absence Alerts, type _bool_ (**Default: false**)
- `method`: _(Optional)_ Method used for the webhook request. Valid options are: `post`, `put`, `patch`, `get`, `delete`. Type _string_ (**Default: "post"**)
- `operator`: _(Optional)_ Whether the Alert will trigger on the presence or absence of logs. Valid options are `presence` and `absence`, type _string_ (**Default: "presence"**)
- `terminal`: _(Optional)_ Whether the Alert will trigger after the `triggerinterval` if the Alert condition is met (e.g., send an Alert after 30s). Valid options are `true` and `false` for presence Alerts and `true` for absence Alerts, type _bool_ (**Default: true**)
- `triggerinterval`: _(Optional)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. For absence Alerts, valid options are: `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. Type _string_ (**Defaults: "30" for presence; "15m" for absence**)
- `triggerlimit`: _(Required)_ Number of lines before the Alert is triggered. (eg. Setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`), type _integer_
- `url`: _(Required)_ URL of the webhook, type _string_


