# Resource: `logdna_index_rate_alert`

Manages [LogDNA Index Rate Alert](https://docs.mezmo.com/docs/index-rate-alerts). Configuring alerts based on the index rate or retention and storage rate of your log data helps you track unusual behavior in your systems. For example, if there's a sudden spike in volume, Mezmo's Index Rate Alert feature tells you which applications or sources produced the data spike. It also shows any recently added sources. Index rate alerts can also help managers who are responsible for budgets to analyze and predict storage costs.

To get started, all you need to do is to specify a configuration and one of our currently supported alerts recipients: email, Slack, or PagerDuty.

Be aware that only one index rate alert configuration is allowed per account

## Example - Index Rate Alert

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_index_rate_alert" "config" {
  max_lines = 3
  max_z_score = 3
  threshold_alert = "separate"
  frequency = "hourly"
  enabled = true
  channels {
    email = ["test@test.com"]
    slack = ["https://slack_url/key"]
    pagerduty = ["service_key"]
  }
  webhook_channel {
    url = "https:/testurl.com"
    method = "POST"
    headers = {
      header1 = "value1"
    }
      bodytemplate = jsonencode({
      something = "something"
    })
  }
}

```

## Destroy
There is not a DELETE endpoint supported by the Index Rate Alert API. For this reason, removing the Index Rate Alert Config effectively disables it. (set enabled to false in DB)

## Import

Index Rate Alert can be imported by static ID "config", which can be found using the [Get Index Rate Alert API](https://docs.mezmo.com/log-analysis-api/ref#get-index-rate-alert):

1. Custom HTTP Headers - `servicekey: <SERVICE_KEY>` or `apikey: <SERVICE_KEY>`
```sh
curl --request GET \
     --url <API_URL>/v1/config/index-rate \
     --header 'Accept: application/json' \
     --header 'servicekey: <SERVICE_KEY>'
```
2. Basic Auth - `Authorization: Basic <encodeInBase64(credentials)>`.<br />
Credentials is a string formatted as `<username>:<password>`. Our usage here entails substituting `<SERVICE_KEY>` as the username and leaving the password blank. The colon separator should still be included in the resulting string `<SERVICE_KEY>:`
```sh
curl --request GET \
     --url <API_URL>/v1/config/index-rate \
     --header 'Accept: application/json' \
     --header 'Authorization: Basic <BASE_64_ENCODED_CREDENTIALS>'
```

```sh
terraform import logdna_index_rate_alert.config config
```

Note that only the alert channels supported by this provider will be imported.

## Argument Reference

The following arguments are supported by `logdna_alert`:

- `max_lines`: The number of lines required in order to set off the alert, type _int_
- `max_z_score`: The number of standard deviations above the 30-day average lines in order to set off the alert, type _int_
- `threshold_alert`: Set if you want alerts to be triggered if one or both of the max lines and standard deviation have been triggered or individually, type _string_ ["separate" | "both"]
- `frequency`: Notify recipients once per hour or once per day (starting from the first passing of the threshold) until the index rate declines back below the thresholds, ceasing all alerts., type _string_ ["hourly" | "daily"]
- `enabled`: (Required) Enable an existing configuration, type _boolean_

### channels

`channels` supports the following arguments:

- `email`: **_[]string_** An array of email addresses (strings) to notify
- `slack`: **_[]string_** An array of slack hook urls (strings) to notify
- `pagerduty`: **_[]string_** An array of pagerduty service integration keys (strings) to notify
