package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/rbac_role"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/datasource_auth_roles"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceAuthRoles{}

func NewDataSourceAuthRoles() datasource.DataSource {
	return &DataSourceAuthRoles{}
}

type DataSourceAuthRoles struct {
	hyperstack *client.HyperstackClient
	client     *rbac_role.ClientWithResponses
}

func (d *DataSourceAuthRoles) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_roles"
}

func (d *DataSourceAuthRoles) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_auth_roles.AuthRolesDataSourceSchema(ctx)
}

func (d *DataSourceAuthRoles) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = rbac_role.NewClientWithResponses(
		d.hyperstack.ApiServer,
		rbac_role.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceAuthRoles) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_auth_roles.AuthRolesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListRBACRolesWithResponse(ctx)
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

	callResult := result.JSON200.Roles
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	data.Roles = d.MapRoles(ctx, resp, *callResult)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceAuthRoles) MapRoles(
	ctx context.Context,
	resp *datasource.ReadResponse,
	data []rbac_role.RBACRoleFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_roles.RolesValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				createdAt := types.StringNull()
				if row.CreatedAt != nil {
					createdAt = types.StringValue(row.CreatedAt.String())
				}

				model, diagnostic := datasource_auth_roles.NewRolesValue(
					datasource_auth_roles.RolesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":          types.Int64Value(int64(*row.Id)),
						"name":        types.StringValue(*row.Name),
						"description": types.StringValue(*row.Description),
						"policies":    d.MapRolesPolicies(ctx, resp, *row.Policies),
						"permissions": d.MapRolesPermissions(ctx, resp, *row.Permissions),
						"created_at":  createdAt,
					},
				)
				resp.Diagnostics.Append(diagnostic...)
				roles = append(roles, model)
			}
			return roles
		}(),
	)
	resp.Diagnostics.Append(diagnostic...)
	return model
}

func (d *DataSourceAuthRoles) MapRolesPolicies(
	ctx context.Context,
	resp *datasource.ReadResponse,
	data []rbac_role.RolePolicyFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_roles.PoliciesValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_roles.NewPoliciesValue(
					datasource_auth_roles.PoliciesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":          types.Int64Value(int64(*row.Id)),
						"name":        types.StringValue(*row.Name),
						"description": types.StringValue(*row.Description),
					},
				)
				resp.Diagnostics.Append(diagnostic...)
				roles = append(roles, model)
			}
			return roles
		}(),
	)
	resp.Diagnostics.Append(diagnostic...)
	return model
}

func (d *DataSourceAuthRoles) MapRolesPermissions(
	ctx context.Context,
	resp *datasource.ReadResponse,
	data []rbac_role.RolePermissionFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_roles.PermissionsValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_roles.NewPermissionsValue(
					datasource_auth_roles.PermissionsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":         types.Int64Value(int64(*row.Id)),
						"resource":   types.StringValue(*row.Resource),
						"permission": types.StringValue(*row.Permission),
					},
				)
				resp.Diagnostics.Append(diagnostic...)
				roles = append(roles, model)
			}
			return roles
		}(),
	)
	resp.Diagnostics.Append(diagnostic...)
	return model
}
