provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_ingestion_exclusion" "http-success" {
  title  = "HTTP 2XX"
  apps   = ["nginx", "apache"]
  query  = "response:(>=200 <300)"
  active = true
  indexonly = false
}

resource "logdna_ingestion_exclusion" "http-noise" {
  title  = "Noisy HTTP Paths"
  apps   = ["nginx", "apache"]
  query  = "robots.txt OR favicon.ico OR .well-known"
  active = true
  indexonly = true
}
