locals {
  flavors = data.hyperstack_core_flavors.this.core_flavors
  flavor = coalescelist(flatten([
    for flavor in local.flavors : coalesce([
        for f in flavor.flavors : f
        if f.cpu == var.flavor_cpus && f.gpu == var.flavor_gpu && f.region_name == var.flavor_region
      ])
    ]))[0]
}