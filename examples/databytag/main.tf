resource "netbox_available_prefixes" "foo" {
  parent_prefix_id = 1
  prefix_length    = 24
  is_pool          = true
  status           = "active"
  role             = "gcp"
  site             = "se1"
  vlan             = "gcp"
  vrf              = "activision"
  tenant           = "cloud"

  description = "testAccDataSourceAvailablePrefixesConfigByParameters ==> foo"
  tags        = ["datasource-adc-accTag01", "datasource-AvailablePrefix-accTag02", "datasource-AvailablePrefix-accTag03"]
  custom_fields {}
}

resource "netbox_available_prefixes" "bar" {
  parent_prefix_id = 1
  prefix_length    = 23
  is_pool          = true
  status           = "active"
  role             = "gcp"
  site             = "se1"
  vlan             = "gcp"
  vrf              = "activision"
  tenant           = "cloud"

  description = "testAccDataSourceAvailablePrefixesConfigByParameters ==> bar"
  tags        = ["datasource-adc-accTag01", "datasource-AvailablePrefix-accTag04", "datasource-AvailablePrefix-accTag05"]
  custom_fields {}
}

resource "netbox_available_prefixes" "neo" {
  parent_prefix_id = 1
  prefix_length    = 26
  is_pool          = true
  status           = "active"
  role             = "cloudera"
  site             = "se1"
  vlan             = "gcp"
  vrf              = "activision"
  tenant           = "cloud"

  description = "testAccDataSourceAvailablePrefixesConfigByParameters ==> neo"
  tags        = ["datasource-adc-accTag06", "datasource-AvailablePrefix-accTag07", "datasource-AvailablePrefix-accTag08"]
  custom_fields {}
}

data "netbox_available_prefixes" "tag" {
  name       = "prefix_lookup_by_tag"
  tag        = lower("datasource-adc-accTag01")
  depends_on = [netbox_available_prefixes.bar, netbox_available_prefixes.foo, netbox_available_prefixes.neo]
}

data "netbox_available_prefixes" "role" {
  name       = "prefix_lookup_by_role"
  role       = lower("cloudera")
  depends_on = [netbox_available_prefixes.bar, netbox_available_prefixes.foo, netbox_available_prefixes.neo]
}

output "available_prefix_tags" {
  value = element(data.netbox_available_prefixes.tag.prefixes, 0).tags
}

output "available_prefix_tags_values" {
 value = data.netbox_available_prefixes.tag.prefixes
}

output "available_prefix_role" {
  value = element(data.netbox_available_prefixes.role.prefixes, 0).tags
}

output "available_prefix_role_vlaues" {
  value = data.netbox_available_prefixes.role.prefixes
}