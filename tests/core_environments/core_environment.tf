resource "hyperstack_core_environment" "test_environment" {
  name   = local.name
  region = var.region
}