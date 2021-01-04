package netbox

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/fenglyu/go-netbox/netbox/client/ipam"
	"github.com/fenglyu/go-netbox/netbox/models"
)

func TestAccAvailablePrefixes_basic(t *testing.T) {
	context := map[string]interface{}{
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
		"parent_prefix_id":     testNetboxParentPrefixId,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		//Providers:         testAccProviders,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAvailablePrefixesDestroyProducer(t),
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
				ImportStateVerifyIgnore: []string{"parent_prefix_id"},
			},
		},
	})
}

func TestAccAvailablePrefixes_basic1(t *testing.T) {
	context := map[string]interface{}{
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
		"parent_prefix_id":     testNetboxParentPrefixIdWithVrf,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAvailablePrefixWithParentPrefixIdExample1(context),
			},
			{
				ResourceName:            "netbox_available_prefixes.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent_prefix_id"},
			},
		},
	})
}

func TestAccAvailablePrefixes_basic2(t *testing.T) {
	context := map[string]interface{}{
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
		"parent_prefix_id":     testNetboxParentPrefixIdWithVrf,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAvailablePrefixWithParentPrefixIdExample2(context),
			},
			{
				ResourceName:            "netbox_available_prefixes.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent_prefix_id"},
			},
		},
	})
}

func TestAccAvailablePrefixes_EmptyCustomFields(t *testing.T) {
	context := map[string]interface{}{
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
		"parent_prefix_id":     testNetboxParentPrefixIdWithVrf,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAvailablePrefixWithParentPrefixEmptyCF(context),
			},
			{
				ResourceName:            "netbox_available_prefixes.custom_fields",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent_prefix_id"},
			},
		},
	})
}

func TestAccAvailablePrefixesMultipleSteps(t *testing.T) {
	context := map[string]interface{}{
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
		"parent_prefix_id":     testNetboxParentPrefixIdWithVrf,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAvailablePrefixWithParentPrefixIdMultipleStep1(context),
			},
			{
				Config: testAccAvailablePrefixWithParentPrefixIdMultipleStep2(context),
			},
			{
				Config: testAccAvailablePrefixWithParentPrefixIdMultipleStep3(context),
			},
			{
				ResourceName:            "netbox_available_prefixes.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent_prefix_id"},
			},
		},
	})
}

func testAccAvailablePrefixWithParentPrefixIdMultipleStep1(context map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "bar" {
	parent_prefix_id = %{parent_prefix_id}
	prefix_length = %{random_prefix_length}
  	is_pool          = true
  	status           = "active"

  	role = "gcp"
  	site = "se1"
  	vlan = "gcp"
  	vrf  = "activision"
  	tenant = "cloud"
	tags = ["AvailablePrefix-acc%{random_suffix}-01", "AvailablePrefix-acc%{random_suffix}-02", "AvailablePrefix-acc%{random_suffix}-03"]

	custom_fields {
		helpers      = "cf-acc%{random_suffix}-01"
		ipv4_acl_in  = "cf-acc%{random_suffix}-02"
		ipv4_acl_out = "cf-acc%{random_suffix}-03"
	}
}`, context)
}

func testAccAvailablePrefixWithParentPrefixIdMultipleStep2(context map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "bar" {
	parent_prefix_id = %{parent_prefix_id}
	prefix_length = %{random_prefix_length}
  	is_pool          = true
  	status           = "reserved"

  	role = "gcp"
  	site = "se1"
  	vlan = "gcp"
  	vrf  = "activision"
  	tenant = "cloud"
	tags = ["AvailablePrefix-acc%{random_suffix}-03", "AvailablePrefix-acc%{random_suffix}-04", "AvailablePrefix-acc%{random_suffix}-05"]

	custom_fields {
		helpers      = "cf-acc%{random_suffix}-01"
		ipv4_acl_in  = "cf-acc%{random_suffix}-02"
		ipv4_acl_out = "cf-acc%{random_suffix}-03"
	}
}`, context)
}

