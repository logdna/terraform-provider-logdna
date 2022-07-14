# Resource: `logdna_key`

This resource allows you to manage ingestion and service keys.

## Example

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
}

resource "logdna_key" "service-key" {
  type = "service"
  name = "terraform-my_service_key"

  lifecycle {
    create_before_destroy = true
  }
}

resource "logdna_key" "ingestion-key" {
  type = "ingestion"
  name = "terraform-my_ingestion_key"

  lifecycle {
    create_before_destroy = true
  }
}
```

The `create_before_destroy` and `lifecycle` meta-argument are not required; however, these options ensure a valid key is always available when a key is being recreated. This helps avoid any disruptions in service.

~> **NOTE:** We recommend prefixing the name of your terraform resources so they can be distinguished from other resources in the UI.

## Key Rotation

This resource can be used in conjuction with automated scripts to perform automatic key rotations, e.g.,

```sh
# Run this every time you want to rotate the key
$ terraform apply -replace="logdna_key.my_key"
```

## Argument Reference

The following arguments are supported:

- `type`: **string** _(Required)_ The type of key to be used. Can be one of either `service` or `ingestion`.
- `name`: **string** _(Optional)_ A non-unique name for the key. If not supplied, a default one is generated.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id`: **string** The unique identifier of this key.
- `key`: **string** The actual key value.
- `name`: **string** The name of the key.
- `type`: **string** The type of key. Can be one of either `service` or `ingestion`.
- `created`: **int** The date the key was created in Unix time milliseconds.

## Import

A key can be imported using the `id`, e.g.,

```sh
$ terraform import logdna_key.my_key <id>
```
