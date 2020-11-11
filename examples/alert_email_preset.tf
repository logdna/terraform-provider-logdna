provider "logdna" {
  servicekey = "Your service key goes here"
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
