output "flavors" {
  value = [
    for v in data.hyperstack_core_flavors.this.core_flavors : {
      gpu         = v.gpu
      region_name = v.region_name
      flavors     = [
        for flavor in v.flavors : {
          id              = flavor.id
          name            = flavor.name
          region_name     = flavor.region_name
          cpu             = flavor.cpu
          ram             = flavor.ram
          disk            = flavor.disk
          ephemeral       = flavor.ephemeral
          gpu             = flavor.gpu
          gpu_count       = flavor.gpu_count
          stock_available = flavor.stock_available
          created_at      = flavor.created_at
        }
      ]
    }
  ]
}
