provider "netbox" {
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  #parent_prefix = "10.0.4.0/24"
  parent_prefix_id = 502
  prefix_length    = 29
  is_pool          = false
  status           = "active"

  description = "foo"
  tags        = ["test01", "test02"]

  /*
  custom_fields  {
      helpers      = "123123"
      ipv4_acl_in  = ""
      ipv4_acl_out = ""
  }
  */
}

output "available_prefix" {
  value = netbox_available_prefixes.gke-pods.prefix
}

output "available_prefix_id" {
  value = netbox_available_prefixes.gke-pods.id
}

