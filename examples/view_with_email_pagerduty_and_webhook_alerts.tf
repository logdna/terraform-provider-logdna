provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_view" "my_view" {
  apps       = ["app1", "app2"]
  categories = ["Demo1", "Demo2"]
  hosts      = ["host1", "host2"]
  levels     = ["fatal", "critical"]
  name       = "Email PagerDuty and Webhook Alerts"
  query      = "test"
  tags       = ["host1", "host2"]
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
    key             = "your PagerDuty key goes here"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
  webhook_channel {
    headers = {
      hello = "test3"
      test  = "test2"
    }
    bodytemplate = jsonencode({
      hello = "test1"
      test  = "test2"
    })
    immediate       = "false"
    method          = "post"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
    url             = "https://yourwebhook/endpoint"
  }
}
