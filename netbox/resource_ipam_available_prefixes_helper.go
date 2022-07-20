package netbox

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

func expandCustomFields(d *schema.ResourceData, v interface{}) (map[string]interface{}, error) {
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

func flatterCustomFields(d *schema.ResourceData, v interface{}) []map[string]interface{} {
	cfs := make([]map[string]interface{}, 0)
	if v == nil {
		return nil
	}
	cf, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]interface{})
	log.Println("flatterCustomFields  ", d.Get("custom_fields"))

	if _, ok := cf["helpers"]; ok {
		result["helpers"] = cf["helpers"]
	}
	if _, ok := cf["ipv4_acl_in"]; ok {
		result["ipv4_acl_in"] = cf["ipv4_acl_in"]
	}
	if _, ok := cf["ipv4_acl_out"]; ok {
		result["ipv4_acl_out"] = cf["ipv4_acl_out"]
	}

	log.Println("result  ", result)
	cfs = append(cfs, result)
	return cfs
}

func flatternDatasourceCF(d *schema.ResourceData, v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	cf, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return []map[string]interface{}{cf}
}

// Same as schema.IsCIDRNetwork
func IsCIDRNetworkDiagFunc(min, max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) (diags diag.Diagnostics) {
		v, ok := i.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected type of %v to be string", path),
				AttributePath: path,
			})
			return diags
		}

		_, ipnet, err := net.ParseCIDR(v)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected %v to contain a valid Value, got: %s with err: %s", path, v, err),
				AttributePath: path,
			})
			return diags
		}

		if ipnet == nil || v != ipnet.String() {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary: fmt.Sprintf("expected %v to contain a valid network Value, expected %s, got %s",
					path, ipnet, v),
				AttributePath: path,
			})
		}

		sigbits, _ := ipnet.Mask.Size()
		if sigbits < min || sigbits > max {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected %v to contain a network Value with between %d and %d significant bits, got: %d", path, min, max, sigbits),
				AttributePath: path,
			})
		}

		return diags
	}
}

func IntAtLeastDiagFunc(min int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) (diags diag.Diagnostics) {
		v, ok := i.(int)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected type of %v to be integer", path),
				AttributePath: path,
			})
			return diags
		}
		if v < min {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected %v to be at least (%d), got %d", path, min, v),
				AttributePath: path,
			})
			return diags
		}

		return diags
	}
}

func IntBetweenDiagFunc(min, max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) (diags diag.Diagnostics) {
		v, ok := i.(int)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected type of %v to be integer", path),
				AttributePath: path,
			})
			return diags
		}

		if v < min || v > max {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected %v to be in the range (%d - %d), got %d", path, min, max, v),
				AttributePath: path,
			})
			return diags
		}

		return diags
	}
}

func StringInSliceDiagFunc(valid []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) (diags diag.Diagnostics) {
		v, ok := i.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected type of %s to be string", path),
				AttributePath: path,
			})
			return diags
		}

		for _, str := range valid {
			if v == str || (ignoreCase && strings.EqualFold(v, str)) {
				return diags
			}
		}

		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("expected %s to be one of %v, got %s", path, valid, v),
			AttributePath: path,
		})
		return diags
	}
}

// StringLenBetweenDiagFunc returns a SchemaValidateFunc which tests if the provided value
// is of type string and has length between min and max (inclusive)
func StringLenBetween(min, max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) (diags diag.Diagnostics) {
		v, ok := i.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected type of %s to be string", path),
				AttributePath: path,
			})
			return diags
		}

		if len(v) < min || len(v) > max {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("expected length of %v to be in the range (%d - %d), got %s", path, min, max, v),
				AttributePath: path,
			})
		}

		return diags
	}
}
