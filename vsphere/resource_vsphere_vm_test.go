package vsphere

import (
	// "fmt"
	// "time"

	"github.com/hashicorp/terraform/helper/resource"
	// "github.com/hashicorp/terraform/terraform"

	"testing"
)

const testAccCheckVsphereVmConfig_basic = `
resource "vsphere_vm" "testingvm1" {
  template_name = "packer-stemcell-vsphere-1432631533-vm"
  vm_name = "testingvm1"
  memory_mb = "8192"
  cpus = "4"
  customization_specification = "Ubuntu"
}
`

func TestAccVsphereVm_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckVsphereVmConfig_basic,
			},
		},
	})
}
