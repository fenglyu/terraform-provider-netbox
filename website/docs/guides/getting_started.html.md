---
layout: "netbox"
page_title: "Getting Started with the Netbox provider"
sidebar_current: "docs-netbox-provider-guides-getting-started"
description: |-
  Getting started with the Netbox provider
---

# Getting started with the Netbox provider

~> **Note:** The provider is compatible with terraform v0.12.*, But it hasn't been tested with terraform `v0.13-dev`, It generally should work.


## Before you begin

* Login into Netbox with your account and set up API token in user profile.
* [Install Terraform](https://www.terraform.io/intro/getting-started/install.html)
and read the Terraform getting started guide that follows. This guide will
assume basic proficiency with Terraform - it is an introduction to the Netbox provider.

## Configuring the Provider

Create a Terraform config file named "main.tf". Inside, Include the following configuration:

```hcl
provider "netbox" {
  api_token = "<authentication token>"
  host      = "netbox.k8s.me"
}
```

> Required field

+ The `api_token` field should the Api token created in netbox user profile.
+ The `host` field should be the address of the netbox service.

> Optional field

+ The `base_path` field is an the base path of netbox api entry point, by default "/api".
+ The `request_timeout` field should be used to configure the API request timeout in
the format of [go time duration string](https://golang.org/src/time/format.go?#L1369)

## Carve an available prefix under a parent prefix
Resource `netbox_available_prefixes` is named following netbox's api schema, Look at the
[`netbox_available_prefixes documentation`](/docs/providers/netbox/r/available_prefixes.html)
for more configurable fields.

```hcl

resource "netbox_available_prefixes" "foo" {
  #parent_prefix = "10.0.4.0/24"
  parent_prefix_id = 9999
  prefix_length    = 29
  status           = "active"

  description = "foo"
  tags        = ["test01", "test02"]
}
```

-> In order to make the result verbose in the first time, You can add an output in the config.

```hcl
output "available_prefix" {
  value = netbox_available_prefixes.foo.prefix
}
```

## Authentication

The Netbox provider provides two means of providing credentials for authentication.
The following methods are supported, in this order, and explained below.

- Static credentials
- Environment variables
 
### Static credentials

!> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding an `api_token` and `host` in-line in the Netbox provider block:

Usage:

```hcl
provider "netbox" {
  api_token = "<authentication token>"
  host      = "netbox.k8s.me"
  base_path = "/api"
}
```

### Environment variables

You can provide your credentials via the `NETBOX_TOKEN` and
`NETBOX_HOST`, environment variables, representing your netbox api token 
and netbox server address, respectively. 
The `NETBOX_BASE_PATH` is an optional environment variable which can be ignored,
It's by default "/api". 
The `NETBOX_API_TOKEN` and `API_TOKEN` environment variables can also be used 
to represent the API Token. if applicable:

```hcl
provider "netbox" {}
```

Usage:

```sh
$ export NETBOX_TOKEN=""                                                                                                                                   ✘ 130 
$ export NETBOX_HOST=""
$ export NETBOX_BASE_PATH="/api"

$ terraform plan
```

## Provisioning your resources
With a Terraform config and with you credentials configured, it's time to provision your resources.

```sh
terraform apply
```

By now, You've gotten started using `Netbox provider` and carved an avalaible prefix in CIDR notation for you GCE/GKE network.
An expected output looks like.

```hcl
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

available_prefix = 10.0.0.0/29
available_prefix_id = 335

```