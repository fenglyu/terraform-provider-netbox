provider "netbox" {
  api_token = "434476c51e79b0badfad4afcd9a64b4dede1adb9"
  host      = "netbox.k8s.me"
  base_path = "/api"
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {

  parent_prefix_id = 125
  prefix_length    = 28
  is_pool          = true
  status           = "active"

  tenant = "foo"
  description = "cidr for gke-pods, Hello, My friend"
  tags        = ["test01", "test02", "test03", "test04"]
}


data "netbox_available_prefixes" "foo"{
  prefix_id = netbox_available_prefixes.gke-pods.id
  //prefix = "10.0.0.0/28"
}

output "available_prefix_json" {
  value = data.netbox_available_prefixes.foo
}

output "available_prefix_tags" {
  value = join("_", data.netbox_available_prefixes.foo.tags)
}