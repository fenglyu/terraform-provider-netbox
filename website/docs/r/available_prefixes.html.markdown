---
subcategory: "Available Prefixes"
layout: "netbox"
page_title: "Netbox: netbox_available_prefixes"
sidebar_current: "docs-netbox-available-prefixes-x"
description: |-
  Manages an available prefix in NETBOX.
---

# netbox\_available\_prefixes
Carve a prefix `A prefix is an IPv4 or IPv6 network and mask expressed in CIDR notation (e.g. 192.0.2.0/24)` based
on a prarent prefix or its ID.

## Example Usage
```hcl
## create a prefix with a prefix which is specified by its id `parent_prefix_id`
resource "netbox_available_prefixes" "default" {
  parent_prefix_id = 1234
  prefix_length    = 28
  is_pool          = true
  status           = "active"
  description = "foo bar"
  tags        = ["foo", "bar"]
}
```

```hcl
## create a prefix with a prefix which is specified by `parent_prefix`
resource "netbox_available_prefixes" "default" {
  parent_prefix    = "10.1.0.0/16"
  prefix_length    = 23
  is_pool          = false
  status           = "active"
  description = "foo bar"
  tags        = ["foo", "bar"]
}
```


## Argument Reference

The following arguments are supported:
* `parent_prefix`       - (Required) Crave available prefixes under the parent_prefix.
* `parent_prefix_id`    - (Required) A UID identifying the prefix under which available prefix is craved.
* `prefix_length`       - (Required) The mask expressed in CIDR notation, E.G. 24 in 192.0.2.0/24.
* `status`              - (Optional) Each prefix can be assigned a status. It's one of statuses "container", "active", "reserved", "deprecated". Defaults to "active".
```hcl
* Container - A summary of child prefixes
* Active - Provisioned and in use
* Reserved - Designated for future use
* Deprecated - No longer in use
```

* `is_pool`             - (Optional) If enabled, NetBox will treat this prefix as a range (such as a NAT pool) wherein every IP address is valid and assignable. This logic is used for identifying available IP addresses within a prefix. If this flag is disabled, NetBox will assume that the first and last (broadcast) address within the prefix are unusable. Defaults to false.
* `role`                - (Optional) A prefix's **role** defines its function. Role assignment is optional and roles are fully customizable.
* `site`                - (Optional) The site the prefix is assigned to.
* `tags`                - (Optional) A list of network tags to attach to the instance.
* `tenant`              - (Optional) A tenant represents a discrete entity for administrative purposes.
* `vlan`                - (Optional) A isolated layber two domain this prefix is on/ or related to.
* `vrf`                 - (Optional) A VRF object in NetBox represents a virtual routing and forwarding (VRF) domain.
* `description`         - (Optional) A brief description of this resource.
* `custom_fields`       - (Optional) Custom fields can be left blank.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:
* `id`      - An identifier for the resource in integer form
* `prefix`  - The available prefix in CIDR notation which is computed
* `family`  - The Ipv4/Ipv6 family
* `created` - The day when the prefix is create
* `last_updated` -  The time when the prefix is last updated

## Import
~> **Note:** The fields `parent_prefix_id` and `vrf` cannot be imported automatically. The API doesn't return this information. If you are setting one of these fields in your config, you will need to update your state manually after importing the resource.

Prefix can be imported by its id, e.g.
```hcl
$ terraform import netbox_available_prefixes.foo 911
```
