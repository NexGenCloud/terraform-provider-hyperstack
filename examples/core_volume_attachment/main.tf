# Create a virtual machine
resource "hyperstack_core_virtual_machine" "example" {
  name             = var.vm_name
  environment_name = var.environment_name
  key_name         = var.keypair_name
  image_name       = var.image_name
  flavor_name      = var.flavor_name
}

# Create multiple data volumes (non-bootable)
resource "hyperstack_core_volume" "data" {
  count            = var.volume_count
  name             = "${var.vm_name}-data-volume-${count.index + 1}"
  environment_name = var.environment_name
  size             = var.volume_size
  volume_type      = var.volume_type
  description      = "Data volume ${count.index + 1} for ${var.vm_name}"
  image_id         = null # Non-bootable data volume
}

# Attach all volumes to the VM
resource "hyperstack_core_volume_attachment" "example" {
  vm_id      = hyperstack_core_virtual_machine.example.id
  volume_ids = [for vol in hyperstack_core_volume.data : vol.id]
  protected  = false

  depends_on = [
    hyperstack_core_volume.data
  ]
}
