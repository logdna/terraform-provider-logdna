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
  bucket = "example"
  endpoint = "example.com"
  apikey = "key"
  resourceinstanceid = "id"
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
  bucket = "example"
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
  accountname = "example name"
  accountkey = "example key"
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
  bucket = "example"
  projectid = "id"
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
  space = "example"
  endpoint = "example.com"
  accesskey = "key"
  secretkey = "key"
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
  authurl = "example.com"
  expires = 5
  username = "example user"
  password = "password"
  tenantname = "example"
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
- `bucket`: **string _(Optional)_** Bucket
- `endpoint`: **string _(Optional)_** Endpoint
- `apikey`: **string _(Optional)_** Api key
- `resourceinstanceid`: **string _(Optional)_** Resource instance ID
- `accountname`: **string _(Optional)_** Account name
- `accountkey`: **string _(Optional)_** Account key
- `projectid`: **string _(Optional)_** Project ID
- `space`: **string _(Optional)_** Space
- `endpoint`: **string _(Optional)_** Endpoint
- `accesskey`: **string _(Optional)_** Access key
- `secretkey`: **string _(Optional)_** Secret key
- `authurl`: **string _(Optional)_** Auth URL
- `expires`: **_integer (Optional)_** Expires
- `username`: **string _(Optional)_** Username
- `password`: **string _(Optional)_** Password
- `tenantname`: **string _(Optional)_** Tenant Name

## IBM COS Archiving

The following properties must be provided to create IBM COS Archiving: `bucket`, `endpoint`, `apikey`, `resourceinstanceid`

## AWS S3 Archiving

The following properties must be provided to create AWS S3 Archiving: `bucket`

## Azure Blob Archiving

The following properties must be provided to create Azure Blob Archiving: `accountname`, `accountkey`

## Google Cloud Services Archiving

The following properties must be provided to create GCS Archiving: `bucket`, `projectid`

## Digital Ocean Archiving

The following properties must be provided to create DOS Archiving: `space`, `endpoint`, `accesskey`, `secretkey`

## OpenStack Swift Archiving

The following properties must be provided to create Swift Archiving: `authurl`, `username`, `password`, `tenantname`. (`expires` is optional)


Note that the provided settings must be valid. The connection to
the archiving integration will be validated before the configuration
can be saved.
