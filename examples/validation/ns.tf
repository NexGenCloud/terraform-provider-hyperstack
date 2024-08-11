resource "kubernetes_namespace" "this" {
  metadata {
    name = var.ns
  }
}

provider "kubernetes" {
  config_path = module.cluster.kube_config_file
  insecure    = true
}

