locals {
  name = "${var.name_prefix}${random_string.this_name.result}"
}
