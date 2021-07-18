provider "logdna" {
  servicekey = "Your service key goes here"
}

resource "logdna_stream_config" "config" {
  user      = var.stream_user
  password  = var.stream_password
  topic     = "example"
  brokers = [
    "broker-1.example.org:9093",
    "broker-2.example.org:9093",
    "broker-3.example.org:9093",
  ]
}
