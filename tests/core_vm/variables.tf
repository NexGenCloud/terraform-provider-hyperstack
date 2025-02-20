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
      count = 1
    }
  }
}
