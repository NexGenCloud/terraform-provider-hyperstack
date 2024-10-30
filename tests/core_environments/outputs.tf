output "id" {
  value = hyperstack_core_environment.test_environment.id
}

output "name" {
  value = hyperstack_core_environment.test_environment.name
}

output "region" {
  value = hyperstack_core_environment.test_environment.region
}

output "created_at" {
  value = hyperstack_core_environment.test_environment.created_at
}

output "environments" {
  value = [
    for v in data.hyperstack_core_environments.this.core_environments : {
      id         = v.id
      name       = v.name
      region     = v.region
      created_at = v.created_at
    }
  ]
}

output "environment" {
  value = {
    id         = data.hyperstack_core_environment.this.id
    name       = data.hyperstack_core_environment.this.name
    region     = data.hyperstack_core_environment.this.region
    created_at = data.hyperstack_core_environment.this.created_at
  }
}