resource "local_sensitive_file" "kubeconfig" {
  content         = jsonencode(local.kubeconfig)
  filename        = "${var.artifacts_dir}/kubeconfig"
  file_permission = "0644"
}
