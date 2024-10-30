data "hyperstack_core_environments" "this" {
  depends_on = [
    hyperstack_core_environment.test_environment,
  ]
}

data "hyperstack_core_environment" "this" {
  id = hyperstack_core_environment.test_environment.id
}
