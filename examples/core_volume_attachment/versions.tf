terraform {
  required_version = ">= 1.0"

  required_providers {
    hyperstack = {
      source  = "nexgencloud/hyperstack"
      version = ">= 0.2.0"
    }
  }
}

provider "hyperstack" {
  # API key is read from HYPERSTACK_API_KEY environment variable
}
