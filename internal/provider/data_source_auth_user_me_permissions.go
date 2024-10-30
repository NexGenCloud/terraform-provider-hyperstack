package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/user_permission"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_auth_user_me_permissions"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceAuthUserMePermissions{}

func NewDataSourceAuthUserMePermissions() datasource.DataSource {
	return &DataSourceAuthUserMePermissions{}
}

type DataSourceAuthUserMePermissions struct {
	hyperstack *client.HyperstackClient
	client     *user_permission.ClientWithResponses
}

func (d *DataSourceAuthUserMePermissions) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_user_me_permissions"
}

func (d *DataSourceAuthUserMePermissions) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_auth_user_me_permissions.AuthUserMePermissionsDataSourceSchema(ctx)
}

func (d *DataSourceAuthUserMePermissions) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = user_permission.NewClientWithResponses(
		d.hyperstack.ApiServer,
		user_permission.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceAuthUserMePermissions) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_auth_user_me_permissions.AuthUserMePermissionsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListMyUserPermissionsWithResponse(ctx)
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

	callResult := result.JSON200.Permissions
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

func (d *DataSourceAuthUserMePermissions) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]user_permission.UserPermissionFields,
) datasource_auth_user_me_permissions.AuthUserMePermissionsModel {
	return datasource_auth_user_me_permissions.AuthUserMePermissionsModel{
		AuthUserMePermissions: func() types.Set {
			return d.MapPermissions(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceAuthUserMePermissions) MapPermissions(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []user_permission.UserPermissionFields,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_auth_user_me_permissions.AuthUserMePermissionsValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_user_me_permissions.NewAuthUserMePermissionsValue(
					datasource_auth_user_me_permissions.AuthUserMePermissionsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":         types.Int64Value(int64(*row.Id)),
						"resource":   types.StringValue(*row.Resource),
						"permission": types.StringValue(*row.Permission),
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
