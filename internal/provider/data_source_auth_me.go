package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"io/ioutil"
	"net/http"
)

var _ datasource.DataSource = &DataSourceAuthMe{}

func NewDataSourceAuthMe() datasource.DataSource {
	return &DataSourceAuthMe{}
}

type AuthMeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	User    struct {
		Email     string `json:"email"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
	} `json:"user"`
}

type DataSourceAuthMe struct {
	hyperstack *client.HyperstackClient
}

type DataSourceAuthMeModel struct {
	Email     types.String `tfsdk:"email"`
	Name      types.String `tfsdk:"name"`
	Username  types.String `tfsdk:"username"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (d *DataSourceAuthMe) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_me"
}

func (d *DataSourceAuthMe) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get me information",

		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				MarkdownDescription: "User email",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "User name",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "User username",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "User created_at",
				Computed:            true,
			},
		},
	}
}

func (d *DataSourceAuthMe) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	hyperstack, ok := req.ProviderData.(*client.HyperstackClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.HyperstackClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.hyperstack = hyperstack
}

func (d *DataSourceAuthMe) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceAuthMeModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// URL of the API
	url := d.hyperstack.ApiServer + "/auth/me"

	// Create a new request
	hreq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	// Add required headers
	hreq.Header.Add("accept", "application/json")
	hreq.Header.Add("api_key", d.hyperstack.ApiToken)

	// Create an HTTP client and send the request
	hresp, err := d.hyperstack.Client.Do(hreq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error",
			fmt.Sprintf("%s", err),
		)
		return
	}
	defer hresp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(hresp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error",
			fmt.Sprintf("%s", err),
		)
		return
	}
	//panic(string(body))

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResponse AuthMeResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	data.Email = types.StringValue(apiResponse.User.Email)
	data.Name = types.StringValue(apiResponse.User.Name)
	data.Username = types.StringValue(apiResponse.User.Username)
	data.CreatedAt = types.StringValue(apiResponse.User.CreatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
