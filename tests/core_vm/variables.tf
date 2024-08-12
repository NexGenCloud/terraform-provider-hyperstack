variable "region" {
  type    = string
}

variable "artifacts_dir" {
  type = string
}

variable "name_prefix" {
  type    = string
}

#noinspection TFIncorrectVariableType
variable "vms" {
  type = map(object({
    enabled = optional(bool, true)

    flavor = object({
      name = optional(string)
      gpu_name = optional(string)
      gpu_count = optional(number)
      cpu_count = optional(number)
    })

    image_type = optional(string, "Ubuntu")
    image_version = optional(string, "Server 20.04 LTS")

    count = optional(number, 1)
  }))

  default = {
    "cpu4-4" = {
      enabled    = true
      flavor = {
        gpu_name  = ""
        cpu_count = 4
      }
      count = 4
    }
    "cpu8-4" = {
      enabled    = false
      flavor = {
        gpu_name  = ""
        cpu_count = 8
      }
      count = 4
    }
    "a6000-2" = {
      enabled    = false
      flavor = {
        gpu_name  = "RTX-A6000"
        gpu_count = 1
      }
      count = 2
      image_version = "Server 22.04 LTS R535 CUDA 12.2"
    }
    "l40-2" = {
      enabled    = false
      flavor = {
        gpu_name  = "L40"
        gpu_count = 1
      }
      count = 2
      image_version = "Server 22.04 LTS R535 CUDA 12.2"
    }
    "a100x1-2" = {
      enabled    = false
      flavor = {
        gpu_name  = "A100-80G-PCIe"
        gpu_count = 1
      }
      count = 2
      image_version = "Server 22.04 LTS R535 CUDA 12.2"
    }
    "h100x1-2" = {
      enabled    = false
      flavor = {
        gpu_name  = "H100-80G-PCIe"
        gpu_count = 1
      }
      count = 2
      image_version = "Server 22.04 LTS R535 CUDA 12.2"
    }
  }
}
