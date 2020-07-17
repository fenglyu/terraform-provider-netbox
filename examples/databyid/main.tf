provider "netbox" {
  //api_token = "<authentication token>"
  host      = "netbox.k8s.me"
  base_path = "/api"
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  #parent_prefix = "10.0.4.0/24"
  parent_prefix_id = 302
  prefix_length    = 29
  is_pool          = false
  status           = "active"
  role = "gcp"
  site = "se1"
  vlan = "gcp"
  #vrf  = "activision"
  tenant = "cloud"
  description = "foo"
  tags        = ["test01", "test02"]
  custom_fields    = {
    helpers      = "helpers"
    ipv4_acl_in  = "ipv4_acl_in"
    ipv4_acl_out = "ipv4_acl_out"
    number       = 123
    required     = true
    test_url     = "https://www.microsoft.com"
    color        = "red"
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