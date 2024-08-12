module "vllm" {
  source = "../../../../examples/k8s/vllm"
  count  = var.validation.gpu ? 1 : 0

  depends_on = [
    module.gpu_validation,
  ]

  name                  = var.name
  namespace             = module.simple[0].namespace
  kube_config_file      = var.kube_config_file
  cluster_name          = var.cluster_name
  artifacts_dir         = var.artifacts_dir
  load_balancer_address = var.load_balancer_address
}
