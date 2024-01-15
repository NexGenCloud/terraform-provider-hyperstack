output "user_email" {
  value = data.hyperstack_auth_me.this.user.email
}

output "user_name" {
  value = data.hyperstack_auth_me.this.user.name
}

output "user_username" {
  value = data.hyperstack_auth_me.this.user.username
}

output "user_created_at" {
  value = data.hyperstack_auth_me.this.user.created_at
}
