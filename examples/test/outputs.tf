output "user_email" {
  value = data.hyperstack_auth_me.this.email
}

output "user_name" {
  value = data.hyperstack_auth_me.this.name
}

output "user_username" {
  value = data.hyperstack_auth_me.this.username
}

output "user_created_at" {
  value = data.hyperstack_auth_me.this.created_at
}
