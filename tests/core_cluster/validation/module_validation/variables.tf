variable "name" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "artifacts_dir" {
  type = string
}

variable "kube_config_file" {
  type = string
}

variable "load_balancer_address" {
  type = string
}

variable "validation" {
  type = object({
    gpu = bool
  })
}
