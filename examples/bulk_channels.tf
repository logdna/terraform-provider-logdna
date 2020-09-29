provider "logdna" {
  servicekey = "your service key goes here"
}

resource "logdna_view" "my_view" {
  name     = "Email Pagerduty and Webhook Alerts"
  query    = "test"
  apps     = ["app1", "app2"]
  levels   = ["fatal", "critical"]
  hosts    = ["host1", "host2"]
  category = ["Demo"]
  tags     = ["host1", "host2"]
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
    immediate       = "false"
    key             = "your pagerduty key goes here"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
  webhook_channel {
    headers = {
      hello = "test3"
      test  = "test2"
    }
    bodytemplate = {
      hello = "test1"
      test  = "test2"
    }
    immediate       = "false"
    method          = "post"
    url             = "https://yourwebhook/endpoint"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
}
