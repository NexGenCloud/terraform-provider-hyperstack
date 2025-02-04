variable "region" {
  type = string
}

variable "artifacts_directory" {
  type = string
}

variable "name" {
  type = string
}

variable "environment_name" {
  type = string
}

variable "flavor_name" {
  type = string
}

variable "image_name" {
  type = string
}

variable "ingress_ports" {
  type = list(number)
  default = []
}

variable "user_data" {
  type    = string
  default = ""
}

variable "callback_url" {
  type    = string
  default = null
}

variable "enable_port_randomization" {
  type    = bool
  default = null
}

variable "assign_floating_ip" {
  type    = bool
  default = null
}

variable "create_bootable_volume" {
  type    = bool
  default = null
}
