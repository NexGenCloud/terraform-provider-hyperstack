package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/auth"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/datasource_auth_me"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceAuthMe{}

func NewDataSourceAuthMe() datasource.DataSource {
	return &DataSourceAuthMe{}
}

type DataSourceAuthMe struct {
	hyperstack *client.HyperstackClient
	client     *auth.ClientWithResponses
}

func (d *DataSourceAuthMe) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_me"
}

func (d *DataSourceAuthMe) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_auth_me.AuthMeDataSourceSchema(ctx)
}

func (d *DataSourceAuthMe) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = auth.NewClientWithResponses(
		d.hyperstack.ApiServer,
		auth.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceAuthMe) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_auth_me.AuthMeModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.AuthUserInformationWithResponse(ctx)
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

	userResult := result.JSON200.User
	if userResult == nil {
		resp.Diagnostics.AddWarning(
			"No user data",
			"",
		)
		return
	}

	user, diag := datasource_auth_me.NewUserValue(
		datasource_auth_me.UserValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"email":      types.StringValue(*userResult.Email),
			"name":       types.StringValue(*userResult.Name),
			"username":   types.StringValue(*userResult.Username),
			"created_at": types.StringValue(userResult.CreatedAt.String()),
		},
	)
	resp.Diagnostics.Append(diag...)
	data.User = user

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
