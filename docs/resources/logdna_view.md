# Resource: `logdna_view`

Manages [LogDNA Views](https://docs.logdna.com/docs/views) as well as
[View-specific Alerts](https://docs.logdna.com/docs/alerts#how-to-attach-an-alert-to-an-existing-view).
These differ from `logdna_alert` which are "preset alerts", while these are specific to
certain views.  To get started, specify a `name` and one of: `apps`, `hosts`, `levels`,
`query` or `tags`. We currently support configuring these view alerts to be sent via
email, webhook, or PagerDuty.

## Example - Basic View

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_view" "my_view" {
  name  = "Basic View"
  query = "level:debug my query"
}
```

## Example - Multi-channel View

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA Region
}

resource "logdna_view" "my_view" {
  apps     = ["app1", "app2"]
  categories = ["Demo1", "Demo2"]
  hosts    = ["host1"]
  levels   = ["warn", "error"]
  name     = "Terraform Multi-channel View"
  query    = "my query"
  tags     = ["tag1", "tag2"]

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

## Import

Views can be imported by `id`, which can be found in the URL when editing the
View in the web UI:

```sh
$ terraform import logdna_view.your-view-name <id>
```

Note that only the alert channels supported by this provider will be imported.

## Argument Reference

The following arguments are supported by `logdna_view`:

_Note:_ A `name` and at least one of the following properties: `apps`, `hosts`, `levels`,
`query`, `tags` must be specified to create a View.

- `apps`: **_string_** _(Optional)_ Array of app names to filter the View by.
- `categories`: **[]string** _(Optional)_ Array of existing category names that this View
  should be nested under. _Note: If the category does not exist, the View will by default
  be created in uncategorized_.
- `hosts`: **[]string** _(Optional)_ Array of host names to filter the View by.
- `levels`: **[]string** _(Optional)_ Array of level names to filter the View by.
- `name`: **string _(Required)_** The name of this View.
- `query`: **string** _(Optional)_  Search query for the View.
- `tags`: **[]string** _(Optional)_ Array of tag names to filter the View by.

### email_channel

`email_channel` supports the following arguments:

- `emails`: **[]string _(Required)_** An array of email addresses (strings) to notify in
  the Alert
- `immediate`: **string** _(Optional; Default: `"false"`)_ Whether the Alert will be
  triggered immediately after the trigger limit is reached. Valid options are `"true"` and
  `"false"` for presence Alerts and `"false"` for absence Alerts.
- `operator`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for
  absence)_ Whether the Alert will trigger on the presence or absence of logs. Valid
  options are `presence` and `absence`.
- `terminal`: **_string_** _(Optional; Default: `"true"`)_ Whether the Alert will trigger
  after the `triggerinterval` if the Alert condition is met (e.g., send an Alert after
  30s). Valid options are `"true"` and `"false"` for presence Alerts, and `"true"` for
  absence Alerts.
- `timezone`: **string** _(Optional)_ Which time zone the log timestamps will be formatted
  in. Timezones are represented as
  [database time zones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).
- `triggerinterval`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for
  absence)_ Interval which the Alert will be looking for presence or absence of log lines.
  For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`,
  `12h`, and `24h`. For absence Alerts, valid options are: `15m`, `30m`, `1h`, `6h`,
  `12h`, and `24h`.
- `triggerlimit`: **_integer (Required)_** Number of lines before the Alert is triggered
  (e.g. setting a value of `10` for an `absence` Alert would alert you if `10` lines were
  not seen in the `triggerinterval`).

### pagerduty_channel

`pagerduty_channel` supports the following arguments:

- `immediate`: **_string_** _(Optional; Default: `"false"`)_ Whether the Alert will be
  triggered immediately after the trigger limit is reached. Valid options are `"true"` and
  `"false"` for presence Alerts, and `"false"` for absence Alerts.
- `key`: **string _(Required)_** The service key used for PagerDuty.
- `operator`: **_string_** _(Optional; Default: `presence`)_ Whether the Alert will
  trigger on the presence or absence of logs. Valid options are `presence` and `absence`.
- `terminal`: **_string_** _(Optional; Default: `"true"`)_ Whether the Alert will trigger
  after the `triggerinterval` if the Alert condition is met (e.g. send an Alert after
  30s). Valid options are `"true"` and `"false"` for presence Alerts, and `"true"` for
  absence Alerts.
- `triggerinterval`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for
  absence)_ Interval which the Alert will be looking for presence or absence of log lines.
  For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`,
  `12h`, and `24h`. For absence Alerts, valid options are: `15m`, `30m`, `1h`, `6h`,
  `12h`, and `24h`.
- `triggerlimit`: **_integer (Required)_** Number of lines before the Alert is triggered
  (e.g. setting a value of `10` for an `absence` Alert would alert you if `10` lines were
  not seen in the `triggerinterval`).

### webhook_channel

`webhook_channel` supports the following arguments:

- `bodytemplate`: **string** _(Optional)_ JSON-formatted string for the body of the
  webhook. We recommend using
  [`jsonencode()`](https://www.terraform.io/docs/configuration/functions/jsonencode.html)
  to easily convert a Terraform map into a JSON string.
- `headers`: **_map<string, string>** _(Optional)_ Key-value pair for webhook request
  headers and header values. Example: `"MyHeader" = "MyValue"`
- `immediate`: **_string_** _(Optional; Default: `"false"`)_ Whether the Alert will
  trigger immediately after the trigger limit is reached. Valid options are `"true"` and
  `"false"` for presence Alerts, and `"false"` for absence Alerts.
- `method`: **_string_** _(Optional; Default: `post`)_ Method used for the webhook
  request. Valid options are: `post`, `put`, `patch`, `get`, `delete`.
- `operator`: **_string_** _(Optional; Default: `presence`)_ Whether the Alert will
  trigger on the presence or absence of logs. Valid options are `presence` and `absence`.
- `terminal`: **_string_** _(Optional; Default: `"true"`)_ Whether the Alert will trigger
  after the `triggerinterval` if the Alert condition is met (e.g. send an Alert after
  30s). Valid options are `"true"` and `"false"` for presence Alerts, and `"true"` for
  absence Alerts.
- `triggerinterval`: **_string_** _(Optional; Defaults: `"30"` for presence; `"15m"` for 
  absence)_ Interval which the Alert will be looking for presence or absence of log lines.
  For presence Alerts, valid options are: `30`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`,
  `12h`, and `24h`. For absence Alerts, valid options are: `15m`, `30m`, `1h`, `6h`,
  `12h`, and `24h`.
- `triggerlimit`: **_integer (Required)_** Number of lines before the Alert is triggered.
  (eg. Setting a value of `10` for an `absence` Alert would alert you if `10` lines were
  not seen in the `triggerinterval`)
- `url`: **_string (Required)_** The URL of the webhook.
