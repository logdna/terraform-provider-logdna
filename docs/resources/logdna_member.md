# Resource: `logdna_member`

This resource allows you to manage the team members of an organization.

## Example

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
}

resource "logdna_member" "admin_user" {
  email = "user@domain.jp.co"
  role = "admin"
}
```

## Argument Reference

The following arguments are supported:

- `email`: **string** _(Required)_ The email of the user. If a user with that email does not exist, they will be invited to join Mezmo.
- `role`: **string** _(Required)_ The role of this user. Can be one of `admin`, `member`, and `readonly`. `owner` roles can only be changed through the UI.
- `groups`: **string[]** _(Optional)_ The id of the groups the user belongs to. Defaults to an empty list.

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `email`: **string** The email of the member.
- `role`: **string** The role of the member.
- `groups`: **string[]** The groups the member belongs to.

## Import

A member can be imported using their `email`, e.g.,

```sh
$ terraform import logdna_member.user1 <email>
```
