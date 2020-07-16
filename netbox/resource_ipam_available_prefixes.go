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

	"github.com/fenglyu/go-netbox/netbox/client/dcim"
	"github.com/fenglyu/go-netbox/netbox/client/ipam"
	"github.com/fenglyu/go-netbox/netbox/client/tenancy"
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
				Description:  "The mask in integer form",
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Role",
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tenant",
			},
			"vlan": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VLAN",
			},
			"vrf": {
				Type:        schema.TypeString,
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
				Description: "IPv4, or Ipv6",
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

	var site *models.Site
	if sites, err := getDcimSites(config, d); err == nil {
		site = sites[0]
		wPrefix.Site = &site.ID
	}

	var tenant int64
	if tenantData, ok := d.GetOk("tenant"); ok {
		if *site.Tenant.Name != tenantData.(string) {
			return fmt.Errorf("Incompatible site %s and the tenant %s, expected tenant %s", *site.Name, tenantData.(string), *site.Tenant.Name)
		}
		tenant = site.Tenant.ID
		wPrefix.Tenant = &tenant
	}

	var vrf *models.VRF
	if vrfs, err := getIpamVrfs(config, d); err == nil {
		vrf = vrfs[0]
		wPrefix.Vrf = &vrf.ID
	}

	var vlan *models.VLAN
	if vlans, err := getIpamVlans(config, d); err == nil {
		vlan = vlans[0]
		wPrefix.Vlan = &vlan.ID
	}

	var status string
	if statusData, ok := d.GetOk("status"); ok {
		status = statusData.(string)
		wPrefix.Status = prefixStatusIDMap[status]
	}

	var role *models.Role
	if roles, err := getIpamRoles(config, d); err == nil {
		role = roles[0]
		wPrefix.Role = &role.ID
	}

	var IsPool bool
	if isPoolData, ok := d.GetOk("is_pool"); ok {
		IsPool = isPoolData.(bool)
		wPrefix.IsPool = &IsPool
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
		customFields = cfData.(map[string]interface{})
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

	//  Duo the bug in API `ipam/prefixes/{ID}/available-prefixes/` which can't setup the right vrf
	//  Here we update its vrf as a fixup
	if vrf != nil {
		vrfP := models.WritablePrefix{
			Vrf:    &vrf.ID,
			Prefix: availablePrefix.Prefix,
		}
		vrfpartialUpdate := ipam.IpamPrefixesPartialUpdateParams{
			ID:   availablePrefix.ID,
			Data: &vrfP,
		}
		vrfpartialUpdate.WithContext(context.Background())
		partialUpdatePrefixRes, _ := json.Marshal(vrfpartialUpdate)
		log.Println("partialUpdatePrefix: ", string(partialUpdatePrefixRes))

		vrfRes, uerr := config.client.Ipam.IpamPrefixesPartialUpdate(&vrfpartialUpdate, nil)
		if uerr != nil {
			d.SetId("")
			return fmt.Errorf("%v %v", vrfRes, uerr)
		}
	}

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
	d.Set("custom_fields", prefix.CustomFields)
	d.Set("is_pool", prefix.IsPool)
	d.Set("created", prefix.Created.String())
	d.Set("family", prefix.Family)
	d.Set("last_updated", prefix.LastUpdated.String())

	if prefix != nil && prefix.Role != nil {
		d.Set("role", prefix.Role.Name)
	}

	if ppid, ok := d.GetOk("parent_prefix_id"); ok {
		d.Set("parent_prefix_id", ppid.(int))
	}

	d.Set("prefix", prefix.Prefix)
	pl := strings.Split(*prefix.Prefix, "/")[1]
	prefixLength, _ := strconv.Atoi(pl)
	d.Set("prefix_length", prefixLength)
	if prefix != nil && prefix.Site != nil {
		d.Set("site", prefix.Site.Name)
	} else {
		d.Set("site", "")
	}
	if prefix != nil && prefix.Status != nil {
		d.Set("status", prefixStatusIDMapReverse[*prefix.Status.Value])
	}
	d.Set("tags", prefix.Tags)
	if prefix != nil && prefix.Tenant != nil {
		d.Set("tenant", prefix.Tenant.Name)
	}
	if prefix != nil && prefix.Vlan != nil {
		d.Set("vlan", prefix.Vlan.Name)
	}
	if prefix != nil && prefix.Vrf != nil {
		d.Set("vrf", prefix.Vrf.Name)
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

	if d.HasChange("is_pool") && !d.IsNewResource() {
		v := d.Get("is_pool").(bool)
		writablePrefix.IsPool = &v
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
	if d.HasChange("site") && !d.IsNewResource() {
		if siteId, err := getModelId(config, d, "site"); err == nil {
			writablePrefix.Site = &siteId
		}
	}
	if d.HasChange("vrf") && !d.IsNewResource() {
		if vrfId, err := getModelId(config, d, "vrf"); err == nil {
			writablePrefix.Vrf = &vrfId
		}
	}
	if d.HasChange("tenant") && !d.IsNewResource() {
		if tenantId, err := getModelId(config, d, "tenant"); err == nil {
			writablePrefix.Tenant = &tenantId
		}
	}
	if d.HasChange("vlan") && !d.IsNewResource() {
		if vlanId, err := getModelId(config, d, "vlan"); err == nil {
			writablePrefix.Vlan = &vlanId
		}
	}
	if d.HasChange("role") && !d.IsNewResource() {
		if roleId, err := getModelId(config, d, "role"); err == nil {
			writablePrefix.Role = &roleId
		}
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	partialUpdatePrefix := ipam.IpamPrefixesPartialUpdateParams{
		ID:      int64(id),
		Data:    &writablePrefix,
		Context: context.Background(),
	}

	partialUpdatePrefixRes, _ := json.Marshal(partialUpdatePrefix)
	log.Println("resourceIpamAvailablePrefixesUpdate partialUpdatePrefix: ", string(partialUpdatePrefixRes))

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
	log.Println("[debug] ", res)
	if ok && resourceSchemaField != "" {
		return res.(string), nil
	}
	return "", fmt.Errorf("Cannot determine %s: set in this resource", resourceSchemaField)

}

func getIpamRoles(config *Config, d *schema.ResourceData) ([]*models.Role, error) {
	roleName, err := getAttrFromSchema("role", d, config)
	if err != nil {
		return nil, err
	}
	roleParam := ipam.IpamRolesListParams{
		Name:    &roleName,
		Limit:   &NetboxApigeneralQueryLimit,
		Context: context.Background(),
	}
	roleRes, err := config.client.Ipam.IpamRolesList(&roleParam, nil)
	if err != nil {
		fmt.Println("IpamRolesList ", err)
	}

	roleReadOKRes, _ := json.Marshal(&roleRes.Payload.Results)
	log.Println("roleReadOKRes ", string(roleReadOKRes))

	if roleRes == nil || *roleRes.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow role %s , not found", roleName)
	}
	return roleRes.Payload.Results, nil
}

func getDcimSites(config *Config, d *schema.ResourceData) ([]*models.Site, error) {
	siteName, err := getAttrFromSchema("site", d, config)
	if err != nil {
		return nil, err
	}
	siteParam := dcim.DcimSitesListParams{
		Name:    &siteName,
		Limit:   &NetboxApigeneralQueryLimit,
		Context: context.Background(),
	}
	siteRes, err := config.client.Dcim.DcimSitesList(&siteParam, nil)
	if err != nil {
		return nil, fmt.Errorf("DcimSitesListParams %s", err.Error())
	}

	if siteRes == nil || *siteRes.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow Site %s , not found", siteName)
	}
	return siteRes.Payload.Results, nil
}

func getIpamVlans(config *Config, d *schema.ResourceData) ([]*models.VLAN, error) {
	vlanName, err := getAttrFromSchema("vlan", d, config)
	if err != nil {
		return nil, err
	}
	vlanParam := ipam.IpamVlansListParams{
		Name:    &vlanName,
		Limit:   &NetboxApigeneralQueryLimit,
		Context: context.Background(),
	}
	vlanData, err := config.client.Ipam.IpamVlansList(&vlanParam, nil)
	if err != nil {
		return nil, fmt.Errorf("IpamVlansList %s", err.Error())
	}
	if vlanData == nil || *vlanData.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow vlan %s , not found", vlanName)
	}
	return vlanData.Payload.Results, nil
}

func getIpamVrfs(config *Config, d *schema.ResourceData) ([]*models.VRF, error) {
	vrfName, err := getAttrFromSchema("vrf", d, config)
	if err != nil {
		return nil, err
	}

	vrfParam := ipam.IpamVrfsListParams{
		Name:    &vrfName,
		Limit:   &NetboxApigeneralQueryLimit,
		Context: context.Background(),
	}
	vrfData, err := config.client.Ipam.IpamVrfsList(&vrfParam, nil)
	if err != nil {
		return nil, fmt.Errorf("IpamVlansList %s", err.Error())
	}
	if vrfData == nil || *vrfData.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow vrf %s , not found", vrfName)
	}
	return vrfData.Payload.Results, nil
}

