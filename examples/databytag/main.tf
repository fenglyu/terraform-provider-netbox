provider "netbox" {
 // api_token = "<authentication token>"
  host      = "netbox.k8s.me"
  base_path      = "/api"
  request_timeout = "4m"
}

data "netbox_available_prefixes" "foo"{
  name = "lookup 001"
  prefix = "10.0.8.0/26"
}


resource "netbox_available_prefixes" "foo" {
  parent_prefix_id 	= 2
prefix_length 	= 23
is_pool          	= true
status          	= "active"
role = "gcp"
site = "se1"
vlan = "gcp"
vrf  = "activision"
tenant = "cloud"

description = "testAccDataSourceAvailablePrefixesConfigByParameters ==> foo"
tags        = ["accTag01", "accTag02", "accTag03"]
custom_fields  {}
}

resource "netbox_available_prefixes" "bar" {
  parent_prefix_id 	= 2
  prefix_length 	= 24
is_pool          	= true
status          	= "active"
role = "gcp"
site = "se1"
vlan = "gcp"
vrf  = "activision"
tenant = "cloud"

description = "testAccDataSourceAvailablePrefixesConfigByParameters ==> bar"
tags        = ["accTag01", "accTag04", "accTag05"]
custom_fields  {}
}

resource "netbox_available_prefixes" "neo" {
  parent_prefix_id 	= 2
  prefix_length 	= 21
is_pool          	= true
status          	= "active"
role = "cloudera"
site = "se1"
vlan = "gcp"
vrf  = "activision"
tenant = "cloud"

description = "testAccDataSourceAvailablePrefixesConfigByParameters ==> neo"
tags        = ["accTag06", "accTag07", "accTag08"]
custom_fields  {}
}

data "netbox_available_prefixes" "tag"{
name = "prefix_lookup_by_tag"
tag = lower("accTag01")
depends_on  = [netbox_available_prefixes.bar, netbox_available_prefixes.foo, netbox_available_prefixes.neo]
}

data "netbox_available_prefixes" "role"{
name = "prefix_lookup_by_role"
role =  "cloudera"
depends_on  = [netbox_available_prefixes.bar, netbox_available_prefixes.foo, netbox_available_prefixes.neo]
}

output "by_role" {
  value = data.netbox_available_prefixes.role
}

output "by_tags" {
  value = data.netbox_available_prefixes.tag
}
