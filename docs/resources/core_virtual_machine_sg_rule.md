---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "hyperstack_core_virtual_machine_sg_rule Resource - terraform-provider-hyperstack"
subcategory: ""
description: |-
  
---

# hyperstack_core_virtual_machine_sg_rule (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `direction` (String) The direction of traffic that the firewall rule applies to.
- `ethertype` (String) The Ethernet type associated with the rule.
- `protocol` (String) The network protocol associated with the rule. Call the [`GET /core/sg-rules-protocols`](https://infrahub-api-doc.nexgencloud.com/#get-/core/sg-rules-protocols) endpoint to retrieve a list of permitted network protocols.
- `remote_ip_prefix` (String) The IP address range that is allowed to access the specified port. Use "0.0.0.0/0" to allow any IP address.
- `virtual_machine_id` (Number)

### Optional

- `port_range_max` (Number)
- `port_range_min` (Number)

### Read-Only

- `created_at` (String)
- `id` (Number) The ID of this resource.
- `status` (String)
