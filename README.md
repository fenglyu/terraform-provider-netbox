# Netbox


## Build netbox provider
```shell script
 % make build-dev13 version=0.1.8                             
```

## start/stop netbox testing environment with initializer data for testing
```shell script
% make test-netbox-env-up

# stop with
% make test-netbox-env-down
```

## Running Terraform Provider Acceptance tests
```shell script
 % export NETBOX_HOST="netbox.k8s.me"                         
export NETBOX_TOKEN=""
export NETBOX_BASE_PATH="/api"
export NETBOX_PARENT_PREFIX_ID=2
export NETBOX_PARENT_PREFIX_WITH_VRF_ID=1

% make testacc                                                                                                                                                                                     ✘ 130  ==> Checking source code against gofmt...                                                                                                                                                                  ==> Checking that code complies with gofmt requirements...                                                                                                                                                 TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(go list ./...) -v -timeout 240m -ldflags="-X=github.com/fenglyu/terraform-provider-netbox/version.ProviderVersion=acc"                                       ?       github.com/fenglyu/terraform-provider-netbox    [no test files]                                                                                                                                    2020/10/08 22:55:18 [INFO] Instantiating http client for host netbox.k8s.me and path /api                                                                                                                  2020/10/08 22:55:18 [INFO] Instantiating http client for host netbox.k8s.me and path /api                                                                                                                  === RUN   TestLoadAndValidate                                                                                                                                                                              2020/10/08 22:55:18 [INFO] Instantiating http client for host netbox.k8s.me and path /api
--- PASS: TestLoadAndValidate (0.02s)           
=== RUN   TestApiAccessTestInHttp                                                                    
    config_test.go:30: scheme http not supported, only https
--- SKIP: TestApiAccessTestInHttp (0.00s)        
=== RUN   TestApiAccessTestInHttps                                                                   
--- PASS: TestApiAccessTestInHttps (0.02s)                                                           
=== RUN   TestAccDataSourceAvailablePrefixesByPrefix                             
=== PAUSE TestAccDataSourceAvailablePrefixesByPrefix     
=== RUN   TestAccDataSourceAvailablePrefixesByPrefixId      
=== PAUSE TestAccDataSourceAvailablePrefixesByPrefixId        
=== RUN   TestAccDataSourceAvailablePrefixesByTag                                                    
=== PAUSE TestAccDataSourceAvailablePrefixesByTag                                                    
=== RUN   TestAccDataSourceAvailablePrefixesByRole 
=== PAUSE TestAccDataSourceAvailablePrefixesByRole                      
=== RUN   TestProvider                                                                               
--- PASS: TestProvider (0.00s)                                                                  
...
```

##



