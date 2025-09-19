package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/NexGenCloud/hyperstack-sdk-go/lib/environment"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/resource_core_environment"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ResourceCoreEnvironment{}
var _ resource.ResourceWithImportState = &ResourceCoreEnvironment{}

func NewResourceCoreEnvironment() resource.Resource {
	return &ResourceCoreEnvironment{}
}

type ResourceCoreEnvironment struct {
	hyperstack *client.HyperstackClient
	client     *environment.ClientWithResponses
}

func (r *ResourceCoreEnvironment) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_environment"
}

func (r *ResourceCoreEnvironment) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_core_environment.CoreEnvironmentResourceSchema(ctx)
}

func (r *ResourceCoreEnvironment) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	r.client, err = environment.NewClientWithResponses(
		r.hyperstack.ApiServer,
		environment.WithRequestEditorFn(r.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (r *ResourceCoreEnvironment) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data resource_core_environment.CoreEnvironmentModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateEnvironmentWithResponse(ctx, func() environment.CreateEnvironmentJSONRequestBody {
		return environment.CreateEnvironmentJSONRequestBody{
			Name:   data.Name.ValueString(),
			Region: data.Region.ValueString(),
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

	callResult := result.JSON200.Environment
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

func (r *ResourceCoreEnvironment) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data resource_core_environment.CoreEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.RetrieveEnvironmentWithResponse(ctx, int(data.Id.ValueInt64()))
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

	callResult := result.JSON200.Environment
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

func (r *ResourceCoreEnvironment) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var dataOld resource_core_environment.CoreEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataOld)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: might not need anymore
	id := int(dataOld.Id.ValueInt64())

	var data resource_core_environment.CoreEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.UpdateEnvironmentWithResponse(ctx, id, func() environment.UpdateEnvironmentJSONRequestBody {
		return environment.UpdateEnvironmentJSONRequestBody{
			Name: data.Name.ValueString(),
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

	callResult := result.JSON200.Environment
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

func (r *ResourceCoreEnvironment) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_core_environment.CoreEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt64())

	// Add delay before environment deletion to ensure resource cleanup
	// This addresses the 405 error where API thinks resources still exist
	tflog.Info(ctx, "Waiting 5 seconds before environment deletion to ensure resource cleanup...")
	time.Sleep(5 * time.Second)

	result, err := r.client.DeleteEnvironmentWithResponse(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	// Handle 404 Not Found - resource already deleted
	if result.StatusCode() == 404 {
		// Resource doesn't exist, consider it successfully deleted
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

	resp.State.RemoveResource(ctx)
}

func (r *ResourceCoreEnvironment) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ResourceCoreEnvironment) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *environment.EnvironmentFields,
) resource_core_environment.CoreEnvironmentModel {
	return resource_core_environment.CoreEnvironmentModel{
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
		Region: func() types.String {
			if response.Region == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Region)
		}(),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
		Features: r.MapFeatures(ctx, diags, response.Features),
	}
}

func (r *ResourceCoreEnvironment) MapFeatures(
	ctx context.Context,
	diags *diag.Diagnostics,
	data *environment.EnvironmentFeatures,
) resource_core_environment.FeaturesValue {
	model, diagnostic := resource_core_environment.NewFeaturesValue(
		resource_core_environment.FeaturesValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"green_status": func() attr.Value {
				if data.GreenStatus == nil {
					return types.StringNull()
				}
				return types.StringValue(string(*data.GreenStatus))
			}(),
			"network_optimised": func() attr.Value {
				if data.NetworkOptimised == nil {
					return types.BoolNull()
				}
				return types.BoolValue(*data.NetworkOptimised)
			}(),
		},
	)
	diags.Append(diagnostic...)

	return model
}
