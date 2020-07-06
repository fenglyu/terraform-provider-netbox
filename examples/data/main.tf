provider "netbox" {
  api_token = "c4a3c627b64fa514e8e0840a94c06b04eb8674d9"
  host      = "netbox.k8s.me"
  base_path      = "/api"
  request_timeout = "4m"
}

data "netbox_available_prefixes" "gke-pods"{
  prefix = "10.0.4.0/26"
}
