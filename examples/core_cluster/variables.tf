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

variable "master_flavor" {
  type = string
}

variable "node_flavor" {
  type = string
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
