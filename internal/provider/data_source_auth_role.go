package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/rbac_role"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/datasource_auth_role"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceAuthRole{}

func NewDataSourceAuthRole() datasource.DataSource {
	return &DataSourceAuthRole{}
}

type DataSourceAuthRole struct {
	hyperstack *client.HyperstackClient
	client     *rbac_role.ClientWithResponses
}

func (d *DataSourceAuthRole) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_role"
}

func (d *DataSourceAuthRole) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_auth_role.AuthRoleDataSourceSchema(ctx)
}

func (d *DataSourceAuthRole) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceAuthRole) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_auth_role.AuthRoleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetARBACRoleDetailWithResponse(ctx, int(data.Id.ValueInt64()))
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

	data.Role = d.MapRole(ctx, resp, *callResult)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceAuthRole) MapRole(
	ctx context.Context,
	resp *datasource.ReadResponse,
	data rbac_role.RBACRoleFields,
) datasource_auth_role.RoleValue {
	createdAt := types.StringNull()
	if data.CreatedAt != nil {
		createdAt = types.StringValue(data.CreatedAt.String())
	}

	model, diagnostic := datasource_auth_role.NewRoleValue(
		datasource_auth_role.RoleValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":          types.Int64Value(int64(*data.Id)),
			"name":        types.StringValue(*data.Name),
			"description": types.StringValue(*data.Description),
			"policies":    d.MapRolesPolicies(ctx, resp, *data.Policies),
			"permissions": d.MapRolesPermissions(ctx, resp, *data.Permissions),
			"created_at":  createdAt,
		},
	)
	resp.Diagnostics.Append(diagnostic...)
	return model
}

func (d *DataSourceAuthRole) MapRolesPolicies(
	ctx context.Context,
	resp *datasource.ReadResponse,
	data []rbac_role.RolePolicyFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_role.PoliciesValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_role.NewPoliciesValue(
					datasource_auth_role.PoliciesValue{}.AttributeTypes(ctx),
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

func (d *DataSourceAuthRole) MapRolesPermissions(
	ctx context.Context,
	resp *datasource.ReadResponse,
	data []rbac_role.RolePermissionFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_role.PermissionsValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_role.NewPermissionsValue(
					datasource_auth_role.PermissionsValue{}.AttributeTypes(ctx),
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
