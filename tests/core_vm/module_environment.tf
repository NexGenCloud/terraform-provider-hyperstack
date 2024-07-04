module "environment" {
  source = "../../examples/core_environment"

  name   = local.name
  region = var.hyperstack_region
}
