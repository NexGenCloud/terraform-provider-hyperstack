resource "hyperstack_core_keypair" "test_keypair" {
  name        = local.name
  environment = hyperstack_core_environment.this.name
  public_key  = var.public_key
}
