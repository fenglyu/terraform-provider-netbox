package netbox

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAvailablePrefixesByPrefix(t *testing.T) {

	context := map[string]interface{}{
		"parent_prefix_id":     testNetboxParentPrefixIdWithVrf,
		"random_prefix_length": randIntRange(t, 16, 30),
	}
	resourceName := "data.netbox_available_prefixes.bar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAvailablePrefixesConfigByPrefix(context),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAvailablePrefixesCheck(resourceName, "netbox_available_prefixes.foo"),
					resource.TestCheckResourceAttr(resourceName, "name", "prefix_lookup"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.created", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.description", "testAccDataSourceComputeInstanceConfig description"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.family", "4"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.is_pool", "true"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.last_updated", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T[\d:.]+Z$`)),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.prefix", regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+/\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.role", "gcp"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.site", "se1"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.status", "active"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.vlan", "gcp"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.vrf", "activision"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.tenant", "cloud"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.tags.#", "3"),
				),
			},
		},
	})
}

func TestAccDataSourceAvailablePrefixesByPrefixId(t *testing.T) {

	context := map[string]interface{}{
		"parent_prefix_id":     testNetboxParentPrefixIdWithVrf,
		"random_prefix_length": randIntRange(t, 16, 30),
	}
	resourceName := "data.netbox_available_prefixes.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAvailablePrefixesConfigByPrefixId(context),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAvailablePrefixesCheck(resourceName, "netbox_available_prefixes.foo"),
					resource.TestCheckResourceAttr(resourceName, "name", "prefix_lookup"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.created", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.description", "testAccDataSourceComputeInstanceConfig description"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.family", "4"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.is_pool", "true"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.last_updated", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T[\d:.]+Z$`)),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.prefix", regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+/\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.role", "gcp"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.site", "se1"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.status", "active"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.vlan", "gcp"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.vrf", "activision"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.tenant", "cloud"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.tags.#", "3"),
				),
			},
		},
	})
}

/*
// Those two tests require terraform 0.13.0 to properly work, Skip them here
func TestAccDataSourceAvailablePrefixesByTag(t *testing.T) {

	context := map[string]interface{}{
		"parent_prefix_id":     testNetboxParentPrefixIdWithVrf,
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
	}
	resourceName := "data.netbox_available_prefixes.tag"
	//resourceRoleName := "data.netbox_available_prefixes.bar"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAvailablePrefixesConfigByParameters(context),
				Check: resource.ComposeTestCheckFunc(
					//	testAccDataSourceAvailablePrefixesCheck(resourceName, "netbox_available_prefixes.foo"),
					resource.TestCheckResourceAttr(resourceName, "name", "prefix_lookup_by_tag"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.#", "2"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.created", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.description", regexp.MustCompile(`^testAccDataSourceAvailablePrefixesConfigByParameters`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.family", "4"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.is_pool", "true"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.last_updated", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T[\d:.]+Z$`)),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.prefix", regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+/\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.role", "gcp"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.site", "se1"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.status", "active"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.vlan", "gcp"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.vrf", "activision"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.tenant", "cloud"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.tags.#", "3"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.1.created", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)),
					resource.TestMatchResourceAttr(resourceName, "prefixes.1.description", regexp.MustCompile(`^testAccDataSourceAvailablePrefixesConfigByParameters`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.1.family", "4"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.1.id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.1.is_pool", "true"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.1.last_updated", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T[\d:.]+Z$`)),
					resource.TestMatchResourceAttr(resourceName, "prefixes.1.prefix", regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+/\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.1.role", "gcp"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.1.site", "se1"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.1.status", "active"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.1.vlan", "gcp"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.1.vrf", "activision"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.1.tenant", "cloud"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.1.tags.#", "3"),
				),
			},
		},
	})
}

func TestAccDataSourceAvailablePrefixesByRole(t *testing.T) {

	context := map[string]interface{}{
		"parent_prefix_id":     testNetboxParentPrefixIdWithVrf,
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
	}

	resourceName := "data.netbox_available_prefixes.role"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAvailablePrefixesConfigByParameters(context),
				Check: resource.ComposeTestCheckFunc(
					//testAccDataSourceAvailablePrefixesCheck(resourceName, "netbox_available_prefixes.neo"),
					resource.TestCheckResourceAttr(resourceName, "name", "prefix_lookup_by_role"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.#", "1"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.created", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.description", regexp.MustCompile(`^testAccDataSourceAvailablePrefixesConfigByParameters`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.family", "4"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.id", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.is_pool", "true"),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.last_updated", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T[\d:.]+Z$`)),
					resource.TestMatchResourceAttr(resourceName, "prefixes.0.prefix", regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+/\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.role", "cloudera"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.site", "se1"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.status", "active"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.vlan", "gcp"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.vrf", "activision"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.tenant", "cloud"),
					resource.TestCheckResourceAttr(resourceName, "prefixes.0.tags.#", "3"),
				),
			},
		},
	})
}
*/

func testAccDataSourceAvailablePrefixesCheck(datasourceName string, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[datasourceName]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", datasourceName)
		} else {
			log.Println(ds)
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
			if datasourceAttributes[fmt.Sprintf("prefixes.0.%s", attrToCheck)] != resourceAttributes[attrToCheck] {
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
  name = "prefix_lookup"
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
  name = "prefix_lookup"
  id = netbox_available_prefixes.foo.id
}
`, config)
}

func testAccDataSourceAvailablePrefixesConfigByParameters(config map[string]interface{}) string {
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

  description = "testAccDataSourceAvailablePrefixesConfigByParameters ==> foo"
  tags        = ["datasource-%{random_suffix}-accTag01", "datasource-AvailablePrefix-accTag02", "datasource-AvailablePrefix-accTag03"]
  custom_fields  {}
}

resource "netbox_available_prefixes" "bar" {
  parent_prefix_id 	= %{parent_prefix_id} 
  prefix_length 	= %{random_prefix_length} - 1
  is_pool          	= true
  status          	= "active"
  role = "gcp"
  site = "se1"
  vlan = "gcp"
  vrf  = "activision"
  tenant = "cloud"

  description = "testAccDataSourceAvailablePrefixesConfigByParameters ==> bar"
  tags        = ["datasource-%{random_suffix}-accTag01", "datasource-AvailablePrefix-accTag04", "datasource-AvailablePrefix-accTag05"]
  custom_fields  {}
}

resource "netbox_available_prefixes" "neo" {
  parent_prefix_id 	= %{parent_prefix_id}
  prefix_length 	= %{random_prefix_length} + 1
  is_pool          	= true
  status          	= "active"
  role = "cloudera"
  site = "se1"
  vlan = "gcp"
  vrf  = "activision"
  tenant = "cloud"

  description = "testAccDataSourceAvailablePrefixesConfigByParameters ==> neo"
  tags        = ["datasource-%{random_suffix}-accTag06", "datasource-AvailablePrefix-accTag07", "datasource-AvailablePrefix-accTag08"]
  custom_fields  {}
}

data "netbox_available_prefixes" "tag"{
  name = "prefix_lookup_by_tag"
  tag = lower("datasource-%{random_suffix}-accTag01")
  depends_on  = [netbox_available_prefixes.bar, netbox_available_prefixes.foo, netbox_available_prefixes.neo]
}

data "netbox_available_prefixes" "role"{
  name = "prefix_lookup_by_role"
  role =  lower("cloudera")
  depends_on  = [netbox_available_prefixes.bar, netbox_available_prefixes.foo, netbox_available_prefixes.neo]
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
