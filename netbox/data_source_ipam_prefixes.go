package netbox

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/fenglyu/go-netbox/netbox/client/ipam"
)

func dataSourceIpamAvailablePrefixes() *schema.Resource {

	prefixSchema := datasourceSchemaFromResourceSchema(resourceIpamAvailablePrefixes().Schema)
	// Add prefix id to prefix output

	prefixSchema["id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Optional:    true,
		Description: "An identifier for the prefix",
	}

	// This is a schema Element that will allow us to read and place all returned prefixes into the
	// `prefixes` attribute.
	return &schema.Resource{
		Read: dataSourceIpamAvailablePrefixesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of data lookup",
			},
			"contains": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "contains query",
			},
			"family": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix family",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search for a prefix by ID",
			},
			"id_in": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search for a prefix with a set of IDs",
			},
			"is_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix is a pool",
			},
			"mask_length": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Max Length of Prefix",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Limit the number of returned results",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The initial index from which to return the results.",
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Full Prefix CIDR to find",
			},
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query String",
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix Role",
			},
			"role_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix Role ID",
			},
			"site": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix Site",
			},
			"site_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix Site ID",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix Status",
			},
			"tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix tag",
			},
			"tenant": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix tenant",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix Tenant ID",
			},
			"vlan_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix Vlan ID",
			},
			"vlan_vid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix Vlan VID",
			},
			"vrf": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix VRF",
			},
			"vrf_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix VRF ID",
			},
			"within": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix Within Query",
			},
			"within_include": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prefix Within Include",
			},
			"prefixes": {
				Type:     schema.TypeList,
				Required: false,
				Computed: true,
				Elem: &schema.Resource{
					Schema: prefixSchema,
				},
			},
		},
	}
}

func dataSourceIpamAvailablePrefixesRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	// construct a prefix query
	param := ipam.IpamPrefixesListParams{}

	if v, ok := d.GetOk("contains"); ok {
		contains := v.(string)
		param.SetContains(&contains)
	}

	if v, ok := d.GetOk("family"); ok {
		family := v.(string)
		param.SetFamily(&family)
	}

	if v, ok := d.GetOk("id"); ok {
		id := v.(string)
		param.SetIDIn(&id)
	}

	if v, ok := d.GetOk("id_in"); ok {
		idIn := v.(string)
		param.SetIDIn(&idIn)
	}

	if v, ok := d.GetOk("is_pool"); ok {
		isPool := v.(string)
		param.SetIsPool(&isPool)
	}

	if v, ok := d.GetOk("limit"); ok {
		limit := v.(int64)
		param.SetLimit(&limit)
	}

	if v, ok := d.GetOk("mask_length"); ok {
		maskLength := v.(float64)
		param.SetMaskLength(&maskLength)
	}

	if v, ok := d.GetOk("offset"); ok {
		offset := v.(int64)
		param.SetOffset(&offset)
	}

	if v, ok := d.GetOk("prefix"); ok {
		prefix := v.(string)

		prefixLength, err := strconv.Atoi(strings.Split(prefix, "/")[1])
		if err != nil {
			return fmt.Errorf("Error parsing prefix parameter %v", err)
		}

		maskLength := float64(prefixLength)

		param.SetWithinInclude(&prefix)
		param.SetMaskLength(&maskLength)
	}

	if v, ok := d.GetOk("q"); ok {
		query := v.(string)
		param.SetQ(&query)
	}

	if v, ok := d.GetOk("role"); ok {
		role := v.(string)
		param.SetRole(&role)
	}

	if v, ok := d.GetOk("role_id"); ok {
		roleID := v.(string)
		param.SetRoleID(&roleID)
	}

	if v, ok := d.GetOk("site"); ok {
		site := v.(string)
		param.SetSite(&site)
	}

	if v, ok := d.GetOk("site_id"); ok {
		siteID := v.(string)
		param.SetSiteID(&siteID)
	}

	if v, ok := d.GetOk("status"); ok {
		status := v.(string)
		param.SetStatus(&status)
	}

	if v, ok := d.GetOk("tag"); ok {
		tag := v.(string)
		param.SetTag(&tag)
	}

	if v, ok := d.GetOk("tenant"); ok {
		tenant := v.(string)
		param.SetTenant(&tenant)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		tenantID := v.(string)
		param.SetTenantID(&tenantID)
	}

	if v, ok := d.GetOk("vlan_id"); ok {
		vlanID := v.(string)
		param.SetVlanID(&vlanID)
	}

	if v, ok := d.GetOk("vlan_vid"); ok {
		vlanVID := v.(float64)
		param.SetVlanVid(&vlanVID)
	}

	if v, ok := d.GetOk("vrf"); ok {
		vrf := v.(string)
		param.SetVrf(&vrf)
	}

	if v, ok := d.GetOk("vrf_id"); ok {
		vrfID := v.(string)
		param.SetVrfID(&vrfID)
	}

	if v, ok := d.GetOk("within"); ok {
		within := v.(string)
		param.SetWithin(&within)
	}

	if v, ok := d.GetOk("within_include"); ok {
		withinInclude := v.(string)
		param.SetWithinInclude(&withinInclude)
	}

	param.WithContext(context.Background())
	ipamPrefixListBody, err := config.client.Ipam.IpamPrefixesList(&param, nil)
	if err != nil {
		return err
	}

	// Container to store results
	prefixes := make([]map[string]interface{}, 0)

	prefixIdList := make([]string, 0)
	for _, prefix := range ipamPrefixListBody.Payload.Results {
		data := map[string]interface{}{}
		data["description"] = prefix.Description
		data["custom_fields"] = flatternDatasourceCF(d, prefix.CustomFields)
		data["is_pool"] = prefix.IsPool
		data["created"] = prefix.Created.String()
		data["family"] = prefix.Family
		data["last_updated"] = prefix.LastUpdated.String()
		data["prefix"] = prefix.Prefix
		data["status"] = prefixStatusIDMapReverse[*prefix.Status.Value]
		data["tags"] = prefix.Tags
		data["id"] = prefix.ID

		if prefix.Site != nil {
			data["site"] = prefix.Site.Name
		}
		if prefix.Tenant != nil {
			data["tenant"] = prefix.Tenant.Name
		}
		if prefix.Role != nil {
			data["role"] = prefix.Role.Name
		}
		if prefix.Vlan != nil {
			data["vlan"] = prefix.Vlan.Name
		}
		if prefix.Vrf != nil {
			data["vrf"] = prefix.Vrf.Name
		}

		pl, err := strconv.Atoi(strings.Split(*prefix.Prefix, "/")[1])
		if err != nil {
			return fmt.Errorf("Error parsing *prefix.Prefix parameter %v", err)
		}
		data["prefix_length"] = pl

		prefixes = append(prefixes, data)
		prefixIdList = append(prefixIdList, fmt.Sprintf("%d", prefix.ID))
	}

	if err := d.Set("prefixes", prefixes); err != nil {
		return fmt.Errorf("Error retrieving prefixes: %s", err)
	}

	if v, ok := d.GetOk("id"); ok {
		id := v.(string)
		d.SetId(id)
	} else {
		//d.SetId(d.Get("name").(string))
		// For multiple reseults, Name in the form of "{name}/id0_id1_id2"
		d.SetId(fmt.Sprintf("%s/%s", d.Get("name").(string), strings.Join(prefixIdList, "_")))
	}

	return nil
}
