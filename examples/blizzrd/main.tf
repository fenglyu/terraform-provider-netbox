provider "netbox" {
}

resource "netbox_available_prefixes" "gke-pods" {
  #parent_prefix = "10.0.4.0/24"
  parent_prefix_id = 9480
  prefix_length    = 29
  is_pool          = false
  status           = "active"

  ##vrf = 0
  tenant           = "cloud"
  role             = "Development"
  site             = "gcp" 
  description = "test/cloud/flv-test-0 || usw2-pri-gke-nodes"
  tags        = ["k8s", "gke", "gke-pods", "test01", "test02"]

  custom_fields    = {
      helpers      = ""
      ipv4_acl_in  = ""
      ipv4_acl_out = ""
  }

}

output "available_prefix" {
  value = netbox_available_prefixes.gke-pods.prefix
}

