package vsphere

import (
	"fmt"
	"github.com/vmware/govmomi"
	"golang.org/x/net/context"
	"net/url"
)

type Config struct {
	Username string
	Password string
	Host     string
	client   *govmomi.Client
}

func (c *Config) Client() (*govmomi.Client, error) {
	if c.client == nil {
		sdkURL, err := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk",
			c.Username,
			c.Password,
			c.Host))

		if err != nil {
			return nil, err
		}
		client, err := govmomi.NewClient(context.TODO(), sdkURL, true)
		if err != nil {
			return nil, err
		} else {
			c.client = client
			return client, nil
		}
	} else {
		return c.client, nil
	}

}
