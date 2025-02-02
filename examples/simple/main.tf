provider "netbox" {
  request_timeout = "4m"
}



resource "netbox_available_prefixes" "gke-pods" {
  parent_prefix = "240.0.0.0/4"
  //parent_prefix_id = 2
  prefix_length = 7
  tags          = ["BasePathTest-acc"]
  /*
    custom_fields   {
      helpers      = "sdfdf"

      ipv4_acl_in  = "ipv4_a343434cl_in"
      ipv4_acl_out = "ipv4_aefserwcl_out"

  }
  */
  custom_fields {}
}

output "available_prefix" {
  value = netbox_available_prefixes.gke-pods.prefix
}

output "available_prefix_id" {
  value = netbox_available_prefixes.gke-pods.id
}

