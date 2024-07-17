resource "hyperstack_core_keypair" "this" {
  name        = local.name
  environment = module.environment.environment.name
  public_key  = tls_private_key.this.public_key_openssh
}
