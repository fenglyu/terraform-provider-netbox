---
layout: "netbox"
page_title: "Getting Started with the Netbox provider"
sidebar_current: "docs-netbox-provider-developing-guide"
description: |-
    Netbox provider developing guide
---

# Run Acceptance Tests

## Set test required environment variables

```bash
export NETBOX_TOKEN="a30439d5093375b36"
export NETBOX_HOST="netbox.k8s.me"
export NETBOX_BASE_PATH="/api"
# 10.0.0.0/8 (VRF: Global), prefix id is 1
export NETBOX_PARENT_PREFIX_ID=1
# 240.0.0.0/4 (VRF: activision (20183) (20183)), prefix id is 16
export NETBOX_PARENT_PREFIX_WITH_VRF_ID=16

```

## Environment Variables Reference

The following variables are required:
+ `NETBOX_TOKEN`       - (Required) The Api token created in netbox user profile.
+ `NETBOX_HOST`        - (Required) The host address of netbox service. Default "localhost:8000".
+ `NETBOX_BASE_PATH`   - (Optional) The base path of netbox api entry point, by default "/api".
+ `NETBOX_PARENT_PREFIX_ID`       - (Required) The Id of a prefix which the acceptance test will create resources on.
+ `NETBOX_PARENT_PREFIX_WITH_VRF_ID`     - (Required) The Id of a prefix with within a vrf which acceptance test will create resources on.


## Running acc test
```bash
% make testacc                                                                                                                                                                                     ✘ 130 
==> Checking source code against gofmt...
==> Checking that code complies with gofmt requirements...
TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(go list ./...) -v -timeout 240m -ldflags="-X=github.com/fenglyu/terraform-provider-netbox/version.ProviderVersion=acc"
?       github.com/fenglyu/terraform-provider-netbox    [no test files]
2020/08/04 23:11:47 [INFO] Instantiating http client for host netbox.k8s.me and path /api
2020/08/04 23:11:47 [INFO] Instantiating http client for host netbox.k8s.me and path /api
=== RUN   TestLoadAndValidate
2020/08/04 23:11:48 [INFO] Instantiating http client for host netbox.k8s.me and path /api
--- PASS: TestLoadAndValidate (0.04s)
=== RUN   TestApiAccessTestInHttps
--- PASS: TestApiAccessTestInHttps (0.04s)
=== RUN   TestAccDataSourceAvailablePrefixes_basic
=== PAUSE TestAccDataSourceAvailablePrefixes_basic
=== RUN   TestProvider
--- PASS: TestProvider (0.00s)
=== RUN   TestAccProviderBasePath_setBasePath
=== PAUSE TestAccProviderBasePath_setBasePath
=== RUN   TestAccProviderBasePath_setInvalidBasePath
=== PAUSE TestAccProviderBasePath_setInvalidBasePath
=== RUN   TestAccAvaliablePrefixes_basic
=== PAUSE TestAccAvaliablePrefixes_basic
=== RUN   TestAccAvaliablePrefixes_basic1
=== PAUSE TestAccAvaliablePrefixes_basic1
=== RUN   TestAccAvaliablePrefixes_basic2
=== PAUSE TestAccAvaliablePrefixes_basic2
=== RUN   TestAccAvaliablePrefixes_EmptyCustomFields
=== PAUSE TestAccAvaliablePrefixes_EmptyCustomFields
=== CONT  TestAccDataSourceAvailablePrefixes_basic
=== CONT  TestAccAvaliablePrefixes_basic1
=== CONT  TestAccProviderBasePath_setInvalidBasePath
=== CONT  TestAccAvaliablePrefixes_EmptyCustomFields
=== CONT  TestAccAvaliablePrefixes_basic
=== CONT  TestAccProviderBasePath_setBasePath
=== CONT  TestAccAvaliablePrefixes_basic2
--- PASS: TestAccProviderBasePath_setInvalidBasePath (0.03s)
testAccCheckAvailablePrefixesDestroyProducer:  There is not prefix with ID 64
--- PASS: TestAccAvaliablePrefixes_basic (5.37s)
testAccCheckAvailablePrefixesDestroyProducer:  There is not prefix with ID 65
--- PASS: TestAccProviderBasePath_setBasePath (6.95s)
testAccCheckAvailablePrefixesDestroyProducer:  There is not prefix with ID 66
--- PASS: TestAccAvaliablePrefixes_basic2 (8.01s)
testAccCheckAvailablePrefixesDestroyProducer:  There is not prefix with ID 67
--- PASS: TestAccAvaliablePrefixes_basic1 (9.81s)
testAccCheckAvailablePrefixesDestroyProducer:  There is not prefix with ID 68
--- PASS: TestAccAvaliablePrefixes_EmptyCustomFields (10.40s)
testAccCheckAvailablePrefixesDestroyProducer:  There is not prefix with ID 69
testAccCheckAvailablePrefixesDestroyProducer:  There is not prefix with ID 69
--- PASS: TestAccDataSourceAvailablePrefixes_basic (12.16s)
PASS
ok      github.com/fenglyu/terraform-provider-netbox/netbox     12.617s
?       github.com/fenglyu/terraform-provider-netbox/scripts/sidebar    [no test files]
?       github.com/fenglyu/terraform-provider-netbox/version    [no test files]
flv@genji ~/dev/go/terraform-plugins/terraform-provider-netbox

```

## Running one specific testcase
```bash

% TF_ACC=1 go test $(go list ./...) -v -run '^(TestAccAvaliablePrefixes_basic2)$'

?       github.com/fenglyu/terraform-provider-netbox    [no test files]
2020/08/04 23:12:50 [INFO] Instantiating http client for host netbox.k8s.me and path /api
2020/08/04 23:12:50 [INFO] Instantiating http client for host netbox.k8s.me and path /api
=== RUN   TestAccAvaliablePrefixes_basic2
=== PAUSE TestAccAvaliablePrefixes_basic2
=== CONT  TestAccAvaliablePrefixes_basic2
testAccCheckAvailablePrefixesDestroyProducer:  There is not prefix with ID 70
--- PASS: TestAccAvaliablePrefixes_basic2 (1.77s)
PASS
ok      github.com/fenglyu/terraform-provider-netbox/netbox     2.104s
?       github.com/fenglyu/terraform-provider-netbox/scripts/sidebar    [no test files]
?       github.com/fenglyu/terraform-provider-netbox/version    [no test files]

```