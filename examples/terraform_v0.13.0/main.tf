provider "netbox" {
  request_timeout = "4m"
}

variable "project_xpn_subnets" {
  default = {
    "usw2-pri" = {
      network_cidr_prefix = 29
      region              = "us-west2"
      secondary_ip_ranges = {
        "gke-pods" = {
          is_routable         = false
          network_cidr_prefix = 24
        }
        "gke-services" = {
          is_routable         = false
          network_cidr_prefix = 27
        }
      }
      reserved_peering_ranges = {
        "master_ipv4_range" = {
          network_cidr_prefix = 28
        }
      }
    }
  }
}

locals {
  ipam_parent_prefixes = {
    corp-xpn-net = {
      global  = 9058
      rfc1112 = 9116
      role    = "Corporation"
    }

    dev-xpn-net = {
      global  = 9480
      rfc1112 = 9182
      role    = "Development"
    }

    prod-xpn-net = {
      apac     = 7852
      americas = 7850
      europe   = 7851
      rfc1112  = 9115
      role     = "Production"
    }

    test-xpn-net = {
      global  = 502
      rfc1112 = 627
      role    = "gcp"
    }


  }

  // placeholder
  xpn_is_valid_cfg      = true
  xpn_host_network_name = "test-xpn-net"

  // Build list of maps of primary ranges requiring a prefix to be generated by ipam resource
  xpn_ipam_primary_ranges_pairs = [
    for subnet_name, subnet_cfg in var.project_xpn_subnets : {
      is_routable   = lookup(subnet_cfg, "is_routable", true)
      prefix_length = lookup(subnet_cfg, "network_cidr_prefix", null)
      range_name    = subnet_name
      region        = lookup(subnet_cfg, "region", null)
    }
    if local.xpn_is_valid_cfg
    && lookup(subnet_cfg, "network_address", "") == ""
    && lookup(subnet_cfg, "network_cidr_prefix", null) != null
  ]

  // Build list of maps of nested ranges (under primary) requiring a prefix to be generated by ipam resource
  xpn_ipam_nested_ranges_pairs = flatten([
    for subnet_name, subnet_cfg in var.project_xpn_subnets : [
      for range_name, range_cfg in merge(
        lookup(subnet_cfg, "secondary_ip_ranges", {}),
        lookup(subnet_cfg, "reserved_peering_ranges", {})
        ) : {
        is_routable   = lookup(range_cfg, "is_routable", true)
        prefix_length = lookup(range_cfg, "network_cidr_prefix", null)
        range_name    = "${subnet_name}-${range_name}"
        region        = lookup(subnet_cfg, "region", null)
      }
      if lookup(range_cfg, "network_address", "") == ""
      && lookup(range_cfg, "network_cidr_prefix", null) != null
    ]
    if local.xpn_is_valid_cfg
  ])

  // Merge primary and nested ranges lists of maps into a single map for ipam resource processing
  xpn_ipam_ranges = {
    for item in concat(
      local.xpn_ipam_nested_ranges_pairs,
      local.xpn_ipam_primary_ranges_pairs
    ) : item.range_name => item
  }

  // Reconstruct a list of reserved_peering_ranges nested maps, following ipam ressource generation
  xpn_reserved_peering_ranges_pairs = flatten([
    for subnet_name, subnet_cfg in var.project_xpn_subnets : [
      for range_name, range_cfg in lookup(subnet_cfg, "reserved_peering_ranges", {}) : {
        ip_cidr_range = (
          lookup(range_cfg, "network_address", "") != "" && lookup(range_cfg, "network_cidr_prefix", null) != null
          ? "${range_cfg.network_address}/${range_cfg.network_cidr_prefix}"
          : netbox_available_prefixes.range["${subnet_name}-${range_name}"].prefix
        )
        range_name  = range_name
        subnet_name = subnet_name
      }
      if lookup(range_cfg, "network_address", "") == ""
      && lookup(range_cfg, "network_cidr_prefix", null) != null
    ]
    if local.xpn_is_valid_cfg
  ])

  // Build map of reserved_peering_ranges nested maps, so the main xpn_subnets map can be rebuilt after ipam processing
  xpn_reserved_peering_ranges = {
    for item in local.xpn_reserved_peering_ranges_pairs :
    "${item.subnet_name}__${item.range_name}" => item.ip_cidr_range
  }

  // Reconstruct a list of reserved_peering_ranges nested maps, following ipam ressource generation
  xpn_secondary_ranges_pairs = flatten([
    for subnet_name, subnet_cfg in var.project_xpn_subnets : [
      for range_name, range_cfg in lookup(subnet_cfg, "secondary_ip_ranges", {}) : {
        ip_cidr_range = (
          lookup(range_cfg, "network_address", "") != "" && lookup(range_cfg, "network_cidr_prefix", null) != null
          ? "${range_cfg.network_address}/${range_cfg.network_cidr_prefix}"
          : netbox_available_prefixes.range["${subnet_name}-${range_name}"].prefix
        )
        range_name  = range_name
        subnet_name = subnet_name
      }
      if lookup(range_cfg, "network_address", "") == ""
      && lookup(range_cfg, "network_cidr_prefix", null) != null
    ]
    if local.xpn_is_valid_cfg
  ])

  // Build map of secondary_ranges nested maps, so the main xpn_subnets map can be rebuilt after ipam processing
  xpn_secondary_ranges = {
    for item in local.xpn_secondary_ranges_pairs :
    "${item.subnet_name}__${item.range_name}" => item.ip_cidr_range
  }
}

