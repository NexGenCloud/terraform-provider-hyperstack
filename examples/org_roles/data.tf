data "hyperstack_auth_roles" "this" {
}

data "hyperstack_auth_role" "this" {
  for_each = toset([for v in data.hyperstack_auth_roles.this.auth_roles.*.id : tostring(v)])

  id = each.value
}

data "hyperstack_keypairs" "this" {

}

data "hyperstack_environments" "this" {
}