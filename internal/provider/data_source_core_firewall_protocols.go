package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/security_rules"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_firewall_protocols"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceCoreFirewallProtocols{}

func NewDataSourceCoreFirewallProtocols() datasource.DataSource {
	return &DataSourceCoreFirewallProtocols{}
}

type DataSourceCoreFirewallProtocols struct {
	hyperstack *client.HyperstackClient
	client     *security_rules.ClientWithResponses
}

func (d *DataSourceCoreFirewallProtocols) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_firewall_protocols"
}

func (d *DataSourceCoreFirewallProtocols) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_firewall_protocols.CoreFirewallProtocolsDataSourceSchema(ctx)
}

func (d *DataSourceCoreFirewallProtocols) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = security_rules.NewClientWithResponses(
		d.hyperstack.ApiServer,
		security_rules.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreFirewallProtocols) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_firewall_protocols.CoreFirewallProtocolsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListFirewallRuleProtocolsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if result.JSON200 == nil {
		bodyBytes, _ := ioutil.ReadAll(result.HTTPResponse.Body)
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	callResult := result.JSON200.Protocols
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No user data",
			"",
		)
		return
	}

	data = d.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceCoreFirewallProtocols) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]string,
) datasource_core_firewall_protocols.CoreFirewallProtocolsModel {
	return datasource_core_firewall_protocols.CoreFirewallProtocolsModel{
		CoreFirewallProtocols: func() types.Set {
			return d.MapProtocols(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceCoreFirewallProtocols) MapProtocols(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []string,
) types.Set {
	model, diagnostic := types.SetValue(
		types.StringType,
		func() []attr.Value {
			protocols := make([]attr.Value, 0)
			for _, row := range data {
				model := types.StringValue(row)
				protocols = append(protocols, model)
			}
			return protocols
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
