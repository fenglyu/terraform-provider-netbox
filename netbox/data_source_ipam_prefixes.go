package netbox

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/fenglyu/go-netbox/netbox/client/ipam"
	"github.com/fenglyu/go-netbox/netbox/models"
)

func dataSourceIpamAvailablePrefixes() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceIpamAvailablePrefixes().Schema)

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "prefix")

	// Add "prefix_id" to support id passing
	dsSchema["prefix_id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "An identifier for the prefix",
	}

	return &schema.Resource{
		Read:   dataSourceIpamAvailablePrefixesRead,
		Schema: dsSchema,
	}
}

func dataSourceIpamAvailablePrefixesRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	var prefix *models.Prefix
	if v, ok := d.GetOk("prefix_id"); ok {

		id := v.(int)
		params := ipam.IpamPrefixesReadParams{
			ID: int64(id),
		}
		params.WithContext(context.Background())

		ipamPrefixesReadOK, err := config.client.Ipam.IpamPrefixesRead(&params, nil)
		if err != nil || ipamPrefixesReadOK == nil {
			return fmt.Errorf("Cannot determine prefix with ID %d", id)
		}
		prefix = ipamPrefixesReadOK.Payload
	} else if v, ok := d.GetOk("prefix"); ok {
		pStr := v.(string)

		withinInclude := pStr
		prefixLength, err := strconv.Atoi(strings.Split(withinInclude, "/")[1])
		if err != nil {
			return fmt.Errorf("Error in [getIpamPrefixes] %v", err)
		}
		maskLength := float64(prefixLength)
		param := ipam.IpamPrefixesListParams{
			MaskLength:    &maskLength,
			WithinInclude: &withinInclude,
		}
		param.WithContext(context.Background())
		ipamPrefixListBody, err := config.client.Ipam.IpamPrefixesList(&param, nil)
		if err != nil {
			return err
		}
		ipamPrefixesReadOKRes, _ := json.Marshal(&ipamPrefixListBody.Payload.Results)
		log.Println("ipamPrefixListBody", string(ipamPrefixesReadOKRes))
		if ipamPrefixListBody == nil || *ipamPrefixListBody.Payload.Count < 1 {
			//return fmt.Errorf("Unknow prefix %s not found", *prefix.Prefix)
			return fmt.Errorf("Unknow prefix %s not found", v)
			d.SetId("")
		}
		prefix = ipamPrefixListBody.Payload.Results[0]
	}

	jsonPrefix, _ := json.Marshal(prefix)
	log.Println("[INFO] dataSourceIpamPrefixesRead ", string(jsonPrefix))
	d.Set("description", prefix.Description)
	d.Set("custom_fields", prefix.CustomFields)
	d.Set("is_pool", prefix.IsPool)
	d.Set("created", prefix.Created.String())
	d.Set("family", prefix.Family)
	d.Set("last_updated", prefix.LastUpdated.String())

	if prefix != nil && prefix.Role != nil {
		d.Set("role", prefix.Role.ID)
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
