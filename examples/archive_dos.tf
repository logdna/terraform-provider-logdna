provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
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
