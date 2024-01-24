package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/keypair"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/datasource_keypairs"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/resource_keypair"
	"io/ioutil"
	"strings"
)

var _ resource.Resource = &ResourceKeypair{}
var _ resource.ResourceWithImportState = &ResourceKeypair{}

func NewResourceKeypair() resource.Resource {
	return &ResourceKeypair{}
}

type ResourceKeypair struct {
	hyperstack *client.HyperstackClient
	client     *keypair.ClientWithResponses
}

func (r *ResourceKeypair) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_keypair.KeypairModel
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
	var id int64
	for _, row := range *searchCallResult {
		if strings.Contains(*row.Name, strings.TrimSpace(data.Name.ValueString())) {
			id = int64(*row.Id)
			break
		}
	}

	// Perform the update operation
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

	data.Keypair = MapKeypairFieldsToKeypairValue(callResult.Keypair)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceKeypair) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keypairs"
}

func (r *ResourceKeypair) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_keypair.KeypairResourceSchema(ctx)
	resp.Schema.Attributes["public_key"] = schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
	resp.Schema.Attributes["environment_name"] = schema.StringAttribute{
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	}
}

func (r *ResourceKeypair) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceKeypair) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data resource_keypair.KeypairModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	result, err := r.client.ImportKeypairWithResponse(ctx, func() keypair.ImportKeypairJSONRequestBody {
		return keypair.ImportKeypairJSONRequestBody{
			EnvironmentName: data.EnvironmentName.ValueString(),
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
	data.Keypair = MapKeypairFieldsToKeypairValue(callResult.Keypair)

	//Set data.Keypairs to new empty list if it's not yet set
	data.Keypairs, _ = types.ListValue(datasource_keypairs.KeypairsValue{}.Type(ctx), func() []attr.Value {
		return make([]attr.Value, 0)
	}())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceKeypair) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data resource_keypair.KeypairModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	data.Keypairs = r.MapKeypairs(ctx, resp, *callResult)
	//set data.Keypair to ( find in data.Keypairs where name == data.Name )
	for _, row := range *callResult {
		if strings.Contains(*row.Name, strings.TrimSpace(data.Name.ValueString())) {
			data.Keypair = MapDataSourceKeypairFieldsToKeypairValue(row)

			break
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceKeypair) MapKeypairs(
	ctx context.Context,
	resp *resource.ReadResponse,
	data []keypair.KeypairFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_keypairs.KeypairsValue{}.Type(ctx),
		func() []attr.Value {
			keypairs := make([]attr.Value, 0)
			for _, row := range data {
				createdAt := types.StringNull()
				if row.CreatedAt != nil {
					createdAt = types.StringValue(row.CreatedAt.String())
				}

				model, diagnostic := datasource_keypairs.NewKeypairsValue(
					datasource_keypairs.KeypairsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":          types.Int64Value(int64(*row.Id)),
						"name":        types.StringValue(*row.Name),
						"public_key":  types.StringValue(*row.PublicKey),
						"fingerprint": types.StringValue(*row.Fingerprint),
						"environment": types.StringValue(*row.Environment),
						"created_at":  createdAt,
					},
				)
				resp.Diagnostics.Append(diagnostic...)
				keypairs = append(keypairs, model)
			}
			return keypairs
		}(),
	)
	resp.Diagnostics.Append(diagnostic...)
	return model
}

func (r *ResourceKeypair) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_keypair.KeypairModel
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

func (r *ResourceKeypair) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func MapKeypairFieldsToKeypairValue(kf *keypair.KeypairFields) resource_keypair.KeypairValue {
	kv := resource_keypair.KeypairValue{
		CreatedAt:   types.StringValue(kf.CreatedAt.String()),
		Environment: types.StringValue(*kf.Environment),
		Fingerprint: types.StringValue(*kf.Fingerprint),
		Id:          types.Int64Value(int64(*kf.Id)),
		Name:        types.StringValue(*kf.Name),
		PublicKey:   types.StringValue(*kf.PublicKey),
	}

	return kv
}

func MapDataSourceKeypairFieldsToKeypairValue(kf keypair.KeypairFields) resource_keypair.KeypairValue {
	kv := resource_keypair.KeypairValue{
		CreatedAt:   types.StringValue(kf.CreatedAt.String()),
		Environment: types.StringValue(*kf.Environment),
		Fingerprint: types.StringValue(*kf.Fingerprint),
		Id:          types.Int64Value(int64(*kf.Id)),
		Name:        types.StringValue(*kf.Name),
		PublicKey:   types.StringValue(*kf.PublicKey),
	}

	return kv
}
