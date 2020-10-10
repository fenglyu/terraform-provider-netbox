# terraform-provider-netbox

## 0.1.8 (Oct 8, 2020)
Add:heavy_check_mark: :
1. The netbox harness test integration finished successfully, huge convenience for Acceptance test.
2. Provider dependency upgrade to[ Terraform Plugin SDK v2 ](https://www.terraform.io/docs/extend/guides/v2-upgrade-guide.html) .

<details>
<summary>3. Support to fetch `prarent_prefix_id ` when running `terraform import `</summary>

```
 % terraform import  netbox_available_prefixes.gke-pods 11                                                                                                                                                 
netbox_available_prefixes.gke-pods: Importing from ID "11"...                                                                                                                                              
netbox_available_prefixes.gke-pods: Import prepared!                                                                                                                                                       
  Prepared netbox_available_prefixes for import                                                                                                                                                            
netbox_available_prefixes.gke-pods: Refreshing state... [id=11]                                                                                                                                            
                                                                                                                                                                                                           
Import successful!                                                                                                                                                                                         
                                                                                                     
The resources that were imported are shown above. These resources are now in                         
your Terraform state and will henceforth be managed by Terraform.                                    
                                                                                                     
 % terraform plan                                                                                                                                                                                          
Refreshing Terraform state in-memory prior to plan...                                                
The refreshed state will be used to calculate this plan, but will not be                             
persisted to local or remote state storage.                                                          
                                                                                                     
netbox_available_prefixes.gke-pods: Refreshing state... [id=11]                                      
                                                                                                     
------------------------------------------------------------------------                             
                                                                                                     
An execution plan has been generated and is shown below.                                             
Resource actions are indicated with the following symbols:                                           
                                                                                                     
Terraform will perform the following actions:                                                        
                                                                                                     
Plan: 0 to add, 0 to change, 0 to destroy.                                                                                                                                                                 
                                                                                                     
Changes to Outputs:                                                                                                                                                                                        
  + available_prefix    = "10.0.0.32/28"                                                             
  + available_prefix_id = "11"                                                                       
                                                                                                     
------------------------------------------------------------------------                             
                                                                                                     
Note: You didn't specify an "-out" parameter to save this plan, so Terraform                         
can't guarantee that exactly these actions will be performed if                                      
"terraform apply" is subsequently run.                                                                                                                                                                     
                                                                                                                                                                                                                   
 % terraform apply  --auto-approve                                                                   
netbox_available_prefixes.gke-pods: Refreshing state... [id=11]                                                                                                                                            
                                                                                                     
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.                                          
                                                                                                     
Outputs:

available_prefix = 10.0.0.32/28
available_prefix_id = 11
```
 </details>
 <details>
<summary>4. Enhanced output when parent prefix has no available prefixes for certain length.
</summary>

```bash
 % terraform apply  --auto-approve                                                                   
netbox_available_prefixes.gke-pods[2]: Creating...                                                   
netbox_available_prefixes.gke-pods[0]: Creating...                                                   
netbox_available_prefixes.gke-pods[1]: Creating...                                                   
netbox_available_prefixes.gke-pods[3]: Creating...                                                   
netbox_available_prefixes.gke-pods[2]: Creation complete after 0s [id=12]                            
netbox_available_prefixes.gke-pods[0]: Creation complete after 0s [id=13]                            
                                                                                                     
Error: Insufficient space is available to accommodate the requested prefix size(s) "/9"                                                                                                                    
                                                                   
Error: Insufficient space is available to accommodate the requested prefix size(s) "/9"     
```
 </details>

Delete:negative_squared_cross_mark: :
1. Remove self-built netbox testing environment

 <details>
<summary>Acceptance test records
</summary>

```bash
flv@genji ~/dev/go/terraform-plugins/terraform-provider-netbox                                                                                                                                             
 % make testacc                                                                                                                                                                                     ✘ 130  
==> Checking source code against gofmt...                                                                                                                                                                  
==> Checking that code complies with gofmt requirements...                                                                                                                                                 
TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(go list ./...) -v -timeout 240m -ldflags="-X=github.com/fenglyu/terraform-provider-netbox/version.ProviderVersion=acc"                                       
?       github.com/fenglyu/terraform-provider-netbox    [no test files]                                                                                                                                    
2020/10/08 22:55:18 [INFO] Instantiating http client for host netbox.k8s.me and path /api                                                                                                                  
2020/10/08 22:55:18 [INFO] Instantiating http client for host netbox.k8s.me and path /api                                                                                                                  
=== RUN   TestLoadAndValidate                                                                                                                                                                              
2020/10/08 22:55:18 [INFO] Instantiating http client for host netbox.k8s.me and path /api                                                                                                                  
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
=== RUN   TestAccProviderBasePath_setBasePath                                                                                                                                                              
=== PAUSE TestAccProviderBasePath_setBasePath                                                                                                                                                              
=== RUN   TestAccProviderBasePath_setInvalidBasePath                                                                                                                                                       
=== PAUSE TestAccProviderBasePath_setInvalidBasePath                                                                                                                                                       
=== RUN   TestAccAvailablePrefixes_basic                                                                                                                                                                   
=== PAUSE TestAccAvailablePrefixes_basic                                                                                                                                                                   
=== RUN   TestAccAvailablePrefixes_basic1                                                                                                                                                                  
=== PAUSE TestAccAvailablePrefixes_basic1                                                                                                                                                                  
=== RUN   TestAccAvailablePrefixes_basic2                                                                                                                                                                  
=== PAUSE TestAccAvailablePrefixes_basic2                                                                                                                                                                  
=== RUN   TestAccAvailablePrefixes_EmptyCustomFields                                                                                                                                                       
=== PAUSE TestAccAvailablePrefixes_EmptyCustomFields                                                                                                                                                       
=== RUN   TestAccAvailablePrefixesMultipleSteps                                                      
=== PAUSE TestAccAvailablePrefixesMultipleSteps                                                                                                                                                            
=== CONT  TestAccDataSourceAvailablePrefixesByPrefix          
=== CONT  TestAccAvailablePrefixes_basic                                                                                                                                                                   
=== CONT  TestAccDataSourceAvailablePrefixesByRole                                                   
=== CONT  TestAccAvailablePrefixes_EmptyCustomFields
=== RUN   TestAccAvailablePrefixesMultipleSteps                                                      
=== PAUSE TestAccAvailablePrefixesMultipleSteps                                                                                                                                                            
=== CONT  TestAccDataSourceAvailablePrefixesByPrefix          
=== CONT  TestAccAvailablePrefixes_basic                                                                                                                                                                   
=== CONT  TestAccDataSourceAvailablePrefixesByRole                                                   
=== CONT  TestAccAvailablePrefixes_EmptyCustomFields                                                                                                                                                       
=== CONT  TestAccAvailablePrefixesMultipleSteps                                                      
=== CONT  TestAccDataSourceAvailablePrefixesByTag                                                                                                                                                          
=== CONT  TestAccAvailablePrefixes_basic2
=== CONT  TestAccAvailablePrefixes_basic1                                                            
=== CONT  TestAccDataSourceAvailablePrefixesByPrefixId                                               
=== CONT  TestAccProviderBasePath_setInvalidBasePath                                                 
=== CONT  TestAccProviderBasePath_setBasePath                                                        
2020/10/08 22:55:18 [WARN] Truncating attribute path of 0 diagnostics for TypeSet                    
2020/10/08 22:55:18 [WARN] Truncating attribute path of 0 diagnostics for TypeSet                    
--- PASS: TestAccProviderBasePath_setInvalidBasePath (0.72s)         
2020/10/08 22:55:19 [getIpamParentPrefixes] ipamPrefixListBody [{"created":"2020-10-08","custom_fields":{"helpers":null,"ipv4_acl_in":null,"ipv4_acl_out":null},"description":"activision vrf prefix","fami
ly":{"label":"IPv4","value":4},"id":1,"is_pool":false,"last_updated":"2020-10-08T14:43:28.694Z","prefix":"10.0.0.0/8","site":{"id":1,"name":"se1","slug":"se1","url":"http://netbox.k8s.me/api/dcim/sites/1
/"},"status":{"label":"Active","value":"active"},"tenant":{"id":1,"name":"cloud","slug":"cloud","url":"http://netbox.k8s.me/api/tenancy/tenants/1/"},"vlan":{"display_name":"1009 (gcp)","id":1,"name":"gcp
","url":"http://netbox.k8s.me/api/ipam/vlans/1/","vid":1009},"vrf":{"id":1,"name":"activision","url":"http://netbox.k8s.me/api/ipam/vrfs/1/"}},{"created":"2020-10-08","custom_fields":{"helpers":"cf-accnv
t1xgnzdl-01","ipv4_acl_in":"cf-accnvt1xgnzdl-02","ipv4_acl_out":"cf-accnvt1xgnzdl-03"},"family":{"label":"IPv4","value":4},"id":113,"is_pool":true,"last_updated":"2020-10-08T14:55:19.257Z","prefix":"10.0
.0.0/19","role":{"id":2,"name":"gcp","slug":"gcp","url":"http://netbox.k8s.me/api/ipam/roles/2/"},"site":{"id":1,"name":"se1","slug":"se1","url":"http://netbox.k8s.me/api/dcim/sites/1/"},"status":{"label
":"Active","value":"active"},"tags":["AvailablePrefix-accnvt1xgnzdl-04","AvailablePrefix-accnvt1xgnzdl-05","AvailablePrefix-accnvt1xgnzdl-03"],"tenant":{"id":1,"name":"cloud","slug":"cloud","url":"http:/
/netbox.k8s.me/api/tenancy/tenants/1/"},"vlan":{"display_name":"1009 (gcp)","id":1,"name":"gcp","url":"http://netbox.k8s.me/api/ipam/vlans/1/","vid":1009},"vrf":{"id":1,"name":"activision","url":"http://
netbox.k8s.me/api/ipam/vrfs/1/"}}]                                                                   
--- PASS: TestAccAvailablePrefixes_basic (2.02s)
--- PASS: TestAccProviderBasePath_setBasePath (2.24s)                                                
--- PASS: TestAccAvailablePrefixes_basic2 (2.50s)                                                                                                                                                          
--- PASS: TestAccAvailablePrefixes_basic1 (2.61s)                                                    
--- PASS: TestAccAvailablePrefixes_EmptyCustomFields (2.86s)                                         
2020/10/08 22:55:21 [WARN] Truncating attribute path of 0 diagnostics for TypeSet                    
2020/10/08 22:55:21 [WARN] Truncating attribute path of 0 diagnostics for TypeSet                    
--- PASS: TestAccDataSourceAvailablePrefixesByTag (3.60s)     
--- PASS: TestAccDataSourceAvailablePrefixesByPrefix (4.41s)                                         
--- PASS: TestAccDataSourceAvailablePrefixesByPrefixId (4.67s)      
--- PASS: TestAccDataSourceAvailablePrefixesByRole (4.91s)                                           
--- PASS: TestAccAvailablePrefixesMultipleSteps (6.86s)                                              
PASS                                                                                                 
ok      github.com/fenglyu/terraform-provider-netbox/netbox     (cached)
?       github.com/fenglyu/terraform-provider-netbox/scripts/sidebar    [no test files]              
?       github.com/fenglyu/terraform-provider-netbox/version    [no test files]                      
flv@genji ~/dev/go/terraform-plugins/terraform-provider-netbox                           
```
 </details>


## 0.1.7 (Sep 18, 2020)

### Update
 Fix a bug which crashes the plugin when site is not bind to the given tenant


## 0.1.6 (Sep 16, 2020)
### Upgrade
Starting from this release `0.1.6`, `terraform-provider-netbox` will only support feature and update request for netbox v2.8.* and above.
netbox v2.4.* users can still use release `0.1.4`, but there won't be any update for it.

### Feature
Port the data `lookup` for netbox 2.8

## 0.1.4 (Aug 7, 2020)
### Update
Fix data source name scalability issue

## 0.1.3 (Aug 7, 2020)
### New Feature

Data source "netbox_available_prefixes" now support data lookup by netbox query parameters, Implemented by @ndowns @ https://ghosthub.corp.blizzard.net/flv/terraform-provider-netbox/pull/7

### Example
An example can be found here https://ghosthub.corp.blizzard.net/flv/terraform-provider-netbox/blob/master/examples/databytag/main.tf.
> This example is tested under Terraform v0.13,  It relies on `depends_on` for execution order, But Terraform v0.12 has well-known bug/issue that `terraform plan` seems to re-read the data source even after applying resources. https://github.com/hashicorp/terraform/issues/11806



## 0.1.2 (July 31, 2020)
### Race Condition fix

Add `mutex lock` for Resource's Create/Update/Delete Operations, This is expected to fix the duplicated/nested available prefixes being created by concurrent operations. Besides, We can also upgrade netbox to a recent release which has mutex lock implemented its server side.

## 0.1.1 (July 21, 2020)
### Feature Update

`custom_fields` has changed from `map` type *Optional* field to *Required* `list` type with 3 elements.
> In practice, Filed/Attribute change should be reflected in the major release version number increasing, Here we only update minor version number.
1. helpers  (Optional) string
2. ipv4_acl_in  (Optional) string
3. ipv4_acl_out (Optional) string

```
  custom_fields  {
      helpers  = "foo"
      ipv4_acl_in  = "foo"
      ipv4_acl_out = "bar"
  }
```


#### Most of our se1 netbox `prefix` has those elements empty, So we can just put a placeholder like.
```
custom_fields  {}
```

### Warning
> By default, There is no need to set `vrf` field, It's "gobal", Once we configure `vrf` field for Non routable network, The `parent_prefix_id` should also be set to the corresponding parent prefix in the same `vrf` field.
> This is due to the nature of `vrf` that same prefixes can exist among differrent vrfs.


## 0.1.0 (July 16, 2020)
### Features
Resource [`netbox_available_prefixes`](https://ghosthub.corp.blizzard.net/flv/terraform-provider-netbox/blob/netbox_v2.4.7/website/docs/r/available_prefixes.html.md)
Data Source [`netbox_available_prefixes`](https://ghosthub.corp.blizzard.net/flv/terraform-provider-netbox/blob/netbox_v2.4.7/website/docs/d/available_prefixes.html.md)

### Simple Example
https://ghosthub.corp.blizzard.net/flv/terraform-provider-netbox/blob/netbox_v2.4.7/examples/blizzrd/main.tf


