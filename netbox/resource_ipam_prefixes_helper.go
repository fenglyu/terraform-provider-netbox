package netbox

import (
	"github.com/fenglyu/go-netbox/netbox/models"
)

func flatternFamily(f *models.PrefixFamily) []map[string]interface{} {
	if f == nil {
		return nil
	}
	return []map[string]interface{}{{
		"label": f.Label,
		"value": f.Value,
	}}
}

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

func flatternNestedVLAN(nv *models.NestedVLAN) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":           nv.ID,
		"display_name": nv.DisplayName,
		"name":         nv.Name,
		"url":          nv.URL.String(),
		"vid":          nv.Vid,
	}}
}

func flatternNestedVRF(nv *models.NestedVRF) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":           nv.ID,
		"prefix_count": nv.PrefixCount,
		"name":         nv.Name,
		"url":          nv.URL.String(),
		"rd":           nv.Rd,
	}}
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
