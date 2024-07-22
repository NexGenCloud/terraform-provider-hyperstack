variable "region" {
  type = string
}

variable "artifacts_dir" {
  type = string
}

variable "name" {
  type = string
}

variable "environment_name" {
  type = string
}

variable "node_count" {
  type = number
}

variable "kubernetes_version" {
  type = string
}

variable "enable_public_ip" {
  type = bool
  # TODO: check with false
  default = true
}

variable "master_instance_gpu" {
  type    = string
  default = ""
}

variable "master_instance_cpus" {
  type    = number
  default = 8
}

variable "node_instance_gpu" {
  type    = string
  default = ""
}

variable "node_instance_cpus" {
  type    = number
  default = 8
}

variable "image_type" {
  type    = string
  default = "Ubuntu"
}

variable "image_version" {
  type    = string
  default = "Server 20.04 LTS"
}

variable "skip_certificate" {
  type    = bool
  default = true
}
