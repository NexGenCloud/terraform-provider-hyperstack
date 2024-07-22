locals {
  name = var.name

  master_flavor_name = module.flavor_master.name
  node_flavor_name   = module.flavor_node.name
  image_name         = module.image.name

  kubeconfig_in = yamldecode(base64decode(hyperstack_core_cluster.this.kube_config))
  cluster_in = local.kubeconfig_in["clusters"][0]

  # This magic reencoding is needed due to static terraform types
  cluster_cert = jsondecode(var.skip_certificate ? jsonencode({
    "insecure-skip-tls-verify" = true
  }) : jsonencode({
    "certificate-authority-data" = local.cluster_in["cluster"]["certificate-authority-data"]
  }))

  kubeconfig = merge(local.kubeconfig_in, {
    clusters = [
      merge(local.cluster_in, {
        cluster = merge(local.cluster_cert, {
          server = hyperstack_core_cluster.this.api_address
        })
      })
    ]
  })
}
