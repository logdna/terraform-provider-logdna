# Terraform Provider for LogDNA

[![CircleCI](https://circleci.com/gh/logdna/terraform-provider-logdna/tree/master.svg?style=svg)](https://app.circleci.com/pipelines/github/logdna/terraform-provider-logdna)
[![Coverage Status](https://coveralls.io/repos/github/logdna/terraform-provider-logdna/badge.svg)](https://coveralls.io/github/logdna/terraform-provider-logdna)
[![Public Beta](https://img.shields.io/badge/-Public%20Beta-404346?style=flat)](#)

[LogDNA](https://logdna.com) is a centralized log management platform. The LogDNA Provider allows organizations to manage Views and Alerts programmatically via Terraform.

The [official docs for the LogDNA terraform provider](https://registry.terraform.io/providers/logdna/logdna/latest/docs) can be found in the Terraform registry.

Using the `logdna_view` resource, a user can create a View with a provided `name`, `query`, `hosts`, `categories`, `tags`, `email_channel`, `pagerduty_channel`, and `webhook_channel`, delete a View with a given `viewid` or update a View using the `viewid` and `name`.

Using the `logdna_alert` resource, a user can create a Preset Alert with a provided `name`, `email_channel`, `pagerduty_channel` and `webhook_channel`, delete a Preset Alert with a given `presetid` or update a Preset Alert using the `presetid` and `name`.

Run `terraform init`, `terraform plan`, and `terraform apply`, refresh your browser and then navigate to the UI to see your updates.

In addition to the examples provided below, sample .tf files can be found [here](https://github.com/logdna/terraform-provider-logdna/tree/main/examples).

## Example Terraform Configuration for Preset Alerts
```
provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_alert" "my_alert" {
  name  = "Email and PagerDuty Preset Alert"
  email_channel {    
    emails         = ["test@logdna.com"]                 
    operator       = "absence"
    timezone       = "Pacific/Samoa"
    triggerlimit   = 15                  
  }

  pagerduty_channel {
    immediate       = "false"
    key             = "Your PagerDuty API key goes here"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }
}
```

## Example Terraform Configuration for View-specific Alerts
```
provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_view" "my_view" {
  apps     = ["app1", "app2"]
  categories = ["Demo1", "Demo2"]
  hosts    = ["host1", "host2"]
  levels   = ["fatal", "critical"]
  name     = "Email PagerDuty and Webhook View-specific Alerts"
  query    = "test"
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
    key             = "Your PagerDuty API key goes here"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
  }

  webhook_channel {
    bodytemplate = jsonencode({
      hello = "test1"
      test  = "test2"
    })
    headers = {
      hello = "test3"
      test  = "test2"
    }
    immediate       = "false"
    method          = "post"
    terminal        = "true"
    triggerinterval = "15m"
    triggerlimit    = 15
    url             = "https://yourwebhook/endpoint"
  }
}
```

## Testing

To run the provider's test suite, add your LogDNA service key to [logdna/provider_test.go](https://github.com/logdna/terraform-provider-logdna/blob/main/logdna/provider_test.go), and then run the following commands. Your service key can be generated or retrieved from **Settings > Organization > API Keys**.

```
make test
```

Alternatively, you can run:

```
TF_ACC=1 go test ./logdna -v
```
