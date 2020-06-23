package netbox

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"API_TOKEN",
					"NETBOX_API_TOKEN",
				}, nil),
				//ValidateFunc: validateCredentials,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"NETBOX_HOST",
				}, NetboxDefaultHost),
			},
			"base_path": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"NETBOX_BASE_PATH",
				}, NetboxDefaultBasePath),
			},
		},
		ResourcesMap: ResourceMap(),
	}
}

func ResourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"ipam_available_prefixes": resourceIpamPrefixes(),
		//"ipam_prefixes_available_ips":
	}
}

/*
func validateCredentials(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil ||v.(string) == ""{
		return
	}

	apiToken :=v.(string)
	// more validate logic
	if _, err := netbox.NewNetboxWithAPIKey(host, apiToken); err != nil{
		errors = append(errors, fmt.Errorf("credentials are not valid: %s", err))
	}
	return
}
*/
