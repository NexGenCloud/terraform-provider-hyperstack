locals {
  mapped_roles = {for v in data.hyperstack_auth_role.this : v.role.name => v.role}
}
