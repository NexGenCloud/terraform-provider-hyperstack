// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_core_virtual_machine_sg_rule

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CoreVirtualMachineSgRuleResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"direction": schema.StringAttribute{
				Required:            true,
				Description:         "The direction of traffic that the firewall rule applies to.",
				MarkdownDescription: "The direction of traffic that the firewall rule applies to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ethertype": schema.StringAttribute{
				Required:            true,
				Description:         "The Ethernet type associated with the rule.",
				MarkdownDescription: "The Ethernet type associated with the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"port_range_max": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"port_range_min": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:            true,
				Description:         "The network protocol associated with the rule. Call the [`GET /core/sg-rules-protocols`](https://infrahub-api-doc.nexgencloud.com/#get-/core/sg-rules-protocols) endpoint to retrieve a list of permitted network protocols.",
				MarkdownDescription: "The network protocol associated with the rule. Call the [`GET /core/sg-rules-protocols`](https://infrahub-api-doc.nexgencloud.com/#get-/core/sg-rules-protocols) endpoint to retrieve a list of permitted network protocols.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"any",
						"ah",
						"dccp",
						"egp",
						"esp",
						"gre",
						"hopopt",
						"icmp",
						"igmp",
						"ip",
						"ipip",
						"ipv6-encap",
						"ipv6-frag",
						"ipv6-icmp",
						"icmpv6",
						"ipv6-nonxt",
						"ipv6-opts",
						"ipv6-route",
						"ospf",
						"pgm",
						"rsvp",
						"sctp",
						"tcp",
						"udp",
						"udplite",
						"vrrp",
					),
				},
			},
			"remote_ip_prefix": schema.StringAttribute{
				Required:            true,
				Description:         "The IP address range that is allowed to access the specified port. Use \"0.0.0.0/0\" to allow any IP address.",
				MarkdownDescription: "The IP address range that is allowed to access the specified port. Use \"0.0.0.0/0\" to allow any IP address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"virtual_machine_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

type CoreVirtualMachineSgRuleModel struct {
	CreatedAt        types.String `tfsdk:"created_at"`
	Direction        types.String `tfsdk:"direction"`
	Ethertype        types.String `tfsdk:"ethertype"`
	Id               types.Int64  `tfsdk:"id"`
	PortRangeMax     types.Int64  `tfsdk:"port_range_max"`
	PortRangeMin     types.Int64  `tfsdk:"port_range_min"`
	Protocol         types.String `tfsdk:"protocol"`
	RemoteIpPrefix   types.String `tfsdk:"remote_ip_prefix"`
	Status           types.String `tfsdk:"status"`
	VirtualMachineId types.Int64  `tfsdk:"virtual_machine_id"`
}
