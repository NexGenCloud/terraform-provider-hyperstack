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

# New variables for enhanced cluster configuration
variable "deployment_mode" {
  description = "Deployment mode for the cluster"
  type        = string
  default     = "full"
  validation {
    condition     = contains(["full", "standard"], var.deployment_mode)
    error_message = "Deployment mode must be either 'full' or 'standard'."
  }
}

variable "master_count" {
  description = "Number of master nodes"
  type        = number
  default     = 2
  validation {
    condition     = var.master_count >= 2 && var.master_count <= 3
    error_message = "Master count must be between 2 and 3."
  }
}
