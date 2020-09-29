# logdna_alert Resource

Manages [LogDNA alert presets](https://docs.logdna.com/docs/alerts). To get started, all you need to do is to specify a `name` and configuration for one of our currently supported alerts- email, webhook, or pagerduty.
## Example - Basic Alert Setup

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) LogDNA Region Instance, needed for IBM-based instances
}

resource "logdna_alert" "my_alert" {
  name = "Terraform Basic Email Preset Example"
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
## Example - In-Depth Alert Setup

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) LogDNA Region Instance, needed for IBM-based instances
}

resource "logdna_alert" "my_alert" {
  name = "Terraform Alert In-Depth Example"
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
    key             = "your pagerduty service key goes here"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
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

- `name`: (Required) The name this view will be given

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
- `immediate`: (Optional) _Default: false_ Whether the alert will trigger immediately after the trigger limit is reached.
- `operator`: (Optional) _Default: Presence_ Whether the alert will trigger on the presence or absence of logs.
- `terminal`: (Optional) If the alert will trigger after the `triggerinterval` if the alert condition is met (ex. Send an alert after 30s).
- `triggerinterval`: (Optional) Interval which the alert will be looking for presence or absence of log lines. For `presence` alerts, valid values are: `30s`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h` (default `30s`). For `absence` alerts, valid values are:  `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. (default `15m`)

### webhook_channel

`webhook_channel` supports the following arguments:

- `url`: (Required): URL of the webhook
- `bodytemplate`: (Required) JSON Object for the body of the webhook
- `triggerlimit`: (Required) Number of lines before the alert is triggered. (ex. Setting a value of `10` for an `absence` alert would alert you if `10` lines were not seen in the `triggerinterval`)
- `headers`: (Optional) Key value pair for webhook request headers and header values
- `method`: (Optional) _Default: POST_ HTTP Method used for the webhook request
- `immediate`: (Optional) _Default: false_ Whether the alert will trigger immediately after the trigger limit is reached.
- `operator`: (Optional) _Default: Presence_ Whether the alert will trigger on the presence or absence of logs.
- `terminal`: (Optional) If the alert will trigger after the `triggerinterval` if the alert condition is met (ex. Send an alert after 30s).
- `triggerinterval`: (Optional) Interval which the alert will be looking for presence or absence of log lines. For `presence` alerts, valid values are: `30s`, `1m`, `5m`, `15m`, `30m`, `1h`, `6h`, `12h`, and `24h` (default `30s`). For `absence` alerts, valid values are:  `15m`, `30m`, `1h`, `6h`, `12h`, and `24h`. (default `15m`)

