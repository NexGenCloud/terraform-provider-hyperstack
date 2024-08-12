module "image" {
  source = "../../examples/core_image"

  for_each = var.vms

  image_region  = var.region
  image_type    = each.value.image_type
  image_version = each.value.image_version
}
