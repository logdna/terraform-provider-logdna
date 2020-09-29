provider "logdna" {
  servicekey = "your service key goes here"
}

resource "logdna_view" "my_view" {
  name     = "Email Alert"
  query    = "test"
  apps     = ["app1", "app2"]
  levels   = ["fatal", "critical"]
  hosts    = ["host1", "host2"]
  category = ["Demo"]
  tags     = ["host1", "host2"]
  email_channel {
    emails          = ["test@logdna.com"]
    operator        = "absence"
    terminal        = "true"
    timezone        = "Pacific/Samoa"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
}
