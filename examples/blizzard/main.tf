provider "netbox" {
}

resource "netbox_available_prefixes" "gke-pods" {
  parent_prefix = "10.0.0.0/8"
  //parent_prefix_id = 1
  prefix_length = 26
  is_pool       = false
  status        = "active"

  description = "test/cloud/flv-test-0 || usw2-pri-gke-nodes"
  tags        = ["k8s", "gke", "gke-pods", "test01", "test02"]

  custom_fields {}

}

output "available_prefix" {
  value = netbox_available_prefixes.gke-pods.prefix
}

