package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/keypair"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_keypair"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceCoreKeypair{}

func NewDataSourceCoreKeypair() datasource.DataSource {
	return &DataSourceCoreKeypair{}
}

func (d *DataSourceCoreKeypair) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_keypair"
}

func (d *DataSourceCoreKeypair) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_keypair.CoreKeypairDataSourceSchema(ctx)
}

type DataSourceCoreKeypair struct {
	hyperstack *client.HyperstackClient
	client     *keypair.ClientWithResponses
}

func (d *DataSourceCoreKeypair) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = keypair.NewClientWithResponses(
		d.hyperstack.ApiServer,
		keypair.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreKeypair) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_keypair.CoreKeypairModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListKeyPairsWithResponse(ctx, func() *keypair.ListKeyPairsParams {
		return &keypair.ListKeyPairsParams{
			Page:     nil,
			PageSize: nil,
			Search:   nil,
		}
	}())
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

	searchCallResult := result.JSON200.Keypairs
	if searchCallResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	// Find the keypair with the matching name and get its ID
	var keypair *keypair.KeypairFields = nil
	for _, row := range *searchCallResult {
		if int64(*row.Id) == data.Id.ValueInt64() {
			keypair = &row
			break
		}
	}

	// Check if id was found
	if keypair == nil {
		resp.Diagnostics.AddError("No keypair found with the id for update: %s", fmt.Sprintf("%d", data.Id.ValueInt64()))
		return
	}

	data = d.ApiToModel(ctx, &resp.Diagnostics, keypair)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceCoreKeypair) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *keypair.KeypairFields,
) datasource_core_keypair.CoreKeypairModel {
	return datasource_core_keypair.CoreKeypairModel{
		Id:          types.Int64Value(int64(*response.Id)),
		Name:        types.StringValue(*response.Name),
		PublicKey:   types.StringValue(*response.PublicKey),
		Fingerprint: types.StringValue(*response.Fingerprint),
		Environment: d.MapEnvironment(ctx, diags, *response.Environment),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
	}
}

func (d *DataSourceCoreKeypair) MapEnvironment(
	ctx context.Context,
	diags *diag.Diagnostics,
	data keypair.KeypairEnvironmentFields,
) datasource_core_keypair.EnvironmentValue {
	model, diagnostic := datasource_core_keypair.NewEnvironmentValue(
		datasource_core_keypair.EnvironmentValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":     types.Int64Value(int64(*data.Id)),
			"name":   types.StringValue(*data.Name),
			"region": types.StringValue(*data.Region),
			//"features": d.MapEnvironmentFeatures(ctx, diags, *data.Features),
		},
	)
	diags.Append(diagnostic...)

	return model
}
