package netbox

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
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
	// Below is a work-around for netbox v2.4.7
	prefixStatusContainer  int64 = 0
	prefixStatusActive     int64 = 1
	prefixStatusReserved   int64 = 2
	prefixStatusDeprecated int64 = 3

	prefixStatusIDMap = map[string]int64{
		"container":  prefixStatusContainer,
		"active":     prefixStatusActive,
		"reserved":   prefixStatusReserved,
		"deprecated": prefixStatusDeprecated,
	}
	prefixStatusIDMapReverse = map[int64]string{
		prefixStatusContainer:  "container",
		prefixStatusActive:     "active",
		prefixStatusReserved:   "reserved",
		prefixStatusDeprecated: "deprecated",
	}
)

// Config support all configurations for provider
type Config struct {
	ApiToken       string
	Host           string
	BasePath       string
	RequestTimeout time.Duration

	// new box client
	client *client.NetBox
	//context context.Context
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

	httpClient, err := runtimeclient.TLSClient(runtimeclient.TLSClientOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Fatal(err)
	}
	t := runtimeclient.NewWithClient(c.Host, c.BasePath, []string{"https", "http"}, httpClient)
	log.Printf("[INFO] Instantiating http client for host %s and path %s", c.Host, c.BasePath)
	if c.ApiToken != "" {
		t.DefaultAuthentication = runtimeclient.APIKeyAuth(AuthHeaderName, "header", fmt.Sprintf(AuthHeaderFormat, c.ApiToken))
	}
	t.SetDebug(false)
	c.client = client.New(t, strfmt.Default)

	if err := ApiAccessTest(c.Host, c.BasePath, c.ApiToken); err != nil {
		return err
	}

	return nil
}

func ApiAccessTest(host, path, token string) error {
	//Test url example: "http://netbox.k8s.me/api/"
	url := fmt.Sprintf("http://%s%s", host, path)
	method := "GET"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))

	res, err := client.Do(req)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf(res.Status)
	}
	return nil
}
