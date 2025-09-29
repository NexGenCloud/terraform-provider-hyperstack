terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.0"
    }
    hyperstack = {
      source  = "nexgencloud/hyperstack"
      version = "~> 1.41"
    }
  }
}
