output "id" {
  value = module.cluster.id
}

output "name" {
  value = module.cluster.name
}

output "environment_name" {
  value = module.cluster.environment_name
}

output "kubernetes_version" {
  value = module.cluster.kubernetes_version
}

output "api_address" {
  value = module.cluster.api_address
}

output "kube_config" {
  value = module.cluster.kube_config
}

output "kube_config_file" {
  value = module.cluster.kube_config_file
}

output "status" {
  value = module.cluster.status
}

output "status_reason" {
  value = module.cluster.status_reason
}

output "node_count" {
  value = module.cluster.node_count
}

output "node_addresses" {
  value = module.cluster.node_addresses
}

output "keypair_name" {
  value = module.cluster.keypair_name
}

output "enable_public_ip" {
  value = module.cluster.enable_public_ip
}

output "created_at" {
  value = module.cluster.created_at
}

output "clusters_versions" {
  value = data.hyperstack_core_clusters_versions.this.core_clusters_versions
}
