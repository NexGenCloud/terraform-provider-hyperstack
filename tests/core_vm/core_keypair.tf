resource "hyperstack_core_keypair" "this" {
  name        = local.name
  environment = hyperstack_core_environment.this.name
  public_key  = tls_private_key.this.public_key_openssh
}
