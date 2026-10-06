terraform {
  # Sourced from the repo root so Terragrunt copies examples/ into the
  # working directory too. Terragrunt v1 always runs terraform from a cache
  # copy of the unit, so the ../../examples module paths only resolve when
  # the whole repo comes along.
  source = "${get_repo_root()}//tests/core_environment"
}

include "root" {
  path = find_in_parent_folders("root.hcl")
}
