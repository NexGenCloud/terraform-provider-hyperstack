terraform {
  # Sourced from the repo root so Terragrunt copies examples/ into the
  # working directory too. Terragrunt v1 always runs terraform from a cache
  # copy of the unit, so the ../../examples module paths only resolve when
  # the whole repo comes along.
  source = "${get_repo_root()}//tests/core_cluster/validation"
}

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
