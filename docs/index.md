# LogDNA Provider

[![Public Beta](https://img.shields.io/badge/-Public%20Beta-404346?style=flat)](#)

[LogDNA](https://logdna.com) is a centralized log management platform. The LogDNA Terraform Provider allows organizations to manage certain LogDNA resources (alerts, views, etc) programmatically via Terraform.

## Example Usage
```hcl
# Configure the LogDNA Provider
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_view" "http500" {
 name     = "HTTP 500s"
 query    = "response:500"
 email_channel {
    emails          = ["you@yourdomain.com"]  # Email address to send alerts to
    operator        = "presence"              # Trigger on the presence of lines
    terminal        = "true"                  # Alert at the end of the trigger interval
    triggerinterval = "15m"                   # Time window for alert (15 minutes)
    triggerlimit    = 15                      # Lines threshold for alert (15 lines)
 }
}
```

## Pre-requirements and considerations
Before using Terraform for creating resources in LogDNA, review the following notes:
- Verify Terraform is [installed](https://learn.hashicorp.com/tutorials/terraform/install-cli). The minimum supported version is 0.12.0 and can be checked by running `terraform version`.
- The configurations seen in the examples will go into a Terraform configuration file such as `main.tf`.
- Have the service key for your Organization available. To obtain the service key for your LogDNA Organization, go to the LogDNA dashboard and navigate to **Settings > Organization > API Keys** or follow this link [here](https://app.logdna.com/manage/api-keys).
- Authentication is handled via the `servicekey` parameter and can be set in the `provider` configuration section in the `.tf` file.
- When using the LogDNA Terraform provider, be aware that there is a rate limit of 50 requests per minute.
- If you do not provide a specific a `url` in the provider configuration, the URL defaults to `https://api.logdna.com` (recommended).
- If you want to create an Alert that uses PagerDuty to notify you, you will need to provide LogDNA with the [PagerDuty API key](https://support.pagerduty.com/docs/generating-api-keys#events-api-keys). To ensure that the LogDNA Dashboard properly displays the PagerDuty alert notification channel, we recommend that you first link the PagerDuty service to LogDNA via the [Dashboard UI](https://docs.logdna.com/docs/pagerduty-alert-integration) before using this plugin to create a PagerDuty Alert. You may choose to create such resources first and then link PagerDuty, but be aware that they will not work as intended until the connection is reconciled.

## Argument Reference

The following arguments are supported by the `provider` section of the `.tf` file:

- `servicekey`: **string _(Required)_** LogDNA Account Service Key. This can be generated or retrieved from Settings > Organization > API Keys.
- `url`: **string** _(Optional; Default: api.logdna.com)_ The LogDNA region URL. If you’re configuring an IBM Log Analysis with LogDNA or IBM Cloud Activity Tracker with LogDNA, you’ll need to ensure `url` is set to the [correct endpoint depending on the IBM region](https://cloud.ibm.com/docs/log-analysis?topic=log-analysis-endpoints#endpoints_api).
