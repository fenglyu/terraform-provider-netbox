terraform {
  required_providers {
    netbox = "~> 0.0.1"
  }
}

provider "netbox" {
  api_token = "c4a3c627b64fa514e8e0840a94c06b04eb8674d9"
  host      = "netbox.k8s.me"
  base_path      = "/api"
  // Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
  // see time.ParseDuration for time unit
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  // example schema attribute
  //parent_prefix = "10.0.4.0/24"
  parent_prefix_id = 31
  prefix_length = 26
  is_pool = true
  status = "active"
  description = "cidr for gke-pods"
  tags = ["k8s", "gke", "gke-pods"]

}
