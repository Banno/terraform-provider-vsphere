# DEPRECATED IN FAVOR OF TERRAFORM'S OFFICIAL VSPHERE PROVIDER: https://www.terraform.io/docs/providers/vsphere/index.html


#VMWare VSphere Provider For Terraform

Allows you to define infrastructure for VMWare VSphere. 

##Requirements: 
* A working VSphere 5.5 setup
* **Govmomi** for interfacing with vsphere cleanly
* An existing "Customization Specification"

##Use Case
Our current use case is using packer to upload a VM to VSphere, turn it on, and
off to make sure that VMware knows that VMWare tools are installed. From there
we mark it as a template and use that template to deploy

## Provider

The provider can be configured manually via
```
provider "vsphere" {
  vsphere_host = "vcenter.my.domain.com"
  vsphere_username = "username"
  vsphere_password = "password"
}

```

or with environmental variables:
```
export VSPHERE_USERNAME="username"
export VSPHERE_PASSWORD="password"
export VSPHERE_HOST="vsphere.my.domain.com" 
```

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

##Optional configuration
```
resource "vsphere_vm" "machine_name" {
  ...
  resource_pool = "Dev Cluster/Resources/Dev Pool"
  ...
}
```

##How to build
* go get
* go install