resource "netbox_available_prefixes" "range" {
  for_each = local.xpn_ipam_ranges

  custom_fields {}
  #  description   = "${module.gcp-project.project_id} || ${each.key}"
  description   = "test 001|| ${each.key}"
  is_pool       = false
  prefix_length = lookup(each.value, "prefix_length", null)
  role          = local.ipam_parent_prefixes[local.xpn_host_network_name]["role"]
  site          = "se1"
  status        = "active"
  #tags          = each.value["tags"]
  tenant = "cloud"


  parent_prefix_id = (
    lookup(each.value, "is_routable", true)
    ? (
      ## if routable and for dev/corp xpn, pick from global prefix
      local.xpn_host_network_name == "dev-xpn-net"
      || local.xpn_host_network_name == "corp-xpn-net"
      || local.xpn_host_network_name == "test-xpn-net"
      ? local.ipam_parent_prefixes[local.xpn_host_network_name]["global"]

      ## if routable but NOT for dev/corp xpn, check if region is in americas
      : (
        length(regexall("^us-", lookup(each.value, "region", ""))) > 0
        || length(regexall("^northamerica-", lookup(each.value, "region", ""))) > 0
        || length(regexall("^southamerica-", lookup(each.value, "region", ""))) > 0
        ? local.ipam_parent_prefixes[local.xpn_host_network_name]["americas"]

        ## if routable but NOT for dev/corp xpn, and not in americas, check if region is in apac
        : (
          length(regexall("^asia-", lookup(each.value, "region", ""))) > 0
          || length(regexall("^australia-", lookup(each.value, "region", ""))) > 0
          ? local.ipam_parent_prefixes[local.xpn_host_network_name]["apac"]

          ## if routable but NOT for dev/corp xpn, and not in americas or apac, check if region is in europe
          : (
            length(regexall("^europe-", lookup(each.value, "region", ""))) > 0
            ? local.ipam_parent_prefixes[local.xpn_host_network_name]["europe"]

            ## something went wront, can't determine a parent_prefix for the provided region... :(
            : null
          )
        )
      )
    )

    ## if not routable, pick from non-routable rfc1112 prefix associated with local.xpn_host_network_name
    : local.ipam_parent_prefixes[local.xpn_host_network_name]["rfc1112"]
  )

  lifecycle {
    ignore_changes = [vrf]
    //, custom_fields]
  }
}
