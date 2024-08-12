output "namespace" {
  value = kubernetes_namespace.this.id
}

output "namespaces" {
  value = data.kubernetes_all_namespaces.this.namespaces
}

output "helm_nginx_name" {
  value = helm_release.nginx.name
}
