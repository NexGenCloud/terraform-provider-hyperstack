module "vm" {
  source = "../../examples/core_vm"

  name                   = local.name
  artifacts_directory    = var.artifacts_directory
  environment_name       = module.environment.environment.name
  flavor_name            = local.flavor_name
  image_name             = local.image_name
  region                 = var.hyperstack_region
  ingress_ports          = [22, 80, 443]
  create_bootable_volume = false
  user_data              = data.cloudinit_config.this.rendered
  // TODO: Setting this to true results in error state
  assign_floating_ip     = true
}
