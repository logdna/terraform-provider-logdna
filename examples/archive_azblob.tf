provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
}

resource "logdna_archive" "config" {
  integration = "azblob"
  azblob_config {
    accountname = "example name"
    accountkey = "example key"
  }
}
