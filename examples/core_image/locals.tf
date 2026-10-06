locals {
  images = data.hyperstack_core_images.this.core_images
  image = flatten([
    for image in local.images : [
      for i in image.images : i
      if i.type == var.image_type && i.version == var.image_version && i.region_name == var.image_region
    ]
  ])[0]
}