module "gpu_validation" {
  source = "../../../../examples/k8s/gpu_validation"
  count  = var.validation.gpu ? 1 : 0

  depends_on = [
    module.simple,
  ]

  name             = var.name
  namespace        = module.simple[0].namespace
  kube_config_file = var.kube_config_file
  cluster_name     = var.cluster_name
  artifacts_dir    = var.artifacts_dir
}
