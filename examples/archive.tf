provider "logdna" {
  servicekey = "xxxxxxxxxxxxxxxxxxxxxxxx"
}

resource "logdna_archive" "config" {
  integration         = "ibm"
  bucket              = "example"
  endpoint            = "example.com"
  apikey              = "key"
  resourceinstanceid  = "id"
}
