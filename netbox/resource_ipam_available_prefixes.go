package netbox

import (
	"context"
	"fmt"
	"github.com/fenglyu/go-netbox/netbox/client/ipam"
	"github.com/fenglyu/go-netbox/netbox/models"
	"log"
	"strconv"
	"strings"

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
			Create:  schema.DefaultTimeout(10 * time.Minute),
			Update:  schema.DefaultTimeout(10 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"parent_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDRNetwork(8, 128),
				Description:  "crave available prefixes under the parent_prefix",
			},
			"prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsCIDRNetwork(8, 128),
				Description:  "craved available prefix",
			},
			"prefix_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  "A unique integer value identifying this prefix under which is used crave available prefix",
			},
			"prefix_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 128),
				Description:  "The netmask in number form",
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The primary function of this prefix  ",
			},
			"site": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Site",
			},
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Optional:    true,
				Description: `The list of tags attached to the available prefix.`,
			},
			"tenant": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Tenant",
			},
			"vlan": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VLAN",
			},
			"vrf": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VRF",
			},
			"is_pool": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "All IP addresses within this prefix are considered usable",
			},
			"status": {
				Type:         schema.TypeString,
				Default:      "active",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(prefixinitializeStatus, false),
				Description:  "Operational status of this prefix",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 200),
				Description:  "Describe the purpose of this prefix",
			},
			"custom_fields": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Custom fields",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		//	CustomizeDiff: nil,
	}
}

func resourceIpamPrefixesCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	wPrefix := models.WritablePrefix{}

	var prefix_id int64
	if pfx_id, ok := d.GetOk("prefix_id"); ok {
		prefix_id = int64(pfx_id.(int))
		//wPrefix.ID = prefix_id
	}

	var prefix string
	if pfix, ok := d.GetOk("prefix"); ok {
		prefix = pfix.(string)
		wPrefix.Prefix = &prefix
	}

	var prefixlength int64
	if pl, ok := d.GetOk("prefix_length"); ok {
		prefixlength = int64(pl.(int))
		wPrefix.PrefixLength = prefixlength
	}

	var site int64
	if siteData, ok := d.GetOk("site"); ok {
		site = int64(siteData.(int))
		wPrefix.Site = &site
	}

	var vrf int64
	if vrfData, ok := d.GetOk("vrf"); ok {
		vrf = int64(vrfData.(int))
		wPrefix.Vrf = &vrf
	}
	var tenant int64
	if tenantData, ok := d.GetOk("tenant"); ok {
		tenant = int64(tenantData.(int))
		wPrefix.Tenant = &tenant
	}
	var vlan int64
	if vlanData, ok := d.GetOk("vlan"); ok {
		vlan = int64(vlanData.(int))
		wPrefix.Vlan = &vlan
	}

	var status string
	if statusData, ok := d.GetOk("status"); ok {
		status = statusData.(string)
		wPrefix.Status = status
	}

	var role int64
	if roleData, ok := d.GetOk("role"); ok {
		role = int64(roleData.(int))
		wPrefix.Role = &role
	}

	var IsPool bool
	if isPoolData, ok := d.GetOk("is_pool"); ok {
		IsPool = isPoolData.(bool)
		wPrefix.IsPool = IsPool
	}

	var description string
	if desc, ok := d.GetOk("description"); ok {
		description = desc.(string)
		wPrefix.Description = description
	}

	var tags []string
	if tagsData, ok := d.GetOk("tags"); ok {
		tags = tagsData.([]string)
		wPrefix.Tags = tags
	}

	var customFields interface{}
	if cfData, ok := d.GetOk("custom_fields"); ok {
		customFields = cfData.(map[string]string)
		wPrefix.CustomFields = customFields
	}
	param := ipam.IpamPrefixesAvailablePrefixesCreateParams{
		ID:   int64(prefix_id),
		Data: &wPrefix,
	}
	param.WithContext(context.Background())

	log.Printf("[INFO] Requesting AvaliablePrefix creation")
	res, err := config.client.Ipam.IpamPrefixesAvailablePrefixesCreate(&param, nil)
	if err != nil {
		// The resource didn't actually create
		log.Fatalln("[Error] Failed to create AvaliablePrefix: ", err)
		d.SetId("")
		return err
	}
	availablePrefix := res.GetPayload()

	d.SetId(fmt.Sprintf("%d", availablePrefix.ID))

	return resourceIpamPrefixesRead(d, m)
}

func resourceIpamPrefixesRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	prefix, err := getIpamPrefix(config, d)
	if err != nil || prefix == nil {
		return err
	}

	log.Println("[INFO] resourceIpamPrefixesRead ", prefix)
	//d.Set("id", prefix.ID)
	d.Set("description", prefix.Description)

	if err := d.Set("custom_fields", flattenCustomFields(prefix)); err != nil {
		return err
	}
	d.Set("is_pool", prefix.IsPool)
	d.Set("created", prefix.Created)
	if prefix != nil && prefix.Family != nil {
		d.Set("family", flatternFamily(prefix.Family))
	}
	if prefix != nil && prefix.Role != nil {
		d.Set("role", flatternRole(prefix.Role))
	}
	d.Set("last_updated", prefix.LastUpdated.String())
	d.Set("prefix", prefix.Prefix)
	d.Set("prefix_length", strings.Split(*prefix.Prefix, "/")[1])
	if prefix != nil && prefix.Site != nil {
		d.Set("site", flatternSite(prefix.Site))
	}
	if prefix != nil && prefix.Status != nil {
		d.Set("status", flatterPrefixStatus(prefix.Status))
	}

	d.Set("tags", prefix.Tags)

	if prefix != nil && prefix.Tenant != nil {
		d.Set("tenant", flatternNestedTenant(prefix.Tenant))
	}
	if prefix != nil && prefix.Vlan != nil {
		d.Set("vlan", flatternNestedVLAN(prefix.Vlan))
	}

	if prefix != nil && prefix.Vrf != nil {
		d.Set("vrf", flatternNestedVRF(prefix.Vrf))
	}

	return nil
}

func resourceIpamPrefixesUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	var writablePrefix models.WritablePrefix

	// required property
	prefixData := d.Get("prefix").(string)
	writablePrefix.Prefix = &prefixData

	if d.HasChange("prefix_length") && !d.IsNewResource() {
		prefixLengthData := d.Get("prefix_length").(int64)
		writablePrefix.PrefixLength = prefixLengthData
	}
	if d.HasChange("site") && !d.IsNewResource() {
		siteId := d.Get("site").(int64)
		writablePrefix.Site = &siteId
	}
	if d.HasChange("vrf") && !d.IsNewResource() {
		vrfData := d.Get("vrf").(int64)
		writablePrefix.Vrf = &vrfData
	}
	if d.HasChange("tenant") && !d.IsNewResource() {
		tenantData := d.Get("tenant").(int64)
		writablePrefix.Tenant = &tenantData
	}
	if d.HasChange("vlan") && !d.IsNewResource() {
		vlanData := d.Get("vlan").(int64)
		writablePrefix.Vlan = &vlanData
	}
	if d.HasChange("status") && !d.IsNewResource() {
		statusData := d.Get("status").(string)
		flag := false
		for _, str := range prefixinitializeStatus {
			if statusData == str || (strings.ToLower(statusData) == strings.ToLower(str)) {
				flag = true
			}
		}
		if !flag {
			return fmt.Errorf("Not a valid status in %v", prefixinitializeStatus)
		}
		writablePrefix.Status = strings.ToLower(statusData)
	}
	if d.HasChange("role") && !d.IsNewResource() {
		roleData := d.Get("role").(int64)
		writablePrefix.Role = &roleData
	}
	if d.HasChange("is_pool") && !d.IsNewResource() {
		isPoolData := d.Get("is_pool").(bool)
		writablePrefix.IsPool = isPoolData
	}
	if d.HasChange("description") && !d.IsNewResource() {
		descriptionData := d.Get("description").(string)
		writablePrefix.Description = descriptionData
	}
	if d.HasChange("tags") && !d.IsNewResource() {
		writablePrefix.Tags = convertStringSet(d.Get("tags").(*schema.Set))
	}
	if d.HasChange("custom_fields") && !d.IsNewResource() {
		cfData := d.Get("custom_fields").(map[string]string)
		writablePrefix.CustomFields = cfData
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	partialUpdatePrefix := ipam.IpamPrefixesPartialUpdateParams{
		ID:   int64(id),
		Data: &writablePrefix,
	}
	partialUpdatePrefix.WithContext(context.Background())
	_, uerr := config.client.Ipam.IpamPrefixesPartialUpdate(&partialUpdatePrefix, nil)
	if uerr != nil {
		return uerr
	}

	return resourceIpamPrefixesRead(d, m)
}

func resourceIpamPrefixesDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	log.Printf("[INFO]Requesting Prefix deletion: %s", d.Get("prefix").(string))
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	params := ipam.IpamPrefixesDeleteParams{
		ID: int64(id),
	}
	params.WithContext(context.Background())
	_, derr := config.client.Ipam.IpamPrefixesDelete(&params, nil)
	if derr != nil {
		return derr
	}

	d.SetId("")
	return nil
}

func getIpamPrefix(config *Config, d *schema.ResourceData) (*models.Prefix, error) {

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, err
	}
	params := ipam.IpamPrefixesReadParams{
		ID: int64(id),
	}
	params.WithContext(context.Background())

	ipamPrefixesReadOK, err := config.client.Ipam.IpamPrefixesRead(&params, nil)
	if err != nil || ipamPrefixesReadOK == nil {
		return nil, fmt.Errorf("Cannot determine prefix with ID %d", id)
	}
	return ipamPrefixesReadOK.Payload, nil
}

func getIpamPrefixes(config *Config, d *schema.ResourceData) ([]*models.Prefix, error) {

	prefix, err := getPrefix(config, d)
	if err != nil {
		return nil, err
	}

	//var limit int64 = 1
	// Compose Parameters for GET: /ipam/prefixes/
	idStr := d.Id()
	param := ipam.IpamPrefixesListParams{
		ID:     &idStr,
		Prefix: &prefix,
		//Limit:  &limit,
	}
	ipamPrefixListBody, err := config.client.Ipam.IpamPrefixesList(&param, nil)
	if err != nil {
		return nil, err
	}

	if ipamPrefixListBody != nil || *ipamPrefixListBody.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow prefix %s with ID %s, not found", prefix, d.Id())
	}

	return ipamPrefixListBody.Payload.Results, nil
}

func getPrefix(config *Config, d *schema.ResourceData) (string, error) {
	return getAttrFromSchema("prefix", d, config)
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