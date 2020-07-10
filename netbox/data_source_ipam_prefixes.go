package netbox

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/fenglyu/go-netbox/netbox/client/ipam"
	"github.com/fenglyu/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIpamAvailablePrefixes() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceIpamAvailablePrefixes().Schema)

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "prefix")

	return &schema.Resource{
		Read:   dataSourceIpamAvailablePrefixesRead,
		Schema: dsSchema,
	}
}

func dataSourceIpamAvailablePrefixesRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	var prefix *models.Prefix
	if v, ok := d.GetOk("prefix"); ok {
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
			return fmt.Errorf("Unknow prefix %s with ID %s, not found", *prefix.Prefix, d.Id())
		}
		prefix = ipamPrefixListBody.Payload.Results[0]
	}

	jsonPrefix, _ := json.Marshal(prefix)
	log.Println("[INFO] dataSourceIpamPrefixesRead ", string(jsonPrefix))
	d.Set("description", prefix.Description)
	d.Set("custom_fields", flattenCustomFields(prefix))
	d.Set("is_pool", prefix.IsPool)
	d.Set("created", prefix.Created)
	d.Set("family", prefix.Family)
	if prefix != nil && prefix.Role != nil {
		d.Set("role", flatternRoleV247(prefix.Role))
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
		d.Set("vrf", flatternNestedVRFV247(prefix.Vrf))
	}
	d.SetId(fmt.Sprintf("%d", prefix.ID))

	return nil
}
