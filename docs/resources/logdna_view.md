# logdna_view Resource

Manages [LogDNA Views](https://docs.logdna.com/docs/views) as well as [View-specific Alerts](https://docs.logdna.com/docs/alerts#how-to-attach-an-alert-to-an-existing-view). To get started, specify a `name` and one of: `apps`, `hosts`, `levels`, `query` or `tags`. We currently support configuring Alerts to be sent via email, webhook, or PagerDuty.

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
    immediate       = false
    operator        = "absence"
    terminal        = true
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

_Note:_ A `name` and at least one of the following properties: `apps`, `hosts`, `levels`, `query`, `tags` must be specified to create a View.

- `apps`: _(Optional)_ Array of names of apps (each app is of type _string_) to filter the View by
- `categories`: _(Optional)_ Array of existing category names (each category is of type _string_) this View should be nested under. _Note: If categories are not provided, the View will by default be created in uncategorized. Additionally, a user needs to have created categories in the account of interest before they are referenced in the Terraform configuration_ 
- `hosts`: _(Optional)_ Array of names of hosts (each host is of type _string_) to filter the View by
- `levels`: _(Optional)_ Array of names of levels (each level is of type _string_) to filter the View by
- `name`: _(Required)_ Name this View will be given, type _string_
- `query`: _(Optional)_  Search query scope for the View, type _string_
- `tags`: _(Optional)_ Array of names of tags (each tag is of type _string_) to filter the View by

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
