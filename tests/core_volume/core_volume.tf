resource "hyperstack_core_volume" "no_image" {
  name = "${local.name}-no_image"

  environment_name = module.environment.environment.name
  description      = "A volume without image_id"
  volume_type      = local.volume_type
  size             = var.volume_size
  image_id         = null
  callback_url     = null
}

resource "hyperstack_core_volume" "image" {
  name = "${local.name}-image"

  environment_name = module.environment.environment.name
  description      = "A volume with image_id"
  volume_type      = local.volume_type
  size             = var.volume_size
  image_id         = module.image.id
  callback_url     = null
}
