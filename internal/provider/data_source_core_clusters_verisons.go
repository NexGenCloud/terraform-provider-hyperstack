package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/clusters"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_clusters_versions"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceCoreClustersVersions{}

func NewDataSourceCoreClustersVersions() datasource.DataSource {
	return &DataSourceCoreClustersVersions{}
}

type DataSourceCoreClustersVersions struct {
	hyperstack *client.HyperstackClient
	client     *clusters.ClientWithResponses
}

func (d *DataSourceCoreClustersVersions) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_clusters_versions"
}

func (d *DataSourceCoreClustersVersions) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_clusters_versions.CoreClustersVersionsDataSourceSchema(ctx)
}

func (d *DataSourceCoreClustersVersions) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = clusters.NewClientWithResponses(
		d.hyperstack.ApiServer,
		clusters.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreClustersVersions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_clusters_versions.CoreClustersVersionsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetClusterVersionsWithResponse(ctx)
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

	callResult := result.JSON200.Versions
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

func (d *DataSourceCoreClustersVersions) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]string,
) datasource_core_clusters_versions.CoreClustersVersionsModel {
	return datasource_core_clusters_versions.CoreClustersVersionsModel{
		CoreClustersVersions: func() types.Set {
			return d.MapProtocols(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceCoreClustersVersions) MapProtocols(
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
