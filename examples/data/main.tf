provider "netbox" {
  base_path       = "/api"
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  #parent_prefix = "10.0.4.0/24"
  # bear in mind, if vrf is specified, the parent_prefix_id should belong to the vrf,
  # otherwise, it wont' set the available prefix's vrf to "global"(by default).
  count            = 8
  parent_prefix_id = 1
  prefix_length    = 9
  is_pool          = true
  status           = "active"
  site             = "se1"
  vlan             = "gcp"
  role             = "gcp"
  tenant           = "cloud"
  vrf              = "activision"
  description      = "foo"
  tags             = ["test01", "test02", "test04", "test07"]
  custom_fields {}
}


data "netbox_available_prefixes" "foo" {
  name = "querybyId"
  id   = netbox_available_prefixes.gke-pods[0].id
  //prefix = "10.0.0.0/28"
}

output "available_prefix_json" {
  value = data.netbox_available_prefixes.foo
}

output "available_prefix_tags" {
  value = element(data.netbox_available_prefixes.foo.prefixes, 0).tags
}
