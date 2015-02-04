#VMWare VSphere Provider For Terraform

Allows you to define infrastructure for VMWare VSphere. 

##Requirements: 
* A working VSphere 5.5 setup
* **Govmomi** for interfacing with vsphere cleanly
* An existing "Customization Specification"

##Use Case
Our current use case is using packer to upload a VM to VSphere, turn it on, and
off to make sure that VMware knows that VMWare tools are installed. From there
we mark it as a template and use that template to deploy from.

##Environment Variables
export VSPHERE_USERNAME="username"
export VSPHERE_PASSWORD="password"
export VSPHERE_URL="vsphere.my.domain.com" 

##Minimal configuration
```
resource "vsphere_vm" "machine_name" {
  template_name = "packer-virtualbox-iso-timestamp-vm"
  vm_name = "machine"
  memory_mb = 1024
  cpus = 2
  ip_address = "127.0.0.1" #This is optional, will use DHCP if unset
  customization_specification = "Ubuntu"
}
```

