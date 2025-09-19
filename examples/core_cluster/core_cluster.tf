resource "hyperstack_core_cluster" "this" {
  name       = local.name
  node_count = var.node_count

  environment_name = var.environment_name
  keypair_name     = hyperstack_core_keypair.this.name

  kubernetes_version = var.kubernetes_version
  image_name         = local.image_name
  master_flavor_name = var.master_flavor
  node_flavor_name   = var.node_flavor

  # New fields with defaults
  deployment_mode = var.deployment_mode
  master_count    = var.master_count

  # Remove node_groups line
}
