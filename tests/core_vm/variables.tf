variable "hyperstack_region" {
  type    = string
}

variable "artifacts_directory" {
  type = string
}

variable "name_prefix" {
  type    = string
}

variable "instance_gpu" {
  type    = string
  default = ""
}

variable "instance_cpus" {
  type    = number
  default = 4
}

variable "image_type" {
  type    = string
  default = "Ubuntu"
}

variable "image_version" {
  type    = string
  default = "Server 20.04 LTS"
}
