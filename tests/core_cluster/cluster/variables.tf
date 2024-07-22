variable "region" {
  type = string
}

variable "artifacts_dir" {
  type = string
}

variable "name_prefix" {
  type = string
}

variable "node_count" {
  type    = number
  default = 2
}

variable "enable_public_ip" {
  type    = bool
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
