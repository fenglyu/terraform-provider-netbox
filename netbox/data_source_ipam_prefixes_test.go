package netbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceAvailablePrefixes_basic(t *testing.T) {

	context := map[string]interface{}{
		// vrf global
		//"parent_prefix_id":     502,
		// vrf "activison", choose the right parent_prefix_id if you enabled vrf
		"parent_prefix_id":     627,
		"random_prefix_length": randIntRange(t, 16, 30),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAvailablePrefixesConfigByPrefix(context),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAvailablePrefixesCheck("data.netbox_available_prefixes.bar", "netbox_available_prefixes.foo"),
					resource.TestCheckResourceAttr("data.netbox_available_prefixes.bar", "is_pool", "true"),
					resource.TestCheckResourceAttr("data.netbox_available_prefixes.bar", "status", "active"),
					resource.TestCheckResourceAttr("data.netbox_available_prefixes.bar", "family", "4"),
					resource.TestCheckResourceAttr("data.netbox_available_prefixes.bar", "description", "testAccDataSourceComputeInstanceConfig description"),
					// Tags to be tested
					//resource.TestCheckResourceAttr("data.netbox_available_prefixes.bar", "tags", "[\"datasource-AvailablePrefix-acc01\", \"datasource-AvailablePrefix-acc02\", \"datasource-AvailablePrefix-acc02\"]"),
				),
			},
			{
				Config: testAccDataSourceAvailablePrefixesConfigByPrefixId(context),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAvailablePrefixesCheck("data.netbox_available_prefixes.bar", "netbox_available_prefixes.foo"),
					resource.TestCheckResourceAttr("data.netbox_available_prefixes.bar", "is_pool", "true"),
					resource.TestCheckResourceAttr("data.netbox_available_prefixes.bar", "status", "active"),
					resource.TestCheckResourceAttr("data.netbox_available_prefixes.bar", "family", "4"),
					resource.TestCheckResourceAttr("data.netbox_available_prefixes.bar", "description", "testAccDataSourceComputeInstanceConfig description"),
				),
			},
		},
	})
}

func testAccDataSourceAvailablePrefixesCheck(datasourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[datasourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", datasourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		datasourceAttributes := ds.Primary.Attributes
		resourceAttributes := rs.Primary.Attributes

		instanceAttrsToTest := []string{
			"prefix",
			"prefix_length",
			"is_pool",
			"status",
			"description",
			"tags",
			"created",
			"custom_fields",
			"family",
			"last_updated",
			"role",
			"site",
			"tenant",
			"vlan",
			"vrf",
		}

		for _, attrToCheck := range instanceAttrsToTest {
			if datasourceAttributes[attrToCheck] != resourceAttributes[attrToCheck] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attrToCheck,
					datasourceAttributes[attrToCheck],
					resourceAttributes[attrToCheck],
				)
			}
		}

		return nil
	}
}

func testAccDataSourceAvailablePrefixesConfigByPrefix(config map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "foo" {
  parent_prefix_id 	= %{parent_prefix_id}
  prefix_length 	= %{random_prefix_length}
  is_pool          	= true
  status          	= "active"
  role = "gcp"
  site = "se1"
  vlan = "gcp"
  vrf  = "activision"
  tenant = "cloud"

  custom_fields  {}
  description = "testAccDataSourceComputeInstanceConfig description"
  tags        = ["datasource-AvailablePrefix-acc01", "datasource-AvailablePrefix-acc02", "datasource-AvailablePrefix-acc03"]
}

data "netbox_available_prefixes" "bar"{
  prefix = netbox_available_prefixes.foo.prefix
}
`, config)
}

func testAccDataSourceAvailablePrefixesConfigByPrefixId(config map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "foo" {
  parent_prefix_id 	= %{parent_prefix_id}
  prefix_length 	= %{random_prefix_length}
  is_pool          	= true
  status          	= "active"
  role = "gcp"
  site = "se1"
  vlan = "gcp"
  vrf  = "activision"
  tenant = "cloud"

  description = "testAccDataSourceComputeInstanceConfig description"
  tags        = ["datasource-AvailablePrefix-acc06", "datasource-AvailablePrefix-acc04", "datasource-AvailablePrefix-acc05"]
  custom_fields  {}


}

data "netbox_available_prefixes" "bar"{
  prefix_id = netbox_available_prefixes.foo.id
}
`, config)
}

func testAccCheckDataSourceAvailablePrefixesDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "netbox_available_prefixes" {
				continue
			}

			config := testAccProvider.Meta().(*Config)
			_, err := sendRequestforPrefix(config, rs)
			fmt.Println("testAccCheckAvailablePrefixesDestroyProducer: ", err)

			if err == nil {
				return fmt.Errorf("Available Prefix still exists")
			}
		}
		return nil
	}
}
