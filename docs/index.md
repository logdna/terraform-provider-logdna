# LogDNA Provider

[LogDNA](https://logdna.com) is a centralized log management platform. The LogDNA Provider allows organizations to manage Views and Alerts programmatically via Terraform.

## Example Usage

```hcl 
# Configure the LogDNA Provider
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) LogDNA Region Instance, needed for IBM-based instances
}

resource "logdna_view" "http500" {
 name     = "HTTP 500s"
 query    = "request:500"
 email_channel {
    emails          = ["test@logdna.com"]  # Email address to send alerts to
    operator        = "presence"           # Trigger on the presence of lines
    terminal        = "true"               # Alert at the end of the trigger interval
    triggerinterval = "15m"                # Time window for alert (15 minutes)
    triggerlimit    = 15                   # Lines threshold for alert (15 lines)
 }
}
```

## Argument Reference

The following arguments are supported:

- `servicekey`: (Required) LogDNA Account Service Key. This can be generated or retrieved from Settings > Organization > API Keys
- `url`: (Optional) _Default: api.logdna.com_ The LogDNA region URL. If you’re configuring an IBM Log Analysis with LogDNA or IBM Cloud Activity Tracker with LogDNA you’ll need to ensure `url` is set to the [right endpoint depending on the IBM region](https://cloud.ibm.com/docs/Log-Analysis-with-LogDNA?topic=Log-Analysis-with-LogDNA-endpoints#endpoints_api)

