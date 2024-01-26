resource "hyperstack_environment" "test-env" {
  name="test-tf-env-${random_string.this_name.result}"
  region="staging-CA-1"
}