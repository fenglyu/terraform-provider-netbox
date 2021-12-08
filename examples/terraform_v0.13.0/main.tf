provider "netbox" {
  request_timeout = "4m"
}
provider "netbox" {
  api_token = "<authentication token>"
  host      = "netbox.k8s.me"
  base_path = "/api"
}

resource "netbox_available_prefixes" "gke-pods" {

  parent_prefix_id = 1
  prefix_length    = 25
  tags             = ["BasePathTest-acc", "flv"]
  vrf              = "activision"
  custom_fields {}
}

output "available_prefix" {
  value = netbox_available_prefixes.gke-pods.prefix
}

output "available_prefix_id" {
  value = netbox_available_prefixes.gke-pods.id
}
