// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_core_firewall_protocols

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func CoreFirewallProtocolsDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"core_firewall_protocols": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

type CoreFirewallProtocolsModel struct {
	CoreFirewallProtocols types.Set `tfsdk:"core_firewall_protocols"`
}
