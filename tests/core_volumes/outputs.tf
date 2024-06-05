output "volumes" {
  value = [
    for v in data.hyperstack_core_volumes.this.core_volumes : {
      id   = v.id
      name = v.name
      environment = {
        name = v.environment.name
      }
      description  = v.description
      volume_type  = v.volume_type
      size         = v.size
      status       = v.status
      bootable     = v.bootable
      image_id     = v.image_id
      callback_url = v.callback_url
      created_at   = v.created_at
      updated_at   = v.updated_at
    }
  ]
}
