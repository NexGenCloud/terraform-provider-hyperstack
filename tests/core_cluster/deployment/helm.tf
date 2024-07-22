resource "helm_release" "nginx" {
  name       = "${local.name}-nginx"
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
