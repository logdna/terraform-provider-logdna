# Resource: `logdna_alert`

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
    immediate       = "false"
    operator        = "presence"
    triggerlimit    = 15
    triggerinterval = "15m"
    terminal        = "true"
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

  webhook_channel {
    bodytemplate = jsonencode({
      message = "Alerts from {{name}}"
    })
    headers = {
      "Authentication" = "auth_header_value"
      "HeaderTwo"      = "ValueTwo"
    }
    immediate       = "false"
    method          = "post"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
    url             = "https://yourwebhook/endpoint"
 }
}
```

## Argument Reference

The following arguments are supported by `logdna_alert`:

- `name`: (Required) The name this Preset Alert will be given, type _string_

### email_channel

`email_channel` supports the following arguments:

- `emails`: **_[]string (Required)_** An array of email addresses (strings) to notify in the Alert
- `immediate`: **_string_** _(Optional; Default: `"false"`)_ Valid options are `"true"` and `"false"` for presence Alerts and `"false"` for absence Alerts.
- `operator`: **_string_** _(Optional; Default: `presence`)_ Whether the Alert will trigger on the presence or absence of logs. Valid options are `presence` and `absence`.
- `terminal`: **_string_** _(Optional; Default: `"true"`)_ Whether the Alert will trigger after the `triggerinterval` if the Alert condition is met (e.g. send an Alert after 30s). Valid options are `"true"` and `"false"` for presence Alerts and `"true"` for absence Alerts.
- `timezone`: **_string_** _(Optional)_ Which time zone the log timestamps will be formatted in. Timezones are represented as [database time zones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).
- `triggerinterval`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for absence)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. For absence Alerts, valid options are: `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`.
- `triggerlimit`: **_integer (Required)_** Number of lines before the Alert is triggered (e.g. setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`).

### pagerduty_channel

`pagerduty_channel` supports the following arguments:

- `immediate`: **_string_** _(Optional; Default: `"false"`)_ Whether the Alert will trigger immediately after the trigger limit is reached. Valid options are `"true"` and `"false"` for presence Alerts and `"false"` for absence Alerts.
- `key`: **_string (Required)_** The PagerDuty service key.
- `operator`: **_string_** _(Optional; Default: `presence`)_ Whether the Alert will trigger on the presence or absence of logs. Valid options are `presence` and `absence`.
- `terminal`: **_string_** _(Optional; Default: `"true"`)_ Whether the Alert will trigger after the `triggerinterval` if the Alert condition is met (e.g., send an Alert after 30s). Valid options are `"true"` and `"false"` for presence Alerts and `"true"` for absence Alerts.
- `triggerinterval`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for absence)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. For absence Alerts, valid options are: `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`.
- `triggerlimit`: **_integer (Required)_** Number of lines before the Alert is triggered (e.g. setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`).

### webhook_channel

`webhook_channel` supports the following arguments:

- `bodytemplate`: **_string_** _(Optional)_ JSON-formatted string for the body of the webhook. We recommend using [`jsonencode()`](https://www.terraform.io/docs/configuration/functions/jsonencode.html) to easily convert a Terraform map into a JSON string.
- `headers`: **_map<string, string>** _(Optional)_ Key-value pair for webhook request headers and header values. Example: `"MyHeader" = "MyValue"`
- `immediate`: **_string_** _(Optional; Default: `"false"`)_ Whether the Alert will trigger immediately after the trigger limit is reached. Valid options are `"true"` and `"false"` for presence Alerts and `"false"` for absence Alerts.
- `method`: **_string_** _(Optional; Default: `post`)_ Method used for the webhook request. Valid options are: `post`, `put`, `patch`, `get`, `delete`.
- `operator`: **_string_** _(Optional; Default: `presence`)_ Whether the Alert will trigger on the presence or absence of logs. Valid options are `presence` and `absence`.
- `terminal`: **_string_** _(Optional; Default: `"true"`)_ Whether the Alert will trigger after the `triggerinterval` if the Alert condition is met (e.g., send an Alert after 30s). Valid options are `"true"` and `"false"` for presence Alerts and `"true"` for absence Alerts.
- `triggerinterval`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for absence)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. For absence Alerts, valid options are: `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`.
- `triggerlimit`: **_integer (Required)_** Number of lines before the Alert is triggered (e.g. setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`).
- `url`: **_string (Required)_** The URL of the webhook.


