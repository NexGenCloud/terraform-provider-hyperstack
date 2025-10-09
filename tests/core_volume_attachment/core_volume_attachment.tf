# Attach normal data volumes to the VM
resource "hyperstack_core_volume_attachment" "data" {
  for_each = local.vms

  vm_id      = module.vms[each.key].vm.id
  volume_ids = [for vol in hyperstack_core_volume.data : vol.id]
  protected  = var.volume_protected

  depends_on = [
    hyperstack_core_volume.data,
    module.vms
  ]
}
