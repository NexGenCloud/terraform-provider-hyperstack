skip = false

include "root" {
  path = find_in_parent_folders("root.hcl")
}

dependency "cluster" {
  config_path = "../cluster"
}

inputs = {
  cluster_name     = dependency.cluster.outputs.name
  api_address      = dependency.cluster.outputs.api_address
  kube_config_file = dependency.cluster.outputs.kube_config_file
}
