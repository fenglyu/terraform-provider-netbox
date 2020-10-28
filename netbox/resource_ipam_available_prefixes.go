package netbox

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

	lockNamePrefix = "availableprefixes"
)

func resourceIpamAvailablePrefixes() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpamAvailablePrefixesCreate,
		ReadContext:   resourceIpamAvailablePrefixesRead,
		UpdateContext: resourceIpamAvailablePrefixesUpdate,
		DeleteContext: resourceIpamAvailablePrefixesDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIpamAvailablePrefixesImportState,
			//StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,

		// TODO when schema change happens
		// StateUpgraders:

		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(10 * time.Minute),
			Update:  schema.DefaultTimeout(10 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"parent_prefix": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				AtLeastOneOf:     availablePrefixesKeys,
				ValidateDiagFunc: IsCIDRNetworkDiagFunc(1, 128),
				Description:      "crave available prefixes under the parent_prefix",
			},
			"prefix": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: IsCIDRNetworkDiagFunc(1, 128),
				Description:      "craved available prefix",
			},
			"parent_prefix_id": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				AtLeastOneOf:     availablePrefixesKeys,
				ValidateDiagFunc: IntAtLeastDiagFunc(0),
				Description:      "A unique integer value identifying this prefix under which is used crave available prefix",
			},
			"prefix_length": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: IntBetweenDiagFunc(1, 128),
				Description:      "The mask in integer form",
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
				Type:             schema.TypeString,
				Default:          "active",
				Optional:         true,
				ValidateDiagFunc: StringInSliceDiagFunc(prefixinitializeStatus, false),
				Description:      "Operational status of this prefix",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: StringLenBetween(0, 200),
				Description:      "Describe the purpose of this prefix",
			},
			// Blizzard's custom_fields
			"custom_fields": {
				Type:     schema.TypeList,
				Required: true,
				//Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				//ForceNew:    true,
				MaxItems:    1,
				Description: "Set of customized key/value pairs created for prefix.",
				/*				*/
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"helpers": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "helpers (to be explained)",
							Default:          "",
							DiffSuppressFunc: emptyOrDefaultStringSuppress(""),
						},
						"ipv4_acl_in": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "ipv4_acl_in (to be explained)",
							Default:          "",
							DiffSuppressFunc: emptyOrDefaultStringSuppress(""),
						},
						"ipv4_acl_out": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "ipv4_acl_out (to be explained)",
							Default:          "",
							DiffSuppressFunc: emptyOrDefaultStringSuppress(""),
						},
					},
				},
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

		CustomizeDiff: customdiff.All(
			customdiff.If(
				func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
					return d.HasChange("custom_fields")
				},
				suppressEmptyCustomFieldsDiff,
			),
		),
	}
}

func resourceIpamAvailablePrefixesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	//prarent_prefix here is only needed to fetch prarent_prefix_id only
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
		if site != nil && site.Tenant != nil {
			//return fmt.Errorf("Site %s or its tenant not exists", d.Get("site").(string))
			if *site.Tenant.Name != tenantData.(string) {
				return diag.Errorf("Incompatible site %s and the tenant %s, expected tenant %s", *site.Name, tenantData.(string), *site.Tenant.Name)
			}
			tenant = site.Tenant.ID
			wPrefix.Tenant = &tenant
		}
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
		wPrefix.Status = status
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
		if cfMap, err := expandCustomFields(d, cfData); err == nil {
			customFields = cfMap
			wPrefix.CustomFields = customFields
		} else {
			log.Println(err)
		}
	}

	// If parent prefix is given
	wPrefixRes, _ := json.Marshal(wPrefix)
	log.Println("[INFO] ", string(wPrefixRes))

	if _, ok := getParentPrefix(config, d); ok == nil {
		if wPrefix.ID == 0 {
			results, err := getIpamPrefixes(config, d)
			if err != nil {
				return diag.FromErr(err)
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
	// Lock/Unlock have been deprecated, Rewrite them after migrated to sdk v2
	mutexKV.Lock(fmt.Sprintf("%s_%d", lockNamePrefix, prefix_id))
	defer mutexKV.Unlock(fmt.Sprintf("%s_%d", lockNamePrefix, prefix_id))
	res, err := config.client.Ipam.IpamPrefixesAvailablePrefixesCreate(&param, nil)
	if err != nil {
		// The resource didn't actually create
		log.Println("[Error] Failed to create AvaliablePrefix: ", err)
		d.SetId("")
		// work around  netbox insufficient space response problem which returns 204 http code with a message body
		// although according to https://tools.ietf.org/html/rfc2616#section-10.2.5, The 204 response MUST NOT include a message-body,
		// and thus is always terminated by the first empty line after the header fields.
		if strings.Contains(err.Error(), "204") {
			log.Printf("[WARN] Insufficient space is available to accommodate the requested prefix size(s) \"/%d\"", prefixlength)
			return diag.Errorf("Insufficient space is available to accommodate the requested prefix size(s) \"/%d\"", prefixlength)
		}
		return diag.FromErr(err)
	}

	availablePrefix := res.GetPayload()
	d.SetId(fmt.Sprintf("%d", availablePrefix.ID))

	return resourceIpamAvailablePrefixesRead(ctx, d, m)
}

func resourceIpamAvailablePrefixesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	log.Println("d state", d.State())
	log.Println("custom_fields.# ", d.Get("custom_fields.#"))

	prefix, err := getIpamPrefix(config, d)
	if err != nil || prefix == nil {
		return diag.FromErr(err)
	}

	log.Println("[INFO] resourceIpamPrefixesRead ", prefix)
	//d.Set("id", prefix.ID)
	d.Set("description", prefix.Description)

	log.Println("CustomFields: ", flatterCustomFields(d, prefix.CustomFields))

	if prefix != nil && prefix.CustomFields != nil {
		d.Set("custom_fields", flatterCustomFields(d, prefix.CustomFields))
	}

	log.Println("d state", d.State())
	log.Println("custom_fields.# ", d.Get("custom_fields.#"))

	d.Set("is_pool", prefix.IsPool)
	d.Set("created", prefix.Created.String())
	d.Set("family", prefix.Family.Value)
	d.Set("last_updated", prefix.LastUpdated.String())

	if prefix != nil && prefix.Role != nil {
		d.Set("role", prefix.Role.Name)
	}

	if prefix.Prefix != nil && *prefix.Prefix != "" {
		parentPrefix, err := getIpamParentPrefixes(config, d, prefix)
		if err != nil || parentPrefix == nil {
			return diag.FromErr(err)
		}
		if parentPrefix != nil && *parentPrefix.Prefix != "" {
			if _, ok := d.GetOk("parent_prefix_id"); ok {
				d.Set("parent_prefix_id", int(parentPrefix.ID))
			}

			if _, ok := d.GetOk("parent_prefix"); ok {
				d.Set("parent_prefix", parentPrefix.Prefix)
			}
		}
	}

	d.Set("prefix", prefix.Prefix)
	if prefix.Prefix != nil && *prefix.Prefix != "" {
		pl := strings.Split(*prefix.Prefix, "/")[1]
		prefixLength, _ := strconv.Atoi(pl)
		d.Set("prefix_length", prefixLength)
	}

	if prefix != nil && prefix.Site != nil {
		d.Set("site", prefix.Site.Name)
	} else {
		d.Set("site", "")
	}
	if prefix != nil && prefix.Status != nil {
		d.Set("status", *prefix.Status.Value)
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

func resourceIpamAvailablePrefixesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
			return diag.Errorf("Not a valid status in %v", prefixinitializeStatus)
		}
		writablePrefix.Status = strings.ToLower(statusData)
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
		cfData := d.Get("custom_fields").([]interface{})
		log.Println("cfData	", cfData)
		if cfMap, err := expandCustomFields(d, cfData); err == nil {
			writablePrefix.CustomFields = cfMap
		} else {
			log.Println(err)
		}
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
		return diag.FromErr(err)
	}
	partialUpdatePrefix := ipam.IpamPrefixesPartialUpdateParams{
		ID:      int64(id),
		Data:    &writablePrefix,
		Context: context.Background(),
	}

	partialUpdatePrefixRes, _ := json.Marshal(partialUpdatePrefix)
	log.Println("resourceIpamAvailablePrefixesUpdate partialUpdatePrefix: ", string(partialUpdatePrefixRes))

	mutexKV.Lock(fmt.Sprintf("%s_%d", lockNamePrefix, id))
	defer mutexKV.Unlock(fmt.Sprintf("%s_%d", lockNamePrefix, id))

	res, uerr := config.client.Ipam.IpamPrefixesPartialUpdate(&partialUpdatePrefix, nil)
	if uerr != nil {
		// TODO Support verbose response body here
		return diag.Errorf("%v %v", res, uerr)
	}

	return resourceIpamAvailablePrefixesRead(ctx, d, m)
}

func resourceIpamAvailablePrefixesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	log.Printf("[INFO]Requesting Prefix deletion: %s", d.Get("prefix").(string))
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	params := ipam.IpamPrefixesDeleteParams{
		ID: int64(id),
	}
	params.WithContext(context.Background())
	mutexKV.Lock(fmt.Sprintf("%s_%d", lockNamePrefix, id))
	defer mutexKV.Unlock(fmt.Sprintf("%s_%d", lockNamePrefix, id))

	_, derr := config.client.Ipam.IpamPrefixesDelete(&params, nil)
	if derr != nil {
		return diag.FromErr(derr)
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
	param := ipam.IpamPrefixesListParams{
		MaskLength:    &maskLength,
		WithinInclude: &withinInclude,
		Prefix:        &prefix,
	}
	param.WithContext(context.Background())
	ipamPrefixListBody, err := config.client.Ipam.IpamPrefixesList(&param, nil)
	if err != nil {
		return nil, err
	}

	if ipamPrefixListBody == nil || ipamPrefixListBody.Payload == nil || *ipamPrefixListBody.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow prefix %s with ID %s, not found", prefix, d.Id())
	}
	// trace level log
	ipamPrefixesReadOKRes, _ := json.Marshal(&ipamPrefixListBody.Payload.Results)
	log.Println("[getIpamPrefixes] ipamPrefixListBody", string(ipamPrefixesReadOKRes))

	return ipamPrefixListBody.Payload.Results, nil
}

