skip = false

include "root" {
  path = find_in_parent_folders("root.hcl")
}

dependency "cluster" {
  config_path = "../cluster"
}

inputs = {
}

generate "modules" {
  path      = "modules.gen.tf"
  if_exists = "overwrite"
  contents = templatefile("data/module.txt", {
    name     = dependency.cluster.outputs.name
    clusters = dependency.cluster.outputs.clusters
  })
}
