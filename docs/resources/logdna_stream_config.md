# Resource: `logdna_stream_config`

> **IBM Log Analysis and Cloud Activity Tracker users only**

Manages [LogDNA Streaming](https://ibm.github.io/cloud-enterprise-examples/log-streaming/content-overview/) configuration for an account.

## Example

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_stream_config" "config" {
  user      = var.stream_user
  password  = var.stream_password
  topic     = "example"
  brokers   = [
    "broker-1.example.org:9093",
    "broker-2.example.org:9093",
    "broker-3.example.org:9093",
  ]
}
```

## Import

Importing an existing configuration is supported:

```sh
$ terraform import logdna_stream_config.config stream
```

## Argument Reference

The following arguments are supported by `logdna_stream_config`:

- `brokers`: **[]string** _(Required)_ List of of broker URLs. 
- `topic`: **string** _(Required)_ The topic that logs will be published on.
- `user`: **string** _(Required)_ The SASL username for the connection.
- `password`: **string** _(Required)_ The SASL password for the connection.

Note that the provided brokers and credentials must be valid, and
the brokers must be reachable when the resource is created or updated.
The connection to the broker will be validated before the configuration
can be saved.
