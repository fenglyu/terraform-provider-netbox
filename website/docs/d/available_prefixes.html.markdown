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

## Argument Reference

The following arguments are supported:

* `self_link` - (Optional) The self link of the instance. One of `name` or `self_link` must be provided.

* `name` - (Optional) The name of the instance. One of `name` or `self_link` must be provided.

---

* `project` - (Optional) The ID of the project in which the resource belongs.
    If `self_link` is provided, this value is ignored.  If neither `self_link`
    nor `project` are provided, the provider project is used.

* `zone` - (Optional) The zone of the instance. If `self_link` is provided, this
    value is ignored.  If neither `self_link` nor `zone` are provided, the
    provider zone is used.

## Attributes Reference