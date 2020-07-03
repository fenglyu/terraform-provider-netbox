terraform {
  required_providers {
    netbox = "~> 0.0.1"
  }
}

provider "netbox" {
  api_token = "c4a3c627b64fa514e8e0840a94c06b04eb8674d9"
  host      = "127.0.0.1:80"
  base_path      = "/api"
  // Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
  // see time.ParseDuration for time unit
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  // example schema attribute
  prefix_id = 103
 // available_prefixes {
  prefix_length = 26
  ispool = false
  status = "active"
  description = "cidr for gke-pods"
 // }
}
