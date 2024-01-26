terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.0"
    }
    hyperstack = {
      source  = "nexgen/hyperstack"
      version = "~> 0.0.1"
    }
  }
}

provider "hyperstack" {
  staging = true
}
