package netbox

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Global MutexKV
var mutexKV = NewMutexKV()

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"NETBOX_TOKEN",
					"NETBOX_API_TOKEN",
					"API_TOKEN",
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

		DataSourcesMap: map[string]*schema.Resource{
			"netbox_available_prefixes": dataSourceIpamAvailablePrefixes(),
		},

		ResourcesMap: ResourceMap(),
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
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
		"netbox_available_prefixes": resourceIpamAvailablePrefixes(),
		//"ipam_prefixes_available_ips":
	}
}

func providerConfigure(d *schema.ResourceData, p *schema.Provider, terraformVersion string) (interface{}, diag.Diagnostics) {
	config := Config{
		ApiToken: d.Get("api_token").(string),
		Host:     d.Get("host").(string),
		BasePath: d.Get("base_path").(string),
	}

	if v, ok := d.GetOk("request_timeout"); ok {
		var err error
		config.RequestTimeout, err = time.ParseDuration(v.(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	if err := config.LoadAndValidate(context.Background()); err != nil {
		return nil, diag.FromErr(err)
	}

	return &config, nil
}
