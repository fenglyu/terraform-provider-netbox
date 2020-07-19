package netbox

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func emptyOrDefaultStringSuppress(defaultVal string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return (old == "" && new == defaultVal) || (new == "" && old == defaultVal)
	}
}

func suppressEmptyCustomFieldsDiff(d *schema.ResourceDiff, meta interface{}) error {
	oldi, newi := d.GetChange("custom_fields")

	old, ok := oldi.([]interface{})
	if !ok {
		return fmt.Errorf("Expected old Custom Fields diff to be a slice")
	}

	new, ok := newi.([]interface{})
	if !ok {
		return fmt.Errorf("Expected new Custom Fields diff to be a slice")
	}

	if len(old) != 0 && len(new) != 1 {
		return nil
	}

	if len(old) != 1 && len(new) != 0 {
		return nil
	}

	if _, ok := new[0].(map[string]interface{}); !ok {
		return fmt.Errorf("Unable to type assert Custom Fields")
	}

	return nil
}
