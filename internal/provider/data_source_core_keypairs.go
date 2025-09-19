package provider

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/NexGenCloud/hyperstack-sdk-go/lib/keypair"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_keypairs"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DataSourceCoreKeypairs{}

func NewDataSourceCoreKeypairs() datasource.DataSource {
	return &DataSourceCoreKeypairs{}
}

func (d *DataSourceCoreKeypairs) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_keypairs"
}

func (d *DataSourceCoreKeypairs) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_keypairs.CoreKeypairsDataSourceSchema(ctx)
}

type DataSourceCoreKeypairs struct {
	hyperstack *client.HyperstackClient
	client     *keypair.ClientWithResponses
}

func (d *DataSourceCoreKeypairs) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceCoreKeypairs) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_keypairs.CoreKeypairsModel

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

	callResult := result.JSON200.Keypairs
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

func (d *DataSourceCoreKeypairs) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]keypair.KeypairFields,
) datasource_core_keypairs.CoreKeypairsModel {
	return datasource_core_keypairs.CoreKeypairsModel{
		CoreKeypairs: func() types.Set {
			return d.MapKeypairs(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceCoreKeypairs) MapKeypairs(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []keypair.KeypairFields,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_core_keypairs.CoreKeypairsValue{}.Type(ctx),
		func() []attr.Value {
			keypairs := make([]attr.Value, 0)
			for _, row := range data {
				environment, diagnostic := d.MapEnvironment(ctx, diags, *row.Environment).ToObjectValue(ctx)
				diags.Append(diagnostic...)
				
				model, diagnostic := datasource_core_keypairs.NewCoreKeypairsValue(
					datasource_core_keypairs.CoreKeypairsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":          types.Int64Value(int64(*row.Id)),
						"name":        types.StringValue(*row.Name),
						"public_key":  types.StringValue(*row.PublicKey),
						"fingerprint": types.StringValue(*row.Fingerprint),
						"environment": environment,
						"created_at": func() attr.Value {
							if row.CreatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.CreatedAt.String())
						}(),
					},
				)
				diags.Append(diagnostic...)
				keypairs = append(keypairs, model)
			}
			return keypairs
		}(),
	)
	diags.Append(diagnostic...)
	return model
}

func (d *DataSourceCoreKeypairs) MapEnvironment(
	ctx context.Context,
	diags *diag.Diagnostics,
	data keypair.KeypairEnvironmentFields,
) datasource_core_keypairs.EnvironmentValue {
	model, diagnostic := datasource_core_keypairs.NewEnvironmentValue(
		datasource_core_keypairs.EnvironmentValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id": func() attr.Value {
				if data.Id == nil {
					return types.Int64Null()
				}
				return types.Int64Value(int64(*data.Id))
			}(),
			"name": func() attr.Value {
				if data.Name == nil {
					return types.StringNull()
				}
				return types.StringValue(*data.Name)
			}(),
			"region": func() attr.Value {
				if data.Region == nil {
					return types.StringNull()
				}
				return types.StringValue(*data.Region)
			}(),
			"created_at": func() attr.Value {
				if data.CreatedAt == nil {
					return types.StringNull()
				}
				return types.StringValue(data.CreatedAt.String())
			}(),
			"features": func() attr.Value {
				if data.Features == nil {
					return types.ObjectNull(map[string]attr.Type{
						"green_status":      types.StringType,
						"network_optimised": types.BoolType,
					})
				}
				return d.MapEnvironmentFeatures(ctx, diags, *data.Features)
			}(),
		},
	)
	diags.Append(diagnostic...)

	return model
}

func (d *DataSourceCoreKeypairs) MapEnvironmentFeatures(
	ctx context.Context,
	diags *diag.Diagnostics,
	data keypair.KeypairEnvironmentFeatures,
) types.Object {
	model, diagnostic := types.ObjectValue(
		map[string]attr.Type{
			"green_status":      types.StringType,
			"network_optimised": types.BoolType,
		},
		map[string]attr.Value{
			"green_status": func() attr.Value {
				if data.GreenStatus == nil {
					return types.StringNull()
				}
				return types.StringValue(string(*data.GreenStatus))
			}(),
			"network_optimised": func() attr.Value {
				if data.NetworkOptimised == nil {
					return types.BoolNull()
				}
				return types.BoolValue(*data.NetworkOptimised)
			}(),
		},
	)
	diags.Append(diagnostic...)

	return model
}
