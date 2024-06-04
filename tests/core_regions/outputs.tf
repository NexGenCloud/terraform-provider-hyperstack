output "regions" {
  value = [
    for v in data.hyperstack_core_regions.this.core_regions : {
      id          = v.id
      name        = v.name
      description = v.description
    }
  ]
}
