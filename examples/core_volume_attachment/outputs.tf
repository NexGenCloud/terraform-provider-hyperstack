output "vm_id" {
  description = "The ID of the virtual machine"
  value       = hyperstack_core_virtual_machine.example.id
}

output "vm_name" {
  description = "The name of the virtual machine"
  value       = hyperstack_core_virtual_machine.example.name
}

output "volume_ids" {
  description = "The IDs of the created volumes"
  value       = [for vol in hyperstack_core_volume.data : vol.id]
}

output "volume_attachment_id" {
  description = "The ID of the volume attachment resource"
  value       = hyperstack_core_volume_attachment.example.id
}

output "volume_attachments" {
  description = "Details of the volume attachments"
  value       = hyperstack_core_volume_attachment.example.volume_attachments
}

output "attached_volumes" {
  description = "List of volume IDs that are attached"
  value       = hyperstack_core_volume_attachment.example.volume_ids
}
