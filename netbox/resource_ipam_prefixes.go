package netbox

import (
	"fmt"
	"github.com/fenglyu/go-netbox/netbox/client/ipam"
	"github.com/fenglyu/go-netbox/netbox/models"
	"strconv"

	//"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	//	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"time"
)

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
						"id": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Description:  "The unique ID of prefix",
						},
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
							ValidateFunc: validation.StringInSlice(prefixinitializeStatus, false),
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

	prefix, err := getIpamPrefix(config, d)
	if err != nil || prefix == nil {
		return err
	}

	d.Set("id", prefix.ID)
	d.Set("description", prefix.Description)
	d.Set("custom_fields", prefix.CustomFields)
	d.Set("is_pool", prefix.IsPool)
	d.Set("created", prefix.Created)
	d.Set("family", flatternFamily(prefix.Family))
	d.Set("prefix", prefix.Prefix)
	d.Set("role", prefix.Created)
	d.Set("created", prefix.Created)
	d.Set("created", prefix.Created)
	d.Set("created", prefix.Created)
	d.Set("created", prefix.Created)
	d.Set("created", prefix.Created)
	d.Set("created", prefix.Created)

	return nil
}

func resourceIpamPrefixesUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceIpamPrefixesRead(d, m)
}

func resourceIpamPrefixesDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func getIpamPrefix(config *Config, d *schema.ResourceData) (*models.Prefix, error) {
	idStr, err := getID(config, d)
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	params := ipam.IpamPrefixesReadParams{
		ID: int64(id),
	}

	prefix, err := config.client.Ipam.IpamPrefixesRead(&params, nil)
	if err != nil || prefix == nil {
		return nil, fmt.Errorf("Cannot determine prefix with ID %s", id)
	}
	return prefix.Payload, nil
}

func getIpamPrefixes(config *Config, d *schema.ResourceData) ([]*models.Prefix, error) {
	id, err := getID(config, d)
	if err != nil {
		return nil, err
	}

	prefix, err := getPrefix(config, d)
	if err != nil {
		return nil, err
	}

	//var limit int64 = 1
	// Compose Parameters for GET: /ipam/prefixes/
	param := ipam.IpamPrefixesListParams{
		ID:     &id,
		Prefix: &prefix,
		//Limit:  &limit,
	}
	ipamPrefixListBody, err := config.client.Ipam.IpamPrefixesList(&param, nil)
	if err != nil {
		return nil, err
	}

	if ipamPrefixListBody != nil || *ipamPrefixListBody.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow prefix %s with ID %s, not found", prefix, id)
	}

	return ipamPrefixListBody.Payload.Results, nil
}

func flatternFamily(f *models.PrefixFamily) []map[string]interface{} {
	if f == nil {
		return nil
	}
	return []map[string]interface{}{{
		"label": f.Label,
		"value": f.Value,
	}}
}

// TODO
func flatternNestedRole(nr *models.NestedRole) []map[string]interface{} {
	return nil
}

func flattenPrefixes(prefixesList []*models.Prefix) ([]map[string]interface{}, error) {
	flattened := make([]map[string]interface{}, len(prefixesList))

	for i, prefix := range prefixesList {
		flattened[i] = map[string]interface{}{
			"description":   prefix.Description,
			"custom_fields": prefix.CustomFields,
			"is_pool":       prefix.IsPool,
		}
	}
	return nil, nil
}

func getPrefix(config *Config, d *schema.ResourceData) (string, error) {
	return getAttrFromSchema("prefix", d, config)
}

func getID(config *Config, d *schema.ResourceData) (string, error) {
	return getAttrFromSchema("id", d, config)
}

func getAttrFromSchema(resourceSchemaField string, d *schema.ResourceData, config *Config) (string, error) {
	res, ok := d.GetOk(resourceSchemaField)
	if !ok {
		return "", fmt.Errorf("Cannot determine %s: set in this resource", resourceSchemaField)
	}
	return res.(string), nil
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
