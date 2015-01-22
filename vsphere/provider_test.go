package vsphere

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"vsphere": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	username := os.Getenv("VSPHERE_USERNAME")
	password := os.Getenv("VSPHERE_PASSWORD")
	url := os.Getenv("VSPHERE_URL")
	if username == "" || password == "" || url == "" {
		t.Fatal("VSPHERE_USERNAME, VSPHERE_PASSWORD and VSPHERE_URL must be set for acceptance tests to work.")
	}
}
