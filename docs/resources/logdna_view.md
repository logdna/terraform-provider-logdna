# logdna_view Resource

Manages [LogDNA Views](https://docs.logdna.com/docs/views) as well as [View-specific Alerts](https://docs.logdna.com/docs/alerts#how-to-attach-an-alert-to-an-existing-view). To get started, specify a `name` and one of: `apps`, `hosts`, `levels`, `query` or `tags`. We currently support configuring Alerts to be sent via email, webhook, or PagerDuty.

## Example - Basic View Setup

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) LogDNA Region Instance, needed for IBM-based instances
}

resource "logdna_view" "my_view" {
  name  = "My View via Terraform"
  query = "level:debug my query"
}
```

## Example - In-Depth View Setup

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) LogDNA Region Instance, needed for IBM-based instances
}

resource "logdna_view" "my_view" {
  apps     = ["app1", "app2"]
  categories = ["Demo1", "Demo2"]
  hosts    = ["host1"]
  levels   = ["warn", "error"]
  name     = "Terraform View-specific Alert Example"
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
    key             = "Your PagerDuty service key goes here"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }

  webhook_channel {
    bodytemplate = {
      message = "Alerts from {{name}}"
    }
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

The following arguments are supported:

_Note:_ At least one of the following properties: `apps`, `hosts`, `levels`, `query`, `tags` must be specified to create a View. Unless otherwise noted, all field values are case in-sensitive.
- `apps`: _(Optional)_ Array of names of apps (each app is of type _string_) to filter the View by
- `categories`: _(Optional)_ Array of existing category names (each category is of type _string_) this View should be nested under. _Note: If the category does not exist, the View will by default be created in uncategorized_
- `hosts`: _(Optional)_ Array of names of hosts (each host is of type _string_) to filter the View by
- `levels`: _(Optional)_ Array of names of levels (each level is of type _string_) to filter the View by
- `name`: _(Required)_ Name this View will be given, type _string_
- `query`: _(Optional)_  Search query scope for the View, type _string_
- `tags`: _(Optional)_ Array of names of tags (each tag is of type _string_) to filter the View by

### email_channel

`email_channel` supports the following arguments:

- `emails`: _(Required)_ An array of email addresses (each email is of type _string_) to notify in the Alert
- `immediate`: _(Optional)_ Whether the Alert will trigger immediately after the trigger limit is reached, type _string_ (**Default: "false"**)
- `operator`: _(Optional)_ Whether the Alert will trigger on the presence or absence of logs, type _string_ (**Default: "presence"**)
- `terminal`: _(Optional)_ Whether the Alert will trigger after the `triggerinterval` if the Alert condition is met (e.g., send an Alert after 30s), type _string_
- `timezone`: _(Optional)_ Which time zone the log timestamps will be formatted in. Timezones are represented as [database time zones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones), type _string_
- `triggerinterval`: _(Optional)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid values are: 30, 1m, 5m, 15m, 30m, 1h, 6h, 12h, and 24h. For absence Alerts, valid values are: 15m, 30m, 1h, 6h, 12h, and 24h. Type _string_ (**Defaults: "30" for presence; "15m" for absence**)
- `triggerlimit`: _(Required)_ Number of lines before the Alert is triggered. (eg. Setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`), type _integer_

### pagerduty_channel

`pagerduty_channel` supports the following arguments:

- `immediate`: _(Optional)_ Whether the Alert will trigger immediately after the trigger limit is reached, type _string_ (**Default: "false"**)
- `key`: _(Required)_ PagerDuty service key, type _string_
- `operator`: _(Optional)_ Whether the Alert will trigger on the presence or absence of logs, type _string_ (**Default: "presence"**)
- `terminal`: _(Optional)_ Whether the Alert will trigger after the `triggerinterval` if the Alert condition is met (e.g., send an Alert after 30s), type _string_
- `triggerinterval`: _(Optional)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid values are: 30, 1m, 5m, 15m, 30m, 1h, 6h, 12h, and 24h. For absence Alerts, valid values are: 15m, 30m, 1h, 6h, 12h, and 24h. Type _string_ (**Defaults: "30" for presence; "15m" for absence**)
- `triggerlimit`: _(Required)_ Number of lines before the Alert is triggered. (eg. Setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`), type _integer_

### webhook_channel

`webhook_channel` supports the following arguments:

- `bodytemplate`: _(Optional)_ JSON Object for the body of the webhook, type Map of _strings_
- `headers`: _(Optional)_ Key-value pair for webhook request headers and header values, type Map of _strings_
- `immediate`: _(Optional)_ Whether the Alert will trigger immediately after the trigger limit is reached, type _string_ (**Default: "false"**)
- `method`: _(Optional)_ Method used for the webhook request, type _string_ (**Default: "POST"**)
- `operator`: _(Optional)_ Whether the Alert will trigger on the presence or absence of logs, type _string_. Two possible values: "absence" and "presence" (**Default: "presence"**)
- `terminal`: _(Optional)_ Whether the Alert will trigger after the `triggerinterval` if the Alert condition is met (e.g., send an Alert after 30s), type _string_
- `triggerinterval`: _(Optional)_ Interval which the Alert will be looking for presence or absence of log lines. For presence Alerts, valid values are: 30, 1m, 5m, 15m, 30m, 1h, 6h, 12h, and 24h. For absence Alerts, valid values are: 15m, 30m, 1h, 6h, 12h, and 24h. Type _string_ (**Defaults: "30" for presence; "15m" for absence**)
- `triggerlimit`: _(Required)_ Number of lines before the Alert is triggered. (eg. Setting a value of `10` for an `absence` Alert would alert you if `10` lines were not seen in the `triggerinterval`), type _integer_
- `url`: _(Required)_ URL of the webhook, type _string_

