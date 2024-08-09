module "flavor_master" {
  source = "../../../examples/core_flavor"

  for_each = var.clusters

  region    = var.region
  name      = each.value.master_flavor.name
  gpu_name  = each.value.master_flavor.gpu_name
  gpu_count = each.value.master_flavor.gpu_count
  cpu_count = each.value.master_flavor.cpu_count
}

module "flavor_node" {
  source = "../../../examples/core_flavor"

  for_each = var.clusters

  region    = var.region
  name      = each.value.node_flavor.name
  gpu_name  = each.value.node_flavor.gpu_name
  gpu_count = each.value.node_flavor.gpu_count
  cpu_count = each.value.node_flavor.cpu_count
}
