package vsphere

import (
	"os"
	"fmt"
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
			"vsphere_host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_HOST", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"vsphere_vm": resourceVsphereVM(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func checkEnvVariable(name string) {
	if os.Getenv(name) == "" {
		fmt.Println("Missing %s environment variable", name)
		os.Exit(3)
	}
	return
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	checkEnvVariable("VSPHERE_USERNAME")
	checkEnvVariable("VSPHERE_PASSWORD")
	checkEnvVariable("VSPHERE_HOST")

	config := Config{
		Username: d.Get("vsphere_username").(string),
		Password: d.Get("vsphere_password").(string),
		Host:     d.Get("vsphere_host").(string),
	}
	return config.Client()
}
