output "roles" {
  value = data.hyperstack_auth_roles.this.roles
}

output "roles_map" {
  # TODO: shows outdated values for changed role
  value = local.mapped_roles
}

output "role_test" {
  value = hyperstack_auth_role.this
}

output "keypair_test" {
  value = data.hyperstack_keypairs.this
}

output "env_test" {
  value = data.hyperstack_environments.this
}