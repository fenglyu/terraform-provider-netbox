provider "netbox" {
  api_token = "434476c51e79b0badfad4afcd9a64b4dede1adb9"
  host      = "netbox.k8s.me"
  base_path = "/api"

  // Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
  // see time.ParseDuration for time unit
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  // example schema attribute
  //parent_prefix = "10.0.4.0/24"
  parent_prefix_id = 125
  prefix_length    = 27
  is_pool          = true
  status           = "active"

  #vrf = 0
  tenant = 1
  description = "cidr for gke-pods, Hello, My friend"
  tags        = ["k8s", "gke", "gke-pods", "test01", "test02"]
}

output "available_prefix" {
  value = netbox_available_prefixes.gke-pods.prefix
}

output "available_prefix_id" {
  value = netbox_available_prefixes.gke-pods.id
}

