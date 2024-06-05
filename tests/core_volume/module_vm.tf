module "vm" {
  source = "../../examples/core_vm"

  name                = local.name
  artifacts_directory = var.artifacts_directory
  environment_name    = module.environment.environment.name
  flavor_name         = local.flavor_name
  image_name          = local.image_name
  region              = var.region

}