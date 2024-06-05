output "id" {
  value = module.environment.environment.id
}

output "name" {
  value = module.environment.environment.name
}

output "region" {
  value = module.environment.environment.region
}

output "created_at" {
  value = module.environment.environment.created_at
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