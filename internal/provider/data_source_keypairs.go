package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/keypair"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/datasource_keypairs"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceKeypairs{}

func NewDataSourceKeypairs() datasource.DataSource {
	return &DataSourceKeypairs{}
}

func (d *DataSourceKeypairs) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs"
}

func (d *DataSourceKeypairs) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_keypairs.KeypairsDataSourceSchema(ctx)
}

type DataSourceKeypairs struct {
	hyperstack *client.HyperstackClient
	client     *keypair.ClientWithResponses
}

func (d *DataSourceKeypairs) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceKeypairs) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_keypairs.KeypairsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.RetrieveUserKeypairsWithResponse(ctx)
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

	callResult := result.JSON200.Keypairs
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	data.Keypairs = d.MapKeypairs(ctx, resp, *callResult)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceKeypairs) MapKeypairs(
	ctx context.Context,
	resp *datasource.ReadResponse,
	data []keypair.KeypairFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_keypairs.KeypairsValue{}.Type(ctx),
		func() []attr.Value {
			keypairs := make([]attr.Value, 0)
			for _, row := range data {
				createdAt := types.StringNull()
				if row.CreatedAt != nil {
					createdAt = types.StringValue(row.CreatedAt.String())
				}

				model, diagnostic := datasource_keypairs.NewKeypairsValue(
					datasource_keypairs.KeypairsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":          types.Int64Value(int64(*row.Id)),
						"name":        types.StringValue(*row.Name),
						"public_key":  types.StringValue(*row.PublicKey),
						"fingerprint": types.StringValue(*row.Fingerprint),
						"environment": types.StringValue(*row.Environment),
						"created_at":  createdAt,
					},
				)
				resp.Diagnostics.Append(diagnostic...)
				keypairs = append(keypairs, model)
			}
			return keypairs
		}(),
	)
	resp.Diagnostics.Append(diagnostic...)
	return model
}
