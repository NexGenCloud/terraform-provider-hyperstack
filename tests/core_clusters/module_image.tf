module "image" {
  source = "../../examples/core_image"

  image_region  = var.hyperstack_region
  image_type    = var.image_type
  image_version = var.image_version
}
