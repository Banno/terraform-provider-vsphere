package vsphere

/*import (
	"github.com/hashicorp/terraform/helper/resource"
	"log"
)*/

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
  "github.com/vmware/govmomi"
  "github.com/vmware/govmomi/find"
  "github.com/vmware/govmomi/vim25/types"
)

func resourceVsphereVm() *schema.Resource {
	return &schema.Resource{
		Create: resourceVsphereVmCreate,
		Read:   resourceVsphereVmRead,
		Update: resourceVsphereVmUpdate,
		Delete: resourceVsphereVmDelete,

		Schema: map[string]*schema.Schema{
			"template_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
      "vm_name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
		},
	}
}

func resourceVsphereVmCreate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*govmomi.Client) 
  
  finder := find.NewFinder(client, false)

  datacenter, err := finder.DefaultDatacenter()

  if err != nil {
    return err
  }

  finder.SetDatacenter(datacenter)

  resourcePool, err := finder.DefaultResourcePool()

  if err != nil {
    return err
  }

  rpRef := resourcePool.Reference()

  vm, err := finder.VirtualMachine(d.Get("template_name").(string))

  if err != nil {
    return err
  }

  folders, err := datacenter.Folders()

  if err != nil {
    return err
  }

  clonespec := types.VirtualMachineCloneSpec{
    Config: &types.VirtualMachineConfigSpec{},
    Location: types.VirtualMachineRelocateSpec{
      Pool: &rpRef,
    },
  }


  task, err := vm.Clone(folders.VmFolder,d.Get("vm_name").(string),clonespec)


  if err != nil {
    return err
  }

  info, err := task.WaitForResult(nil)
  
  if err != nil {
    return err
  }

  if info.State == "success" {
    fmt.Printf("%s Registered!", d.Get("vm_name").(string))
  }

  fmt.Printf("%s", d.Id())



  return nil
}

func resourceVsphereVmRead(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourceVsphereVmUpdate(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourceVsphereVmDelete(d *schema.ResourceData, meta interface{}) error {
  return nil
}
