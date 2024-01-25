package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/gosdk/environment"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/datasource_environments"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceEnvironments{}

func NewDataSourceEnvironments() datasource.DataSource {
	return &DataSourceEnvironments{}
}

type DataSourceEnvironments struct {
	hyperstack *client.HyperstackClient
	client     *environment.ClientWithResponses
}

func (d *DataSourceEnvironments) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environments"
}

func (d *DataSourceEnvironments) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_environments.EnvironmentsDataSourceSchema(ctx)
}

func (d *DataSourceEnvironments) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceEnvironments) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_environments.EnvironmentsModel

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

	data.Environments = d.MapEnvironments(ctx, resp, *callResult)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceEnvironments) MapEnvironments(
	ctx context.Context,
	resp *datasource.ReadResponse,
	data []environment.EnvironmentFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_environments.EnvironmentsValue{}.Type(ctx),
		func() []attr.Value {
			envs := make([]attr.Value, 0)
			for _, row := range data {
				createdAt := types.StringNull()
				if row.CreatedAt != nil {
					createdAt = types.StringValue(row.CreatedAt.String())

				}

				model, diagnostic := datasource_environments.NewEnvironmentsValue(
					datasource_environments.EnvironmentsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":         types.Int64Value(int64(*row.Id)),
						"name":       types.StringValue(*row.Name),
						"region":     types.StringValue(*row.Region),
						"created_at": createdAt,
					},
				)
				resp.Diagnostics.Append(diagnostic...)
				envs = append(envs, model)
			}
			return envs
		}(),
	)
	resp.Diagnostics.Append(diagnostic...)
	return model
}
