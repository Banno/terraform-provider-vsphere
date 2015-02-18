package vsphere

import (
	"fmt"
	"github.com/vmware/govmomi"
	"net/url"
)

type Config struct {
	Username string
	Password string
	Host     string
}

func (c *Config) Client() (*govmomi.Client, error) {
	sdkURL, err := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk",
		c.Username,
		c.Password,
		c.Host))

	if err != nil {
		return nil, err
	}

	client, err := govmomi.NewClient(*sdkURL, true)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}

}
