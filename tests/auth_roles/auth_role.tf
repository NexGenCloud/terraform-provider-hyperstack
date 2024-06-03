resource "hyperstack_auth_role" "test_role" {
  name        = local.name
  description = var.role_description

  policies    = [for v in var.role_policies : local.mapped_policies_name[v].id]
  permissions = [for v in var.role_permissions : local.mapped_permissions_name[v].id]
}
