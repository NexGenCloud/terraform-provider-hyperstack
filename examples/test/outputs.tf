output "user" {
  value = data.hyperstack_auth_me.this.user
}

output "organization" {
  value = data.hyperstack_auth_organizations.this.organization
}
