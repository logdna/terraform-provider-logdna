# Resource: `logdna_enterprise_key`

This resource allows you to manage enterprise keys.

## Example

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
}

resource "logdna_enterprise_key" "enterprise_service_key" {
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
$ terraform apply -replace="logdna_enterprise_key.my_key"
```

## Argument Reference

This resource does not support any input arguments.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `id`: **string** The unique identifier of this key.
- `key`: **string** The actual key value.
- `name`: **string** The name of the key.
- `created`: **int** The date the key was created in Unix time milliseconds.

## Import

A key can be imported using the `id`, e.g.,

```sh
$ terraform import logdna_enterprise_key.my_key <id>
```
