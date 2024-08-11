terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0.5"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5.1"
    }
    hyperstack = {
      source  = "nexgencloud/hyperstack"
      version = "~> 0.1"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.14.0"
    }
    null = {
      source = "hashicorp/null"
      version = "3.2.2"
    }  
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.31.0"
    }        
  }
}