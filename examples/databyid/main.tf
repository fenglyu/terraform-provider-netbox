provider "netbox" {
  //api_token = "<authentication token>"
 // host      = "netbox.k8s.me"
  base_path = "/api"
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  #parent_prefix = "10.0.4.0/24"
  # bear in mind, if vrf is specified, the parent_prefix_id should belong to the vrf,
  # otherwise, it wont' set the available prefix's vrf to "global"(by default).
  parent_prefix_id = 249
  prefix_length    = 29
  is_pool          = true
  status           = "active"
  site = "hgh3"
  vlan = "HGH3A OS"
  role = "Production"
  description = "foo"
  tags        = ["test01", "test02"]
  custom_fields {}
}


data "netbox_available_prefixes" "foo"{
  name = "querybyId"
  id = netbox_available_prefixes.gke-pods.id
  //prefix = "10.0.0.0/28"
}

output "available_prefix_json" {
  value = data.netbox_available_prefixes.foo
}

output "available_prefix_tags" {
  value = element(data.netbox_available_prefixes.foo.prefixes, 0).tags
}
