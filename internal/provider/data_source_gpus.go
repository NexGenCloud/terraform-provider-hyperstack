package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/gpu"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_gpus"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_regions"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceCoreGpus{}

func NewDataSourceCoreGpus() datasource.DataSource {
	return &DataSourceCoreGpus{}
}

func (d *DataSourceCoreGpus) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_gpus"
}

func (d *DataSourceCoreGpus) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_gpus.CoreGpusDataSourceSchema(ctx)
}

type DataSourceCoreGpus struct {
	hyperstack *client.HyperstackClient
	client     *gpu.ClientWithResponses
}

func (d *DataSourceCoreGpus) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = gpu.NewClientWithResponses(
		d.hyperstack.ApiServer,
		gpu.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreGpus) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_gpus.CoreGpusModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetGPUListWithResponse(ctx)
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

	callResult := result.JSON200.GpuList
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

func (d *DataSourceCoreGpus) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]gpu.GPUFields,
) datasource_core_gpus.CoreGpusModel {
	return datasource_core_gpus.CoreGpusModel{
		CoreGpus: func() types.Set {
			return d.MapGpus(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceCoreGpus) MapRegion1(
	ctx context.Context,
	diags *diag.Diagnostics,
	data gpu.GPURegionFields, // Assuming this is the type of each region
) datasource_core_gpus.RegionsValue {
	model, diagnostic := datasource_core_gpus.NewRegionsValue(
		datasource_core_gpus.RegionsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":   types.Int64Value(int64(*data.Id)), // Replace with actual fields
			"name": types.StringValue(*data.Name),     // Replace with actual fields
			// Add other fields as necessary
			// "field": types.SomeTypeValue(*data.Field),
		},
	)
	diags.Append(diagnostic...)
	return model
}

func (d *DataSourceCoreGpus) MapRegions(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []gpu.GPURegionFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_core_gpus.RegionsValue{}.Type(ctx),
		func() []attr.Value {
			regions := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_core_regions.NewCoreRegionsValue(
					datasource_core_gpus.RegionsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":   types.Int64Value(int64(*row.Id)),
						"name": types.StringValue(*row.Name),
					},
				)
				diags.Append(diagnostic...)
				regions = append(regions, model)
			}
			return regions
		}(),
	)
	diags.Append(diagnostic...)
	return model
}

func (d *DataSourceCoreGpus) MapGpus(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []gpu.GPUFields,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_core_gpus.CoreGpusValue{}.Type(ctx),
		func() []attr.Value {
			gpus := make([]attr.Value, 0)
			for _, row := range data {
				regions := d.MapRegions(ctx, diags, *row.Regions)
				model, diagnostic := datasource_core_gpus.NewCoreGpusValue(
					datasource_core_gpus.CoreGpusValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"created_at": func() attr.Value {
							if row.CreatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.CreatedAt.String())
						}(),
						"example_metadata": types.StringValue(*row.ExampleMetadata),
						"id":               types.Int64Value(int64(*row.Id)),
						"name":             types.StringValue(*row.Name),
						"regions":          regions,
						"updated_at": func() attr.Value {
							if row.UpdatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.UpdatedAt.String())
						}(),
					},
				)
				diags.Append(diagnostic...)
				gpus = append(gpus, model)
			}
			return gpus
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
