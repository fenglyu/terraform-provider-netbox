package netbox

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/fenglyu/go-netbox/netbox/client/ipam"
	"github.com/fenglyu/go-netbox/netbox/models"
)

func TestAccAvaliablePrefixes_basic(t *testing.T) {
	context := map[string]interface{}{
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAvailablePrefixWithParentPrefixIdExample(context),
			},
			{
				ResourceName:      "netbox_available_prefixes.gke-test",
				ImportState:       true,
				ImportStateVerify: true,
				// parent_prefix_id is the parent which we can't easily get without hacking
				// vrf is not returned via GET "http://netbox.k8s.me/api/ipam/prefixes/{ID}/"
				ImportStateVerifyIgnore: []string{"parent_prefix_id", "vrf"},
			},
		},
	})
}

func TestAccAvaliablePrefixes_basic1(t *testing.T) {
	context := map[string]interface{}{
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAvailablePrefixWithParentPrefixIdExample1(context),
			},
			{
				ResourceName:      "netbox_available_prefixes.foo",
				ImportState:       true,
				ImportStateVerify: true,
				// parent_prefix_id is the parent which we can't easily get without hacking
				// vrf is not returned via GET "http://netbox.k8s.me/api/ipam/prefixes/{ID}/"
				ImportStateVerifyIgnore: []string{"parent_prefix_id", "vrf"},
			},
		},
	})
}

func TestAccAvaliablePrefixes_basic2(t *testing.T) {
	context := map[string]interface{}{
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAvailablePrefixWithParentPrefixIdExample2(context),
			},
			{
				ResourceName:      "netbox_available_prefixes.bar",
				ImportState:       true,
				ImportStateVerify: true,
				// parent_prefix_id is the parent which we can't easily get without hacking
				// vrf is not returned via GET "http://netbox.k8s.me/api/ipam/prefixes/{ID}/"
				ImportStateVerifyIgnore: []string{"parent_prefix_id", "vrf"},
			},
		},
	})
}

func testAccAvailablePrefixWithParentPrefixIdExample(context map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "gke-test" {
	parent_prefix_id = 302
	prefix_length = %{random_prefix_length}
  	is_pool          = true
  	status           = "active"

	tags = ["AvailablePrefix-acc%{random_suffix}-01", "AvailablePrefix-acc%{random_suffix}-02", "AvailablePrefix-acc%{random_suffix}-03"]
}`, context)
}

func testAccAvailablePrefixWithParentPrefixIdExample1(context map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "foo" {
	parent_prefix_id = 302
	prefix_length = %{random_prefix_length}
  	is_pool          = true
  	status           = "active"

  	role = "gcp"
  	site = "se1"
  	vlan = "gcp"
  	#vrf  = "activision"
  	tenant = "cloud"
	tags = ["AvailablePrefix-acc%{random_suffix}-03", "AvailablePrefix-acc%{random_suffix}-04", "AvailablePrefix-acc%{random_suffix}-05"]
}`, context)
}

func testAccAvailablePrefixWithParentPrefixIdExample2(context map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "bar" {
	parent_prefix_id = 302
	prefix_length = %{random_prefix_length}
  	is_pool          = true
  	status           = "active"

  	role = "gcp"
  	site = "se1"
  	vlan = "gcp"
  	#vrf  = "activision"
  	tenant = "cloud"
	tags = ["AvailablePrefix-acc%{random_suffix}-03", "AvailablePrefix-acc%{random_suffix}-04", "AvailablePrefix-acc%{random_suffix}-05"]

	custom_fields    = {
		helpers      = "cf-acc%{random_suffix}-01"
		ipv4_acl_in  = "cf-acc%{random_suffix}-02"
		ipv4_acl_out = "cf-acc%{random_suffix}-03"
	}

}`, context)
}

func testAccCheckAvailablePrefixesDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "netbox_available_prefixes" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
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

func sendRequestforPrefix(config *Config, rs *terraform.ResourceState) (*models.Prefix, error) {
	id, err := strconv.Atoi(rs.Primary.Attributes["id"])
	if err != nil {
		return nil, err
	}
	params := ipam.IpamPrefixesReadParams{
		ID: int64(id),
	}
	params.WithContext(context.Background())

	ipamPrefixesReadOK, err := config.client.Ipam.IpamPrefixesRead(&params, nil)
	if err != nil || ipamPrefixesReadOK == nil {
		return nil, fmt.Errorf("There is not prefix with ID %d", id)
	}

	return ipamPrefixesReadOK.Payload, nil
}
