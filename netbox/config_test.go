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

	err := ApiAccessTest(config.Host, config.BasePath, config.ApiToken, []string{"http"}, true)
	if err != nil {
		t.Fatalf("error %v", err)
	}
}

func TestApiAccessTestInHttps(t *testing.T) {
	config := &Config{
		ApiToken: getTestNetboxApiTokenFromEnv(t),
		Host:     getTestNetboxHostFromEnv(t),
		BasePath: getTestNetboxBasePathFromEnv(t),
	}

	err := ApiAccessTest(config.Host, config.BasePath, config.ApiToken, []string{"https"}, true)
	if err != nil {
		t.Fatalf("error %v", err)
	}
}
