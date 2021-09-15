# Resource: `logdna_archive`

Manages [LogDNA Archiving](https://docs.logdna.com/docs/archiving) configuration for an account.

## Example IBM COS Archive

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_archive" "config" {
  integration = "ibm"
  ibm_config {
    bucket = "example"
    endpoint = "example.com"
    apikey = "key"
    resourceinstanceid = "id"
  }
}

```

## Example AWS S3 Archive

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_archive" "config" {
  integration = "s3"
  s3_config {
    bucket = "example"
  }
}

```

## Example Azure Blob Storage Archive

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_archive" "config" {
  integration = "azblob"
  azblob_config {
    accountname = "example name"
    accountkey = "example key"
  }
}

```

## Example Google Cloud Services Archive

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_archive" "config" {
  integration = "gcs"
  gcs_config {
    bucket = "example"
    projectid = "id"
  }
}

```

## Example Digital Ocean Archive

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_archive" "config" {
  integration = "dos"
  dos_config {
    space = "example"
    endpoint = "example.com"
    accesskey = "key"
    secretkey = "key"
  }
}

```

## Example OpenStack Swift Archive

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_archive" "config" {
  integration = "swift"
  swift_config {
    authurl = "example.com"
    expires = 5
    username = "example user"
    password = "password"
    tenantname = "example"
  }
}

```

## Import

Importing an existing configuration is supported:

```sh
$ terraform import logdna_archive.config archive
```

## Argument Reference

The following arguments are supported by `logdna_archive`:

_Note:_ `integration` field must be specified along with the other fields required for the integration to create an archiving configuration.

- `integration`: **string _(Required)_** Archiving integration. Valid values are `ibm`, `s3`, `azblob`, `gcs`, `dos`, `swift`

### ibm_config

`ibm_config` supports the following arguments:

- `bucket`: **string _(Required)_** Bucket
- `endpoint`: **string _(Required)_** Public endpoint
- `apikey`: **string _(Required)_** API key
- `resourceinstanceid`: **string _(Required)_** Resource Instance ID

### azblob_config

`azblob_config` supports the following arguments:

- `accountname`: **string _(Required)_** Storage account name
- `accountkey`: **string _(Required)_** Storage account key

### gcs_config

`gcs_config` supports the following arguments:

- `bucket`: **string _(Required)_** Bucket
- `projectid`: **string _(Required)_** Project ID

### dos_config

`dos_config` supports the following arguments:

- `space`: **string _(Required)_** Space for Digital Ocean Spaces
- `endpoint`: **string _(Required)_** Public endpoint
- `accesskey`: **string _(Required)_** Spaces Access key
- `secretkey`: **string _(Required)_** Spaces Secret key

### swift_config

`swift_config` supports the following arguments:

- `authurl`: **string _(Required)_** Auth URL
- `expires`: **_integer (Optional)_** Days till expiry
- `username`: **string _(Required)_** Username
- `password`: **string _(Required)_** Password
- `tenantname`: **string _(Required)_** Tenant

Note that the provided settings must be valid. The connection to
the archiving integration will be validated before the configuration
can be saved.
