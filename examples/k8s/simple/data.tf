data "kubernetes_all_namespaces" "this" {
  depends_on = [
    kubernetes_namespace.this,
  ]
}
