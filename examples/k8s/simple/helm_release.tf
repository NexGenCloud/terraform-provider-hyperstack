resource "helm_release" "nginx" {
  name      = "${var.name}-nginx"
  namespace = kubernetes_namespace.this.id

  repository = "https://charts.bitnami.com/bitnami"
  chart      = "nginx"
  version    = "18.1.6"

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
