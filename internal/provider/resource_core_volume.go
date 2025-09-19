package provider

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/NexGenCloud/hyperstack-sdk-go/lib/volume"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/resource_core_volume"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ResourceCoreVolume{}
var _ resource.ResourceWithImportState = &ResourceCoreVolume{}

func NewResourceCoreVolume() resource.Resource {
	return &ResourceCoreVolume{}
}

type ResourceCoreVolume struct {
	hyperstack *client.HyperstackClient
	client     *volume.ClientWithResponses
}

func (r *ResourceCoreVolume) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_volume"
}

func (r *ResourceCoreVolume) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_core_volume.CoreVolumeResourceSchema(ctx)
}

func (r *ResourceCoreVolume) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	r.client, err = volume.NewClientWithResponses(
		r.hyperstack.ApiServer,
		volume.WithRequestEditorFn(r.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

// TODO: deduplicate
func (r *ResourceCoreVolume) WaitForResult(
	ctx context.Context,
	pollInterval,
	timeout time.Duration,
	checker func(ctx context.Context) (bool, error),
) error {
	timeoutTimer := time.NewTimer(timeout)
	pollTicker := time.NewTicker(pollInterval)

	defer timeoutTimer.Stop()
	defer pollTicker.Stop()

	for {
		select {
		case <-timeoutTimer.C:
			return fmt.Errorf(
				"Timeout %s reached waiting for resource state change",
				timeout.String(),
			)
		case <-pollTicker.C:
			status, err := checker(ctx)
			if err != nil {
				return fmt.Errorf(
					"Error calling checker while waiting: %e",
					err,
				)
			}
			if status {
				return nil // Resource is created
			}
			// Continue waiting
		}
	}
}

func (r *ResourceCoreVolume) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported for volume resources",
		"Error",
	)
}

