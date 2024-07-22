variable "name_prefix" {
  type = string
}

# Comes from cluster:
# ---

variable "cluster_name" {
  type = string
}

variable "api_address" {
  type = string
}

variable "kube_config_file" {
  type = string
}
