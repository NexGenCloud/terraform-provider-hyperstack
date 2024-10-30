locals {
  flavors = data.hyperstack_core_flavors.this.core_flavors
  flavor = coalescelist(flatten([
    for flavor in local.flavors : coalesce([
        for f in flavor.flavors : f
        if true
          && (var.name == null || f.name == var.name)
          && (var.cpu_count == null || f.cpu == var.cpu_count)
          && (var.gpu_name == null || f.gpu == var.gpu_name)
          && (var.gpu_count == null || f.gpu_count == var.gpu_count)
          && (f.region_name == var.region)
      ])
    ]))[0]
}