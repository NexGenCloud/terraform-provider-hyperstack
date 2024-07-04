resource "hyperstack_core_environment" "this" {
  name   = local.name
  region = var.hyperstack_region
}
