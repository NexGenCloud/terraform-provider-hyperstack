locals {
  name = "${var.name_prefix}${random_string.this_name.result}"

  master_flavor_name = module.flavor_master.name
  node_flavor_name = module.flavor_node.name
  image_name  = module.image.name
}