func (r *ResourceCoreVolume) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "=== DEBUG: Entered ResourceCoreVolume.Create ===")
	var dataOld resource_core_volume.CoreVolumeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &dataOld)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// DEBUG: Print the payload that will be sent to the API
	payload := volume.CreateVolumeJSONRequestBody{
		CallbackUrl:     dataOld.CallbackUrl.ValueStringPointer(),
		Description:     dataOld.Description.ValueStringPointer(),
		EnvironmentName: dataOld.EnvironmentName.ValueString(),
		ImageId: func() *int {
			if dataOld.ImageId.IsNull() {
				return nil
			}
			val := int(dataOld.ImageId.ValueInt64())
			return &val
		}(),
		Name:       dataOld.Name.ValueString(),
		Size:       int(dataOld.Size.ValueInt64()),
		VolumeType: dataOld.VolumeType.ValueString(),
	}
	tflog.Debug(ctx, "[DEBUG] CreateVolume payload", map[string]interface{}{
		"payload": payload,
	})

	result, err := r.client.CreateVolumeWithResponse(ctx, func() volume.CreateVolumeJSONRequestBody {
		return payload
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

	volumeModel := callResult.Volume
	id := *volumeModel.Id
	err = r.WaitForResult(
		ctx,
		3*time.Second,
		300*time.Second,
		func(ctx context.Context) (bool, error) {
			result, err := r.client.ListVolumesWithResponse(ctx, func() *volume.ListVolumesParams {
				return &volume.ListVolumesParams{
					Page:     nil,
					PageSize: nil,
					Search:   nil,
				}
			}())
			if err != nil {
				return false, err
			}

			fmt.Printf("[DEBUG] Polling for volume id=%v\n", id)
			if result.JSON200 == nil {
				fmt.Printf("[DEBUG] API result is nil, status=%v, body=%s\n", result.StatusCode(), string(result.Body))
			}
			if result.JSON200 != nil {
				fmt.Printf("[DEBUG] API result: status=%v, volumes field present=%v\n", result.JSON200.Status, result.JSON200.Volumes != nil)
			}

			if result.JSON200.Volumes == nil {
				// Treat as empty list, not an error
				fmt.Printf("[DEBUG] Volumes list is nil (SDK struct field 'Volumes', API field 'volumes'), will keep polling.\n")
				return false, nil
			}

			var volumeResult *volume.VolumesFields
			for _, row := range *result.JSON200.Volumes {
				if row.Id == nil {
					continue
				}
				if *row.Id == id {
					volumeResult = &row
					break
				}
			}

			if volumeResult == nil {
				return false, fmt.Errorf("Volume not found")
			}

			if volumeResult.Status == nil {
				return false, fmt.Errorf("volume status is nil")
			}
			status := strings.ToLower(*volumeResult.Status)
			fmt.Printf("Polling volume %d: status=%s\n", id, status)
			fmt.Printf("[DEBUG] Found volume id=%v, status=%v\n", id, *volumeResult.Status)
			switch status {
			case "creating", "downloading":
				return false, nil // keep waiting
			case "available":
				// Convert VolumesFields to VolumeFields for compatibility
				volumeModel = &volume.VolumeFields{
					Attachments: volumeResult.Attachments,
					Bootable:    volumeResult.Bootable,
					CallbackUrl: volumeResult.CallbackUrl,
					CreatedAt:   volumeResult.CreatedAt,
					Description: volumeResult.Description,
					Environment: volumeResult.Environment,
					Id:          volumeResult.Id,
					ImageId:     volumeResult.ImageId,
					Name:        volumeResult.Name,
					OsImage:     nil, // This field doesn't exist in VolumesFields
					Size:        volumeResult.Size,
					Status:      volumeResult.Status,
					UpdatedAt:   volumeResult.UpdatedAt,
					VolumeType:  volumeResult.VolumeType,
				}
				return true, nil  // done!
			case "error", "deleted", "deleteed":
				return false, fmt.Errorf("Volume in error or deleted state: %s", status)
			default:
				return false, fmt.Errorf("Unexpected volume status: %s", status)
			}
		},
	)

	// TODO: doesn't save resource info in state
	if err != nil {
		resp.Diagnostics.AddError(
			"Waiting for state change error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	// TODO: doesn't save resource info in state
	if volumeModel.Status == nil {
		resp.Diagnostics.AddError(
			"Failed creating volume: status is nil",
			"Volume status is nil",
		)
		return
	}

	if strings.ToLower(*volumeModel.Status) != "available" {
		resp.Diagnostics.AddError(
			"Failed creating volume: status %s",
			*volumeModel.Status,
		)
		return
	}

	data := r.ApiToModel(ctx, &resp.Diagnostics, r.convertVolumeFieldsToVolumesFields(callResult.Volume))
	r.MergeData(&data, &dataOld)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVolume) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var dataOld resource_core_volume.CoreVolumeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &dataOld)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(dataOld.Id.ValueInt64())

	// Perform a Read operation to get all volumes
	searchResult, err := r.client.ListVolumesWithResponse(ctx, func() *volume.ListVolumesParams {
		return &volume.ListVolumesParams{
			Page:     nil,
			PageSize: nil,
			Search:   nil,
		}
	}())
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

	searchCallResult := searchResult.JSON200.Volumes
	if searchCallResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	// Find the volume with the matching ID
	var volumeResult *volume.VolumesFields
	for _, row := range *searchCallResult {
		if *row.Id == id {
			volumeResult = &row
			break
		}
	}

	// Nothing found
	if volumeResult == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data := r.ApiToModel(ctx, &resp.Diagnostics, volumeResult)
	r.MergeData(&data, &dataOld)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVolume) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_core_volume.CoreVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt64())

	// Perform a Read operation to get all volumes
	result, err := r.client.ListVolumesWithResponse(ctx, func() *volume.ListVolumesParams {
		return &volume.ListVolumesParams{
			Page:     nil,
			PageSize: nil,
			Search:   nil,
		}
	}())
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

	callResult := result.JSON200.Volumes
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	// Find the volume with the matching ID
	var volumeResult *volume.VolumesFields
	for _, row := range *callResult {
		if *row.Id == id {
			volumeResult = &row
			break
		}
	}

	// Nothing found
	if volumeResult == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Now proceed with the Delete operation using the ID
	resultDelete, err := r.client.DeleteVolumeWithResponse(ctx, id)
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

	err = r.WaitForResult(
		ctx,
		3*time.Second,
		120*time.Second,
		func(ctx context.Context) (bool, error) {
			result, err := r.client.ListVolumesWithResponse(ctx, func() *volume.ListVolumesParams {
				return &volume.ListVolumesParams{
					Page:     nil,
					PageSize: nil,
					Search:   nil,
				}
			}())
			if err != nil {
				return false, err
			}

			if result.JSON200 == nil {
				return false, fmt.Errorf("Wrong API result: %d", result.StatusCode())
			}

			var volumeResult *volume.VolumesFields
			for _, row := range *result.JSON200.Volumes {
				if *row.Id == id {
					volumeResult = &row
					break
				}
			}

			return volumeResult == nil, nil
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Waiting for state change error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ResourceCoreVolume) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ResourceCoreVolume) MergeData(
	data *resource_core_volume.CoreVolumeModel,
	dataOld *resource_core_volume.CoreVolumeModel,
) {
	data.EnvironmentName = dataOld.EnvironmentName
}

