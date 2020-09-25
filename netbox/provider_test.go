package netbox

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-random/random"

	"github.com/fenglyu/go-netbox/netbox/client/ipam"
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

// Declare an environment parent prefix id for ACC test
var netboxParentPrefixIdForTestingVars = []string{
	"NETBOX_PARENT_PREFIX_ID",
}

// Parent prefix has a vrf
var netboxParentPrefixWithVrfIdForTestingVars = []string{
	"NETBOX_PARENT_PREFIX_WITH_VRF_ID",
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

func getNetBoxParentPrefixIdForTestingVarsFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, netboxParentPrefixIdForTestingVars...)
	return multiEnvSearch(netboxParentPrefixIdForTestingVars)
}

func getNetboxParentPrefixWithVrfIdForTestingVarsFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, netboxParentPrefixWithVrfIdForTestingVars...)
	return multiEnvSearch(netboxParentPrefixWithVrfIdForTestingVars)
}

func skipIfEnvNotSet(t *testing.T, envs ...string) {
	if t == nil {
		log.Printf("[DEBUG] Not running inside of test - skip skipping")
		return
	}

	for i, k := range envs {
		switch {
		case os.Getenv(k) == "" && i == len(envs)-1:
			t.Skipf("Environment variable %s is not set", strings.Join(envs, ", "))
		case os.Getenv(k) == "":
			log.Printf("[%d] Environment variable %s is not set", i, k)
			break
		case os.Getenv(k) != "":
			return
		default:
			return
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

var (
	testAccProviders      map[string]terraform.ResourceProvider
	testAccProvider       *schema.Provider
	testAccRandomProvider *schema.Provider

	checkPrefixIdOnce               sync.Once
	testNetboxParentPrefixId        int
	testNetboxParentPrefixIdWithVrf int
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccRandomProvider = random.Provider().(*schema.Provider)

	testAccProviders = map[string]terraform.ResourceProvider{
		"netbox": testAccProvider,
		"random": testAccRandomProvider,
	}
	// check the existance of two test parent prefix id
	checkPrefixId()
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func checkPrefixId() {
	checkPrefixIdOnce.Do(initPrefixId)
}

func initPrefixId() {
	// Test if given parent prefix id is available
	if err := testNetboxParentPrefixDefault(); err != nil {
		return
	}
	// 	Test if given parent prefix(with a vrf) id is available
	if err := testNetboxParentPrefixWithVrf(); err != nil {
		return
	}
}

func testPrefixExistWithID(parentPrefixId int) error {

	params := ipam.IpamPrefixesReadParams{
		ID:      int64(parentPrefixId),
		Context: context.Background(),
	}

	config := &Config{
		ApiToken: multiEnvSearch(netboxApiTokenEnvVars),
		Host:     multiEnvSearch(netboxHostEnvVars),
		BasePath: multiEnvSearch(netboxBasePathEnvVars),
	}

	err := config.LoadAndValidate(context.Background())
	if err != nil {
		return err
	}

	ipamPrefixesReadOK, err := config.client.Ipam.IpamPrefixesRead(&params, nil)
	//t.Log(fmt.Sprintf("testPrefixExistWithID %v", ipamPrefixesReadOK))
	if err != nil || ipamPrefixesReadOK == nil {
		return fmt.Errorf("Cannot determine prefix with ID %d", parentPrefixId)
	}
	return nil
}

func testNetboxParentPrefixDefault() error {

	if ppid, err := strconv.Atoi(multiEnvSearch(netboxParentPrefixIdForTestingVars)); err == nil {
		testNetboxParentPrefixId = ppid
	}
	return testPrefixExistWithID(testNetboxParentPrefixId)
}

func testNetboxParentPrefixWithVrf() error {
	if ppid, err := strconv.Atoi(multiEnvSearch(netboxParentPrefixWithVrfIdForTestingVars)); err == nil {
		testNetboxParentPrefixIdWithVrf = ppid
	}
	return testPrefixExistWithID(testNetboxParentPrefixIdWithVrf)
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
	if v := multiEnvSearch(netboxParentPrefixIdForTestingVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(netboxParentPrefixIdForTestingVars, ", "))
	}
	if v := multiEnvSearch(netboxParentPrefixWithVrfIdForTestingVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(netboxParentPrefixWithVrfIdForTestingVars, ", "))
	}
}

func TestAccProviderBasePath_setBasePath(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_prefix_length": randIntRange(t, 16, 30),
		"random_suffix":        randString(t, 10),
		"basePath":             "/api",
		"host":                 "localhost:8080",
		"parent_prefix_id":     testNetboxParentPrefixId,
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
	parent_prefix_id = %{parent_prefix_id}
	prefix_length = %{random_prefix_length}
	tags = ["BasePathTest-acc%{random_suffix}-01", "BasePathTest-acc%{random_suffix}-02", "BasePathTest-acc%{random_suffix}-03"]

  	custom_fields  {}
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
