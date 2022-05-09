# Resource: `logdna_stream_exclusion`

> **IBM Log Analysis and Cloud Activity Tracker users only**

Manages exclusion rules for [LogDNA Streaming](https://ibm.github.io/cloud-enterprise-examples/log-streaming/content-overview/).
Stream exclusion rules define the applications, hostnames, and patterns within
log lines that should exclude a given line from the stream.

## Example

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_stream_exclusion" "http-success" {
  title  = "HTTP 2XX"
  apps   = ["nginx", "apache"]
  query  = "response:(>=200 <300) request:*"
  active = true
}

resource "logdna_stream_exclusion" "http-noise" {
  title  = "Noisy HTTP Paths"
  apps   = ["nginx", "apache"]
  query  = "robots.txt OR favicon.ico OR .well-known"
  active = true
}
```

## Import

Stream Exclusions can be imported by `id`, which can be found
in the URL when editing the Stream Exclusion in the web UI:

Stream Exclusions can be imported by `id`, which can be found using the List Stream Exclusions API:

1. Custom HTTP Headers - `servicekey: <SERVICE_KEY>` or `apikey: <SERVICE_KEY>`
```sh
curl --request GET \
     --url <API_URL>/v1/config/stream/exclusions \
     --header 'Accept: application/json' \
     --header 'servicekey: <SERVICE_KEY>'
```
2. Basic Auth - `Authorization: Basic <encodeInBase64(credentials)>`.<br />
Credentials is a string composed of formatted as `<username>:<password>`, our usage here entails substituting `<SERVICE_KEY>` as the username and leaving the password blank. The colon separator should still included in the resulting string `<SERVICE_KEY>:`
```sh
curl --request GET \
     --url <API_URL>/v1/config/stream/exclusions \
     --header 'Accept: application/json' \
     --header 'Authorization: Basic <BASE_64_ENCODED_CREDENTIALS>'
```

```sh
$ terraform import logdna_stream_exclusion.your-rule-name <id>
```

## Argument Reference

The following arguments are supported by `logdna_stream_exclusion`:

_Note:_ A `title` and at least one of the following properties: `apps`, `hosts`, `query` must be specified to create this resource.

- `title`: **string** _(Optional)_ Title of this exclusion rule that will appear in the UI.
- `active`: **_bool_** _(Optional; Default: false)_ Whether the rule should be active.
- `apps`: **_[]string_** _(Optional)_ Array of app names to exclude.
- `hosts`: **_[]string_** _(Optional)_ Array of hosts to exclude.
- `query`: **_string_** _(Optional)_ A search query to match lines to exclude
