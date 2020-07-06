package netbox

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAvaliablePrefixes_basic(t *testing.T) {
	/* ... potentially existing acceptance testing logic ... */

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckPrefixDestroyProducer(t),
		Steps: []resource.TestStep{
			/* ... existing TestStep ... */
			{
				Config: nil,
			},
			{
				ResourceName:      "example_thing.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

/*
func testAccCheckPrefixDestroyProducer(t *testing.T) func(s *terraform.State) error {
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

*/
