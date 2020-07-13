---
layout: "netbox"
page_title: "Provider: Netbox"
sidebar_current: "docs-netbox-provider-x"
description: |-
    NetBox is an open source web application designed to help manage and document computer networks.
---

# Netbox Provider

-> Try out Terraform 0.12 with the Netbox provider! `netbox`  are 0.12-compatible.

The Netbox provider is used to configure your Netbox IP address management (IPAM).


A typical provider configuration will look something like:

```hcl
provider "netbox" {
  api_token = ""
  host      = "127.0.0.1"
  # request_timeout = "4m"
}
```


## Features and Bug Requests

* If you have a bug or feature request without an existing issue and an existing resource 
* Or existing field is working in an unexpected way
* You'd like the provider to support a new resource or field

Use one of following method to reach us,
1. Create a ghosthup repo issues
2. Leave #cloud foundation team a message in slack channel `gcp-alpha`
3. Email/Slack to maintainer *flv_blizzard.com* directly


