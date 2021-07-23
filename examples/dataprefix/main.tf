provider "netbox" {
  api_token       = "<authentication token>"
  host            = "netbox.k8s.me"
  base_path       = "/api"
  request_timeout = "4m"
}

data "netbox_available_prefixes" "foo" {
  name   = "dataprefix"
  prefix = "10.0.0.0/28"
}


output "available_prefix" {
  value = data.netbox_available_prefixes.foo.prefix
}

output "available_prefix_id" {
  value = data.netbox_available_prefixes.foo.id
}

output "available_prefix_json" {
  value = data.netbox_available_prefixes.foo
}

output "available_prefix_tags" {
  value = element(data.netbox_available_prefixes.foo.prefixes, 0).tags
}