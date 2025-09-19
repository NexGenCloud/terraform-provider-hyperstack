package provider

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/NexGenCloud/hyperstack-sdk-go/lib/rbac_role"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_auth_role"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	result, err := d.client.RetrieveRBACRoleDetailsWithResponse(ctx, int(data.Id.ValueInt64()))
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

func (d *DataSourceAuthRole) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *rbac_role.RbacRoleFields,
) datasource_auth_role.AuthRoleModel {
	return datasource_auth_role.AuthRoleModel{
		Id: func() types.Int64 {
			return types.Int64Value(int64(*response.Id))
		}(),
		Name: func() types.String {
			if response.Name == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Name)
		}(),
		Description: func() types.String {
			if response.Description == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Description)
		}(),
		Policies: func() types.List {
			return d.MapRolePolicies(ctx, diags, *response.Policies)
		}(),
		Permissions: func() types.List {
			return d.MapRolePermissions(ctx, diags, *response.Permissions)
		}(),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
	}
}

func (d *DataSourceAuthRole) MapRolePolicies(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []rbac_role.RolePolicyFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_role.PoliciesValue{}.Type(ctx),
		func() []attr.Value {
			list := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_role.NewPoliciesValue(
					datasource_auth_role.PoliciesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":          types.Int64Value(int64(*row.Id)),
						"name":        types.StringValue(*row.Name),
						"description": types.StringValue(*row.Description),
					},
				)
				diags.Append(diagnostic...)
				list = append(list, model)
			}
			return list
		}(),
	)
	diags.Append(diagnostic...)
	return model
}

func (d *DataSourceAuthRole) MapRolePermissions(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []rbac_role.RolePermissionFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_role.PermissionsValue{}.Type(ctx),
		func() []attr.Value {
			list := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_role.NewPermissionsValue(
					datasource_auth_role.PermissionsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":         types.Int64Value(int64(*row.Id)),
						"resource":   types.StringValue(*row.Resource),
						"permission": types.StringValue(*row.Permission),
					},
				)
				diags.Append(diagnostic...)
				list = append(list, model)
			}
			return list
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
