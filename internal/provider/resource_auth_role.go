package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/nexgen/hyperstack-sdk-go/lib/rbac_role"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/resource_auth_role"
)

var _ resource.Resource = &ResourceAuthRole{}
var _ resource.ResourceWithImportState = &ResourceAuthRole{}

func NewResourceAuthRole() resource.Resource {
	return &ResourceAuthRole{}
}

type ResourceAuthRole struct {
	hyperstack *client.HyperstackClient
	client     *rbac_role.ClientWithResponses
}

func (r *ResourceAuthRole) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_role"
}

func (r *ResourceAuthRole) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_auth_role.AuthRoleResourceSchema(ctx)
}

func (r *ResourceAuthRole) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	r.client, err = rbac_role.NewClientWithResponses(
		r.hyperstack.ApiServer,
		rbac_role.WithRequestEditorFn(r.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (r *ResourceAuthRole) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data resource_auth_role.AuthRoleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateRBACRoleWithResponse(ctx, func() rbac_role.CreateRBACRoleJSONRequestBody {
		return rbac_role.CreateRBACRoleJSONRequestBody{
			Description: data.Description.ValueString(),
			Name:        data.Name.ValueString(),
			Policies: func() *[]int {
				items := make([]int, 0)
				for _, row := range data.Policies.Elements() {
					items = append(items, int(row.(basetypes.Int64Value).ValueInt64()))
				}
				return &items
			}(),
			Permissions: func() *[]int {
				items := make([]int, 0)
				for _, row := range data.Permissions.Elements() {
					items = append(items, int(row.(basetypes.Int64Value).ValueInt64()))
				}
				return &items
			}(),
		}
	}())
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if result.JSON201 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	callResult := result.JSON201.Role
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	data = r.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceAuthRole) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data resource_auth_role.AuthRoleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetARBACRoleDetailWithResponse(ctx, int(data.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if result.JSON404 != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if result.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	callResult := result.JSON200.Roles
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	data = r.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceAuthRole) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var dataOld resource_auth_role.AuthRoleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataOld)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(dataOld.Id.ValueInt64())

	var data resource_auth_role.AuthRoleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.UpdateARBACRoleWithResponse(ctx, id, func() rbac_role.UpdateARBACRoleJSONRequestBody {
		return rbac_role.UpdateARBACRoleJSONRequestBody{
			Description: data.Description.ValueString(),
			Name:        data.Name.ValueString(),
			Policies: func() *[]int {
				items := make([]int, 0)
				for _, row := range data.Policies.Elements() {
					items = append(items, int(row.(basetypes.Int64Value).ValueInt64()))
				}
				return &items
			}(),
			Permissions: func() *[]int {
				items := make([]int, 0)
				for _, row := range data.Permissions.Elements() {
					items = append(items, int(row.(basetypes.Int64Value).ValueInt64()))
				}
				return &items
			}(),
		}
	}())
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if result.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	callResult := result.JSON200.Role
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	data = r.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceAuthRole) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_auth_role.AuthRoleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt64())

	result, err := r.client.DeleteARBACRoleWithResponse(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if result.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ResourceAuthRole) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ResourceAuthRole) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *rbac_role.RBACRoleFields,
) resource_auth_role.AuthRoleModel {
	return resource_auth_role.AuthRoleModel{
		Id: func() types.Int64 {
			if response.Id == nil {
				return types.Int64Null()
			}
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
			return r.MapRolesPolicies(ctx, diags, *response.Policies)
		}(),
		Permissions: func() types.List {
			return r.MapRolesPermissions(ctx, diags, *response.Permissions)
		}(),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
	}
}

func (r *ResourceAuthRole) MapRolesPolicies(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []rbac_role.RolePolicyFields,
) types.List {
	model, diagnostic := types.ListValue(
		types.Int64Type,
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model := types.Int64Value(int64(*row.Id))
				roles = append(roles, model)
			}
			return roles
		}(),
	)
	diags.Append(diagnostic...)
	return model
}

func (r *ResourceAuthRole) MapRolesPermissions(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []rbac_role.RolePermissionFields,
) types.List {
	model, diagnostic := types.ListValue(
		types.Int64Type,
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model := types.Int64Value(int64(*row.Id))
				roles = append(roles, model)
			}
			return roles
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
