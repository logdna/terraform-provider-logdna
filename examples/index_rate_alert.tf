provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_index_rate_alert" "config" {
  max_lines = 3
  max_z_score = 3
  threshold_alert = "separate"
  frequency = "hourly"
  enabled = true
  channels {
    email = ["test@test.com"]
    slack = ["https://slack_url/key"]
    pagerduty = ["service_key"]
  }
}