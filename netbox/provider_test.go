package netbox

import (
	"fmt"
	"log"
	"math/rand"
	"os"
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

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

func getTestNetboxApiTokenFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, netboxApiTokenEnvVars...)
	return multiEnvSearch(netboxApiTokenEnvVars)
}

func getTestNetboxHostFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, netboxHostEnvVars...)
	return multiEnvSearch(netboxHostEnvVars)
}

func getTestNetboxBasePathFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, netboxBasePathEnvVars...)
	return multiEnvSearch(netboxBasePathEnvVars)
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

/*
func randInt(t *testing.T) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}*/

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
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(netboxHostEnvVars, ", "))
	}
	if v := multiEnvSearch(netboxBasePathEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(netboxBasePathEnvVars, ", "))
	}
}

func TestAccProviderBasePath_setBasePath(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
		"basePath":             "/api",
		"host":                 "netbox.k8s.me",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAvailablePrefixesDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderBasePath_setBasePath(context),
			},
			{
				ResourceName:            "netbox_available_prefixes.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent_prefix_id"},
			},
		},
	})
}

func TestAccProviderBasePath_setInvalidBasePath(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"basePath": fmt.Sprintf("/%s", randString(t, 10)),
		"host":     "www.example.com",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderBasePath_setInvalidBasePath(context),
				//ExpectError: regexp.MustCompile("404 Not Found"),
			},
		},
	})
}

func testAccProviderBasePath_setBasePath(context map[string]interface{}) string {
	return Nprintf(`
provider "netbox" {
  	host      = "%{host}"
	base_path = "%{basePath}"
	request_timeout = "4m"
}

resource "netbox_available_prefixes" "default" {
	parent_prefix_id = 302
	prefix_length = %{random_prefix_length}
	tags = ["BasePathTest-acc%{random_suffix}-01", "BasePathTest-acc%{random_suffix}-02", "BasePathTest-acc%{random_suffix}-03"]
}
`, context)
}

func testAccProviderBasePath_setInvalidBasePath(context map[string]interface{}) string {
	return Nprintf(`
provider "netbox" {
  	host      = "%{host}"
	base_path = "%{basePath}"
}
`, context)
}
