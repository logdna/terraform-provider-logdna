provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_category" "my_category" {
  name = "My Category"
  type = "views"
}
