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
  name = "prefix_lookup"
  prefix = "10.0.0.0/28"
}
```

```hcl
resource "netbox_available_prefixes" "foo" {
  parent_prefix_id 	= 125
  prefix_length 	= 27
}

data "netbox_available_prefixes" "bar"{
  name = "prefix_lookup"
  prefix_id = netbox_available_prefixes.foo.id
}
```

## Query with certain tag or role 
```hcl
resource "netbox_available_prefixes" "foo" {
  ...
  tags        = ["datasource-%{random_suffix}-accTag01", "datasource-AvailablePrefix-accTag02", "datasource-AvailablePrefix-accTag03"]
  custom_fields  {}
}

resource "netbox_available_prefixes" "bar" {
  ...
  tags        = ["datasource-%{random_suffix}-accTag01", "datasource-AvailablePrefix-accTag04", "datasource-AvailablePrefix-accTag05"]
  custom_fields  {}
}

resource "netbox_available_prefixes" "neo" {
  ...
  tags        = ["datasource-%{random_suffix}-accTag06", "datasource-AvailablePrefix-accTag07", "datasource-AvailablePrefix-accTag08"]
  custom_fields  {}
}

data "netbox_available_prefixes" "tag"{
  name = "prefix_lookup_by_tag"
  tag = lower("datasource-%{random_suffix}-accTag01")
  depends_on  = [netbox_available_prefixes.bar, netbox_available_prefixes.foo, netbox_available_prefixes.neo]
}

data "netbox_available_prefixes" "role"{
  name = "prefix_lookup_by_role"
  role =  lower("cloudera")
  depends_on  = [netbox_available_prefixes.bar, netbox_available_prefixes.foo, netbox_available_prefixes.neo]
}
```

## Argument Reference

The following arguments are supported:
* `name`   - (Required) A unique dedicated name for the data resource. 
* `prefix` - (Optional) The prefix in CIDR notation. One of `prefix` or `prefix_id` must be provided.
* `id` - (Optional) The Id of prefix. One of `prefix` or `id` must be provided.

## Other arguments also supported in pfix query includes
```
  is_pool
  tenant_group_id
  tenant_group
  tenant_id
  tenant
  q
  family
  prefix
  within
  within_include
  contains
  mask_length
  vrf_id
  vrf
  region_id
  region
  site_id
  site
  vlan_id
  vlan_vid
  role_id
  role
  status
  tag
  id__n
  id__lte
  id__lt
  id__gte
  id__gt
  tenant_group_id__n
  tenant_group__n
  tenant_id__n
  tenant__n
  vrf_id__n
  vrf__n
  region_id__n
  region__n
  site_id__n
  site__n
  vlan_id__n
  role_id__n
  role__n
  status__n
  tag__n
```

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
---
The `custom_fields` field might include fields,
* `helpers` - (Optional) Blizzard Customized Field.
* `ipv4_acl_in` - (Optional) Blizzard Customized Field.
* `ipv4_acl_out` - (Optional) Blizzard Customized Field.

