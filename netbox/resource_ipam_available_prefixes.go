package netbox

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/fenglyu/go-netbox/netbox/client/ipam"
	"github.com/fenglyu/go-netbox/netbox/models"
)

var (
	availablePrefixesKeys = []string{
		"parent_prefix",
		"parent_prefix_id",
	}
)

func resourceIpamAvailablePrefixes() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpamAvailablePrefixesCreate,
		Read:   resourceIpamAvailablePrefixesRead,
		Update: resourceIpamAvailablePrefixesUpdate,
		Delete: resourceIpamAvailablePrefixesDelete,

		Importer: &schema.ResourceImporter{
			//	State: resourceIpamAvailablePrefixesImportState,
			State: schema.ImportStatePassthrough,
		},
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
				ForceNew:     true,
				AtLeastOneOf: availablePrefixesKeys,
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
			"parent_prefix_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: availablePrefixesKeys,
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
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created date",
			},
			"family": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IP family, IPv4, or Ipv6",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last updated timestamp",
			},
		},
		//	CustomizeDiff: nil,
	}
}

func resourceIpamAvailablePrefixesCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	wPrefix := models.WritablePrefix{}

	var prefix_id int64
	if pfx_id, ok := d.GetOk("parent_prefix_id"); ok {
		prefix_id = int64(pfx_id.(int))
		// not necessary
		//wPrefix.ID = prefix_id
	}

	//As of version 2.8, Netbox doesn't require "prefix" in post data,
	//only prefix_length and a parent_id is mandatory
	//prerent_prefix here is prepared to fetch prrent_prefix_id only
	//wPrefix.Prefix = &prefix

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
		wPrefix.Status = prefixStatusIDMap[status]
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
		tags = convertStringSet(tagsData.(*schema.Set))
		wPrefix.Tags = tags
	}

	var customFields interface{}
	if cfData, ok := d.GetOk("custom_fields"); ok {
		customFields = cfData.(map[string]string)
		wPrefix.CustomFields = customFields
	}

	// If parent prefix is given
	wPrefixRes, _ := json.Marshal(wPrefix)
	log.Println("[INFO] ", string(wPrefixRes))

	if _, ok := getParentPrefix(config, d); ok == nil {
		if wPrefix.ID == 0 {
			results, err := getIpamPrefixes(config, d)
			if err != nil {
				return err
			}
			wPrefix.ID = results[0].ID
			prefix_id = results[0].ID
		}
	}
	param := ipam.IpamPrefixesAvailablePrefixesCreateParams{
		ID:   int64(prefix_id),
		Data: &wPrefix,
	}
	param.WithContext(context.Background())

	paramRes, _ := json.Marshal(param)
	log.Printf("[INFO] Requesting AvaliablePrefix creation %s", string(paramRes))

	res, err := config.client.Ipam.IpamPrefixesAvailablePrefixesCreate(&param, nil)
	if err != nil {
		// The resource didn't actually create
		log.Fatalln("[Error] Failed to create AvaliablePrefix: ", err)
		d.SetId("")
		return err
	}
	availablePrefix := res.GetPayload()
	d.SetId(fmt.Sprintf("%d", availablePrefix.ID))

	return resourceIpamAvailablePrefixesRead(d, m)
}

func resourceIpamAvailablePrefixesRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	prefix, err := getIpamPrefix(config, d)
	if err != nil || prefix == nil {
		return err
	}

	log.Println("[INFO] resourceIpamPrefixesRead ", prefix)
	//d.Set("id", prefix.ID)
	d.Set("description", prefix.Description)
	d.Set("custom_fields", flattenCustomFields(prefix))
	d.Set("is_pool", prefix.IsPool)
	d.Set("created", prefix.Created.String())
	d.Set("family", prefix.Family)
	d.Set("last_updated", prefix.LastUpdated.String())

	if prefix != nil && prefix.Role != nil {
		d.Set("role", prefix.Role.ID)
	}

	if ppid, ok := d.GetOk("parent_prefix_id"); ok {
		d.Set("parent_prefix_id", ppid.(int))
	}

	d.Set("prefix", prefix.Prefix)
	pl := strings.Split(*prefix.Prefix, "/")[1]
	prefixLength, _ := strconv.Atoi(pl)
	d.Set("prefix_length", prefixLength)
	if prefix != nil && prefix.Site != nil {
		d.Set("site", prefix.Site.ID)
	}
	if prefix != nil && prefix.Status != nil {
		d.Set("status", prefixStatusIDMapReverse[*prefix.Status.Value])
	}
	d.Set("tags", prefix.Tags)
	if prefix != nil && prefix.Tenant != nil {
		d.Set("tenant", prefix.Tenant.ID)
	}
	if prefix != nil && prefix.Vlan != nil {
		d.Set("vlan", prefix.Vlan.ID)
	}
	if prefix != nil && prefix.Vrf != nil {
		d.Set("vrf", prefix.Vrf.ID)
	}

	d.SetId(fmt.Sprintf("%d", prefix.ID))
	return nil
}

func resourceIpamAvailablePrefixesUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	var writablePrefix models.WritablePrefix

	// required property
	prefixData := d.Get("prefix").(string)
	writablePrefix.Prefix = &prefixData

	if d.HasChange("prefix_length") && !d.IsNewResource() {
		prefixLengthData := int64(d.Get("prefix_length").(int))
		writablePrefix.PrefixLength = prefixLengthData
	}
	if d.HasChange("site") && !d.IsNewResource() {
		siteId := int64(d.Get("site").(int))
		writablePrefix.Site = &siteId
	}
	if d.HasChange("vrf") && !d.IsNewResource() {
		vrfData := int64(d.Get("vrf").(int))
		writablePrefix.Vrf = &vrfData
	}
	if d.HasChange("tenant") && !d.IsNewResource() {
		tenantData := int64(d.Get("tenant").(int))
		writablePrefix.Tenant = &tenantData
	}
	if d.HasChange("vlan") && !d.IsNewResource() {
		vlanData := int64(d.Get("vlan").(int))
		writablePrefix.Vlan = &vlanData
	}
	if d.HasChange("status") && !d.IsNewResource() {
		statusData := d.Get("status").(string)
		flag := false
		for _, str := range prefixinitializeStatus {
			if statusData == str || (strings.EqualFold(statusData, str)) {
				flag = true
			}
		}
		if !flag {
			return fmt.Errorf("Not a valid status in %v", prefixinitializeStatus)
		}
		writablePrefix.Status = prefixStatusIDMap[strings.ToLower(statusData)]
	}
	if d.HasChange("role") && !d.IsNewResource() {
		roleData := int64(d.Get("role").(int))
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
	partialUpdatePrefixRes, _ := json.Marshal(partialUpdatePrefix)
	log.Println("partialUpdatePrefix: ", string(partialUpdatePrefixRes))
	res, uerr := config.client.Ipam.IpamPrefixesPartialUpdate(&partialUpdatePrefix, nil)
	if uerr != nil {
		// TODO Support verbose response body here
		return fmt.Errorf("%v %v", res, uerr)
	}

	return resourceIpamAvailablePrefixesRead(d, m)
}

func resourceIpamAvailablePrefixesDelete(d *schema.ResourceData, m interface{}) error {
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

	prefix, err := getParentPrefix(config, d)
	if err != nil {
		return nil, err
	}

	//var limit int64 = 1
	// Compose Parameters for GET: /ipam/prefixes/
	withinInclude := prefix
	prefixLength, err := strconv.Atoi(strings.Split(withinInclude, "/")[1])
	if err != nil {
		return nil, fmt.Errorf("Error in [getIpamPrefixes] %v", err)
	}
	maskLength := float64(prefixLength)
	// v2.4.7 api query fomat
	// http://netbox.k8s.me/api/ipam/prefixes/?mask_length=24&within_include=10.247.5.0/24
	param := ipam.IpamPrefixesListParams{
		MaskLength:    &maskLength,
		WithinInclude: &withinInclude,
	}
	param.WithContext(context.Background())
	ipamPrefixListBody, err := config.client.Ipam.IpamPrefixesList(&param, nil)
	if err != nil {
		return nil, err
	}
	ipamPrefixesReadOKRes, _ := json.Marshal(&ipamPrefixListBody.Payload.Results)
	log.Println("ipamPrefixListBody", string(ipamPrefixesReadOKRes))
	if ipamPrefixListBody == nil || *ipamPrefixListBody.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow prefix %s with ID %s, not found", prefix, d.Id())
	}

	return ipamPrefixListBody.Payload.Results, nil
}

func getParentPrefix(config *Config, d *schema.ResourceData) (string, error) {
	return getAttrFromSchema("parent_prefix", d, config)
}

func getAttrFromSchema(resourceSchemaField string, d *schema.ResourceData, config *Config) (string, error) {
	res, ok := d.GetOk(resourceSchemaField)
	if !ok {
		return "", fmt.Errorf("Cannot determine %s: set in this resource", resourceSchemaField)
	}
	return res.(string), nil
}
