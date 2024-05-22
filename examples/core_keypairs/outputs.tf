output "id" {
  value = hyperstack_core_keypair.test_keypair.id
}

output "name" {
  value = hyperstack_core_keypair.test_keypair.name
}

output "environment" {
  value = hyperstack_core_keypair.test_keypair.environment
}

output "public_key" {
  value = hyperstack_core_keypair.test_keypair.public_key
}

output "fingerprint" {
  value = hyperstack_core_keypair.test_keypair.fingerprint
}

output "created_at" {
  value = hyperstack_core_keypair.test_keypair.created_at
}

output "keypairs" {
  value = [
    for v in data.hyperstack_core_keypairs.this.core_keypairs : {
      id          = v.id
      name        = v.name
      environment = v.environment
      public_key  = v.public_key
      fingerprint = v.fingerprint
      created_at  = v.created_at
    }
  ]
}

output "keypair" {
  value = {
    id          = data.hyperstack_core_keypair.this.id
    name        = data.hyperstack_core_keypair.this.name
    environment = data.hyperstack_core_keypair.this.environment
    public_key  = data.hyperstack_core_keypair.this.public_key
    fingerprint = data.hyperstack_core_keypair.this.fingerprint
    created_at  = data.hyperstack_core_keypair.this.created_at
  }
}