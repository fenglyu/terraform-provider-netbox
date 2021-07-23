provider "netbox" {
}

resource "netbox_available_prefixes" "foo" {
  parent_prefix_id = 627
  prefix_length    = 29
  is_pool          = true
  status           = "active"
  tags             = ["AvailablePrefix-acc-01", "AvailablePrefix-acc-02", "AvailablePrefix-acc-03"]
  vrf              = "activision"

  custom_fields { /*
     helpers      = "abcd"
*/
    ipv4_acl_in = "sdf"

    ipv4_acl_out = "sdfd"

  }

}


output "available_prefix" {
  value = netbox_available_prefixes.foo.prefix
}

output "available_prefix_id" {
  value = netbox_available_prefixes.foo.id
}
