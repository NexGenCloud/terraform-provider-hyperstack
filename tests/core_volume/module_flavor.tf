module "flavor" {
  source = "../../examples/core_flavor"

  region      = var.region
  gpu_name    = var.instance_gpu
  cpu_count   = var.instance_cpus
}
