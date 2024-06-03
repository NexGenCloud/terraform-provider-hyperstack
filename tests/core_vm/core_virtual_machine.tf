resource "hyperstack_core_virtual_machine" "this" {
  name = local.name

  environment_name = hyperstack_core_environment.this.name

  flavor_name = local.flavor_name
  image_name  = local.image_name

  user_data = data.template_cloudinit_config.this.rendered

  key_name = hyperstack_core_keypair.this.name

  callback_url = null

  assign_floating_ip = true

  // TODO: Setting this to true results in error state
  create_bootable_volume = false
}
