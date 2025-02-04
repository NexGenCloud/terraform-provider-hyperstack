package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/stock"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_stocks"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceCoreStocks{}

func NewDataSourceCoreStocks() datasource.DataSource {
	return &DataSourceCoreStocks{}
}

type DataSourceCoreStocks struct {
	hyperstack *client.HyperstackClient
	client     *stock.ClientWithResponses
}

func (d *DataSourceCoreStocks) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_stocks"
}

func (d *DataSourceCoreStocks) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_stocks.CoreStocksDataSourceSchema(ctx)
}

func (d *DataSourceCoreStocks) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = stock.NewClientWithResponses(
		d.hyperstack.ApiServer,
		stock.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreStocks) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_stocks.CoreStocksModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.RetrieveGpuStocksWithResponse(ctx)
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

	callResult := result.JSON200.Stocks
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

func (d *DataSourceCoreStocks) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]stock.NewStockResponse,
) datasource_core_stocks.CoreStocksModel {
	return datasource_core_stocks.CoreStocksModel{
		Stocks: func() types.List {
			return d.MapStocks(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceCoreStocks) MapStocks(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []stock.NewStockResponse,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_core_stocks.StocksValue{}.Type(ctx),
		func() []attr.Value {
			stocks := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_core_stocks.NewStocksValue(
					datasource_core_stocks.StocksValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"region":    types.StringPointerValue(row.Region),
						"stocktype": types.StringPointerValue(row.StockType),
						"models":    d.MapStockModels(ctx, diags, *row.Models),
					},
				)
				diags.Append(diagnostic...)
				stocks = append(stocks, model)
			}
			return stocks
		}(),
	)
	diags.Append(diagnostic...)
	return model
}

func (d *DataSourceCoreStocks) MapStockModels(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []stock.NewModelResponse,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_core_stocks.ModelsValue{}.Type(ctx),
		func() []attr.Value {
			models := make([]attr.Value, 0)
			for _, row := range data {
				// TODO: simplify
				conf, diagnostic := types.ObjectValue(
					map[string]attr.Type{
						"n10x": basetypes.Int64Type{},
						"n1x":  basetypes.Int64Type{},
						"n2x":  basetypes.Int64Type{},
						"n4x":  basetypes.Int64Type{},
						"n8x":  basetypes.Int64Type{},
					},
					map[string]attr.Value{
						"n10x": types.Int64Value(int64(*row.Configurations.N10x)),
						"n1x":  types.Int64Value(int64(*row.Configurations.N1x)),
						"n2x":  types.Int64Value(int64(*row.Configurations.N2x)),
						"n4x":  types.Int64Value(int64(*row.Configurations.N4x)),
						"n8x":  types.Int64Value(int64(*row.Configurations.N8x)),
					},
				)
				diags.Append(diagnostic...)
				model, diagnostic := datasource_core_stocks.NewModelsValue(
					datasource_core_stocks.ModelsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"available":        types.StringPointerValue(row.Available),
						"configurations":   conf,
						"model":            types.StringPointerValue(row.Model),
						"planned_100_days": types.StringPointerValue(row.Planned100Days),
						"planned_30_days":  types.StringPointerValue(row.Planned30Days),
						"planned_7_days":   types.StringPointerValue(row.Planned7Days),
					},
				)
				diags.Append(diagnostic...)
				models = append(models, model)
			}
			return models
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
