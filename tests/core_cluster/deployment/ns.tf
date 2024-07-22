resource "kubernetes_namespace" "this" {
  metadata {
    name = local.name
  }
}
