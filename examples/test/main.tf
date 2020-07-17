provider "netbox" {
}

resource "netbox_available_prefixes" "gke-test" {
  parent_prefix_id = 302
  prefix_length = 29
  is_pool          = true
  status           = "active"
  tags = ["AvailablePrefix-acc-01", "AvailablePrefix-acc-02", "AvailablePrefix-acc-03"]
}