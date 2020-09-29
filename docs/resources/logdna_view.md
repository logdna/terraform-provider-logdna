# logdna_view Resource

Manages [LogDNA Views](https://docs.logdna.com/docs/views) as well as [View-Specific Alerts](https://docs.logdna.com/docs/alerts#how-to-attach-an-alert-to-an-existing-view). To get started, all you need to do is to specify a `name` and one of: `query`, `apps`, `levels`, `hosts`, or `tags`. We currently support configuring alerts to be sent via email, webhook, or pagerduty.

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
  name     = "Terraform Alert Specific View Example"
  query    = "my query"
  apps     = ["app1", "app2"]
  levels   = ["warn", "error"]
  hosts    = ["host1"]
  category = ["Demo"]
  tags     = ["tag1", "tag2"]

  email_channel {
    emails          = ["test@logdna.com"]
    triggerlimit    = 15
    immediate       = "false"
    operator        = "absence"
    terminal        = "true"
    triggerinterval = "15m"
    timezone        = "Pacific/Samoa"
  }
  
  pagerduty_channel {
    key             = "your pagerduty service key goes here"
    triggerlimit    = 15
    immediate       = "false"
    terminal        = "true"
    triggerinterval = "15m"
  }

  webhook_channel {
    url             = "https://yourwebhook/endpoint"
    triggerlimit    = 15
    headers = {
      "Authentication" = "auth_header_value"
      "HeaderTwo"      = "ValueTwo"
    }
    bodytemplate = {
      message = "Alerts from {{name}}"
    }
    immediate       = "false"
    method          = "post"
    terminal        = "true"
    triggerinterval = "15m"
  }
}
```

## Argument Reference

The following arguments are supported:

_Note:_ At least one of the following properties: `query`, `host`, `app`, `tag`, `level` must be specified to create a view.

- `name`: (Required) The name this view will be given
- `query`: (Optional) The search query scope for the view
- `category`: (Optional) An array of existing category names this view should be nested under. Not case sensitive. Note: if the category does not exist- the view will by default be created in uncategorized
- `hosts`: (Optional) Array of names of hosts to filter the view by
- `apps`: (Optional) Array of names of apps to filter the view by
- `levels`: (Optional) Array of names of levels to filter the view by
- `tags`: (Optional) Array of names of tags to filter the view by

### email_channel

`email_channel` supports the following arguments:

- `emails`: (Required) An array of emails to notify in the alert
- `triggerlimit`: (Required) Number of lines before the alert is triggered. (ex. Setting a value of `10` for an `absence` alert would alert you if `10` lines were not seen in the `triggerinterval`)
- `timezone`: (Optional) Which time zone the log timestamps will be formatted in.
- `immediate`: (Optional) _Default: false_ Whether the alert will trigger immediately after the trigger limit is reached.
- `operator`: (Optional) _Default: Presence_ Whether the alert will trigger on the presence or absence of logs.
- `terminal`: (Optional) If the alert will trigger after the `triggerinterval` if the alert condition is met (ex. Send an alert after 30s).
- `triggerinterval`: (Optional) Interval which the alert will be looking for presence or absence of log lines. For `presence` alerts, valid values are: `30s`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h` (default `30s`). For `absence` alerts, valid values are:  `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. (default `15m`)

### pagerduty_channel

`pagerduty_channel` supports the following arguments:

- `key`: (Required) Pagerduty service key
- `triggerlimit`: (Required) Number of lines before the alert is triggered. (ex. Setting a value of `10` for an `absence` alert would alert you if `10` lines were not seen in the `triggerinterval`)
- `immediate`: (Optional) _Default: false_ Whether the alert will trigger immediately after the trigger limit is reached
- `operator`: (Optional) _Default: Presence_ Whether the alert will trigger on the presence or absence of logs
- `terminal`: (Optional) If the alert will trigger after the `triggerinterval` if the alert condition is met (ex. Send an alert after 30s).
- `triggerinterval`: (Optional) Interval which the alert will be looking for presence or absence of log lines. For `presence` alerts, valid values are: `30s`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h` (default `30s`). For `absence` alerts, valid values are:  `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. (default `15m`)

### webhook_channel

`webhook_channel` supports the following arguments:

- `url`: (Required): URL of the webhook
- `bodytemplate`: (Required) JSON Object for the body of the webhook
- `triggerlimit`: (Required) Number of lines before the alert is triggered. (ex. Setting a value of `10` for an `absence` alert would alert you if `10` lines were not seen in the `triggerinterval`)
- `headers`: (Optional) Key value pair for webhook request headers and header values
- `method`: (Optional) _Default: POST_ HTTP Method used for the webhook request
- `immediate`: (Optional) _Default: false_ Whether the alert will trigger immediately after the trigger limit is reached
- `operator`: (Optional) _Default: Presence_ Whether the alert will trigger on the presence or absence of logs
- `terminal`: (Optional) If the alert will trigger after the `triggerinterval` if the alert condition is met (ex. Send an alert after 30s)
- `triggerinterval`: (Optional) Interval which the alert will be looking for presence or absence of log lines. For `presence` alerts, valid values are: `30s`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h` (default `30s`). For `absence` alerts, valid values are:  `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. (default `15m`)

