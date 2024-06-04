data "hyperstack_auth_roles" "this" {
}

data "hyperstack_auth_policies" "this" {
}

data "hyperstack_auth_permissions" "this" {
}

data "hyperstack_auth_organization" "this" {
}

data "hyperstack_auth_user_me_permissions" "this" {
}

data "hyperstack_auth_user_permissions" "this" {
  id = data.hyperstack_auth_organization.this.users[0].id
}

data "hyperstack_auth_role" "this" {
  id = hyperstack_auth_role.test_role.id
}
