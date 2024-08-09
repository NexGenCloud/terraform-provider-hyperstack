module "cluster" {
  source = "../../../examples/core_cluster"

  for_each = var.clusters

  region             = var.region
  artifacts_dir      = "${var.artifacts_dir}/${each.key}"
  name               = "${local.name}-${each.key}"
  node_count         = each.value.node_count
  environment_name   = module.environment.environment.name
  kubernetes_version = tolist(data.hyperstack_core_clusters_versions.this.core_clusters_versions)[0]

  master_flavor = module.flavor_master[each.key].name
  node_flavor   = module.flavor_node[each.key].name
  image_type    = each.value.image_type
  image_version = each.value.image_version
}
