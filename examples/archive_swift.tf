provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
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
