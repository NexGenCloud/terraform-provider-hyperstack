package provider

import (
	"context"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	API_SERVER         = "https://infrahub-api.nexgencloud.com/v1"
	API_SERVER_STAGING = "https://infrahub-api-stg.ngbackend.cloud/v1"
)

var _ provider.Provider = &hyperstackProvider{}

type hyperstackProvider struct {
	version string
}

type hyperstackProviderModel struct {
	ApiToken types.String `tfsdk:"api_key"`
	Staging  types.Bool   `tfsdk:"staging"`
}

func (p *hyperstackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hyperstack"
	resp.Version = p.version
}

func (p *hyperstackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"staging": schema.BoolAttribute{
				MarkdownDescription: "If staging server should be used",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Hyperstack API token",
				Optional:            true,
			},
		},
	}
}

func (p *hyperstackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	apiToken := os.Getenv("HYPERSTACK_API_KEY")
	staging := os.Getenv("HYPERSTACK_STAGING") == "true"

	var data hyperstackProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	}

	if !data.Staging.IsNull() {
		staging = data.Staging.ValueBool()
	}

	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing API Token Configuration",
			"While configuring the provider, the API token was not found in "+
				"the HYPERSTACK_API_KEY environment variable or provider "+
				"configuration block api_key attribute.",
		)
	}

	apiServer := API_SERVER
	if staging {
		apiServer = API_SERVER_STAGING
	}

	hyperstack := client.NewHyperstackClient(
		apiToken,
		apiServer,
	)
	resp.DataSourceData = hyperstack
	resp.ResourceData = hyperstack
}

func (p *hyperstackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewResourceAuthRole,
		NewResourceCoreEnvironment,
		NewResourceCoreKeypair,
		NewResourceCoreVirtualMachine,
		NewResourceCoreVirtualMachineSgRule,
	}
}

func (p *hyperstackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataSourceAuthMe,
		NewDataSourceAuthOrganizations,
		NewDataSourceAuthRole,
		NewDataSourceAuthRoles,
		NewDataSourceCoreEnvironments,
		NewDataSourceCoreKeypairs,
		NewDataSourceCoreVirtualMachines,
		NewDataSourceCoreRegions,
		NewDataSourceCoreGpus,
		NewDataSourceCoreFlavors,
		NewDataSourceCoreImages,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hyperstackProvider{
			version: version,
		}
	}
}
