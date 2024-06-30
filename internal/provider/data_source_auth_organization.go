package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/organization"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_auth_organization"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceAuthOrganization{}

func NewDataSourceAuthOrganization() datasource.DataSource {
	return &DataSourceAuthOrganization{}
}

type DataSourceAuthOrganization struct {
	hyperstack *client.HyperstackClient
	client     *organization.ClientWithResponses
}

func (d *DataSourceAuthOrganization) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_organization"
}

func (d *DataSourceAuthOrganization) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_auth_organization.AuthOrganizationDataSourceSchema(ctx)
}

func (d *DataSourceAuthOrganization) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = organization.NewClientWithResponses(
		d.hyperstack.ApiServer,
		organization.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceAuthOrganization) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_auth_organization.AuthOrganizationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.GetOrganizationInfoWithResponse(ctx)
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

	callResult := result.JSON200.Organization
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No organization data",
			"",
		)
		return
	}

	data = d.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceAuthOrganization) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *organization.OrganizationInfoModel,
) datasource_auth_organization.AuthOrganizationModel {
	return datasource_auth_organization.AuthOrganizationModel{
		Id: func() types.Int64 {
			return types.Int64Value(int64(response.Id))
		}(),
		Name: func() types.String {
			return types.StringValue(response.Name)
		}(),
		Users: func() types.List {
			return d.MapUsers(ctx, diags, response.Users)
		}(),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
	}
}

func (d *DataSourceAuthOrganization) MapUsers(
	ctx context.Context,
	diags *diag.Diagnostics,
	data *[]organization.OrganizationUserModel,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_organization.UsersValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range *data {
				model, diagnostic := datasource_auth_organization.NewUsersValue(
					datasource_auth_organization.UsersValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id": func() types.Int64 {
							return types.Int64Value(int64(*row.Id))
						}(),
						"sub": func() types.String {
							return types.StringValue(*row.Sub)
						}(),
						"email": func() types.String {
							return types.StringValue(*row.Email)
						}(),
						"username": func() types.String {
							return types.StringValue(*row.Username)
						}(),
						// Not available on staging yet ??
						//"last_login": func() types.String {
						//	if row.LastLogin == nil {
						//		return types.StringNull()
						//	}
						//	return types.StringValue(row.LastLogin.String())
						//}(),
						"name": func() types.String {
							return types.StringValue(*row.Name)
						}(),
						"role":       types.StringValue(*row.Role),
						"rbac_roles": d.MapUsersRoles(ctx, diags, row.RbacRoles),
						"joined_at": func() types.String {
							if row.JoinedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.JoinedAt.String())
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

func (d *DataSourceAuthOrganization) MapUsersRoles(
	ctx context.Context,
	diags *diag.Diagnostics,
	data *[]organization.RBACRoleFieldForOrganization,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_organization.RbacRolesValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range *data {
				model, diagnostic := datasource_auth_organization.NewUsersValue(
					datasource_auth_organization.RbacRolesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"name": func() types.String {
							return types.StringValue(*row.Name)
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
