package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/keypair"
	"github.com/NexGenCloud/hyperstack-terraform-provider/internal/client"
	"github.com/NexGenCloud/hyperstack-terraform-provider/internal/genprovider/resource_core_keypair"
	"io/ioutil"
	"strings"
)

var _ resource.Resource = &ResourceCoreKeypair{}
var _ resource.ResourceWithImportState = &ResourceCoreKeypair{}

func NewResourceCoreKeypair() resource.Resource {
	return &ResourceCoreKeypair{}
}

type ResourceCoreKeypair struct {
	hyperstack *client.HyperstackClient
	client     *keypair.ClientWithResponses
}

func (r *ResourceCoreKeypair) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_keypair"
}

func (r *ResourceCoreKeypair) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_core_keypair.CoreKeypairResourceSchema(ctx)
	resp.Schema.Attributes["public_key"] = schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	resp.Schema.Attributes["environment"] = schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}

func (r *ResourceCoreKeypair) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	r.client, err = keypair.NewClientWithResponses(
		r.hyperstack.ApiServer,
		keypair.WithRequestEditorFn(r.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (r *ResourceCoreKeypair) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_core_keypair.CoreKeypairModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Perform a Read operation to get all keypairs
	searchResult, err := r.client.RetrieveUserKeypairsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(searchResult.HTTPResponse.Body)
	if searchResult.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	searchCallResult := searchResult.JSON200.Keypairs
	if searchCallResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	// Find the keypair with the matching name and get its ID
	var id int64 = -1
	for _, row := range *searchCallResult {
		if strings.Contains(*row.Name, strings.TrimSpace(data.Name.ValueString())) {
			id = int64(*row.Id)
			break
		}
	}

	// Check if id was found
	if id == -1 {
		resp.Diagnostics.AddError("No keypair found with the name for update: %s", data.Name.ValueString())
		return
	}

	result, err := r.client.UpdateKeypairNameWithResponse(ctx, int(id), func() keypair.UpdateKeypairNameJSONRequestBody {
		return keypair.UpdateKeypairNameJSONRequestBody{
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

	callResult := result.JSON200
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	data = r.ApiToModel(ctx, &resp.Diagnostics, callResult.Keypair)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreKeypair) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data resource_core_keypair.CoreKeypairModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	result, err := r.client.ImportKeypairWithResponse(ctx, func() keypair.ImportKeypairJSONRequestBody {
		return keypair.ImportKeypairJSONRequestBody{
			EnvironmentName: data.Environment.ValueString(),
			PublicKey:       data.PublicKey.ValueString(),
			Name:            data.Name.ValueString(),
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

	callResult := result.JSON200
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	data = r.ApiToModel(ctx, &resp.Diagnostics, callResult.Keypair)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreKeypair) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data resource_core_keypair.CoreKeypairModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Perform a Read operation to get all keypairs
	searchResult, err := r.client.RetrieveUserKeypairsWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(searchResult.HTTPResponse.Body)
	if searchResult.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	searchCallResult := searchResult.JSON200.Keypairs
	if searchCallResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	// Find the keypair with the matching name and get its ID
	var keypairResult *keypair.KeypairFields
	for _, row := range *searchCallResult {
		if strings.Contains(*row.Name, strings.TrimSpace(data.Name.ValueString())) {
			keypairResult = &row
			break
		}
	}

	// Nothing found
	if keypairResult == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data = r.ApiToModel(ctx, &resp.Diagnostics, keypairResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreKeypair) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_core_keypair.CoreKeypairModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Perform a Read operation to get all keypairs
	result, err := r.client.RetrieveUserKeypairsWithResponse(ctx)
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

	callResult := result.JSON200.Keypairs
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	// Find the keypair with the matching name and get its ID
	var id int64
	for _, row := range *callResult {
		if strings.Contains(*row.Name, strings.TrimSpace(data.Name.ValueString())) {
			id = int64(*row.Id)
			break
		}
	}

	// Now proceed with the Delete operation using the ID
	resultDelete, err := r.client.DeleteKeypairWithResponse(ctx, int(id))
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if resultDelete.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("%s", string(resultDelete.Body)),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ResourceCoreKeypair) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ResourceCoreKeypair) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *keypair.KeypairFields,
) resource_core_keypair.CoreKeypairModel {
	return resource_core_keypair.CoreKeypairModel{
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
		PublicKey: func() types.String {
			if response.PublicKey == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.PublicKey)
		}(),
		Environment: func() types.String {
			if response.Environment == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Environment)
		}(),
		Fingerprint: func() types.String {
			if response.Fingerprint == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Fingerprint)
		}(),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
	}
}
