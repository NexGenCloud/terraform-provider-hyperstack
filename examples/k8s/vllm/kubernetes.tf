# TODO: change to daemonset
resource "kubernetes_deployment" "vllm" {
  metadata {
    name      = "${var.name}-vllm"
    namespace = local.namespace
    labels = {
      app = "${var.name}-vllm-app"
    }
  }

  spec {
    # TODO: change to daemonset and remove
    replicas = 2

    selector {
      match_labels = {
        app = "${var.name}-vllm-app"
      }
    }

    strategy {
      rolling_update {
        max_surge       = "25%"
        max_unavailable = "25%"
      }
      type = "RollingUpdate"
    }

    template {
      metadata {
        labels = {
          app = "${var.name}-vllm-app"
        }
      }

      spec {
        container {
          name              = "vllm-openai"
          image             = "vllm/vllm-openai:latest"
          image_pull_policy = "Always"

          command = [
            "python3",
            "-m",
            "vllm.entrypoints.openai.api_server",
            "--model",
            var.vllm_model_name
          ]

          port {
            container_port = 8000
            protocol       = "TCP"
          }

          liveness_probe {
            http_get {
              path   = "/health"
              port   = 8000
              scheme = "HTTP"
            }
            initial_delay_seconds = 120
            period_seconds        = 5
            timeout_seconds       = 1
            success_threshold     = 1
            failure_threshold     = 3
          }

          readiness_probe {
            http_get {
              path   = "/health"
              port   = 8000
              scheme = "HTTP"
            }
            initial_delay_seconds = 60
            period_seconds        = 5
            timeout_seconds       = 1
            success_threshold     = 1
            failure_threshold     = 3
          }

          resources {
            limits = {
              "nvidia.com/gpu" = "1"
            }
            requests = {
              "nvidia.com/gpu" = "1"
            }
          }

          volume_mount {
            name       = "cache-volume"
            mount_path = "/root/.cache/huggingface"
          }
        }

        volume {
          name = "cache-volume"

          empty_dir {}
        }
      }
    }
  }
}

resource "kubernetes_service" "vllm" {
  metadata {
    name      = "${var.name}-vllm"
    namespace = local.namespace
    labels = {
      app = "${var.name}-vllm-app"
    }
  }

  spec {
    selector = {
      app = "${var.name}-vllm-app"
    }

    port {
      port        = 8000
      target_port = 8000
      protocol    = "TCP"
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_ingress_v1" "vllm" {
  wait_for_load_balancer = false

  metadata {
    name      = "${var.name}-vllm"
    namespace = local.namespace
    annotations = {
      "nginx.ingress.kubernetes.io/rewrite-target" = "/"
    }
  }

  spec {
    ingress_class_name = "nginx"
    rule {
      http {
        path {
          path = "/"
          backend {
            service {
              name = kubernetes_service.vllm.metadata[0].name
              port {
                number = 8000
              }
            }
          }
        }
      }
    }
  }
}
