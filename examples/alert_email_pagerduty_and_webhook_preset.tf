provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_alert" "my_alert" {
  name = "Email PagerDuty and Webhook Preset Alert"
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
    key             = "Your PagerDuty service key goes here"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }

  webhook_channel {
    bodytemplate = jsonencode({
      fields = {
        description = "{{ matches }} matches found for {{ name }}"
        issuetype = {
          name = "Bug"
        }
        project = {
          key = "test"
        },
        summary = "Alert From {{ name }}"
      }
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
