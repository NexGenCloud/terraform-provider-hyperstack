output "gpus" {
  value = [
    for v in data.hyperstack_core_gpus.this.core_gpus : {
      id      = v.id
      name    = v.name
      regions = [
        for region in v.regions : {
          id   = region.id
          name = region.name
        }
      ]
      example_metadata = v.example_metadata
      created_at       = v.created_at
      updated_at       = v.updated_at
    }
  ]
}
