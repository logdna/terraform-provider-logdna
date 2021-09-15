provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
}

resource "logdna_archive" "config" {
  integration = "gcs"
  gcs_config {
    bucket = "example"
    projectid = "id"
  }
}
