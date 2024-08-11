output "id" {
  value = hyperstack_core_cluster.this.id
}

output "name" {
  value = hyperstack_core_cluster.this.name
}

output "environment_name" {
  value = hyperstack_core_cluster.this.environment_name
}

output "kubernetes_version" {
  value = hyperstack_core_cluster.this.kubernetes_version
}

output "api_address" {
  value = hyperstack_core_cluster.this.api_address
}

output "kube_config" {
  value = hyperstack_core_cluster.this.kube_config
}

output "kube_config_file" {
  value = local_sensitive_file.kubeconfig.filename
}

output "status" {
  value = hyperstack_core_cluster.this.status
}

output "status_reason" {
  value = hyperstack_core_cluster.this.status_reason
}

output "node_count" {
  value = hyperstack_core_cluster.this.node_count
}

output "keypair_name" {
  value = hyperstack_core_cluster.this.keypair_name
}

# output "enable_public_ip" {
#   value = hyperstack_core_cluster.this.enable_public_ip
# }

output "created_at" {
  value = hyperstack_core_cluster.this.created_at
}
