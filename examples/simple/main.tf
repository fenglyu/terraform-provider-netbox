provider "netbox" {
  api_token = ""
  host      = "127.0.0.1:80"
  base_path      = "/api"
}

resource "ipam_prefixes_available_prefixes" "gke-pods" {
  // example schema attribute
  address = "127.0.0.1"
}
