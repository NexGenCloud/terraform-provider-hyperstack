terraform {
  required_providers {
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0.5"
    }
    hyperstack = {
      source  = "nexgencloud/hyperstack"
      version = "~> 1.41"
    }
  }
}
