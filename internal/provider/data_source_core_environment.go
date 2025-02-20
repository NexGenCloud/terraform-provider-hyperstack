package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/environment"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_environment"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceCoreEnvironment{}

func NewDataSourceCoreEnvironment() datasource.DataSource {
	return &DataSourceCoreEnvironment{}
}

type DataSourceCoreEnvironment struct {
	hyperstack *client.HyperstackClient
	client     *environment.ClientWithResponses
}

func (d *DataSourceCoreEnvironment) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_environment"
}

func (d *DataSourceCoreEnvironment) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_environment.CoreEnvironmentDataSourceSchema(ctx)
}

func (d *DataSourceCoreEnvironment) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceCoreEnvironment) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_environment.CoreEnvironmentModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.RetrieveEnvironmentWithResponse(ctx, int(data.Id.ValueInt64()))
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

	callResult := result.JSON200.Environment
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

func (d *DataSourceCoreEnvironment) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *environment.EnvironmentFields,
) datasource_core_environment.CoreEnvironmentModel {
	return datasource_core_environment.CoreEnvironmentModel{
		Id: func() types.Int64 {
			return types.Int64Value(int64(*response.Id))
		}(),
		Name: func() types.String {
			if response.Name == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Name)
		}(),
		Region: func() types.String {
			if response.Region == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Region)
		}(),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
		Features: d.MapFeatures(ctx, diags, response.Features),
	}
}

func (d *DataSourceCoreEnvironment) MapFeatures(
	ctx context.Context,
	diags *diag.Diagnostics,
	data *environment.EnvironmentFeatures,
) datasource_core_environment.FeaturesValue {
	model, diagnostic := datasource_core_environment.NewFeaturesValue(
		datasource_core_environment.FeaturesValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"network_optimised": types.BoolValue(*data.NetworkOptimised),
		},
	)
	diags.Append(diagnostic...)

	return model
}
