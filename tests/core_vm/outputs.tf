output "id" {
  value = module.vm.vm.id
}

output "name" {
  value = module.vm.vm.name
}

output "status" {
  value = module.vm.vm.status
}

output "environment" {
  value = module.vm.vm.environment
}

output "image" {
  value = module.vm.vm.image
}

output "flavor" {
  value = module.vm.vm.flavor
}

output "keypair" {
  value = module.vm.vm.keypair
}

output "volume_attachments" {
  value = module.vm.vm.volume_attachments
}

output "security_rules" {
  value = module.vm.vm.security_rules
}

output "power_state" {
  value = module.vm.vm.power_state
}

output "vm_state" {
  value = module.vm.vm.vm_state
}

output "fixed_ip" {
  value = module.vm.vm.fixed_ip
}

output "floating_ip" {
  value = module.vm.vm.floating_ip
}

output "floating_ip_status" {
  value = module.vm.vm.floating_ip_status
}

output "created_at" {
  value = module.vm.vm.created_at
}
