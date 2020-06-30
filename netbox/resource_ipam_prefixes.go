package netbox

import (
	"fmt"
	"github.com/fenglyu/go-netbox/netbox/client/ipam"
	"github.com/fenglyu/go-netbox/netbox/models"
	"log"
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
							Default:      "activIPv4 or IPv6 network with mask",
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
	config := m.(*Config)

	var id int64
	if idstr, ok := d.GetOk("id"); ok {
		id = idstr.(int64)
	}

	var family string
	if familyData, ok := d.GetOk("family"); ok {
		family = familyData.(string)
	}

	var prefix string
	if pfix, ok := d.GetOk("prefix"); ok {
		prefix = pfix.(string)
	}

	var prefixlength int64
	if pl, ok := d.GetOk("prefix_length"); ok {
		prefixlength = pl.(int64)
	}

	var site int64
	if siteData, ok := d.GetOk("site"); ok {
		site = siteData.(int64)
	}

	var vrf int64
	if vrfData, ok := d.GetOk("vrf"); ok {
		vrf = vrfData.(int64)
	}
	var tenant int64
	if tenantData, ok := d.GetOk("tenant"); ok {
		tenant = tenantData.(int64)
	}
	var vlan int64
	if vlanData, ok := d.GetOk("vlan"); ok {
		vlan = vlanData.(int64)
	}

	var status string
	if statusData, ok := d.GetOk("status"); ok {
		status = statusData.(string)
	}

	var role int64
	if roleData, ok := d.GetOk("role"); ok {
		role = roleData.(int64)
	}

	var IsPool bool
	if isPoolData, ok := d.GetOk("is_pool"); ok {
		IsPool = isPoolData.(bool)
	}

	var description string
	if desc, ok := d.GetOk("description"); ok {
		description = desc.(string)
	}

	var tags []string
	if tagsData, ok := d.GetOk("tags"); ok {
		tags = tagsData.([]string)
	}

	var customFields interface{}
	if cfData, ok := d.GetOk("custom_fields"); ok {
		customFields = cfData.(string)
	}

	wPrefix := models.WritablePrefix{
		Family:       family,
		Prefix:       &prefix,
		PrefixLength: prefixlength,
		Site:         &site,
		Vrf:          &vrf,
		Tenant:       &tenant,
		Vlan:         &vlan,
		Status:       status,
		Role:         &role,
		IsPool:       IsPool,
		Description:  description,
		Tags:         tags,
		CustomFields: customFields,
	}
	param := ipam.IpamPrefixesAvailablePrefixesCreateParams{
		ID:   id,
		Data: &wPrefix,
	}

	log.Printf("[INFO] Requesting AvaliablePrefix creation")
	res, err := config.client.Ipam.IpamPrefixesAvailablePrefixesCreate(&param, nil)
	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return err
	}
	availablePrefix := res.GetPayload()

	d.SetId(fmt.Sprintf("%s", availablePrefix.ID))

	return resourceIpamPrefixesRead(d, m)
}

func resourceIpamPrefixesRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	prefix, err := getIpamPrefix(config, d)
	if err != nil || prefix == nil {
		return err
	}

	//d.Set("id", prefix.ID)
	d.Set("description", prefix.Description)

	if err := d.Set("custom_fields", prefix.CustomFields); err != nil {
		return err
	}
	d.Set("is_pool", prefix.IsPool)
	d.Set("created", prefix.Created)
	d.Set("family", flatternFamily(prefix.Family))
	d.Set("role", flatternRole(prefix.Role))
	d.Set("last_updated", prefix.LastUpdated.String())
	d.Set("prefix", prefix.Prefix)
	d.Set("site", flatternSite(prefix.Site))
	d.Set("status", flatterPrefixStatus(prefix.Status))
	d.Set("tags", prefix.Tags)
	d.Set("tenant", flatternNestedTenant(prefix.Tenant))
	d.Set("vlan", flatternNestedVLAN(prefix.Vlan))
	d.Set("vrf", flatternNestedVRF(prefix.Vrf))

	return nil
}

func resourceIpamPrefixesUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	prefix, err := getIpamPrefix(config, d)
	if err != nil {
		return err
	}
	wp, err := getWritePrefixFromPrefix(prefix)
	if err != nil {
		return err
	}
	if d.HasChange("prefix") && !d.IsNewResource() {
		o, n := d.GetChange("prefix")
		wp.Prefix = n.(string)
	}

	if d.HasChange("site") && !d.IsNewResource() {
		o, n := d.GetChange("site")
		prefix.Site = n
	}

	return resourceIpamPrefixesRead(d, m)
}

func resourceIpamPrefixesDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func getIpamPrefix(config *Config, d *schema.ResourceData) (*models.Prefix, error) {
	/*
		idStr, err := getID(config, d)
		if err != nil {
			return nil, err
		}
	*/
	id, err := strconv.Atoi(d.Id())
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

func getWritePrefixFromPrefix(p *models.Prefix) (*models.WritablePrefix, error) {
	if p == nil {
		return nil, fmt.Errorf("nil pointer")
	}
	var wp models.WritablePrefix
	if p.Prefix != nil {
		wp.Prefix = p.Prefix
	}
	if p.Family != nil {
		wp.Family = *p.Family.Label
	}
	if p.Site != nil {
		wp.Site = &p.Site.ID
	}
	if p.Vrf != nil {
		wp.Vrf = &p.Vrf.ID
	}
	if p.Vlan != nil {
		wp.Vlan = &p.Vlan.ID
	}
	if p.Tenant != nil {
		wp.Tenant = &p.Tenant.ID
	}
	if p.Status != nil {
		wp.Status = *p.Status.Label
	}
	if p.Role != nil {
		wp.Role = &p.Role.ID
	}

	wp.IsPool = p.IsPool
	wp.Description = p.Description
	wp.Tags = p.Tags
	wp.CustomFields = p.CustomFields

	return &wp, nil
}

func getIpamPrefixes(config *Config, d *schema.ResourceData) ([]*models.Prefix, error) {
	/*
		id, err := getID(config, d)
		if err != nil {
			return nil, err
		}
	*/
	id, err := strconv.Atoi(d.Id())
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
