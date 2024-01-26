resource "hyperstack_core_environment" "this" {
  name   = local.name
  region = "staging-CA-1" // TODO: from datasource
}