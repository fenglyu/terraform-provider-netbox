provider "netbox" {
  //api_token = "<authentication token>"
  host      = "netbox.k8s.me"
  base_path = "/api"
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  #parent_prefix = "10.0.4.0/24"
  # bear in mind, if vrf is specified, the parent_prefix_id should belong to the vrf,
  # otherwise, it wont' set the available prefix's vrf to "global"(by default).
  parent_prefix_id = 2
  prefix_length    = 29
  is_pool          = true
  status           = "active"
  role = "gcp"
  site = "se1"
  vlan = "gcp"
  vrf  = "activision"
  tenant = "cloud"
  description = "foo"
  tags        = ["test01", "test02"]
  custom_fields   {
    helpers      = "atetste"
    ipv4_acl_in  = "ipv4_a343434cl_in"
    ipv4_acl_out = "ipv4_aefserwcl_out"
  }
}


data "netbox_available_prefixes" "foo"{
  name = "lookup"
  id = netbox_available_prefixes.gke-pods.id
}

output "available_prefix_json" {
  value = data.netbox_available_prefixes.foo
}

output "available_prefix_tags" {
  value = element(data.netbox_available_prefixes.foo.prefixes, 0).tags
}