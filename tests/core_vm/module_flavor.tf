module "flavor" {
  source = "../../examples/core_flavor"

  for_each = var.vms

  region    = var.region
  name      = each.value.flavor.name
  gpu_name  = each.value.flavor.gpu_name
  gpu_count = each.value.flavor.gpu_count
  cpu_count = each.value.flavor.cpu_count
}
