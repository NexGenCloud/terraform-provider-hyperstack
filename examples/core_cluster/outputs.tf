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
  value = try(hyperstack_core_cluster.this.api_address, null)
}

output "kube_config" {
  value = try(hyperstack_core_cluster.this.kube_config, null)
}

output "kube_config_file" {
  value = try(local_sensitive_file.kubeconfig.filename, null)
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

output "created_at" {
  value = hyperstack_core_cluster.this.created_at
}

output "artifacts_dir" {
  value = var.artifacts_dir
}

output "load_balancer_address" {
  value = local.load_balancer_address
}

# New outputs for enhanced cluster configuration
output "deployment_mode" {
  description = "Deployment mode of the cluster"
  value       = hyperstack_core_cluster.this.deployment_mode
}

output "master_count" {
  description = "Number of master nodes"
  value       = hyperstack_core_cluster.this.master_count
}
