package netbox

import (
	"context"
	"testing"
)

func TestLoadAndValidate(t *testing.T) {

	config := &Config{
		ApiToken: getTestNetboxApiTokenFromEnv(t),
		Host:     getTestNetboxHostFromEnv(t),
		BasePath: getTestNetboxBasePathFromEnv(t),
	}

	err := config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestApiAccessTestInHttp(t *testing.T) {
	config := &Config{
		ApiToken: getTestNetboxApiTokenFromEnv(t),
		Host:     getTestNetboxHostFromEnv(t),
		BasePath: getTestNetboxBasePathFromEnv(t),
	}
	host, schemes := getHost(config.Host)
	if schemes[0] != "http" {
		t.Skipf("scheme http supported only %s", schemes[0])
	}
	err := ApiAccessTest(host, config.BasePath, config.ApiToken, []string{"http"}, true)
	if err != nil {
		t.Skipf("error %v", err)
	}
}

func TestApiAccessTestInHttps(t *testing.T) {
	config := &Config{
		ApiToken: getTestNetboxApiTokenFromEnv(t),
		Host:     getTestNetboxHostFromEnv(t),
		BasePath: getTestNetboxBasePathFromEnv(t),
	}

	host, schemes := getHost(config.Host)
	if schemes[0] != "https" {
		t.Skipf("scheme https supported only %s", schemes[0])
	}
	err := ApiAccessTest(host, config.BasePath, config.ApiToken, []string{"https"}, true)
	if err != nil {
		t.Fatalf("error %v", err)
	}
}
