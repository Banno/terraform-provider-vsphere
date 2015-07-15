provider "vsphere" {
  vsphere_host = ""
  vsphere_username = ""
  vsphere_password = ""
}

resource "vsphere_vm" "testing_vsphere_instance" {
  template_name = "packer-base-vsphere-1424186188-vm"
  vm_name = "testingvm1"
  memory_mb = "8192"
  cpus = "4"
  customization_specification = "Ubuntu"
  # resource_pool = "Dev Cluster/Resources/Dev Pool" # Optional
}
