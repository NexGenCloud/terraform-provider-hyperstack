resource "hyperstack_core_virtual_machine_sg_rule" "ingress_ports" {
  for_each = toset([for p in var.ingress_ports : tostring(p)])

  virtual_machine_id = hyperstack_core_virtual_machine.this.id

  direction        = "ingress"
  ethertype        = "IPv4"
  port_range_min   = tonumber(each.value)
  port_range_max   = tonumber(each.value)
  protocol         = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
}
