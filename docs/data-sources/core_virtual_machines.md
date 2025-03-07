---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "hyperstack_core_virtual_machines Data Source - terraform-provider-hyperstack"
subcategory: ""
description: |-
  
---

# hyperstack_core_virtual_machines (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `page` (String) Page Number
- `page_size` (String) Data Per Page
- `search` (String)

### Read-Only

- `core_virtual_machines` (Attributes Set) (see [below for nested schema](#nestedatt--core_virtual_machines))

<a id="nestedatt--core_virtual_machines"></a>
### Nested Schema for `core_virtual_machines`

Read-Only:

- `contract_id` (Number)
- `created_at` (String)
- `environment` (Attributes) (see [below for nested schema](#nestedatt--core_virtual_machines--environment))
- `fixed_ip` (String)
- `flavor` (Attributes) (see [below for nested schema](#nestedatt--core_virtual_machines--flavor))
- `floating_ip` (String)
- `floating_ip_status` (String)
- `id` (Number)
- `image` (Attributes) (see [below for nested schema](#nestedatt--core_virtual_machines--image))
- `keypair` (Attributes) (see [below for nested schema](#nestedatt--core_virtual_machines--keypair))
- `labels` (List of String)
- `locked` (Boolean)
- `name` (String)
- `os` (String)
- `power_state` (String)
- `security_rules` (Attributes List) (see [below for nested schema](#nestedatt--core_virtual_machines--security_rules))
- `status` (String)
- `vm_state` (String)
- `volume_attachments` (Attributes List) (see [below for nested schema](#nestedatt--core_virtual_machines--volume_attachments))

<a id="nestedatt--core_virtual_machines--environment"></a>
### Nested Schema for `core_virtual_machines.environment`

Read-Only:

- `id` (Number)
- `name` (String)
- `org_id` (Number)
- `region` (String)


<a id="nestedatt--core_virtual_machines--flavor"></a>
### Nested Schema for `core_virtual_machines.flavor`

Read-Only:

- `cpu` (Number)
- `disk` (Number)
- `ephemeral` (Number)
- `gpu` (String)
- `gpu_count` (Number)
- `id` (Number)
- `name` (String)
- `ram` (Number)


<a id="nestedatt--core_virtual_machines--image"></a>
### Nested Schema for `core_virtual_machines.image`

Read-Only:

- `name` (String)


<a id="nestedatt--core_virtual_machines--keypair"></a>
### Nested Schema for `core_virtual_machines.keypair`

Read-Only:

- `name` (String)


<a id="nestedatt--core_virtual_machines--security_rules"></a>
### Nested Schema for `core_virtual_machines.security_rules`

Read-Only:

- `created_at` (String)
- `direction` (String)
- `ethertype` (String)
- `id` (Number)
- `port_range_max` (Number)
- `port_range_min` (Number)
- `protocol` (String)
- `remote_ip_prefix` (String)
- `status` (String)


<a id="nestedatt--core_virtual_machines--volume_attachments"></a>
### Nested Schema for `core_virtual_machines.volume_attachments`

Read-Only:

- `created_at` (String)
- `device` (String)
- `status` (String)
- `volume` (Attributes) (see [below for nested schema](#nestedatt--core_virtual_machines--volume_attachments--volume))

<a id="nestedatt--core_virtual_machines--volume_attachments--volume"></a>
### Nested Schema for `core_virtual_machines.volume_attachments.volume`

Read-Only:

- `description` (String)
- `id` (Number)
- `name` (String)
- `size` (Number)
- `volume_type` (String)
