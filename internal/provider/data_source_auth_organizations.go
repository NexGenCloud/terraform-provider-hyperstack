package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/organization"
	"github.com/NexGenCloud/hyperstack-terraform-provider/internal/client"
	"github.com/NexGenCloud/hyperstack-terraform-provider/internal/genprovider/datasource_auth_organizations"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceAuthOrganizations{}

func NewDataSourceAuthOrganizations() datasource.DataSource {
	return &DataSourceAuthOrganizations{}
}

type DataSourceAuthOrganizations struct {
	hyperstack *client.HyperstackClient
	client     *organization.ClientWithResponses
}

func (d *DataSourceAuthOrganizations) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_organizations"
}

func (d *DataSourceAuthOrganizations) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_auth_organizations.AuthOrganizationsDataSourceSchema(ctx)
}

func (d *DataSourceAuthOrganizations) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSourceAuthOrganizations) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_auth_organizations.AuthOrganizationsModel

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

func (d *DataSourceAuthOrganizations) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *organization.OrganizationInfoModel,
) datasource_auth_organizations.AuthOrganizationsModel {
	return datasource_auth_organizations.AuthOrganizationsModel{
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

func (d *DataSourceAuthOrganizations) MapUsers(
	ctx context.Context,
	diags *diag.Diagnostics,
	data *[]organization.OrganizationUserModel,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_organizations.UsersValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range *data {
				model, diagnostic := datasource_auth_organizations.NewUsersValue(
					datasource_auth_organizations.UsersValue{}.AttributeTypes(ctx),
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
						"name": func() types.String {
							return types.StringValue(*row.Name)
						}(),
						"role": types.StringValue(*row.Role),
						"joined_at": func() types.String {
							if row.JoinedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.JoinedAt.String())
						}(),

						// TODO: implement
						"rbac_roles": types.ListNull(datasource_auth_organizations.RbacRolesValue{}.Type(ctx)),
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
