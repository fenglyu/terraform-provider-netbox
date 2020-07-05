package netbox

import (
	"log"
	"os"
	"testing"
)

var netboxApiTokenEnvVars = []string{
	"NETBOX_API_TOKEN",
	"API_TOKEN",
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
