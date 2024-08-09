variable "region" {
  type = string
}

variable "name" {
  type    = string
  default = null
}

variable "gpu_name" {
  type    = string
  default = null
}

variable "gpu_count" {
  type    = number
  default = null
}

variable "cpu_count" {
  type    = number
  default = null
}
