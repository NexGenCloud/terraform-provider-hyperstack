variable "region" {
  type = string
}

variable "artifacts_dir" {
  type = string
}

variable "name_prefix" {
  type = string
}

#noinspection TFIncorrectVariableType
variable "clusters" {
  type = map(object({
    enabled = optional(bool, true)

    master_flavor = optional(object({
      name = optional(string)
      gpu_name = optional(string)
      gpu_count = optional(number)
      cpu_count = optional(number)
    }), {
      name = "n1-cpu-medium"
    })

    node_count = optional(number, 1)
    node_flavor = object({
      name = optional(string)
      gpu_name = optional(string)
      gpu_count = optional(number)
      cpu_count = optional(number)
    })

    validation = object({
      gpu = optional(bool, false)
    })

    image_type = optional(string, "Ubuntu")
    image_version = optional(string, "Server 22.04 LTS R535 CUDA 12.2")
  }))
  default = {
    "cpu-2" = {
      enabled    = true
      validation = {}
      node_count = 1
      node_flavor = {
        gpu_name  = ""
        cpu_count = 4
      }
    }
    "a6000-2" = {
      enabled    = true
      validation = {
        gpu = true
      }
      node_count = 2
      node_flavor = {
        gpu_name  = "RTX-A6000"
        gpu_count = 1
      }
      image_version = "Server 22.04 LTS R535 CUDA 12.2"
    }
  }
}
