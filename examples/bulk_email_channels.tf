provider "logdna" {
  servicekey = "your service key goes here"
}

resource "logdna_view" "my_view" {
  name     = "Two Email Alerts"
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
    triggerinterval = "15m"
    triggerlimit    = 15
    timezone        = "Pacific/Samoa"
  }
  email_channel {
    emails          = ["test@logdna.com"]
    immediate       = "false"
    operator        = "absence"
    terminal        = "true"
    timezone        = "Pacific/Samoa"
    triggerlimit    = 15
    triggerinterval = "15m"
  }
}
