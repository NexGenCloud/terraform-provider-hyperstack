output "user" {
  value = data.hyperstack_auth_me.this
}

output "organization" {
  value = data.hyperstack_auth_organizations.this
}
