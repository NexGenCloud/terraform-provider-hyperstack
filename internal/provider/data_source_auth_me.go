package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/auth"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_auth_me"
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

	callResult := result.JSON200.User
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No user data",
			"",
		)
		return
	}

	data = d.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceAuthMe) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *auth.UserFields,
) datasource_auth_me.AuthMeModel {
	return datasource_auth_me.AuthMeModel{
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
		Email: func() types.String {
			if response.Email == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Email)
		}(),
		Name: func() types.String {
			if response.Name == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Name)
		}(),
		Username: func() types.String {
			if response.Username == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Username)
		}(),
	}
}
