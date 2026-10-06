resource "hyperstack_core_keypair" "this" {
  name             = var.name
  environment_name = var.environment_name
  public_key       = tls_private_key.this.public_key_openssh
}
