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

Categories can be imported by `type` and `id`, which can be found using the [List Categories API](https://docs.logdna.com/reference/list-categories):

1. Custom HTTP Headers - `servicekey: <SERVICE_KEY>` or `apikey: <SERVICE_KEY>`
```sh
curl --request GET \
     --url <API_URL>/v1/config/categories/<CATEGORY_TYPE> \
     --header 'Accept: application/json' \
     --header 'servicekey: <SERVICE_KEY>'
```
2. Basic Auth - `Authorization: Basic <encodeInBase64(credentials)>`.<br />
Credentials is a string composed of formatted as `<username>:<password>`, our usage here entails substituting `<SERVICE_KEY>` as the username and leaving the password blank. The colon separator should still included in the resulting string `<SERVICE_KEY>:`
```sh
curl --request GET \
     --url <API_URL>/v1/config/categories/<CATEGORY_TYPE> \
     --header 'Accept: application/json' \
     --header 'Authorization: Basic <BASE_64_ENCODED_CREDENTIALS>'
```

```sh
terraform import logdna_category.your-category-name <type>:<id>
```

## Argument Reference

The following arguments are supported by `logdna_category`:

- `name`: **string (Required)** The name this Category will be given
- `type`: **string (Required)** The type this Category belongs to, valid options are: `views`, `boards`, `screens`