func (r *ResourceCoreVolume) convertVolumeFieldsToVolumesFields(vf *volume.VolumeFields) *volume.VolumesFields {
	if vf == nil {
		return nil
	}
	return &volume.VolumesFields{
		Attachments: vf.Attachments,
		Bootable:    vf.Bootable,
		CallbackUrl: vf.CallbackUrl,
		CreatedAt:   vf.CreatedAt,
		Description: vf.Description,
		Environment: vf.Environment,
		Id:          vf.Id,
		ImageId:     vf.ImageId,
		Name:        vf.Name,
		Size:        vf.Size,
		Status:      vf.Status,
		UpdatedAt:   vf.UpdatedAt,
		VolumeType:  vf.VolumeType,
		// Note: OsImage field from VolumeFields is not copied as it doesn't exist in VolumesFields
	}
}

func (r *ResourceCoreVolume) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *volume.VolumesFields,
) resource_core_volume.CoreVolumeModel {
	return resource_core_volume.CoreVolumeModel{
		Bootable:        types.BoolPointerValue(response.Bootable),
		CallbackUrl:     types.StringPointerValue(response.CallbackUrl),
		Description:     types.StringPointerValue(response.Description),
		Environment:     r.MapEnvironment(ctx, diags, *response.Environment),
		EnvironmentName: types.StringNull(),
		Id: func() types.Int64 {
			if response.Id == nil {
				return types.Int64Null()
			}
			return types.Int64Value(int64(*response.Id))
		}(),
		ImageId: func() types.Int64 {
			if response.ImageId == nil {
				return types.Int64Null()
			}
			return types.Int64Value(int64(*response.ImageId))
		}(),
		Name: types.StringPointerValue(response.Name),
		Size: func() types.Int64 {
			if response.Size == nil {
				return types.Int64Null()
			}
			return types.Int64Value(int64(*response.Size))
		}(),
		Status:     types.StringPointerValue(response.Status),
		VolumeType: types.StringPointerValue(response.VolumeType),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
		UpdatedAt: func() types.String {
			if response.UpdatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.UpdatedAt.String())
		}(),
	}
}

func (r *ResourceCoreVolume) MapEnvironment(
	ctx context.Context,
	diags *diag.Diagnostics,
	data volume.EnvironmentFieldsForVolume,
) resource_core_volume.EnvironmentValue {
	model, diagnostic := resource_core_volume.NewEnvironmentValue(
		resource_core_volume.EnvironmentValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"name": types.StringValue(*data.Name),
		},
	)
	diags.Append(diagnostic...)

	return model
}
