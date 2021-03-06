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
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_USERNAME", ""),
			},
			"vsphere_password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_PASSWORD", ""),
			},
			"vsphere_host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_HOST", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vsphere_vm": resourceVsphereVM(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	c := Config{
		Username: d.Get("vsphere_username").(string),
		Password: d.Get("vsphere_password").(string),
		Host:     d.Get("vsphere_host").(string),
	}
	return &c, nil
}
