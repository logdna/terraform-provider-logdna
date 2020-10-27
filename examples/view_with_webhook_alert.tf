provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_view" "my_view" {
  apps       = ["app1", "app2"]
  categories = ["Demo1", "Demo2"]
  hosts      = ["host1", "host2"]
  levels     = ["fatal", "critical"]
  name       = "Webhook Alert"
  query      = "test"
  tags       = ["host1", "host2"]
  webhook_channel {
    headers = {
      hello = "test3"
      test  = "test2"
    }
    bodytemplate = jsonencode({
      fields = {
        description = "{{ matches }} matches found for {{ name }}"
        issuetype = {
          name = "Bug"
        }
        project = {
          key = "test"
        },
        summary = "Alert From {{ name }}"
      }
    })
    immediate       = "false"
    method          = "post"
    url             = "https://yourwebhook/endpoint"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
}
