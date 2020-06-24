package netbox

import (
	//"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

var initializeStatus = []string{
	"container",  "active", "reserved", "deprecated",
}

func dataSourceIpamPrefixes() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema :=
}