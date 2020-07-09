package netbox

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-random/random"
)

var netboxApiTokenEnvVars = []string{
	"NETBOX_TOKEN",
	"NETBOX_API_TOKEN",
	"API_TOKEN",
}

var netboxHostEnvVars = []string{
	"NETBOX_HOST",
}

var netboxBasePathEnvVars = []string{
	"NETBOX_BASE_PATH",
}

func getTestNetboxApiTokenFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, netboxApiTokenEnvVars...)
	return multiEnvSearch(netboxApiTokenEnvVars)
}

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func skipIfEnvNotSet(t *testing.T, envs ...string) {
	if t == nil {
		log.Printf("[DEBUG] Not running inside of test - skip skipping")
		return
	}

	for _, k := range envs {
		if os.Getenv(k) == "" {
			t.Skipf("Environment variable %s is not set", k)
		}
	}
}

func randInt(t *testing.T) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}

func randIntRange(t *testing.T, min, max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max-min) + min
}

func randString(t *testing.T, length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	set := "abcdefghijklmnopqrstuvwxyz012346789"
	for i := 0; i < length; i++ {
		result[i] = set[r.Intn(len(set))]
	}
	return string(result)
}

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccRandomProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccRandomProvider = random.Provider().(*schema.Provider)

	testAccProviders = map[string]terraform.ResourceProvider{
		"netbox": testAccProvider,
		"random": testAccRandomProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := multiEnvSearch(netboxApiTokenEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(netboxApiTokenEnvVars, ", "))
	}
	if v := multiEnvSearch(netboxHostEnvVars); v == "" {
		t.Fatalf("One of %s must be set to us-central1-a for acceptance tests", strings.Join(netboxHostEnvVars, ", "))
	}
	if v := multiEnvSearch(netboxBasePathEnvVars); v == "" {
		t.Fatalf("One of %s must be set to us-central1-a for acceptance tests", strings.Join(netboxBasePathEnvVars, ", "))
	}
}

func TestAccProviderBasePath_setBasePath(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrefixDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderBasePath_setBasePath("c4a3c627b64fa514e8e0840a94c06b04eb8674d9", "netbox.k8s.me", "/api", "gke-test", randIntRange(t, 18, 30)),
			},
			{
				ResourceName:      "netbox_available_prefixes.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProviderBasePath_setInvalidBasePath(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrefixDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderBasePath_setBasePath("https://www.example.com/compute/beta/", randIntRange(t, 18, 30)),
				ExpectError: regexp.MustCompile("got HTTP response code 404 with body"),
			},
		},
	})
}

func testAccProviderBasePath_setBasePath(apiToken, host, basePath, name string, prefixLength int) string {
	return fmt.Sprintf(`
	provider "netbox" {
		api_token = "%s"
		host      = "%s"
		base_path = "%s"
	}
`, apiToken, basePath, name)
}