func getIpamParentPrefixes(config *Config, d *schema.ResourceData, prefix *models.Prefix) (*models.Prefix, error) {
	// Compose Parameters for GET: /ipam/prefixes/
	// The api call to get parent prefix is like: /api/ipam/prefixes/?contains=10.1.0.0/16&vrf_id=null
	// The WebUI parent prefix fetch pretty much uses the same DB query logic: https://github.com/netbox-community/netbox/blob/develop/netbox/ipam/views.py#L342-L349
	param := ipam.IpamPrefixesListParams{
		Context:  context.Background(),
		Contains: prefix.Prefix,
	}

	var vrfID string
	if prefix.Vrf != nil {
		vrfID = strconv.FormatInt(prefix.Vrf.ID, 10)
		param.VrfID = &vrfID
	} else {
		param.VrfID = &vrfID
	}

	ipamPrefixListBody, err := config.client.Ipam.IpamPrefixesList(&param, nil)
	if err != nil {
		return nil, err
	}

	if ipamPrefixListBody == nil || ipamPrefixListBody.Payload == nil || *ipamPrefixListBody.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow prefix %s with ID %s, not found", *prefix.Prefix, d.Id())
	} else if *ipamPrefixListBody.Payload.Count < 2 {
		return nil, fmt.Errorf("prefix %s with ID %s has no parent prefix", *prefix.Prefix, d.Id())
	}
	// trace level log
	ipamPrefixesReadOKRes, _ := json.Marshal(&ipamPrefixListBody.Payload.Results)
	log.Println("[getIpamParentPrefixes] ipamPrefixListBody", string(ipamPrefixesReadOKRes))

	var parent *models.Prefix
	for _, p := range ipamPrefixListBody.Payload.Results {
		if !strings.EqualFold(*p.Prefix, *prefix.Prefix) {
			parent = p
		}
	}
	return parent, nil
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
		Limit:   &NetboxApiGeneralQueryLimit,
		Context: context.Background(),
	}
	roleRes, err := config.client.Ipam.IpamRolesList(&roleParam, nil)
	if err != nil {
		fmt.Println("IpamRolesList ", err)
	}

	if roleRes == nil || roleRes.Payload == nil || *roleRes.Payload.Count < 1 {
		return nil, fmt.Errorf("Unknow role %s , not found", roleName)
	}
	// trace level log
	roleReadOKRes, _ := json.Marshal(&roleRes.Payload.Results)
	log.Println("roleReadOKRes ", string(roleReadOKRes))

	return roleRes.Payload.Results, nil
}

func getDcimSites(config *Config, d *schema.ResourceData) ([]*models.Site, error) {
	siteName, err := getAttrFromSchema("site", d, config)
	if err != nil {
		return nil, err
	}
	siteParam := dcim.DcimSitesListParams{
		Name:    &siteName,
		Limit:   &NetboxApiGeneralQueryLimit,
		Context: context.Background(),
	}
	siteRes, err := config.client.Dcim.DcimSitesList(&siteParam, nil)
	if err != nil {
		return nil, fmt.Errorf("DcimSitesListParams %s", err.Error())
	}

	if siteRes == nil || siteRes.Payload == nil || *siteRes.Payload.Count < 1 {
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
		Limit:   &NetboxApiGeneralQueryLimit,
		Context: context.Background(),
	}
	vlanData, err := config.client.Ipam.IpamVlansList(&vlanParam, nil)
	if err != nil {
		return nil, fmt.Errorf("IpamVlansList %s", err.Error())
	}
	if vlanData == nil || vlanData.Payload == nil || *vlanData.Payload.Count < 1 {
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
		Limit:   &NetboxApiGeneralQueryLimit,
		Context: context.Background(),
	}
	vrfData, err := config.client.Ipam.IpamVrfsList(&vrfParam, nil)
	if err != nil {
		return nil, fmt.Errorf("IpamVlansList %s", err.Error())
	}
	if vrfData == nil || vrfData.Payload == nil || *vrfData.Payload.Count < 1 {
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
		Limit:   &NetboxApiGeneralQueryLimit,
		Context: context.Background(),
	}
	tenantData, err := config.client.Tenancy.TenancyTenantsList(&tenantParam, nil)
	if err != nil {
		return nil, fmt.Errorf("TenancyTenantsList %s", err.Error())
	}

	if tenantData == nil || tenantData.Payload == nil || *tenantData.Payload.Count < 1 {
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

func resourceIpamAvailablePrefixesImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// config := meta.(*Config)
	log.Println("resourceIpamAvailablePrefixesImportState ", d.Get("custom_fields"))
	if _, ok := d.GetOk("custom_fields"); !ok {
		d.Set("custom_fields", nil)
	}

	return []*schema.ResourceData{d}, nil
}
