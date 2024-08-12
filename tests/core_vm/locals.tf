locals {
  name = "${var.name_prefix}${random_string.this_name.result}"

  vms_tmp = flatten([
    for name, value in var.vms : [
      for i in range(value.count) : {
        name = name
        key  = "${name}-${i}"
      }
    ]
    if value.enabled
  ])
  vms = {
    for value in local.vms_tmp : value.key => merge(var.vms[value.name], {
      name = value.name
    })
  }
}