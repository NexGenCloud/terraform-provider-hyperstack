# Hyperstack Terraform Provider (Alpha)

Documentation: https://registry.terraform.io/providers/NexGenCloud/hyperstack/latest/docs

Hyperstack provides a Terraform provider that enables infrastructure-as-code management of your cloud resources. Similar to other major cloud providers, you can use familiar Terraform workflows to provision and manage your infrastructure on Hyperstack, including virtual machines, Kubernetes clusters, and network resources.

This repository contains the source code for the Hyperstack Terraform provider. The provider is currently in **alpha** and is under active development. All features are subject to change and is provided as-is with no guarantees.

We welcome contributions from the community to help improve the provider and add new features. Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file for more information on how to contribute.

## Prerequisites

Before you begin using the Hyperstack Terraform provider, ensure you have completed the following steps:

1. Install Terraform by following the official [HashiCorp documentation](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
2. Create a new API key or copy your existing API key. For additional information on accessing your API key, see [Hyperstack API Keys Documentation](https://infrahub-doc.nexgencloud.com/docs/api-reference/getting-started-api/authentication/).
3. Set up your environment variables with your API key:

```bash
export HYPERSTACK_API_KEY=your-api-key-here
```

## Usage

To use the Hyperstack Terraform provider, add the following code to your Terraform configuration:

```hcl
terraform {
  required_providers {
    hyperstack = {
      source  = "NexGenCloud/hyperstack"
      version = "1.43.0-alpha"
    }
  }
}
```

## Creating your first Hyperstack VM

To create your first Hyperstack VM, add the following code to your Terraform configuration:

```hcl
terraform {
  required_providers {
    hyperstack = {
      source = "app.terraform.io/nexgencloud/hyperstack"
      version = "1.43.0-alpha"
    }
    tls = {
      source = "hashicorp/tls"
      version = "4.0.5"
    }
  }
}
provider "hyperstack" {}
provider "tls" {}

# Create an environment
resource "hyperstack_core_environment" "canada" {
  name   = "terraform-environment"
  region = "CANADA-1"
}

# Generate a keypair for SSH access
resource "tls_private_key" "ed25519" {
  algorithm = "ED25519"
}

# Create a keypair resource
resource "hyperstack_core_keypair" "ed25519" {
  name        = "terraform-keypair"
  environment = hyperstack_core_environment.canada.name
  public_key  = tls_private_key.ed25519.public_key_openssh
}

# Save the private key content
resource "local_sensitive_file" "ssh" {
  filename = "./llm-inference-benchmarking-keypair.pem"
  content  = tls_private_key.ed25519.private_key_openssh
}

# Create a virtual machine
resource "hyperstack_core_virtual_machine" "example-vm" {
  name               = "terraform-example"
  environment_name   = hyperstack_core_environment.canada.name
  key_name           = hyperstack_core_keypair.ed25519.name
  image_name         = "Ubuntu Server 22.04 LTS R535 CUDA 12.2 with Docker"
  flavor_name        = "n1-cpu-small"
  user_data          = ""
  assign_floating_ip = true
}

# Allow port 22 to SSH into this machine
resource "hyperstack_core_virtual_machine_sg_rule" "ssh_access" {
  virtual_machine_id = hyperstack_core_virtual_machine.example-vm.id
  direction        = "ingress"
  ethertype        = "IPv4"
  port_range_min   = 22
  port_range_max   = 22
  protocol         = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
}

# Output
output "ssh_private_key" {
  description = "SSH private key of the e2e test"
  value       = tls_private_key.ed25519.private_key_openssh
  sensitive   = true
}

output "ssh_public_key" {
  description = "SSH public key of the e2e test"
  value       = tls_private_key.ed25519.public_key_openssh
}

output "vm_floating_ip" {
  description = "The floating IP for the basic test VM"
  value       = hyperstack_core_virtual_machine.example-vm.floating_ip
}
```

You can now run `terraform plan` and `terraform apply` to create your first Hyperstack VM. If you do not see the outputs, you may need to run `terraform apply` again.

Warning: this will incur charges on your Hyperstack account.

If you do not have a Hyperstack account, you can sign up for a free account at [Hyperstack](https://nexgencloud.com/).

Please refer to the [Hyperstack Terraform Provider Documentation](https://infrahub-doc.nexgencloud.com/docs/libraries/terraform) for more examples.

## License

This project is licensed under MIT license - see the [LICENSE](LICENSE) file for details.