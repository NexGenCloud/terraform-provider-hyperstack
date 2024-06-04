data "hyperstack_core_keypairs" "this" {
  depends_on = [
    hyperstack_core_keypair.test_keypair,
  ]
}

data "hyperstack_core_keypair" "this" {
  id = hyperstack_core_keypair.test_keypair.id
}
