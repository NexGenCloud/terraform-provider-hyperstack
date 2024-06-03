resource "hyperstack_core_virtual_machine_sg_rule" "ssh_access" {
  virtual_machine_id = hyperstack_core_virtual_machine.this.id

  direction        = "ingress"
  ethertype        = "IPv4"
  port_range_min   = 22
  port_range_max   = 22
  protocol         = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
}

resource "hyperstack_core_virtual_machine_sg_rule" "http_access" {
  virtual_machine_id = hyperstack_core_virtual_machine.this.id

  direction        = "ingress"
  ethertype        = "IPv4"
  port_range_min   = 80
  port_range_max   = 80
  protocol         = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
}

resource "hyperstack_core_virtual_machine_sg_rule" "https_access" {
  virtual_machine_id = hyperstack_core_virtual_machine.this.id

  direction        = "ingress"
  ethertype        = "IPv4"
  port_range_min   = 443
  port_range_max   = 443
  protocol         = "tcp"
  remote_ip_prefix = "0.0.0.0/0"
}
