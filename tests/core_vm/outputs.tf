output "vms" {
  value = {
    for k, v in module.vms : k => {
      id                  = v.vm.id
      name                = v.vm.name
      status              = v.vm.status
      environment         = v.vm.environment
      image               = v.vm.image
      flavor              = v.vm.flavor
      keypair             = v.vm.keypair
      volume_attachments  = v.vm.volume_attachments
      security_rules      = v.vm.security_rules
      power_state         = v.vm.power_state
      vm_state            = v.vm.vm_state
      fixed_ip            = v.vm.fixed_ip
      floating_ip         = v.vm.floating_ip
      floating_ip_status  = v.vm.floating_ip_status
      created_at          = v.vm.created_at
    }
  }
}
