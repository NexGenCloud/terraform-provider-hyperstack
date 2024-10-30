locals {
  name = "${var.name_prefix}${random_string.this_name.result}"

  mapped_policies_name    = {for v in data.hyperstack_auth_policies.this.auth_policies : v.name => v}
  mapped_policies_id    = {for v in data.hyperstack_auth_policies.this.auth_policies : v.id => v}
  mapped_permissions_name = {for v in data.hyperstack_auth_permissions.this.auth_permissions : v.permission => v}
  mapped_permissions_id = {for v in data.hyperstack_auth_permissions.this.auth_permissions : v.id => v}
}
