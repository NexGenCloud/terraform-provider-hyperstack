module "flavor_master" {
  source = "../../examples/core_flavor"

  flavor_region = var.region
  flavor_gpu    = var.master_instance_gpu
  flavor_cpus   = var.master_instance_cpus
}

module "flavor_node" {
  source = "../../examples/core_flavor"

  flavor_region = var.region
  flavor_gpu    = var.node_instance_gpu
  flavor_cpus   = var.node_instance_cpus
}