func testAccAvailablePrefixWithParentPrefixIdMultipleStep3(context map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "bar" {
	parent_prefix_id = %{parent_prefix_id}
	prefix_length = %{random_prefix_length}
  	is_pool          = false
  	status           = "active"

  	role = "gcp"
  	site = "se1"
  	vlan = "gcp"
  	vrf  = "activision"
  	tenant = "cloud"
	tags = ["AvailablePrefix-acc%{random_suffix}-06", "AvailablePrefix-acc%{random_suffix}-07", "AvailablePrefix-acc%{random_suffix}-08"]

	custom_fields {
		helpers      = "cf-acc%{random_suffix}-01"
		ipv4_acl_in  = "cf-acc%{random_suffix}-02"
		ipv4_acl_out = "cf-acc%{random_suffix}-03"
	}
}`, context)
}

func testAccAvailablePrefixWithParentPrefixIdExample1(context map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "foo" {
	parent_prefix_id = %{parent_prefix_id}
	prefix_length = %{random_prefix_length}
  	is_pool          = true
  	status           = "active"

  	role = "gcp"
  	site = "se1"
  	vlan = "gcp"
  	vrf  = "activision"
  	tenant = "cloud"
	tags = ["AvailablePrefix-acc%{random_suffix}-03", "AvailablePrefix-acc%{random_suffix}-04", "AvailablePrefix-acc%{random_suffix}-05"]

  	custom_fields  {}
}`, context)
}

func testAccAvailablePrefixWithParentPrefixIdExample2(context map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "bar" {
	parent_prefix_id = %{parent_prefix_id}
	prefix_length = %{random_prefix_length}
  	is_pool          = true
  	status           = "active"

  	role = "gcp"
  	site = "se1"
  	vlan = "gcp"
  	vrf  = "activision"
  	tenant = "cloud"
	tags = ["AvailablePrefix-acc%{random_suffix}-03", "AvailablePrefix-acc%{random_suffix}-04", "AvailablePrefix-acc%{random_suffix}-05"]

	custom_fields {
		helpers      = "cf-acc%{random_suffix}-01"
		ipv4_acl_in  = "cf-acc%{random_suffix}-02"
		ipv4_acl_out = "cf-acc%{random_suffix}-03"
	}
}`, context)
}

func testAccAvailablePrefixWithParentPrefixEmptyCF(context map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "custom_fields" {
	parent_prefix_id = %{parent_prefix_id}
	prefix_length = %{random_prefix_length}
  	is_pool          = true
  	status           = "active"

  	role = "gcp"
  	site = "se1"
  	vlan = "gcp"
  	vrf  = "activision"
  	tenant = "cloud"
	tags = ["AvailablePrefix-acc%{random_suffix}-06"]

	custom_fields {
		helpers      = ""
		ipv4_acl_in  = ""
		ipv4_acl_out = ""
	}
}`, context)
}

func testAccAvailablePrefixWithParentPrefixIdExample(context map[string]interface{}) string {
	return Nprintf(`
resource "netbox_available_prefixes" "gke-test" {
	parent_prefix_id = %{parent_prefix_id}
	prefix_length = %{random_prefix_length}
  	is_pool          = true
  	status           = "active"

	tags = ["AvailablePrefix-acc%{random_suffix}-01", "AvailablePrefix-acc%{random_suffix}-02", "AvailablePrefix-acc%{random_suffix}-03"]
	custom_fields {}
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

func sendRequestforPrefix(config *Config, rs *terraform.ResourceState) ([]*models.Prefix, error) {
	idStr := rs.Primary.Attributes["id"]

	idList := make([]string, 0)
	// datasource
	re := regexp.MustCompile(`[a-zA-Z-_]+`)
	if re.MatchString(idStr) {
		prefixesLen, err := strconv.Atoi(rs.Primary.Attributes["prefixes.#"])
		if err != nil {
			return nil, err
		}
		i := 0
		for i < prefixesLen {
			id := rs.Primary.Attributes[fmt.Sprintf("prefixes.%d.id", i)]
			idList = append(idList, id)
			i++
		}
	} else {
		idList = append(idList, idStr)
	}

	prefixList := make([]*models.Prefix, 0)
	var mulError *multierror.Error

	for _, ids := range idList {
		id, err := strconv.Atoi(ids)
		if err != nil {
			log.Println("err ==> ", ids)
			return nil, err
		}
		params := ipam.IpamPrefixesReadParams{
			ID: int64(id),
		}
		params.WithContext(context.Background())

		ipamPrefixesReadOK, err := config.client.Ipam.IpamPrefixesRead(&params, nil)
		if err != nil {
			mulError = multierror.Append(mulError, fmt.Errorf("Lookup Prefix ID %s", idStr))
		}
		if ipamPrefixesReadOK != nil && ipamPrefixesReadOK.Payload != nil {
			prefixList = append(prefixList, ipamPrefixesReadOK.Payload)
		}
	}

	return prefixList, mulError
}
