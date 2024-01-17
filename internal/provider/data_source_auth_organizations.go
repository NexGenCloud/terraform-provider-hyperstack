package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/organization"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/datasource_auth_organizations"
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
			"No user data",
			"",
		)
		return
	}

	users := make([]attr.Value, 0)

	for _, userModel := range *callResult.Users {
		user, diag := datasource_auth_organizations.NewUsersValue(
			datasource_auth_organizations.UsersValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":       types.Int64Value(int64(*userModel.Id)),
				"sub":      types.StringValue(*userModel.Sub),
				"email":    types.StringValue(*userModel.Email),
				"username": types.StringValue(*userModel.Username),
				"name":     types.StringValue(*userModel.Name),
				"role":     types.StringValue(*userModel.Role),

				// TODO: implement
				"rbac_roles": types.ListNull(datasource_auth_organizations.RbacRolesValue{}.Type(ctx)),
				"joined_at":  types.StringNull(),
			},
		)
		resp.Diagnostics.Append(diag...)

		users = append(users, user)
	}

	usersValue, diag := types.ListValue(
		datasource_auth_organizations.UsersValue{}.Type(ctx),
		users,
	)
	resp.Diagnostics.Append(diag...)

	org, diag := datasource_auth_organizations.NewOrganizationValue(
		datasource_auth_organizations.OrganizationValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":         types.Int64Value(int64(callResult.Id)),
			"name":       types.StringValue(callResult.Name),
			"users":      usersValue,
			"created_at": types.StringValue(callResult.CreatedAt.String()),
		},
	)
	resp.Diagnostics.Append(diag...)
	data.Organization = org

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
