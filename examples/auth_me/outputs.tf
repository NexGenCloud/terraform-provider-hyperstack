output "email" {
  value = data.hyperstack_auth_me.this.email
}

output "name" {
  value = data.hyperstack_auth_me.this.name
}

output "username" {
  value = data.hyperstack_auth_me.this.username
}

output "created_at" {
  value = data.hyperstack_auth_me.this.created_at
}
