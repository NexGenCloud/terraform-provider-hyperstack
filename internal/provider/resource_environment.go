package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/environment"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/resource_environments"
)

var _ resource.Resource = &ResourceEnvironment{}
var _ resource.ResourceWithImportState = &ResourceEnvironment{}

func NewResourceEnvironment() resource.Resource {
	return &ResourceEnvironment{}
}

type ResourceEnvironment struct {
	hyperstack *client.HyperstackClient
	client     *environment.ClientWithResponses
}

func (r *ResourceEnvironment) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *ResourceEnvironment) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_environments.EnvironmentsResourceSchema(ctx)
	//resp.Schema.Attributes["id"] = schema.Int64Attribute{
	//	Computed: true,
	//}
}

func (r *ResourceEnvironment) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceEnvironment) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data resource_environments.EnvironmentsModel

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

	callResultMapEnv := r.MapEnvironmentFieldsToEnvironment(ctx, resp.State, resp.Diagnostics, *result.JSON200.Environment)
	data.Environment = r.MapEnvironment(ctx, resp.State, resp.Diagnostics, callResultMapEnv)
	data.Id = types.Int64Value(data.Environment.Id.ValueInt64())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceEnvironment) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data resource_environments.EnvironmentsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetAnEnvironmentDetailsWithResponse(ctx, int(data.Id.ValueInt64()))
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

	callResultMapEnv := r.MapEnvironmentFieldsToEnvironment(ctx, resp.State, resp.Diagnostics, *result.JSON200.Environment)
	data.Environment = r.MapEnvironment(ctx, resp.State, resp.Diagnostics, callResultMapEnv)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceEnvironment) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var dataOld resource_environments.EnvironmentsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataOld)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(dataOld.Id.ValueInt64())

	var data resource_environments.EnvironmentsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.UpdateAnEnvironmentWithResponse(ctx, id, func() environment.UpdateAnEnvironmentJSONRequestBody {
		return environment.UpdateAnEnvironmentJSONRequestBody{
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
	callResultMapEnv := r.MapEnvironmentFieldsToEnvironment(ctx, resp.State, resp.Diagnostics, *result.JSON200.Environment)
	data.Environment = r.MapEnvironment(ctx, resp.State, resp.Diagnostics, callResultMapEnv)
	data.Id = data.Environment.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceEnvironment) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_environments.EnvironmentsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt64())

	result, err := r.client.DeleteAnEnvironmentWithResponse(ctx, id)
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

	resp.State.RemoveResource(ctx)
}

func (r *ResourceEnvironment) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ResourceEnvironment) MapEnvironment(
	ctx context.Context,
	state tfsdk.State,
	diags diag.Diagnostics,
	data environment.Environment,
) resource_environments.EnvironmentValue {
	createdAt := types.StringNull()
	if data.Environment.CreatedAt != nil {
		createdAt = types.StringValue(data.Environment.CreatedAt.String())
	}

	model, diagnostic := resource_environments.NewEnvironmentValue(
		resource_environments.EnvironmentValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":         types.Int64Value(int64(*data.Environment.Id)),
			"name":       types.StringValue(*data.Environment.Name),
			"region":     types.StringValue(*data.Environment.Region),
			"created_at": createdAt,
		},
	)
	diags.Append(diagnostic...)
	return model
}

func (r *ResourceEnvironment) MapEnvironmentFieldsToEnvironment(ctx context.Context, state tfsdk.State, diags diag.Diagnostics, data environment.EnvironmentFields) environment.Environment {
	env := environment.Environment{
		Environment: &environment.EnvironmentFields{
			Id:        data.Id,
			Name:      data.Name,
			Region:    data.Region,
			CreatedAt: data.CreatedAt,
		},
	}
	return env
}
