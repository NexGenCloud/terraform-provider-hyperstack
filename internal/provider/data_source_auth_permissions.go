package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/permission"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_auth_permissions"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceAuthPermissions{}

func NewDataSourceAuthPermissions() datasource.DataSource {
	return &DataSourceAuthPermissions{}
}

type DataSourceAuthPermissions struct {
	hyperstack *client.HyperstackClient
	client     *permission.ClientWithResponses
}

func (d *DataSourceAuthPermissions) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_permissions"
}

func (d *DataSourceAuthPermissions) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_auth_permissions.AuthPermissionsDataSourceSchema(ctx)
}

func (d *DataSourceAuthPermissions) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = permission.NewClientWithResponses(
		d.hyperstack.ApiServer,
		permission.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceAuthPermissions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_auth_permissions.AuthPermissionsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListPermissionsWithResponse(ctx)
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

	callResult := result.JSON200.Permissions
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

func (d *DataSourceAuthPermissions) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]permission.PermissionFields,
) datasource_auth_permissions.AuthPermissionsModel {
	return datasource_auth_permissions.AuthPermissionsModel{
		AuthPermissions: func() types.Set {
			return d.MapPermissions(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceAuthPermissions) MapPermissions(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []permission.PermissionFields,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_auth_permissions.AuthPermissionsValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_permissions.NewAuthPermissionsValue(
					datasource_auth_permissions.AuthPermissionsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":         types.Int64Value(int64(*row.Id)),
						"resource":   types.StringValue(*row.Resource),
						"permission": types.StringValue(*row.Permission),
						"method":     types.StringValue(*row.Method),
						"endpoint":   types.StringValue(*row.Endpoint),
						"created_at": func() attr.Value {
							if row.CreatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.CreatedAt.String())
						}(),
					},
				)
				diags.Append(diagnostic...)
				roles = append(roles, model)
			}
			return roles
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
