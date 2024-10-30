locals {
  name = "${var.name_prefix}${random_string.this_name.result}"

  clusters = {
    for name, value in var.clusters : name => value
    if value.enabled
  }
}

locals {
  ns = "${var.name_prefix}${random_string.this_name.result}"
}
