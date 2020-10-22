# LogDNA Provider

[![Public Beta](https://img.shields.io/badge/-Public%20Beta-404346?style=flat)](#)

[LogDNA](https://logdna.com) is a centralized log management platform. The LogDNA Provider allows organizations to manage Views and Alerts programmatically via Terraform.

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
    emails          = ["test@logdna.com"]  # Email address to send alerts to
    operator        = "presence"           # Trigger on the presence of lines
    terminal        = "true"               # Alert at the end of the trigger interval
    triggerinterval = "15m"                # Time window for alert (15 minutes)
    triggerlimit    = 15                   # Lines threshold for alert (15 lines)
 }
}
```

## Pre-requirements and considerations
Before using Terraform for creating resources in LogDNA, review the following notes:
- Verify Terraform is [installed](https://learn.hashicorp.com/tutorials/terraform/install-cli). The minimum supported version is 0.12.0 and can be checked by running `terraform version`.
- Have the service key for your Organization available. To obtain the service key for your LogDNA Organization, go to the LogDNA dashboard and navigate to **Settings > Organization > API Keys** or follow this link [here](https://app.logdna.com/manage/api-keys).
- Authentication is handled via the `servicekey` parameter and can be set in the provider configuration.
- Be aware that the underlying LogDNA Configuration API has a rate limit of 50 requests per minute; therefore, when using the LogDNA Terraform provider, there is also a limit of 50 resource operations per minute.
- If you do not provide a specific base URL in the provider configuration, the base url defaults to `https://api.logdna.com`.
- If you want to create an Alert that uses PagerDuty to notify you, you will need to provide LogDNA with the [PagerDuty API key](https://support.pagerduty.com/docs/generating-api-keys#events-api-keys). To ensure that the LogDNA Dashboard properly displays the PagerDuty alert notification channel, we recommend that you first link the PagerDuty service to LogDNA via the [Dashboard UI](https://docs.logdna.com/docs/pagerduty-alert-integration), before using the Configuration API to create a PagerDuty Alert. However, not doing so doesn't in any way prevent the use of the Configuration API to create Alerts that use PagerDuty.

## Argument Reference

The following arguments are supported:

- `servicekey`: _(Required)_ LogDNA Account Service Key. This can be generated or retrieved from Settings > Organization > API Keys. Type _string_
- `url`: _(Optional)_ The LogDNA region URL. If you’re configuring an IBM Log Analysis with LogDNA or IBM Cloud Activity Tracker with LogDNA you’ll need to ensure `url` is set to the [right endpoint depending on the IBM region](https://cloud.ibm.com/docs/Log-Analysis-with-LogDNA?topic=Log-Analysis-with-LogDNA-endpoints#endpoints_api). Type _string_ (**_Default: api.logdna.com_**)
