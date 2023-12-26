terraform {
  required_providers {
    hyperstack = {
      source  = "nexgen/hyperstack"
      version = "~> 0.0.1"
    }
  }
}

provider "hyperstack" {}