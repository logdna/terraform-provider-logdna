provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
}

resource "logdna_archive" "config" {
  integration = "s3"
  s3_config {
    bucket = "example"
  }
}
