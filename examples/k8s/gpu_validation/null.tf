resource "null_resource" "probe_gpu_operator" {
  triggers = {
    script = md5(file("${path.module}/data/probe_gpu_operator.sh"))
  }

  provisioner "local-exec" {
    environment = {
      KUBECONFIG   = var.kube_config_file
      CLUSTER_NAME = var.cluster_name
    }

    command = "${path.module}/data/probe_gpu_operator.sh > ${var.artifacts_dir}/probe_gpu_operator.log"
  }
}

resource "null_resource" "cuda_verify" {
  triggers = {
    script = md5(file("${path.module}/data/cuda_verify.sh"))
  }

  provisioner "local-exec" {
    environment = {
      NAME         = kubernetes_pod_v1.cuda_vectoradd.metadata[0].name
      NAMESPACE    = kubernetes_pod_v1.cuda_vectoradd.metadata[0].namespace
      KUBECONFIG   = var.kube_config_file
      CLUSTER_NAME = var.cluster_name
    }
    command = "${path.module}/data/cuda_verify.sh > ${var.artifacts_dir}/cuda_verify.log"
  }

  depends_on = [
    null_resource.probe_gpu_operator,
    kubernetes_pod_v1.cuda_vectoradd,
  ]
}
