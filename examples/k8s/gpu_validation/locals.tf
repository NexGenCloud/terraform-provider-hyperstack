locals {
  namespace = data.kubernetes_namespace.this.metadata[0].name
}