resource "helm_release" "nginx" {
  name       = "${var.ns}-nginx"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "nginx"
  version    = "18.1.6"

  namespace = kubernetes_namespace.this.id

  values = [
    <<EOF
extraEnvVars:
  - name: LOG_LEVEL
    value: error
EOF
  ]

  set {
    name  = "ingress.enabled"
    value = "false"
  }

  set {
    name  = "service.type"
    value = "ClusterIP"
  }
}
provider "helm" {
  kubernetes {
    config_path = module.cluster.kube_config_file
    insecure    = true
  }
}
