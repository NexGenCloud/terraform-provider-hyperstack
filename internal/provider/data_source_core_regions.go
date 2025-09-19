package provider

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/NexGenCloud/hyperstack-sdk-go/lib/region"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_regions"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DataSourceCoreRegions{}

func NewDataSourceCoreRegions() datasource.DataSource {
	return &DataSourceCoreRegions{}
}

func (d *DataSourceCoreRegions) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_regions"
}

func (d *DataSourceCoreRegions) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_regions.CoreRegionsDataSourceSchema(ctx)
}

type DataSourceCoreRegions struct {
	hyperstack *client.HyperstackClient
	client     *region.ClientWithResponses
}

func (d *DataSourceCoreRegions) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = region.NewClientWithResponses(
		d.hyperstack.ApiServer,
		region.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreRegions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_regions.CoreRegionsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListRegionsWithResponse(ctx)
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

	callResult := result.JSON200.Regions
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

func (d *DataSourceCoreRegions) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]region.RegionFields,
) datasource_core_regions.CoreRegionsModel {
	return datasource_core_regions.CoreRegionsModel{
		CoreRegions: func() types.Set {
			return d.MapRegions(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceCoreRegions) MapRegions(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []region.RegionFields,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_core_regions.CoreRegionsValue{}.Type(ctx),
		func() []attr.Value {
			regions := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_core_regions.NewCoreRegionsValue(
					datasource_core_regions.CoreRegionsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id": func() attr.Value {
							if row.Id == nil {
								return types.Int64Null()
							}
							return types.Int64Value(int64(*row.Id))
						}(),
						"name": func() attr.Value {
							if row.Name == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.Name)
						}(),
						"description": func() attr.Value {
							if row.Description == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.Description)
						}(),
						"country": func() attr.Value {
							if row.Country == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.Country)
						}(),
						"green_status": func() attr.Value {
							if row.GreenStatus == nil {
								return types.StringNull()
							}
							return types.StringValue(string(*row.GreenStatus))
						}(),
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
