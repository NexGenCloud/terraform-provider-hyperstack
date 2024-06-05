resource "hyperstack_core_keypair" "this" {
  name        = var.name
  environment = var.environment_name
  public_key  = tls_private_key.this.public_key_openssh
}
