variable "region" {
  type    = string
  default = "staging-CA-1"
}

variable "artifacts_directory" {
  type = string
}

variable "name_prefix" {
  type    = string
  default = "test-"
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

variable "volume_size" {
  type    = number
  default = 50
}
