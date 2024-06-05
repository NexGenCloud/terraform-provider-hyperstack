data "hyperstack_core_environments" "this" {
  depends_on = [
    module.environment.environment,
  ]
}

data "hyperstack_core_environment" "this" {
  id = module.environment.environment.id
}
