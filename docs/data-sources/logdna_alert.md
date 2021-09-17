# Data Source: `logdna_alert`

Pulls in the relevant details from existing [LogDNA Preset Alerts](https://docs.logdna.com/docs/alerts). The `logdna_alert` _data source_ remotely fetches the information from a preset alert using its ID. This preset alert may or may not be directly managed by Terraform. For the management of preset alerts as Terraform _resources_, refer to the documentation [here](../resources/logdna_alert.md).

Users can opt to reference certain preset alerts as data sources rather than definining them as resources. In this scenario, the preset alerts are not be managed as Terraform _resources_ but the data can still be accessed and referenced in other modules.

To create a `logdna_alert` data source, the `presetid` argument must be provided in its declaration to ensure the data being read remotely is correct.

## Example Usage

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_alert" "managed" {
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

  pagerduty_channel {
    immediate       = "true"
    key             = "Your PagerDuty API key goes here"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
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

# create data source by referencing the ID from an alert declared in the same config
data "logdna_alert" "managed_remote" {
  presetid = logdna_alert.managed.id
}

# create data source by pulling in the details from an alert associated with the ID
data "logdna_alert" "external_remote" {
  presetid = "xxxxxxxxxx" # the associated ID can be grabbed from the Web UI, API calls, Terraform config, etc
}

# pass in data source attributes as arguments for module(s) declared in the same config
resource "logdna_view" "test" {
  name  = "Basic View"
  query = "level:debug my query"

  email_channel = data.logdna_alert.managed_remote.email_channel
  pagerduty_channel = data.logdna_alert.external_remote.pagerduty_channel
  slack_channel = data.logdna_alert.external_remote.slack_channel
  webhook_channel = data.logdna_alert.external_remote.webhook_channel
}
```

## Argument Reference

The `logdna_alert` data source supports (and requires) the following argument:

- `presetid`: (Required) The ID associated with a specific preset alert from which we will be pulling details

## Attribute Reference

The `logdna_alert` data source exposes the same attributes supported as arguments in the managed resource. For more detailed descriptions, refer to the documentation [here](../resources/logdna_alert.md#Argument+Reference).

The following attributes (if they exist) can be referenced in the `logdna_alert` data source:

- `name`: Name of the given preset alert
- `email_channel`: List of notifications configured via email in the given preset alert
- `pagerduty_channel`: List of notifications configured via PagerDuty in the given preset alert
- `slack_channel`: List of notifications configured via Slack in the given preset alert
- `webhook_channel`: List of notifications configured via webhook(s) in the given preset alert
