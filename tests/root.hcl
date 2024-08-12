skip = true

locals {
  artifacts_dir = "${get_env("ARTIFACTS_DIR")}/${path_relative_to_include("root")}"
  log_suffix = get_env("LOG_SUFFIX", "default")
  log_level     = "DEBUG"
}

inputs = {
  region = get_env("HYPERSTACK_REGION")
  artifacts_dir = local.artifacts_dir
  name_prefix   = "tf-"
}

terraform {
  extra_arguments "tf_config" {
    commands = [
      "apply",
      "destroy",
      "import",
      "init",
      "plan",
      "refresh",
      "taint",
      "untaint",
      "output",
    ]

    env_vars = {
      TF_CLI_CONFIG_FILE = "provider-mirror.tfrc"
      TF_DATA_DIR        = "${local.artifacts_dir}/.terraform"
      TF_LOG_PATH        = "${local.artifacts_dir}/terraform_${local.log_suffix}.log"
      TF_IN_AUTOMATION   = "1"
      TF_LOG             = local.log_level
    }

    arguments = [
      #"-compact-warnings",
    ]
  }

  extra_arguments "init_upgrade" {
    commands = [
      "init",
    ]

    arguments = [
      "-upgrade",
      "-migrate-state",
    ]
  }

  extra_arguments "apply_parallel" {
    commands = [
      "apply",
      "destroy",
    ]

    arguments = [
      "-parallelism=50",
    ]
  }

  before_hook "rm_lock" {
    commands = ["init"]
    execute = ["rm", "-f", ".terraform.lock.hcl"]
  }

  before_hook "create_dir" {
    commands = ["init"]
    execute = ["mkdir", "-p", local.artifacts_dir]
  }

  after_hook "outputs" {
    commands     = ["apply", "plan"]
    execute      = ["bash", "-c", "terraform output -json > '${local.artifacts_dir}/outputs.json'"]
  }
}

remote_state {
  backend = "local"
  config = {
    path = "${local.artifacts_dir}/terraform.tfstate"
  }

  disable_dependency_optimization = true

  generate = {
    path      = "backend.tf"
    if_exists = "overwrite"
  }
}

generate "local-provider" {
  path      = "provider-mirror.tfrc"
  if_exists = "overwrite"
  contents  = <<EOF
provider_installation {
  filesystem_mirror {
    # Path from subfolder
    path    = "${path_relative_from_include()}/../artifacts/provider-mirror"
    #path    = "${get_repo_root()}/artifacts/provider-mirror"
    include = ["nexgencloud/*"]
  }

  direct {
    exclude = ["nexgencloud/*"]
  }
}
EOF
}
