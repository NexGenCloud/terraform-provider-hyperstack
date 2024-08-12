resource "kubernetes_pod_v1" "cuda_vectoradd" {
  metadata {
    name      = "${var.name}-cuda-vectoradd"
    namespace = local.namespace
  }

  spec {
    restart_policy = "OnFailure"

    container {
      name  = "cuda-vectoradd"
      image = "nvcr.io/nvidia/k8s/cuda-sample:vectoradd-cuda11.7.1-ubuntu20.04"
    }
  }

  target_state = ["Succeeded", "Running"]

  depends_on = [
    null_resource.probe_gpu_operator,
  ]
}