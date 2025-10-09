# Create multiple normal (non-bootable) volumes to attach to the VM
resource "hyperstack_core_volume" "data" {
  count = var.volume_count

  name = "${local.name}-data-volume-${count.index + 1}"

  environment_name = module.environment.environment.name
  description      = "Test data volume ${count.index + 1} for volume attachment test"
  volume_type      = local.volume_type
  size             = var.volume_size
  image_id         = null
  callback_url     = null
}
