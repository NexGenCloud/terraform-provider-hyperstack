# Only org owner can get roles
skip = true

include "root" {
  path = find_in_parent_folders("root.hcl")
}
