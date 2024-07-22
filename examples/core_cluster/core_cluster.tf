resource "hyperstack_core_cluster" "this" {
  name       = local.name
  node_count = var.node_count

  environment_name = var.environment_name
  keypair_name     = hyperstack_core_keypair.this.name

  kubernetes_version = var.kubernetes_version
  image_name         = local.image_name
  master_flavor_name = local.master_flavor_name
  node_flavor_name   = local.node_flavor_name

  enable_public_ip = var.enable_public_ip
}