func getTenancyTenant(config *Config, d *schema.ResourceData) ([]*models.Tenant, error) {
	tenantName, err := getAttrFromSchema("tenant", d, config)
	if err != nil {
		return nil, err
	}
	tenantParam := tenancy.TenancyTenantsListParams{
		Name:    &tenantName,
		Limit:   &NetboxApigeneralQueryLimit,
		Context: context.Background(),
	}
	tenantData, err := config.client.Tenancy.TenancyTenantsList(&tenantParam, nil)
	if err != nil {
		return nil, fmt.Errorf("TenancyTenantsList %s", err.Error())
	}

	if tenantData == nil || *tenantData.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow Tenant %s , not found", tenantName)
	}
	return tenantData.Payload.Results, nil
}

func getModelId(config *Config, d *schema.ResourceData, key string) (int64, error) {
	switch key {
	case "site":
		sites, err := getDcimSites(config, d)
		if err != nil {
			return 0, err
		}
		return sites[0].ID, nil
	case "role":
		roles, err := getIpamRoles(config, d)
		if err != nil {
			return 0, err
		}
		return roles[0].ID, nil
	case "vlan":
		vlans, err := getIpamVlans(config, d)
		if err != nil {
			return 0, err
		}
		return vlans[0].ID, nil
	case "vrf":
		vrfs, err := getIpamVrfs(config, d)
		if err != nil {
			return 0, err
		}
		return vrfs[0].ID, nil
	case "tenant":
		tenants, err := getTenancyTenant(config, d)
		if err != nil {
			return 0, err
		}
		return tenants[0].ID, nil
	default:
		return -1, fmt.Errorf("Uknown key %s", key)
	}
}
