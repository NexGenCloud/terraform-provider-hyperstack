locals {
  //name = "${var.name_prefix}${random_string.this_name.result}"
  name = "${var.name_prefix}predefined-vol-attach"

  vms_types = {
    for name, value in var.vms : name => value
    if value.enabled
  }

  vms_tmp = flatten([
    for name, value in local.vms_types : [
      for i in range(value.count) : {
        name = name
        key  = "${name}-${i}"
      }
    ]
  ])
  vms = {
    for value in local.vms_tmp : value.key => merge(var.vms[value.name], {
      name = value.name
    })
  }

  volume_type = tolist(data.hyperstack_core_volume_types.this.core_volume_types)[0]
}
