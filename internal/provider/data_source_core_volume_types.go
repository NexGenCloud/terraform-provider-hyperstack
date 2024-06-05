package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/volume"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_volume_types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceCoreVolumeTypes{}

func NewDataSourceCoreVolumeTypes() datasource.DataSource {
	return &DataSourceCoreVolumeTypes{}
}

type DataSourceCoreVolumeTypes struct {
	hyperstack *client.HyperstackClient
	client     *volume.ClientWithResponses
}

func (d *DataSourceCoreVolumeTypes) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_volume_types"
}

func (d *DataSourceCoreVolumeTypes) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_volume_types.CoreVolumeTypesDataSourceSchema(ctx)
}

func (d *DataSourceCoreVolumeTypes) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = volume.NewClientWithResponses(
		d.hyperstack.ApiServer,
		volume.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreVolumeTypes) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_volume_types.CoreVolumeTypesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListVolumeTypesWithResponse(ctx)
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

	callResult := result.JSON200.VolumeTypes
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

func (d *DataSourceCoreVolumeTypes) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]string,
) datasource_core_volume_types.CoreVolumeTypesModel {
	return datasource_core_volume_types.CoreVolumeTypesModel{
		CoreVolumeTypes: func() types.Set {
			return d.MapVolumeTypes(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceCoreVolumeTypes) MapVolumeTypes(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []string,
) types.Set {
	model, diagnostic := types.SetValue(
		types.StringType,
		func() []attr.Value {
			volume_types := make([]attr.Value, 0)
			for _, row := range data {
				model := types.StringValue(row)
				volume_types = append(volume_types, model)
			}
			return volume_types
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
