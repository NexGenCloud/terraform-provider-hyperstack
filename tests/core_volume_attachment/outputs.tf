output "vms" {
  value = {
    for k, v in module.vms : k => {
      id                 = v.vm.id
      name               = v.vm.name
      status             = v.vm.status
      environment        = v.vm.environment
      image              = v.vm.image
      flavor             = v.vm.flavor
      keypair            = v.vm.keypair
      volume_attachments = v.vm.volume_attachments
      security_rules     = v.vm.security_rules
      power_state        = v.vm.power_state
      vm_state           = v.vm.vm_state
      fixed_ip           = v.vm.fixed_ip
      floating_ip        = v.vm.floating_ip
      floating_ip_status = v.vm.floating_ip_status
      created_at         = v.vm.created_at
    }
  }
}

output "volumes" {
  value = {
    for i, vol in hyperstack_core_volume.data : i => {
      id          = vol.id
      name        = vol.name
      size        = vol.size
      volume_type = vol.volume_type
      status      = vol.status
      bootable    = vol.bootable
    }
  }
}

output "volume_attachments" {
  value = {
    for k, v in hyperstack_core_volume_attachment.data : k => {
      id                 = v.id
      vm_id              = v.vm_id
      volume_ids         = v.volume_ids
      protected          = v.protected
      volume_attachments = v.volume_attachments
    }
  }
}
