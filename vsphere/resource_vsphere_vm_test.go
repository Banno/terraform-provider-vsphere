package vsphere

import (
	"fmt"

	"github.com/vmware/govmomi/find"
	"golang.org/x/net/context"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"
)

const testAccCheckVsphereVmConfig_basic = `
resource "vsphere_vm" "testingvm1" {
  vm_name = "testingvm1"

  template_name = "packer-stemcell-vsphere-1432631533-vm"
  memory_mb = "8192"
  cpus = "4"
  customization_specification = "Ubuntu"
}
`

func TestAccVsphereVm_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVsphereVmDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckVsphereVmConfig_basic,
			},
		},
	})
}

func testAccCheckVsphereVmDestroy(s *terraform.State) error {
	client, err := testAccProvider.Meta().(*Config).Client()
	if err != nil {
		return err
	}
	finder := find.NewFinder(client.Client, false)
	datacenter, err := finder.DefaultDatacenter(context.TODO())
	if err != nil {
		return err
	}
	finder.SetDatacenter(datacenter)
	found, err := finder.VirtualMachine(context.TODO(), "testingvm1")
	if found != nil {
		return fmt.Errorf("Expected to not find %s: #%v, #%v", "testingvm1", found, err)
	}

	return nil
}
