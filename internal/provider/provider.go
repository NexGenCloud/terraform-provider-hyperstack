package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &hyperstackProvider{}

type hyperstackProvider struct {
	version string
}

type hyperstackProviderModel struct {
	ApiToken   types.String `tfsdk:"api_key"`
	ApiAddress types.String `tfsdk:"api_address"`
}

func (p *hyperstackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hyperstack"
	resp.Version = p.version
}

func (p *hyperstackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for Nexgen Hyperstack platform",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Hyperstack API token",
				Optional:    true,
				Sensitive:   true,
			},
			"api_address": schema.StringAttribute{
				Description: "Hyperstack API address",
				Optional:    true,
			},
		},
	}
}

func (p *hyperstackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var data hyperstackProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiTokenEnv := fmt.Sprintf("%sAPI_KEY", EnvPrefix)
	apiToken := os.Getenv(apiTokenEnv)
	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	}
	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing API token Configuration",
			fmt.Sprintf(
				"While configuring the provider, the API token was not found in "+
					"the %s environment variable or provider "+
					"configuration block api_key attribute.",
				apiTokenEnv,
			),
		)
	}

	apiAddressEnv := fmt.Sprintf("%sAPI_ADDRESS", EnvPrefix)
	apiAddress := os.Getenv(apiAddressEnv)
	if !data.ApiAddress.IsNull() {
		apiAddress = data.ApiAddress.ValueString()
	}
	if apiAddress == "" {
		apiAddress = ApiAddress
	}
	if apiAddress == "" {
		resp.Diagnostics.AddError(
			"Missing API server Configuration",
			fmt.Sprintf(
				"While configuring the provider, the API server was not found in "+
					"the %s environment variable or provider "+
					"configuration block api_key attribute.",
				apiAddressEnv,
			),
		)
	}

	hyperstack := client.NewHyperstackClient(
		apiToken,
		apiAddress,
	)
	resp.DataSourceData = hyperstack
	resp.ResourceData = hyperstack
}

func (p *hyperstackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewResourceAuthRole,
		NewResourceCoreCluster,
		NewResourceCoreEnvironment,
		NewResourceCoreKeypair,
		NewResourceCoreVirtualMachine,
		NewResourceCoreVirtualMachineSgRule,
		NewResourceCoreVolume,
		NewResourceCoreVolumeAttachment,
	}
}

func (p *hyperstackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDataSourceAuthMe,
		NewDataSourceAuthOrganization,
		NewDataSourceAuthPermissions,
		NewDataSourceAuthPolicies,
		NewDataSourceAuthRole,
		NewDataSourceAuthRoles,
		NewDataSourceAuthUserMePermissions,
		NewDataSourceAuthUserPermissions,
		NewDataSourceCoreClustersVersions,
		NewDataSourceCoreDashboard,
		NewDataSourceCoreEnvironment,
		NewDataSourceCoreEnvironments,
		NewDataSourceCoreFirewallProtocols,
		NewDataSourceCoreFlavors,
		NewDataSourceCoreGpus,
		NewDataSourceCoreImages,
		NewDataSourceCoreKeypair,
		NewDataSourceCoreKeypairs,
		NewDataSourceCoreRegions,
		NewDataSourceCoreStocks,
		NewDataSourceCoreVirtualMachines,
		NewDataSourceCoreVolumeTypes,
		NewDataSourceCoreVolumes,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hyperstackProvider{
			version: version,
		}
	}
}
