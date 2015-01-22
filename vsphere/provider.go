package vsphere

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"vsphere_username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_USERNAME", nil),
			},
			"vsphere_password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_PASSWORD", nil),
			},
			"vsphere_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_URL", nil),
			},
		},
		ResourcesMap:  map[string]*schema.Resource{
      "vsphere_vm": resourceVsphereVm(),
    },
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username: d.Get("vsphere_username").(string),
		Password: d.Get("vsphere_password").(string),
		URL:      d.Get("vsphere_url").(string),
	}
  return config.Client()
}
