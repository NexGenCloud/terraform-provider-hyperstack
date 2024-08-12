locals {
  name = var.name

  master_flavor_name = var.master_flavor
  node_flavor_name   = var.node_flavor
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

  # Extract the IP address using a regular expression that captures text between https:// and :6443
  # TODO: this should come from api
  load_balancer_address = regex("https?://([0-9\\.]+):6443", hyperstack_core_cluster.this.api_address)[0]
}