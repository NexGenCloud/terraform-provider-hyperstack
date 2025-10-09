# Volume Attachment Example

This example demonstrates how to attach volumes to a virtual machine in Hyperstack.

## Overview

This example creates:
- A virtual machine
- Multiple data volumes
- Volume attachments to connect the volumes to the VM

## Prerequisites

- Terraform >= 1.0
- Hyperstack API key
- An existing Hyperstack environment
- An existing SSH keypair in your Hyperstack account

## Usage

1. Set your Hyperstack API key:
```bash
export HYPERSTACK_API_KEY="your-api-key-here"
```

2. Initialize Terraform:
```bash
terraform init
```

3. Review the plan:
```bash
terraform plan
```

4. Apply the configuration:
```bash
terraform apply
```

## Configuration

You can customize the deployment by setting variables:

```bash
terraform apply \
  -var="environment_name=my-env" \
  -var="vm_name=my-vm" \
  -var="keypair_name=my-key" \
  -var="volume_count=3" \
  -var="volume_size=200"
```

## Variables

| Name | Description | Default |
|------|-------------|---------|
| `environment_name` | The name of the Hyperstack environment | `my-environment` |
| `vm_name` | The name of the virtual machine | `example-vm` |
| `keypair_name` | The name of the SSH keypair | `my-keypair` |
| `image_name` | The name of the OS image | `Ubuntu 22.04 LTS` |
| `flavor_name` | The flavor name for the VM | `n1-cpu-small` |
| `volume_size` | Size of each volume in GB | `100` |
| `volume_type` | Type of volume to create | `ssd` |
| `volume_count` | Number of volumes to create and attach | `2` |

## Outputs

- `vm_id` - The ID of the virtual machine
- `vm_name` - The name of the virtual machine
- `volume_ids` - The IDs of the created volumes
- `volume_attachment_id` - The ID of the volume attachment resource
- `volume_attachments` - Details of the volume attachments
- `attached_volumes` - List of volume IDs that are attached

## Important Notes

1. **Volume Attachment**: Volumes must exist before they can be attached to a VM.
2. **Protected Flag**: Setting `protected = true` prevents accidental detachment of volumes.
3. **Cleanup**: When destroying the infrastructure, volumes will be detached but not deleted.
4. **Dependencies**: The `depends_on` ensures volumes are created before attachment.

## Clean Up

To destroy all resources created by this example:

```bash
terraform destroy
```

## Advanced Usage

### Single Volume Attachment

```terraform
resource "hyperstack_core_volume" "single" {
  name             = "single-volume"
  environment_name = "my-environment"
  size             = 100
  volume_type      = "ssd"
}

resource "hyperstack_core_volume_attachment" "single" {
  vm_id      = hyperstack_core_virtual_machine.example.id
  volume_ids = [hyperstack_core_volume.single.id]
  protected  = true
}
```

### Protected Volumes

```terraform
resource "hyperstack_core_volume_attachment" "protected" {
  vm_id      = hyperstack_core_virtual_machine.example.id
  volume_ids = [hyperstack_core_volume.data.id]
  protected  = true  # Prevents accidental detachment
}
```

## Troubleshooting

### Volume Already Attached

If you see an error about a volume already being attached, ensure:
- The volume is not attached to another VM
- You're not trying to attach the same volume multiple times

### Volume Not Found

If you see an error about a volume not being found:
- Verify the volume exists in the same environment
- Check that the volume ID is correct
- Ensure the volume is in an `available` state

### Attachment Timeout

If attachment times out:
- Check the volume status in the Hyperstack console
- Verify the VM is in a running state
- Try again after a few minutes
