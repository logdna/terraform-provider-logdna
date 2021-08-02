# Slack Alerts can be configured by using the Webhook channel in conjunction with
# a channel-specific webhook URL
# https://slack.com/help/articles/115005265063-Incoming-webhooks-for-Slack

provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_view" "my_view_with_slack_alert" {
  name       = "Slack Alert"
  query      = "test"

  webhook_channel {
    immediate       = "false"
    method          = "post"
    # Your unique Slack webhook URL
    url             = "https://hooks.slack.com/services/XXXXX/XXXXX/XXXXX"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
}
