resource "hyperstack_core_virtual_machine" "this" {
  name = local.name

  environment_name = hyperstack_core_environment.this.name
  flavor_name      = "s" // TODO: from datasource

  image_name = "Ubuntu Server 22.04 LTS (Jammy Jellyfish)" // TODO: from datasource

  key_name = hyperstack_core_keypair.this.name
}
