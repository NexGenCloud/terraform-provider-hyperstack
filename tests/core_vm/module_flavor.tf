module "flavor" {
  source = "../../examples/core_flavor"

  flavor_region = var.hyperstack_region
  flavor_gpu    = var.instance_gpu
  flavor_cpus   = var.instance_cpus
}
