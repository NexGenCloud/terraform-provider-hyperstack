resource "hyperstack_core_virtual_machine" "this" {
  name = var.name

  environment_name = var.environment_name

  flavor_name = var.flavor_name
  image_name  = var.image_name

  user_data = var.user_data

  key_name = hyperstack_core_keypair.this.name

  callback_url = var.callback_url

  assign_floating_ip = var.assign_floating_ip

  // TODO: Setting this to true results in error state
  create_bootable_volume = var.create_bootable_volume
}
