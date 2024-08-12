module "simple" {
  source = "../../../../examples/k8s/simple"
  count  = 1

  name = var.name
}
