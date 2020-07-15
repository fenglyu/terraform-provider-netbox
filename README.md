# Netbox


## Build netbox provider
```shell script
make build-dev version=v0.0.1 
```

## Run an simple example
```shell script
 % cd examples/simple
 % terraform init

```

## Build for release
```
 % make  build version=v0.0.1
==> Checking source code against gofmt...
==> Checking that code complies with gofmt requirements...
go generate  ./...
==> Installing gox...
==> Building...
CGO_ENABLED=0 gox -os="linux darwin windows" -arch="amd64" -ldflags "-s -w" -output "pkg/{{.OS}}_{{.Arch}}/terraform-provider-netbox_v0.0.1" .
Number of parallel builds: 15

-->   windows/amd64: github.com/fenglyu/terraform-provider-netbox
-->     linux/amd64: github.com/fenglyu/terraform-provider-netbox
-->    darwin/amd64: github.com/fenglyu/terraform-provider-netbox

 % tree pkg
pkg
├── darwin_amd64
│   └── terraform-provider-netbox_v0.0.1
├── linux_amd64
│   └── terraform-provider-netbox_v0.0.1
└── windows_amd64
    └── terraform-provider-netbox_v0.0.1.exe

3 directories, 3 files

```


## Test
```shell script

export NETBOX_TOKEN=""                                                                                                                                   ✘ 130 
export NETBOX_HOST=""
export NETBOX_BASE_PATH="/api"

flv@genji ~/dev/go/terraform-plugins/terraform-provider-netbox
 % make testacc                                                                                                                                                                                       ✘ 2 
==> Checking source code against gofmt...
==> Checking that code complies with gofmt requirements...
TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(go list ./...) -v  -timeout 240m -ldflags="-X=github.com/fenglyu/terraform-provider-netbox/version.ProviderVersion=acc"
?       github.com/fenglyu/terraform-provider-netbox    [no test files]
=== RUN   TestProvider
--- PASS: TestProvider (0.00s)
=== RUN   TestAccProviderBasePath_setBasePath
=== PAUSE TestAccProviderBasePath_setBasePath
=== RUN   TestAccProviderBasePath_setInvalidBasePath
=== PAUSE TestAccProviderBasePath_setInvalidBasePath
=== RUN   TestAccAvaliablePrefixes_basic
=== PAUSE TestAccAvaliablePrefixes_basic
=== CONT  TestAccProviderBasePath_setBasePath
=== CONT  TestAccAvaliablePrefixes_basic
=== CONT  TestAccProviderBasePath_setInvalidBasePath
--- PASS: TestAccProviderBasePath_setInvalidBasePath (0.01s)
testAccCheckAvailablePrefixesDestroyProducer:  There is not prefix with ID 187
--- PASS: TestAccAvaliablePrefixes_basic (1.03s)
testAccCheckAvailablePrefixesDestroyProducer:  There is not prefix with ID 186
--- PASS: TestAccProviderBasePath_setBasePath (1.03s)
PASS
ok      github.com/fenglyu/terraform-provider-netbox/netbox     1.055s
?       github.com/fenglyu/terraform-provider-netbox/version    [no test files]
flv@genji ~/dev/go/terraform-plugins/terraform-provider-netbox

```
