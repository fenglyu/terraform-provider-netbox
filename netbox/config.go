package netbox

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
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

	NetboxApiGeneralQueryLimit int64 = 0
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

	if c.Host == "" {
		c.Host = NetboxDefaultHost
	}

	InsecureSkipVerify := true
	if v := os.Getenv("HTTPS_INSECURE_SKIP_VERIFY"); v != "" {
		switch v {
		case "true":
			InsecureSkipVerify = true
		case "false":
			InsecureSkipVerify = false
		default:
			break
		}
	}

	host, schemes := getHost(c.Host)

	if err := ApiAccessTest(host, c.BasePath, c.ApiToken, schemes, InsecureSkipVerify); err != nil {
		return err
	}
	httpClient, err := runtimeclient.TLSClient(runtimeclient.TLSClientOptions{InsecureSkipVerify: InsecureSkipVerify})
	if err != nil {
		log.Fatal(err)
	}
	t := runtimeclient.NewWithClient(host, c.BasePath, schemes, httpClient)

	log.Printf("[INFO] Instantiating http client for host %s and path %s", host, c.BasePath)
	if c.ApiToken != "" {
		t.DefaultAuthentication = runtimeclient.APIKeyAuth(AuthHeaderName, "header", fmt.Sprintf(AuthHeaderFormat, c.ApiToken))
	}
	//t.SetDebug(true)
	c.client = client.New(t, strfmt.Default)

	return nil
}

func getHost(Host string) (string, []string) {
	schemes := []string{} //, "http" "https"
	host := ""
	u, err := url.Parse(Host)
	if err != nil {
		log.Fatal(err)
	}

	reg := regexp.MustCompile(`^(?:https?://)(?:[\w]+)(?:\.?[\w]{2,})+:(?:[0-8]{2,5})$`)
	if !reg.MatchString(Host) {
		host = Host
		if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
			schemes = append(schemes, "http")
		} else {
			schemes = append(schemes, "https")
		}
	} else {
		switch {
		case strings.EqualFold(u.Scheme, "https"):
			schemes = append(schemes, "https")
			host = fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
			break
		case strings.EqualFold(u.Scheme, "http"):
			schemes = append(schemes, "http")
			host = fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
			break
		default:
			// use https by default
			// when schema is empty, then host is c.Host
			schemes = append(schemes, "https")
			break
		}
	}
	return host, schemes
}

func selectScheme(schemes []string) string {
	schLen := len(schemes)
	if schLen == 0 {
		return ""
	}

	scheme := schemes[0]
	// prefer https, but skip when not possible
	if scheme != "https" && schLen > 1 {
		for _, sch := range schemes {
			if sch == "https" {
				scheme = sch
				break
			}
		}
	}
	return scheme
}

func ApiAccessTest(host, path, token string, schemes []string, InsecureSkipVerify bool) error {
	//Test url example: "http://netbox.k8s.me/api/"
	schema := selectScheme(schemes)
	url := fmt.Sprintf("%s://%s%s", schema, host, path)
	method := "GET"
	client := &http.Client{}

	if schema != "" && strings.EqualFold(schema, "https") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: InsecureSkipVerify},
		}
		client = &http.Client{Transport: tr}
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(res.Status)
	}
	return nil
}
