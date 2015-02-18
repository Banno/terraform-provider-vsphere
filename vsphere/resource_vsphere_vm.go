package vsphere

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func resourceVsphereVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceVsphereVMCreate,
		Read:   resourceVsphereVMRead,
		Update: resourceVsphereVMUpdate,
		Delete: resourceVsphereVMDelete,

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
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
			"cpus": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"memory_mb": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"customization_specification": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVsphereVMCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*govmomi.Client)
	if client == nil {
		return fmt.Errorf("client is nil")
	}

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
			NumCPUs:             d.Get("cpus").(int),
			MemoryMB:            int64(d.Get("memory_mb").(int)),
			CpuHotAddEnabled:    true,
			CpuHotRemoveEnabled: true,
			MemoryHotAddEnabled: true,
		},
		Location: types.VirtualMachineRelocateSpec{
			Pool: &rpRef,
		},
		PowerOn: true,
	}

	ipAddress := d.Get("ip_address").(string)

	specManager := client.CustomizationSpecManager()
	specItem, err := specManager.GetCustomizationSpec(d.Get("customization_specification").(string))
	if err != nil {
		return err
	}

	if ipAddress != "" {
		ip := types.CustomizationFixedIp{
			IpAddress: ipAddress,
		}
		specItem.Spec.NicSettingMap[0].Adapter.Ip = &ip
	} else {
		ip := types.CustomizationDhcpIpGenerator{}
		specItem.Spec.NicSettingMap[0].Adapter.Ip = &ip
	}

	clonespec.Customization = &specItem.Spec

	task, err := vm.Clone(folders.VmFolder, d.Get("vm_name").(string), clonespec)

	if err != nil {
		return err
	}

	_, err = task.WaitForResult(nil)

	if err != nil {
		return err
	}

	d.SetId(d.Get("vm_name").(string))

	return resourceVsphereVMRead(d, meta)
}

func resourceVsphereVMRead(d *schema.ResourceData, meta interface{}) error {
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

	ip, err := vm.WaitForIP()
	if err != nil {
		return err
	}
	d.Set("ip_address", ip)

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

func resourceVsphereVMUpdate(d *schema.ResourceData, meta interface{}) error {
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
		NumCPUs:  d.Get("cpus").(int),
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

	return resourceVsphereVMRead(d, meta)
}

func resourceVsphereVMDelete(d *schema.ResourceData, meta interface{}) error {
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
