skip =  true

terraform {
  source = ".."
}

include "root" {
  path = find_in_parent_folders("root.hcl")
}

dependency "cluster" {
  config_path = "../../cluster"
}

inputs = {
  cluster_name     = dependency.cluster.outputs.clusters["cpu_2"].name
  api_address      = dependency.cluster.outputs.clusters["cpu_2"].api_address
  kube_config_file = dependency.cluster.outputs.clusters["cpu_2"].kube_config_file
}
