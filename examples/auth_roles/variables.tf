variable "name_prefix" {
  type    = string
  default = "test-"
}

variable "role_source" {
  type    = string
  default = "admins"
}

variable "role_description" {
  type    = string
  default = "Test role"
}

variable "role_policies" {
  type    = list(string)
  default = [
    "policy:ReadPermissions",
    "policy:VirtualMachinePermissions",
  ]
}

variable "role_permissions" {
  type    = list(string)
  default = [
    "environment:update",
  ]
}
