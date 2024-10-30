module "vms" {
  source = "../../examples/core_vm"

  for_each = local.vms

  name                   = "${local.name}-${each.key}"
  artifacts_directory    = "${var.artifacts_dir}/${each.key}"
  environment_name       = module.environment.environment.name
  flavor_name            = module.flavor[each.value.name].name
  image_name             = module.image[each.value.name].name
  region                 = var.region
  ingress_ports          = [22, 80, 443]
  create_bootable_volume = false
  user_data              = data.cloudinit_config.this.rendered
  // TODO: Setting this to true results in error state
  assign_floating_ip     = true
}
