# Resource: `logdna_key`

This resource allows you to manage ingestion and service keys.

## Example

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
}

resource "logdna_key" "service-key" {
  type = "service"

  lifecycle {
    create_before_destroy = true
  }
}

resource "logdna_key" "ingestion-key" {
  type = "ingestion"

  lifecycle {
    create_before_destroy = true
  }
}
```

The `create_before_destroy` and `lifecycle` meta-argument are not required, but ensure a valid key is always available so there's no disruption of service.

## Key Rotation

This resource can be used in conjuction with automated scripts to perform automatic key rotations, e.g.,

```sh
# Run this every time you want to rotate the key
$ terraform apply -replace="logdna_key.my_key"
```

## Argument Reference

The following arguments are supported:

- `type`: **string** _(Required)_ The type of key to be used. Should be either `service` or `ingestion`.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id`: **string** The unique identifier of this key.
- `key`: **string** The actual key value.
- `type`: **string** The type of key.
- `created`: **int** The date the key was created in Unix time milliseconds.

## Import

A key can be imported using the `id`, e.g.,

```sh
$ terraform import logdna_key.my_key <id>
```
