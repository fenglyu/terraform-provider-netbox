package main

import (
	"github.com/fenglyu/terraform-provider-netbox/netbox"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: netbox.Provider})
}
