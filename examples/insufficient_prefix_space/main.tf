provider "netbox" {
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  count            = 4
  parent_prefix_id = 3
  prefix_length    = 9
  tags             = ["BasePathTest-acc"]
  custom_fields {}
}

output "available_prefix" {
  value = netbox_available_prefixes.gke-pods[0].prefix
}

output "available_prefix_id" {
  value = netbox_available_prefixes.gke-pods[0].id
}
