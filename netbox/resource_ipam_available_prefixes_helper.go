package netbox

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/fenglyu/go-netbox/netbox/models"
)

/*
func flatternFamily(f *models.PrefixFamily) []map[string]interface{} {
	if f == nil {
		return nil
	}
	return []map[string]interface{}{{
		"label": f.Label,
		"value": f.Value,
	}}
}

// nested_role struct for netbox v2.8.6
func flatternRole(nr *models.NestedRole) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":           nr.ID,
		"name":         nr.Name,
		"prefix_count": nr.PrefixCount,
		"slug":         nr.Slug,
		"url":          nr.URL.String(),
		"vlan_count":   nr.VlanCount,
	}}
}
*/

// NestedRole struct for netbox v2.4.7
func flatternRoleV247(nr *models.NestedRole) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":   nr.ID,
		"name": nr.Name,
		"slug": nr.Slug,
		"url":  nr.URL.String(),
	}}
}

func flatternSite(ns *models.NestedSite) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":   ns.ID,
		"name": ns.Name,
		"slug": ns.Slug,
		"url":  ns.URL.String(),
	}}
}

func flatterPrefixStatus(ps *models.PrefixStatus) []map[string]interface{} {
	return []map[string]interface{}{{
		"label": ps.Label,
		"value": ps.Value,
	}}
}

func flatternNestedTenant(nt *models.NestedTenant) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"id":   nt.ID,
			"name": nt.Name,
			"slug": nt.Slug,
			"url":  nt.URL.String(),
		},
	}
}

func jsonfy(rs []map[string]interface{}) string {
	st, err := json.Marshal(rs)
	if err != nil {
		return ""
	}
	return string(st)
}

func flatternNestedVLAN(nv *models.NestedVLAN) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":           nv.ID,
		"display_name": nv.DisplayName,
		"name":         nv.Name,
		"url":          nv.URL.String(),
		"vid":          nv.Vid,
	}}
}

// NestedVRF struct for netbox v2.4.7
func flatternNestedVRFV247(nv *models.NestedVRF) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":   nv.ID,
		"name": nv.Name,
		"url":  nv.URL.String(),
		"rd":   nv.Rd,
	}}
}

/*
func flatternNestedVRF(nv *models.NestedVRF) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":           nv.ID,
		"prefix_count": nv.PrefixCount,
		"name":         nv.Name,
		"url":          nv.URL.String(),
		"rd":           nv.Rd,
	}}
}


func flattenCustomFields(p *models.Prefix) map[string]string {
	cf := p.CustomFields.(map[string]interface{})
	cfMap := make(map[string]string)
	for k, v := range cf {
		cfMap[k] = v.(string)
	}
	return cfMap
}


type CustomFields struct {
	Helpers      string `json:"helpers"`
	Ipv4_acl_in  string `json:"ipv4_acl_in"`
	Ipv4_acl_out string `json:"ipv4_acl_out"`
}

*/
func convertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	sort.Strings(s)

	return s
}

func expandCustomFields(v interface{}) (map[string]interface{}, error) {
	if v == nil {
		// We can't set default values for lists.
		return nil, nil
	}

	ls := v.([]interface{})
	cf := make(map[string]interface{}, len(ls))

	if len(ls) == 0 {
		// We can't set default values for lists
		return cf, nil
	}

	if len(ls) > 1 || ls[0] == nil {
		return nil, fmt.Errorf("expected exactly one custom field")
	}

	original := ls[0].(map[string]interface{})
	if v, ok := original["helpers"]; ok {
		cf["helpers"] = v.(string)
	}

	if v, ok := original["ipv4_acl_in"]; ok {
		cf["ipv4_acl_in"] = v.(string)
	}

	if v, ok := original["ipv4_acl_out"]; ok {
		cf["ipv4_acl_out"] = v.(string)
	}
	return cf, nil
}

func flatterCustomFields(v interface{}) []map[string]interface{} {
	cf := make([]map[string]interface{}, 0)
	if v == nil {
		return nil
	}
	cf = append(cf, v.(map[string]interface{}))
	return cf
}
