output "roles" {
  value = data.hyperstack_auth_roles.this.auth_roles
}

output "roles_map" {
  # TODO: shows outdated values for changed role
  value = local.mapped_roles
}

output "role_test" {
  value = hyperstack_auth_role.this
}
