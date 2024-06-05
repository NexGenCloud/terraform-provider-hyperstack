terraform {
  required_providers {
    hyperstack = {
      source  = "nexgencloud/hyperstack"
      version = "~> 0.1"
    }
  }

  backend "local" {}
}
