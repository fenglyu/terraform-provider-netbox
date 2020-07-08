terraform {
  required_providers {
    netbox = "~> 0.0.1"
  }
}

provider "netbox" {
  request_timeout = "4m"
}

resource "netbox_available_prefixes" "gke-pods" {
  parent_prefix_id = 9075
  prefix_length = 29
  is_pool = false
  status = "active"
  description = "test/cloud/flv-test-0 || usw2-pri-gke-nodes"
  tags = ["k8s", "gke", "gke-pods"]
}
