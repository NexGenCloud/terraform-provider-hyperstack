terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.31.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.14.0"
    }
  }
}

provider "kubernetes" {
  config_path = var.kube_config_file
  insecure    = true
}

provider "helm" {
  kubernetes {
    config_path = var.kube_config_file
    insecure    = true
  }
}
