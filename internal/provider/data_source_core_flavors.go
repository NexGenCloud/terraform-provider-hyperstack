package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/flavor"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_flavors"
	"io/ioutil"
	"math/big"
)

var _ datasource.DataSource = &DataSourceCoreFlavors{}

func NewDataSourceCoreFlavors() datasource.DataSource {
	return &DataSourceCoreFlavors{}
}

func (d *DataSourceCoreFlavors) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_flavors"
}

func (d *DataSourceCoreFlavors) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_flavors.CoreFlavorsDataSourceSchema(ctx)
}

type DataSourceCoreFlavors struct {
	hyperstack *client.HyperstackClient
	client     *flavor.ClientWithResponses
}

func (d *DataSourceCoreFlavors) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = flavor.NewClientWithResponses(
		d.hyperstack.ApiServer,
		flavor.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreFlavors) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_flavors.CoreFlavorsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Initialize the parameters as nil
	params := (*flavor.RetrieveFlavorsParams)(nil)
	result := (*flavor.RetrieveFlavorsResponse)(nil)
	err := error(nil)

	// If data.Region is not nil or empty, construct the parameters
	if !data.Region.IsNull() && data.Region.String() != "" {
		stringRegion := data.Region.String()

		params = &flavor.RetrieveFlavorsParams{
			Region: &stringRegion,
		}
		result, err = d.client.RetrieveFlavorsWithResponse(ctx, params)
	} else {
		result, err = d.client.RetrieveFlavorsWithResponse(ctx, nil)
	}

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

	callResult := result.JSON200.Data
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	data = d.ApiToModel(ctx, &resp.Diagnostics, callResult)
	// Convert the types.Set to a slice of attr.Value

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceCoreFlavors) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]flavor.FlavorItemGetResponse,
) datasource_core_flavors.CoreFlavorsModel {
	return datasource_core_flavors.CoreFlavorsModel{
		CoreFlavors: d.MapFlavors(ctx, diags, *response),
	}
}

func (d *DataSourceCoreFlavors) MapFlavors(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []flavor.FlavorItemGetResponse,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_core_flavors.CoreFlavorsValue{}.Type(ctx),
		func() []attr.Value {
			coreFlavors := make([]attr.Value, 0)
			for _, row := range data {
				gpu := ""
				if row.Gpu != nil {
					gpu = *row.Gpu
				}
				regionName := ""
				if row.RegionName != nil {
					regionName = *row.RegionName
				}
				flavorsList, _ := types.ListValue(
					datasource_core_flavors.FlavorsValue{}.Type(ctx),
					func() []attr.Value {
						flavors := make([]attr.Value, 0)

						for _, flavorItem := range *row.Flavors {
							cpu := int64(0)
							if flavorItem.Cpu != nil {
								cpu = int64(*flavorItem.Cpu)
							}
							createdAt := ""
							if flavorItem.CreatedAt != nil {
								createdAt = flavorItem.CreatedAt.String()
							}
							disk := int64(0)
							if flavorItem.Disk != nil {
								disk = int64(*flavorItem.Disk)
							}
							gpuCount := int64(0)
							if flavorItem.GpuCount != nil {
								gpuCount = int64(*flavorItem.GpuCount)
							}
							id := int64(0)
							if flavorItem.Id != nil {
								id = int64(*flavorItem.Id)
							}
							name := ""
							if flavorItem.Name != nil {
								name = *flavorItem.Name
							}
							ram := big.NewFloat(0)
							if flavorItem.Ram != nil {
								ram = big.NewFloat(float64(*flavorItem.Ram))
							}
							stockAvailable := false
							if flavorItem.StockAvailable != nil {
								stockAvailable = *flavorItem.StockAvailable
							}

							modelFlavor, diagnostic := datasource_core_flavors.NewFlavorsValue(
								datasource_core_flavors.FlavorsValue{}.AttributeTypes(ctx),
								map[string]attr.Value{
									"cpu":             types.Int64Value(cpu),
									"created_at":      types.StringValue(createdAt),
									"disk":            types.Int64Value(disk),
									"gpu":             types.StringValue(gpu),
									"gpu_count":       types.Int64Value(gpuCount),
									"id":              types.Int64Value(id),
									"name":            types.StringValue(name),
									"ram":             types.NumberValue(ram),
									"region_name":     types.StringValue(regionName),
									"stock_available": types.BoolValue(stockAvailable),
								},
							)
							flavors = append(flavors, modelFlavor)
							diags.Append(diagnostic...)
						}
						return flavors
					}())
				modelCoreFlavor, _ := datasource_core_flavors.NewCoreFlavorsValue(
					datasource_core_flavors.CoreFlavorsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"flavors":     flavorsList,
						"gpu":         types.StringValue(gpu),
						"region_name": types.StringValue("test"),
					})
				coreFlavors = append(coreFlavors, modelCoreFlavor)

			}

			return coreFlavors
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
