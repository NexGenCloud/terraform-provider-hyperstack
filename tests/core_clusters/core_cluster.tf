resource "hyperstack_core_cluster" "this" {
  name       = local.name
  node_count = var.node_count

  environment_name = module.environment.environment.name
  keypair_name     = hyperstack_core_keypair.this.name

  kubernetes_version = tolist(data.hyperstack_core_clusters_versions.this.core_clusters_versions)[0]
  image_name         = local.image_name
  master_flavor_name = local.master_flavor_name
  node_flavor_name   = local.node_flavor_name

  enable_public_ip = var.enable_public_ip
}
