module "validation" {
  source = "../../../examples/validation"
  for_each = var.clusters
  kube_config_file = each.value.kube_config_file
  api_address      = each.value.api_address
  cluster_name     = each.value.cluster_name
  artifacts_dir    = "${var.artifacts_dir}/${each.key}"
  ns               = local.ns
}