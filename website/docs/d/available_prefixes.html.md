---
subcategory: "Available Prefixes"
layout: "netbox"
page_title: "Netbox: netbox_available_prefixes"
sidebar_current: "docs-netbox-datasource-available-prefixes-x"
description: |-
  Manages an available prefix in NETBOX.
---

# netbox\_available\_prefixes
Get information about a prefix

## Example Usage

```hcl
data "netbox_available_prefixes" "foo"{
  prefix = "10.0.0.0/28"
}
```

```hcl
resource "netbox_available_prefixes" "foo" {
  parent_prefix_id 	= 125
  prefix_length 	= 27
}

data "netbox_available_prefixes" "bar"{
  prefix_id = netbox_available_prefixes.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `prefix` - (Optional) The prefix in CIDR notation. One of `prefix` or `prefix_id` must be provided.
* `prefix_id` - (Optional) The Id of prefix. One of `prefix` or `prefix_id` must be provided.


## Attributes Reference
* `prefix`  - The available prefix in CIDR notation which is computed.
* `family`  - The Ipv4/Ipv6 family.
* `created` - The day when the prefix is create.
* `last_updated` -  The time when the prefix is last updated.
* `status`  - The status of prefix, in one of **"container", "active", "reserved", "deprecated"**.

```
* Container - A summary of child prefixes
* Active - Provisioned and in use
* Reserved - Designated for future use
* Deprecated - No longer in use
```
* `is_pool`             - Whether the prefix is pool.
* `role`                - Role of the prefix.
* `site`                - The site the prefix is assigned to.
* `tags`                - A list of network tags to attach to the prefix.
* `tenant`              - The tenant for the prefix
* `vlan`                - The vlan this prefix is on/ or related to.
* `vrf`                 - The VRF this prefix is on.
* `description`         - A brief description of this resource.
* `custom_fields`       - Customized fields for prefix

