provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_view" "my_view" {
  apps       = ["app1", "app2"]
  categories = ["Demo1", "Demo2"]
  hosts      = ["host1", "host2"]
  levels     = ["fatal", "critical"]
  name       = "Email Alert"
  query      = "test"
  tags       = ["host1", "host2"]
  email_channel {
    emails          = ["test@logdna.com"]
    operator        = "absence"
    terminal        = "true"
    timezone        = "Pacific/Samoa"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
}
