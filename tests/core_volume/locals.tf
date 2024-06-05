locals {
  name = "${var.name_prefix}${random_string.this_name.result}"

  flavor_name = module.flavor.name
  image_name  = module.image.name

  volume_type = tolist(data.hyperstack_core_volume_types.this.core_volume_types)[0]
}