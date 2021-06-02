provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_view" "kube_probes" {
  name  = "kube_probes"
  query = "agent:kube-probe"
  pagerduty_channel {
    immediate         = "true"
    triggerinterval   = "30m"
    triggerlimit      = 30
    key               = "Your PagerDuty API key goes here"
  }
}

resource "logdna_view" "e2e_metrics_reporter" {
  name  = "e2e_metrics_reporter"
  query = "app:e2e-latency-metrics-reporter"
  apps  = ["reporting", "metrics"]
}

resource "logdna_view" "not_info_logs" {
  name  = "not_info_logs"
  query = "-level:info"
  webhook_channel {
    bodytemplate = jsonencode({
      fields = {
        description = "{{ matches }} matches found for {{ name }}"
        headers = {
          x-reported-from = "terraform"
        }
        project = {
          key = "not_info_logs"
        },
        summary = "Non-info log entries from {{ name }}"
      }
    })
    immediate       = "false"
    terminal        = "true"
    method          = "post"
    url             = "https://ourwebhook/log_responses/not_info"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
}

resource "logdna_view" "exception_errors" {
  name    = "exception_errors"
  query   = "ClassCastException OR NullPointerException"
  levels  = ["fatal", "critical"]
}
