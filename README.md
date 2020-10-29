# Terraform Provider for LogDNA

[![CircleCI](https://circleci.com/gh/logdna/terraform-provider-logdna/tree/master.svg?style=svg)](https://app.circleci.com/pipelines/github/logdna/terraform-provider-logdna)
[![Coverage Status](https://coveralls.io/repos/github/logdna/terraform-provider-logdna/badge.svg)](https://coveralls.io/github/logdna/terraform-provider-logdna)

🚧 In public beta 🚧

[LogDNA](https://logdna.com) is a centralized log management platform. The LogDNA Provider allows organizations to manage Views and Alerts programmatically via Terraform.

The [official docs for the LogDNA terraform provider](https://registry.terraform.io/providers/logdna/logdna/latest/docs) can be found in the Terraform registry.

## Example Terraform Configuration
```
provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_view" "my_view" {
  name     = "Email PagerDuty and Webhook View-specific Alerts"
  query    = "test"
  apps     = ["app1", "app2"]
  levels   = ["fatal", "critical"]
  hosts    = ["host1", "host2"]
  categories = ["Demo1", "Demo2"]
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
    key             = "Your PagerDuty key goes here"
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
    url             = "https://yourwebhook/endpoint"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
}
```

Using the `logdna_view` resource, a user can create a View with provided `name`, `query`, `hosts`, `categories`, `tags`, `email_channel`, `pagerduty_channel`, and `webhook_channel`, delete a View with a given `viewid` or update a View using the `viewid` and `name`.

Run `terraform init`, `terraform plan`, and `terraform apply`, refresh your browser and then navigate to the UI to see your updates.

## Testing

To run the provider's test suite, add your LogDNA service key to [logdna/provider_test.go](https://github.com/logdna/terraform-provider-logdna/blob/main/logdna/provider_test.go), and then run the following commands. Your service key can be generated or retrieved from **Settings > Organization > API Keys**.

```
make test
```

Alternatively, you can run:

```
TF_ACC=1 go test ./logdna -v
```
