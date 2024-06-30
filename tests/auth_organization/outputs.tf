output "id" {
  value = data.hyperstack_auth_organization.this.id
}

output "name" {
  value = data.hyperstack_auth_organization.this.name
}

// Not available on staging yet ??
# output "credit" {
#   value = data.hyperstack_auth_organization.this.credit
# }
#
# output "threshold" {
#   value = data.hyperstack_auth_organization.this.threshold
# }
#
# output "total_instances" {
#   value = data.hyperstack_auth_organization.this.total_instances
# }
#
# output "total_volumes" {
#   value = data.hyperstack_auth_organization.this.total_volumes
# }
#
# output "total_containers" {
#   value = data.hyperstack_auth_organization.this.total_containers
# }
#
# output "total_clusters" {
#   value = data.hyperstack_auth_organization.this.total_clusters
# }

output "users" {
  value = {
    for v in data.hyperstack_auth_organization.this.users : v.id => {
      id         = v.id
      sub        = v.sub
      email      = v.email
      username   = v.username
      name       = v.name
      role       = v.role
      rbac_roles = [
        for role in v.rbac_roles : {
          name = role.name
        }
      ]
      joined_at  = v.joined_at
      // Not available on staging yet ??
      #last_login = v.last_login
    }
  }
}

output "created_at" {
  value = data.hyperstack_auth_organization.this.created_at
}
