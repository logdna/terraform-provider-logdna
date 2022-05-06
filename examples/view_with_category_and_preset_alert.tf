provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_category" "my_category" {
  name = "My Category"
  type = "views"
}

resource "logdna_alert" "my_alert" {
  name = "Email Preset Alert"
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

resource "logdna_view" "my_view" {
  apps       = ["app1", "app2"]
  categories = ["Demo1", "Demo2"]
  hosts      = ["host1", "host2"]
  levels     = ["fatal", "critical"]
  name       = "Email Alert"
  query      = "test"
  tags       = ["host1", "host2"]
  categories = [logdna_category.my_category.name]
  presetid   = logdna_alert.my_alert.id

  depends_on = ["logdna_alert.my_alert","logdna_category.my_category"]
}
