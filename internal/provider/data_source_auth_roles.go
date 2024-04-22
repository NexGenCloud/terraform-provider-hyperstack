package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/rbac_role"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_auth_roles"
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

	data = d.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceAuthRoles) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]rbac_role.RBACRoleFields,
) datasource_auth_roles.AuthRolesModel {
	return datasource_auth_roles.AuthRolesModel{
		AuthRoles: func() types.Set {
			return d.MapRoles(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceAuthRoles) MapRoles(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []rbac_role.RBACRoleFields,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_auth_roles.AuthRolesValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_roles.NewAuthRolesValue(
					datasource_auth_roles.AuthRolesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":          types.Int64Value(int64(*row.Id)),
						"name":        types.StringValue(*row.Name),
						"description": types.StringValue(*row.Description),
						"policies":    d.MapRolesPolicies(ctx, diags, *row.Policies),
						"permissions": d.MapRolesPermissions(ctx, diags, *row.Permissions),
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

func (d *DataSourceAuthRoles) MapRolesPolicies(
	ctx context.Context,
	diags *diag.Diagnostics,
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
				diags.Append(diagnostic...)
				roles = append(roles, model)
			}
			return roles
		}(),
	)
	diags.Append(diagnostic...)
	return model
}

func (d *DataSourceAuthRoles) MapRolesPermissions(
	ctx context.Context,
	diags *diag.Diagnostics,
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
				diags.Append(diagnostic...)
				roles = append(roles, model)
			}
			return roles
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
