package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/dashboard"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_dashboard"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
	"math/big"
)

var _ datasource.DataSource = &DataSourceCoreDashboard{}

func NewDataSourceCoreDashboard() datasource.DataSource {
	return &DataSourceCoreDashboard{}
}

type DataSourceCoreDashboard struct {
	hyperstack *client.HyperstackClient
	client     *dashboard.ClientWithResponses
}

func (d *DataSourceCoreDashboard) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_dashboard"
}

func (d *DataSourceCoreDashboard) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_dashboard.CoreDashboardDataSourceSchema(ctx)
}

func (d *DataSourceCoreDashboard) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = dashboard.NewClientWithResponses(
		d.hyperstack.ApiServer,
		dashboard.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreDashboard) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_dashboard.CoreDashboardModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetDashboardWithResponse(ctx)
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

	callResult := result.JSON200.Overview
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

func (d *DataSourceCoreDashboard) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *dashboard.OverviewInfo,
) datasource_core_dashboard.CoreDashboardModel {
	return datasource_core_dashboard.CoreDashboardModel{
		Instance: datasource_core_dashboard.InstanceValue{
			CostPerHour: types.NumberValue(big.NewFloat(float64(*response.Instance.CostPerHour))),
			Count:       types.Int64Value(int64(*response.Instance.Count)),
			Gpus:        types.Int64Value(int64(*response.Instance.Gpus)),
			Ram:         types.NumberValue(big.NewFloat(float64(*response.Instance.Ram))),
			Vcpus:       types.Int64Value(int64(*response.Instance.Vcpus)),
		},
		Container: datasource_core_dashboard.ContainerValue{
			CostPerHour: types.NumberValue(big.NewFloat(float64(*response.Container.CostPerHour))),
			Count:       types.Int64Value(int64(*response.Container.Count)),
			Gpus:        types.Int64Value(int64(*response.Container.Gpus)),
			Ram:         types.NumberValue(big.NewFloat(float64(*response.Container.Ram))),
			Vcpus:       types.Int64Value(int64(*response.Container.Vcpus)),
		},
		Volume: datasource_core_dashboard.VolumeValue{
			CostPerHour: types.NumberValue(big.NewFloat(float64(*response.Volume.CostPerHour))),
			Count:       types.Int64Value(int64(*response.Volume.Count)),
			Using:       types.Int64Value(int64(*response.Volume.Using)),
		},
	}
}
