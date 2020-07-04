package netbox

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fenglyu/go-netbox/netbox/client"
	runtimeclient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

var (
	NetboxDefaultHost     = "localhost:8000"
	NetboxDefaultBasePath = "/api"
	AuthHeaderName        = "Authorization"
	AuthHeaderFormat      = "Token %v"

	prefixinitializeStatus = []string{
		"container", "active", "reserved", "deprecated",
	}
)

// Config support all configurations for provider
type Config struct {
	ApiToken       string
	Host           string
	BasePath       string
	RequestTimeout time.Duration

	// new box client
	client  *client.NetBox
	context context.Context
}

func (c *Config) LoadAndValidate(ctx context.Context) error {

	if c.BasePath == "" {
		c.BasePath = NetboxDefaultBasePath
	}

	if c.BasePath == "" {
		c.BasePath = NetboxDefaultBasePath
	}
	if c.Host == "" {
		c.Host = NetboxDefaultHost
	}

	t := runtimeclient.New(c.Host, c.BasePath, client.DefaultSchemes)
	log.Printf("[INFO] Instantiating http client for host %s and path %s", c.Host, c.BasePath)
	if c.ApiToken != "" {
		t.DefaultAuthentication = runtimeclient.APIKeyAuth(AuthHeaderName, "header", fmt.Sprintf(AuthHeaderFormat, c.ApiToken))
	}
	//t.SetDebug(true)
	c.client = client.New(t, strfmt.Default)

	return nil
}
