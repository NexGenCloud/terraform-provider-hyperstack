resource "null_resource" "vllm_validate" {
  triggers = {
    cluster_instance_ids = md5(file("${path.module}/data/vllm_validate.sh"))
  }
  provisioner "local-exec" {
    environment = {
      KUBECONFIG   = var.kube_config_file
      MODEL        = var.vllm_model_name
      SERVER_IP    = var.load_balancer_address
      CLUSTER_NAME = var.cluster_name
    }
    command = "${path.module}/data/vllm_validate.sh > ${var.artifacts_dir}/vllm_validate.log"
  }

  depends_on = [
    kubernetes_deployment.vllm,
    kubernetes_ingress_v1.vllm,
  ]
}
