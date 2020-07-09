package netbox

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAvaliablePrefixes_basic(t *testing.T) {
	/* ... potentially existing acceptance testing logic ... */
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
			},
		},
	})
}

func testAccAvailablePrefixWithParentPrefixIdExample(context map[string]interface{}) string {
	return Nprintf(`
	resource "netbox_available_prefixes" "gke-test" {
		parent_prefix_id = 1
		prefix_length = %{random_prefix_length}
		tags = ["test-acc%{random_suffix}-01", "test-acc%{random_suffix}-02", "test-acc%{random_suffix}-03"]
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

			_, err = sendRequest(config, "GET", "", url, nil)
			if err == nil {
				return fmt.Errorf("ComputeAddress still exists at %s", url)
			}
		}

		return nil
	}
}
