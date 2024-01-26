resource "hyperstack_auth_role" "this" {
  name        = "testrole-${random_string.this_name.result}"
  description = "Test admin role"
  policies    = [for v in local.mapped_roles["admins"].policies : v.id]
  permissions = [for v in local.mapped_roles["admins"].permissions : v.id]
}
