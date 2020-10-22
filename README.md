# terraform-provider-logdna 

🚧 Work in progress 🚧

## Example Terraform Configuration
```
provider "logdna" {
  servicekey = "your service key goes here"
}

resource "logdna_view" "my_view" {
  name     = "email pagerduty and webhook"
  query    = "test"
  apps     = ["app1", "app2"]
  levels   = ["fatal", "critical"]
  hosts    = ["host1", "host2"]
  category = ["Demo"]
  tags     = ["tag1", "tag2"]
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
    url             = ""https://yourwebhook/endpoint""
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
}
```

Using the logdna_view resource- a user can create a view with a provided name query, hosts, category, tags, email, pagerduty and webhook channels, delete a view with a given viewid or update a view's properties.

Run go build, terraform init, terraform plan and terraform apply and then navigate to the UI to see your updates!

## Test

In order to run the provider's test suite add your LogDNA servicekey to logdna/provider_test.go and then run the following. Your servicekey can be generated or retrieved from Settings > Organization > API Keys.

```
make test
```

Alternatively, you can run:

```
TF_ACC=1 go test ./logdna -v
```
