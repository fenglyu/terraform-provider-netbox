package main


import (
 "github.com/hashicorp/terraform-plugin-sdk/plugin"
 "github.com/fenglyu/terraform-provider-netbox"
)
//

func main() {
 plugin.Serve(&plugin.ServeOpts{
  ProviderFunc: google.Provider})
}
