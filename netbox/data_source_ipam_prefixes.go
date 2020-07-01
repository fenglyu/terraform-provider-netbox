package netbox

import (
	//"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var initializeStatus = []string{
	"container", "active", "reserved", "deprecated",
}

func dataSourceIpamPrefixes() *schema.Resource {
	// Generate datasource schema from resource
	return nil
}
