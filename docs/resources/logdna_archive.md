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

_Note:_ `integration` field must be specified alongside its associated config arguments (ex: integration: "s3" must include s3_config{<args>})

- `integration`: **string _(Required)_** Archiving integration. Valid values are `ibm`, `s3`, `azblob`, `gcs`, `dos`, `swift`

### ibm_config

`ibm_config` supports the following arguments:

- `bucket`: **string _(Required)_** IBM COS storage bucket name
- `endpoint`: **string _(Required)_** IBM COS public (region) endpoint
- `apikey`: **string _(Required)_** IBM COS API key
- `resourceinstanceid`: **string _(Required)_** IBM COS instance identifier

### azblob_config

`azblob_config` supports the following arguments:

- `accountname`: **string _(Required)_** Azure Blob Storage account name
- `accountkey`: **string _(Required)_** Azure Blob Storage account access key

### gcs_config

`gcs_config` supports the following arguments:

- `bucket`: **string _(Required)_** Google Cloud Storage bucket name
- `projectid`: **string _(Required)_** Google Cloud project identifier

### dos_config

`dos_config` supports the following arguments:

- `space`: **string _(Required)_** DigitalOcean Spaces storage "bucket" name
- `endpoint`: **string _(Required)_** DigitalOcean Spaces (region) endpoint
- `accesskey`: **string _(Required)_** DigitalOcean Spaces API access key
- `secretkey`: **string _(Required)_** DigitalOcean Spaces API secret key

### swift_config

`swift_config` supports the following arguments:

- `authurl`: **string _(Required)_** OpenStack Swift authentication URL
- `expires`: **_integer (Optional)_** OpenStack Swift storage object days till expiry
- `username`: **string _(Required)_** OpenStack Swift user name
- `password`: **string _(Required)_** OpenStack Swift user password
- `tenantname`: **string _(Required)_** OpenStack Swift tenant/project/account name

Note that the provided settings must be valid. The connection to
the archiving integration will be validated before the configuration
can be saved.
