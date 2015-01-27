package vsphere

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
  "github.com/vmware/govmomi/vim25/mo"
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
      "ip_address": &schema.Schema{
        Type:   schema.TypeString,
        Computed: true,
      },
      "cpus": &schema.Schema{
        Type: schema.TypeInt,
        Required: true,
      },
      "memory_mb": &schema.Schema{
        Type: schema.TypeInt,
        Required: true,
      },
      "static_ip": &schema.Schema{
        Type: schema.TypeString,
        Optional: true,
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
		Config: &types.VirtualMachineConfigSpec{
      NumCPUs: d.Get("cpus").(int),
      MemoryMB: int64(d.Get("memory_mb").(int)),
      CpuHotAddEnabled: true,
      CpuHotRemoveEnabled: true,
      MemoryHotAddEnabled: true,
    },
		Location: types.VirtualMachineRelocateSpec{
			Pool: &rpRef,
		},
    PowerOn: true,
	}

  ipAddress := d.Get("static_ip").(string)
  
  if ipAddress != "" {
    ip := types.CustomizationFixedIp{
      IpAddress: ipAddress,
    }
    specManager := client.CustomizationSpecManager()
    specItem, _ := specManager.GetCustomizationSpec("Ubuntu 1")
    specItem.Spec.NicSettingMap[0].Adapter.Ip = &ip
    clonespec.Customization = &specItem.Spec
  }

	task, err := vm.Clone(folders.VmFolder, d.Get("vm_name").(string), clonespec)

	if err != nil {
		return err
	}

	_, err = task.WaitForResult(nil)

	if err != nil {
		return err
	}

  d.SetId(d.Get("vm_name").(string))

	return resourceVsphereVmRead(d, meta)
}

func resourceVsphereVmRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*govmomi.Client)

  finder := find.NewFinder(client, false)

  datacenter, err := finder.DefaultDatacenter()

  if err != nil {
    return err
  }

  finder.SetDatacenter(datacenter)

  if err != nil {
    return err
  }

  vm, err := finder.VirtualMachine(d.Get("vm_name").(string))

  if err != nil {
    if err.Error() == fmt.Sprintf("vm '%s' not found", d.Get("vm_name").(string)) {
      d.SetId("")
      return nil
    }
  }

  props := []string{"summary"}
  

  var mvm mo.VirtualMachine

  err = client.Properties(vm.Reference(), props, &mvm)

  if err != nil {
    return err
  }

 
  d.Set("memory_mb", mvm.Summary.Config.MemorySizeMB)
  d.Set("cpus", mvm.Summary.Config.NumCpu)

	return nil
}

func resourceVsphereVmUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*govmomi.Client)

  finder := find.NewFinder(client, false)

  datacenter, err := finder.DefaultDatacenter()

  if err != nil {
    return err
  }

  finder.SetDatacenter(datacenter)

  vm, err := finder.VirtualMachine(d.Get("vm_name").(string))

  if err != nil {
    return err
  }

  configspec := types.VirtualMachineConfigSpec{
      NumCPUs: d.Get("cpus").(int),
      MemoryMB: int64(d.Get("memory_mb").(int)),
  }

  task, err := vm.Reconfigure(configspec)

  if err != nil {
    return err
  }

  _, err = task.WaitForResult(nil)


  if err != nil {
    return err
  }

	return resourceVsphereVmRead(d, meta)
}

func resourceVsphereVmDelete(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*govmomi.Client)

  finder := find.NewFinder(client, false)

  datacenter, err := finder.DefaultDatacenter()

  if err != nil {
    return err
  }

  finder.SetDatacenter(datacenter)

  if err != nil {
    return err
  }

  vm, err := finder.VirtualMachine(d.Get("vm_name").(string))

  if err != nil {
    return err
  }

  task, err := vm.PowerOff()

  if err != nil {
    return err
  }

  _, err = task.WaitForResult(nil)

  if err != nil {
    return err
  }

  task, err = vm.Destroy()

  if err != nil {
    return err
  }

  _, err = task.WaitForResult(nil)

  if err != nil {
    return err
  }

  return nil

}
