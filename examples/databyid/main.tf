provider "netbox" {
  //api_token = "<authentication token>"
  host      = "netbox.k8s.me"
  base_path = "/api"
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  #parent_prefix = "10.0.4.0/24"
  parent_prefix_id = 371
  prefix_length    = 29
  is_pool          = true
  status           = "active"
  role = "gcp"
  site = "se1"
  vlan = "gcp"
  #vrf  = "activision"
  tenant = "cloud"
  description = "foo"
  tags        = ["test01", "test02"]
  custom_fields    = {
    helpers      = "atetste"
    ipv4_acl_in  = "ipv4_a343434cl_in"
    ipv4_acl_out = "ipv4_aefserwcl_out"
  }
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