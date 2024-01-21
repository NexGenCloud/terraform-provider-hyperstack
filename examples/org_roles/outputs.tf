output "roles" {
  value = data.hyperstack_auth_roles.this.roles
}

output "roles_map" {
  value = data.hyperstack_auth_role.this
}
