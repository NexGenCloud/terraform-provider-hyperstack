package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/policy"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_auth_policies"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceAuthPolicies{}

func NewDataSourceAuthPolicies() datasource.DataSource {
	return &DataSourceAuthPolicies{}
}

type DataSourceAuthPolicies struct {
	hyperstack *client.HyperstackClient
	client     *policy.ClientWithResponses
}

func (d *DataSourceAuthPolicies) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_policies"
}

func (d *DataSourceAuthPolicies) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_auth_policies.AuthPoliciesDataSourceSchema(ctx)
}

func (d *DataSourceAuthPolicies) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = policy.NewClientWithResponses(
		d.hyperstack.ApiServer,
		policy.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceAuthPolicies) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_auth_policies.AuthPoliciesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListPoliciesWithResponse(ctx)
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

	callResult := result.JSON200.Policies
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

func (d *DataSourceAuthPolicies) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]policy.PolicyFields,
) datasource_auth_policies.AuthPoliciesModel {
	return datasource_auth_policies.AuthPoliciesModel{
		AuthPolicies: func() types.Set {
			return d.MapPolicies(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceAuthPolicies) MapPolicies(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []policy.PolicyFields,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_auth_policies.AuthPoliciesValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_policies.NewAuthPoliciesValue(
					datasource_auth_policies.AuthPoliciesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":          types.Int64Value(int64(*row.Id)),
						"name":        types.StringValue(*row.Name),
						"description": types.StringValue(*row.Description),
						"permissions": d.MapPoliciesPermissions(ctx, diags, *row.Permissions),
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

func (d *DataSourceAuthPolicies) MapPoliciesPermissions(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []policy.PolicyPermissionFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_auth_policies.PermissionsValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_auth_policies.NewPermissionsValue(
					datasource_auth_policies.PermissionsValue{}.AttributeTypes(ctx),
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
