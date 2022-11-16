# Resource: `logdna_ingestion_exclusion`

The resource allows you to filter out logs that you don't need to store, preventing lines from being ingested
in our searchable database. You can define exclusion rules by application, hostname, and patterns within log
lines.

## Example

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
}

resource "logdna_ingestion_exclusion" "http-success" {
  title  = "HTTP 2XX"
  apps   = ["nginx", "apache"]
  query  = "response:(>=200 <300)"
  active = true
  indexonly = false
}

resource "logdna_ingestion_exclusion" "http-noise" {
  title  = "Noisy HTTP Paths"
  apps   = ["nginx", "apache"]
  query  = "robots.txt OR favicon.ico OR .well-known"
  active = true
  indexonly = true
}
```

## Argument Reference

The following arguments are supported by `logdna_ingestion_exclusion`:

_Note:_ A `title` and at least one of the following properties: `apps`, `hosts`, `query` must be specified to create this resource.

- `title`: **string** _(Optional)_ Title of this exclusion rule that will appear in the UI.
- `active`: **_bool_** _(Optional; Default: false)_ Whether the rule should be active.
- `indexonly`: **_bool_** _(Optional; Default: true)_ Live-tail and alerting will be preserved when `false`.
- `apps`: **_[]string_** _(Optional)_ Array of app names to exclude.
- `hosts`: **_[]string_** _(Optional)_ Array of hosts to exclude.
- `query`: **_string_** _(Optional)_ A search query to match lines to exclude
