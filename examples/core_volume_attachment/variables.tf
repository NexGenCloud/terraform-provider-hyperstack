variable "environment_name" {
  description = "The name of the Hyperstack environment"
  type        = string
  default     = "my-environment"
}

variable "vm_name" {
  description = "The name of the virtual machine"
  type        = string
  default     = "example-vm"
}

variable "keypair_name" {
  description = "The name of the SSH keypair"
  type        = string
  default     = "my-keypair"
}

variable "image_name" {
  description = "The name of the OS image"
  type        = string
  default     = "Ubuntu 22.04 LTS"
}

variable "flavor_name" {
  description = "The flavor name for the VM"
  type        = string
  default     = "n1-cpu-small"
}

variable "volume_size" {
  description = "Size of each volume in GB"
  type        = number
  default     = 100
}

variable "volume_type" {
  description = "Type of volume to create"
  type        = string
  default     = "ssd"
}

variable "volume_count" {
  description = "Number of volumes to create and attach"
  type        = number
  default     = 2
}
