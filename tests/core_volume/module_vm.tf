module "vm" {
  source = "../../examples/core_vm"

  name                = local.name
  artifacts_directory = var.artifacts_dir
  environment_name    = module.environment.environment.name
  flavor_name         = local.flavor_name
  image_name          = local.image_name
  region              = var.region

}
