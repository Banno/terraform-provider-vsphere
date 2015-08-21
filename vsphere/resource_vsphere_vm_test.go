package vsphere

import (
	"fmt"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/mo"
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
				Check:  testAccCheckVsphereVmApp_basic,
			},
		},
	})
}

func testAccCheckVsphereVmApp_basic(s *terraform.State) error {
	found, err := testAccGetFinder().VirtualMachine(context.TODO(), "testingvm1")
	if found == nil || err != nil {
		return fmt.Errorf("Expected to not find %s: #%v, #%v", "testingvm1", found, err)
	}

	props := []string{"summary"}
	var mvm mo.VirtualMachine
	found.Properties(context.TODO(), found.Reference(), props, &mvm)

	memoryMb := mvm.Summary.Config.MemorySizeMB
	if memoryMb != 8192 {
		return fmt.Errorf("memory: expected %d to equal 8192", memoryMb)
	}

	cpus := mvm.Summary.Config.NumCpu
	if cpus != 4 {
		return fmt.Errorf("cpus: expected %d to equal 4", cpus)
	}

	return nil
}

func testAccCheckVsphereVmDestroy(s *terraform.State) error {
	found, err := testAccGetFinder().VirtualMachine(context.TODO(), "testingvm1")
	if found != nil {
		return fmt.Errorf("Expected to not find %s: #%v, #%v", "testingvm1", found, err)
	}

	return nil
}

func testAccGetFinder() *find.Finder {
	client, err := testAccProvider.Meta().(*Config).Client()
	if err != nil {
		panic(fmt.Sprintf("getting finder, got err: #%v", err))
	}
	finder := find.NewFinder(client.Client, false)
	datacenter, err := finder.DefaultDatacenter(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("getting finder, got err: #%v", err))
	}
	finder.SetDatacenter(datacenter)
	return finder
}
