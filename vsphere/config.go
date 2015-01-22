package vsphere

import (
	"fmt"
	"github.com/vmware/govmomi"
	"net/url"
)

type Config struct {
	Username string
	Password string
	URL      string

	vsphereClient *govmomi.Client
}

func (c *Config) Client() (*govmomi.Client, error) {
	client := newClient(c)

	return client, nil
}

func newClient(c *Config) *govmomi.Client {
	sdk_url, _ := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk",
		c.Username,
		c.Password,
		c.URL))

	client, _ := govmomi.NewClient(*sdk_url, true)

	return client

}
