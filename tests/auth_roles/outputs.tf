output "id" {
  value = hyperstack_auth_role.test_role.id
}

output "name" {
  value = hyperstack_auth_role.test_role.name
}

output "description" {
  value = hyperstack_auth_role.test_role.description
}

output "policies" {
  value = [
    for policy in hyperstack_auth_role.test_role.policies : {
      id          = local.mapped_policies_id[tostring(policy)].id
      name        = local.mapped_policies_id[tostring(policy)].name
      description = local.mapped_policies_id[tostring(policy)].description
    }
  ]
}

output "permissions" {
  value = [
    for permission in hyperstack_auth_role.test_role.permissions : {
      id         = local.mapped_permissions_id[tostring(permission)].id
      resource   = local.mapped_permissions_id[tostring(permission)].resource
      permission = local.mapped_permissions_id[tostring(permission)].permission
    }
  ]
}

output "created_at" {
  value = hyperstack_auth_role.test_role.created_at
}

output "roles" {
  value = [
    for v in data.hyperstack_auth_roles.this.auth_roles : {
      id          = v.id
      name        = v.name
      description = v.description
      policies    = [
        for policy in v.policies : {
          id          = policy.id
          name        = policy.name
          description = policy.description
        }
      ]
      permissions = [
        for permission in v.permissions : {
          id         = permission.id
          resource   = permission.resource
          permission = permission.permission
        }
      ]
      created_at = v.created_at
    }
  ]
}

output "role" {
  value = {
    id          = data.hyperstack_auth_role.this.id
    name        = data.hyperstack_auth_role.this.name
    description = data.hyperstack_auth_role.this.description
    policies    = [
      for policy in data.hyperstack_auth_role.this.policies : {
        id          = policy.id
        name        = policy.name
        description = policy.description
      }
    ]
    permissions = [
      for permission in data.hyperstack_auth_role.this.permissions : {
        id         = permission.id
        resource   = permission.resource
        permission = permission.permission
      }
    ]
    created_at = data.hyperstack_auth_role.this.created_at
  }
}

output "user_me_permissions" {
  value = [
    # TODO: auth_user_me_permissions -> permissions
    for permission in data.hyperstack_auth_user_me_permissions.this.auth_user_me_permissions : {
      id         = permission.id
      resource   = permission.resource
      permission = permission.permission
    }
  ]
}

output "user_permissions" {
  value = [
    # TODO: auth_user_permissions -> permissions
    for permission in data.hyperstack_auth_user_permissions.this.auth_user_permissions : {
      id         = permission.id
      resource   = permission.resource
      permission = permission.permission
    }
  ]
}
