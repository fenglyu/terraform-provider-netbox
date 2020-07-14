provider "netbox" {
  api_token = "434476c51e79b0badfad4afcd9a64b4dede1adb9"
  host      = "netbox.k8s.me"
  base_path      = "/api"
  request_timeout = "4m"
}

data "netbox_available_prefixes" "foo"{
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
  value = join("_", data.netbox_available_prefixes.foo.tags)
}