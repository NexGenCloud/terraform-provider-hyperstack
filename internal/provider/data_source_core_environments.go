package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/environment"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/datasource_core_environments"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceCoreEnvironments{}

func NewDataSourceCoreEnvironments() datasource.DataSource {
	return &DataSourceCoreEnvironments{}
}

type DataSourceCoreEnvironments struct {
	hyperstack *client.HyperstackClient
	client     *environment.ClientWithResponses
}

func (d *DataSourceCoreEnvironments) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_environments"
}

func (d *DataSourceCoreEnvironments) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_environments.CoreEnvironmentsDataSourceSchema(ctx)
}

func (d *DataSourceCoreEnvironments) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = environment.NewClientWithResponses(
		d.hyperstack.ApiServer,
		environment.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreEnvironments) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_environments.CoreEnvironmentsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListEnvironmentsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(result.HTTPResponse.Body)
	if result.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	callResult := result.JSON200.Environments
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	data = d.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceCoreEnvironments) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]environment.EnvironmentFields,
) datasource_core_environments.CoreEnvironmentsModel {
	return datasource_core_environments.CoreEnvironmentsModel{
		CoreEnvironments: func() types.Set {
			return d.MapEnvironments(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceCoreEnvironments) MapEnvironments(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []environment.EnvironmentFields,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_core_environments.CoreEnvironmentsValue{}.Type(ctx),
		func() []attr.Value {
			envs := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_core_environments.NewCoreEnvironmentsValue(
					datasource_core_environments.CoreEnvironmentsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":     types.Int64Value(int64(*row.Id)),
						"name":   types.StringValue(*row.Name),
						"region": types.StringValue(*row.Region),
						"created_at": func() attr.Value {
							if row.CreatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.CreatedAt.String())
						}(),
					},
				)
				diags.Append(diagnostic...)
				envs = append(envs, model)
			}
			return envs
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
