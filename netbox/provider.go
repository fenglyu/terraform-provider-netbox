package netbox

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"time"
)

func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
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
			"request_timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: ResourceMap(),
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, provider, terraformVersion)
	}
	return provider
}

func ResourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"ipam_available_prefixes": resourceIpamPrefixes(),
		//"ipam_prefixes_available_ips":
	}
}

func providerConfigure(d *schema.ResourceData, p *schema.Provider, terraformVersion string) (interface{}, error) {
	config := Config{
		ApiToken: d.Get("api_token").(string),
		Host:     d.Get("host").(string),
		BasePath: d.Get("base_path").(string),
	}

	if v, ok := d.GetOk("request_timeout"); ok {
		var err error
		config.RequestTimeout, err = time.ParseDuration(v.(string))
		if err != nil {
			return nil, err
		}
	}

	if err := config.LoadAndValidate(p.StopContext()); err != nil {
		return nil, err
	}

	return &config, nil
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
