module "cluster" {
  source = "../../../examples/core_cluster"

  region             = var.region
  artifacts_dir      = var.artifacts_dir
  name               = local.name
  node_count         = var.node_count
  environment_name   = module.environment.environment.name
  enable_public_ip   = var.enable_public_ip
  kubernetes_version = tolist(data.hyperstack_core_clusters_versions.this.core_clusters_versions)[0]

  master_instance_gpu  = var.master_instance_gpu
  master_instance_cpus = var.master_instance_cpus
  node_instance_gpu    = var.node_instance_gpu
  node_instance_cpus   = var.node_instance_cpus
  image_type           = var.image_type
  image_version        = var.image_version
}
