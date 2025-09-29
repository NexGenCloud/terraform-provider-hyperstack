terraform {
  required_version = "~> 1.7"

  required_providers {
    hyperstack = {
      source  = "nexgencloud/hyperstack"
      version = "~> 1.41"
    }
  }
}
