terraform {
  required_version = "~> 1.7"

  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0.5"
    }
    template = {
      source  = "hashicorp/template"
      version = "~> 2.2.0"
    }
    hyperstack = {
      source  = "nexgencloud/hyperstack"
      version = "~> 0.1"
    }
  }

  backend "local" {}
}

provider "hyperstack" {
  staging = true
}
