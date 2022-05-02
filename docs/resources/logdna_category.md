# Resource: `logdna_category`

Manages [LogDNA Categories](https://docs.logdna.com/reference/getting-started-with-the-configuration-api). Categories are designed to organize views, boards and screens. Categories can be created standalone and then attached to any views, boards, or screens.

To get started, all you need to do is to specify one of supported `type`: views, boards, or screens and `name`.

## Example - Basic Category

```hcl
provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
  url = "https://api.logdna.com" # (Optional) specify a LogDNA region
}

resource "logdna_category" "my_category" {
  type = "views"
  name = "My Category via Terraform"
}
```

## Import

Preset Categories can be imported by `type` and `id`:

```sh
terraform import logdna_category.your-category-name <type>:<id>
```