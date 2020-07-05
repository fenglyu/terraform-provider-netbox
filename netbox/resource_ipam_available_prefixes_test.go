package netbox

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAvaliablePrefixes_basic(t *testing.T) {
	/* ... potentially existing acceptance testing logic ... */

	resource.ParallelTest(t, resource.TestCase{
		/* ... existing TestCase functions ... */
		Steps: []resource.TestStep{
			/* ... existing TestStep ... */
			{
				ResourceName:      "example_thing.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
