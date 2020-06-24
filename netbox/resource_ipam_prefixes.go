package netbox

import (
	//"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	//	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"time"
)

var initializeStatus = []string{
	"container", "active", "reserved", "deprecated",
}

func resourceIpamPrefixes() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpamPrefixesCreate,
		Read:   resourceIpamPrefixesRead,
		Update: resourceIpamPrefixesUpdate,
		Delete: resourceIpamPrefixesDelete,
		/*
			Importer: &schema.ResourceImporter{
				State: resourceIpamPrefixesImportState,
			},
		*/
		SchemaVersion: 1,
		// TODO after test coverage finished
		//MigrateState:
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The IPAM prefix in netbox",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsCIDRNetwork(8, 32),
							Description:  "IPv4 or IPv6 network with mask",
						},

						"role": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The primary function of this prefix  ",
						},
						"site": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Site",
						},
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "tags",
						},
						"tenant": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Tenant",
						},
						"vlan": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "VLAN",
						},
						"vrf": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "VRF",
						},
						"ispool": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "All IP addresses within this prefix are considered usable",
						},
						"status": {
							Type:         schema.TypeString,
							Default:      "activIPv4 or IPv6 network with maske",
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"container", "active", "reserved", "deprecated"}, false),
							Description:  "Operational status of this prefix",
						},
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(0, 200),
							Description:  "Describe the purpose of this prefix",
						},
					},
				},
			},
		},

		//	CustomizeDiff: nil,
	}
}

func resourceIpamPrefixesCreate(d *schema.ResourceData, m interface{}) error {
	return resourceIpamPrefixesRead(d, m)
}

func resourceIpamPrefixesRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	return nil
}

func resourceIpamPrefixesUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceIpamPrefixesRead(d, m)
}

func resourceIpamPrefixesDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

/*
func resourceIpamPrefixesImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/zones/(?P<zone>[^/]+)/instances/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<zone>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/zones/{{zone}}/instances/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
*/
