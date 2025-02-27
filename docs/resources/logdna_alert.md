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
    immediate           = "true"
    key                 = "Your PagerDuty API key goes here"
    terminal            = "true"
    triggerinterval     = "15m"
    triggerlimit        = 15
    autoresolve         = true
    autoresolvelimit    = 10
    autoresolveinterval = "15m"
  }

  slack_channel {
    immediate       = "false"
    operator        = "absence"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
    url             = "https://hooks.slack.com/services/identifier/secret"
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

## Import

Preset Alerts can be imported by `presetid`, which can be found using the [List Preset Alerts API](https://docs.logdna.com/reference/list-alerts):

1. Custom HTTP Headers - `servicekey: <SERVICE_KEY>` or `apikey: <SERVICE_KEY>`
```sh
curl --request GET \
     --url <API_URL>/v1/config/presetalert \
     --header 'Accept: application/json' \
     --header 'servicekey: <SERVICE_KEY>'
```
2. Basic Auth - `Authorization: Basic <encodeInBase64(credentials)>`.<br />
Credentials is a string composed of formatted as `<username>:<password>`, our usage here entails substituting `<SERVICE_KEY>` as the username and leaving the password blank. The colon separator should still included in the resulting string `<SERVICE_KEY>:`
```sh
curl --request GET \
     --url <API_URL>/v1/config/presetalert \
     --header 'Accept: application/json' \
     --header 'Authorization: Basic <BASE_64_ENCODED_CREDENTIALS>'
```

```sh
terraform import logdna_alert.your-alert-name <presetid>
```

Note that only the alert channels supported by this provider will be imported.

## Argument Reference

The following arguments are supported by `logdna_alert`:

- `name`: (Required) The name this Preset Alert will be given, type _string_

### email_channel

`email_channel` supports the following arguments:

- `emails`: **_[]string (Required)_** An array of email addresses (strings) to notify in the Alert
- `immediate`: **_string_** _(Optional; Default: `"false"`)_ If set to `"true"`, an alert will be sent immediately after the `triggerlimit` is met. For absence alerts, this field must be `"false"`. For presence alerts, at least one of `immediate` or `terminal` must be `"true"`.
- `operator`: **_string_** _(Optional; Default: `presence`)_ Whether the Alert will trigger on the presence or absence of logs. Valid options are `presence` and `absence`.
- `terminal`: **_string_** _(Optional; Default: `"true"`)_ If set to `"true"`, an alert will be sent after both the `triggerlimit` and `triggerinterval` are met. For absence alerts, this field must be `"true"`. For presence alerts, at least one of `immediate` or `terminal` must be `"true"`.
- `timezone`: **_string_** _(Optional)_ Which time zone the log timestamps will be formatted in. Timezones are represented as [database time zones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).
- `triggerinterval`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for absence)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h`, and `25h`. For absence Alerts, valid options are: `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h`, and `25h`.
- `triggerlimit`: **_integer (Required)_** Number of lines before the Alert is triggered (e.g. setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`).

### pagerduty_channel

`pagerduty_channel` supports the following arguments:

- `immediate`: **_string_** _(Optional; Default: `"false"`)_ If set to `"true"`, an alert will be sent immediately after the `triggerlimit` is met. For absence alerts, this field must be `"false"`. For presence alerts, at least one of `immediate` or `terminal` must be `"true"`.
- `key`: **_string (Required)_** The PagerDuty service key.
- `operator`: **_string_** _(Optional; Default: `presence`)_ Whether the Alert will trigger on the presence or absence of logs. Valid options are `presence` and `absence`.
- `terminal`: **_string_** _(Optional; Default: `"true"`)_ If set to `"true"`, an alert will be sent after both the `triggerlimit` and `triggerinterval` are met. For absence alerts, this field must be `"true"`. For presence alerts, at least one of `immediate` or `terminal` must be `"true"`.
- `triggerinterval`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for absence)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h`, and `25h`. For absence Alerts, valid options are: `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h`, and `25h`.
- `triggerlimit`: **_integer (Required)_** Number of lines before the Alert is triggered (e.g. setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`).
- `autoresolve`: **_boolean_** Set to true if you want the set a condition to resolve the incident that was raised by this alert.
- `autoresolveinterval`: **_string_** _(Required if autoresolve is set to true)_ Interval of time to aggregate and check # of matched lines against the auto resolve limit. For absence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h`, and `25h`. For presence Alerts, valid options are: `15m`, `30m`, `1h`, `6h`, `12h`, `24h`, and `25h`.
- `autoresolvelimit`: **_integer_** _(Required if autoresolve is set to true)_ Specify the number of log lines that match the view's filtering and search criteria. When the number of log lines is reached, this incident will be set to resolved in PagerDuty.

### slack_channel

`slack_channel` supports the following arguments:

- `immediate`: **_string_** _(Optional; Default: `"false"`)_ If set to `"true"`, an alert will be sent immediately after the `triggerlimit` is met. For absence alerts, this field must be `"false"`. For presence alerts, at least one of `immediate` or `terminal` must be `"true"`.
- `operator`: **_string_** _(Optional; Default: `presence`)_ Whether the Alert will trigger on the presence or absence of logs. Valid options are `presence` and `absence`.
- `terminal`: **_string_** _(Optional; Default: `"true"`)_ If set to `"true"`, an alert will be sent after both the `triggerlimit` and `triggerinterval` are met. For absence alerts, this field must be `"true"`. For presence alerts, at least one of `immediate` or `terminal` must be `"true"`.
- `triggerinterval`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for absence)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h`, and `25h`. For absence Alerts, valid options are: `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h`, and `25h`.
- `triggerlimit`: **_integer (Required)_** Number of lines before the Alert is triggered (e.g. setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`).
- `url`: **_string (Required)_** The URL of the webhook for a given Slack application/integration (& channel).

### webhook_channel

`webhook_channel` supports the following arguments:

- `bodytemplate`: **_string_** _(Optional)_ JSON-formatted string for the body of the webhook. We recommend using [`jsonencode()`](https://www.terraform.io/docs/configuration/functions/jsonencode.html) to easily convert a Terraform map into a JSON string.
- `headers`: **_map<string, string>** _(Optional)_ Key-value pair for webhook request headers and header values. Example: `"MyHeader" = "MyValue"`
- `immediate`: **_string_** _(Optional; Default: `"false"`)_ If set to `"true"`, an alert will be sent immediately after the `triggerlimit` is met. For absence alerts, this field must be `"false"`. For presence alerts, at least one of `immediate` or `terminal` must be `"true"`.
- `method`: **_string_** _(Optional; Default: `post`)_ Method used for the webhook request. Valid options are: `post`, `put`, `patch`, `get`, `delete`.
- `operator`: **_string_** _(Optional; Default: `presence`)_ Whether the Alert will trigger on the presence or absence of logs. Valid options are `presence` and `absence`.
- `terminal`: **_string_** _(Optional; Default: `"true"`)_ If set to `"true"`, an alert will be sent after both the `triggerlimit` and `triggerinterval` are met. For absence alerts, this field must be `"true"`. For presence alerts, at least one of `immediate` or `terminal` must be `"true"`.
- `triggerinterval`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for absence)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h`, and `25h`. For absence Alerts, valid options are: `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, `24h`, and `25h`.
- `triggerlimit`: **_integer (Required)_** Number of lines before the Alert is triggered (e.g. setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`).
- `url`: **_string (Required)_** The URL of the webhook.
