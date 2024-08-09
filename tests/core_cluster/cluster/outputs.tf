output "clusters" {
  value = {
    for k, v in module.cluster : k => {
      id                 = v.id
      name               = v.name
      environment_name   = v.environment_name
      kubernetes_version = v.kubernetes_version
      api_address        = v.api_address
      kube_config        = v.kube_config
      kube_config_file   = v.kube_config_file
      status             = v.status
      status_reason      = v.status_reason
      node_count         = v.node_count
      keypair_name       = v.keypair_name
      created_at         = v.created_at
    }
  }
}

output "clusters_versions" {
  value = data.hyperstack_core_clusters_versions.this.core_clusters_versions
}
