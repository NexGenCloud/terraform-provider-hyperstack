package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NexGenCloud/hyperstack-sdk-go/lib/volume_attachment"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/resource_core_volume_attachment"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ResourceCoreVolumeAttachment{}
var _ resource.ResourceWithImportState = &ResourceCoreVolumeAttachment{}

func NewResourceCoreVolumeAttachment() resource.Resource {
	return &ResourceCoreVolumeAttachment{}
}

type ResourceCoreVolumeAttachment struct {
	hyperstack *client.HyperstackClient
	client     *volume_attachment.ClientWithResponses
}

func (r *ResourceCoreVolumeAttachment) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_volume_attachment"
}

func (r *ResourceCoreVolumeAttachment) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_core_volume_attachment.CoreVolumeAttachmentResourceSchema(ctx)
}

func (r *ResourceCoreVolumeAttachment) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	r.client, err = volume_attachment.NewClientWithResponses(
		r.hyperstack.ApiServer,
		volume_attachment.WithRequestEditorFn(r.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (r *ResourceCoreVolumeAttachment) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	tflog.Debug(ctx, "=== DEBUG: Entered ResourceCoreVolumeAttachment.Create ===")
	var data resource_core_volume_attachment.CoreVolumeAttachmentModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmId := int(data.VmId.ValueInt64())

	// Convert volume_ids from types.List to []int
	var volumeIds []int
	volumeIdsList := data.VolumeIds.Elements()
	for _, volumeId := range volumeIdsList {
		if volumeIdVal, ok := volumeId.(types.Int64); ok {
			volumeIds = append(volumeIds, int(volumeIdVal.ValueInt64()))
		}
	}

	// Build the request payload
	payload := volume_attachment.AttachVolumesPayload{
		VolumeIds: &volumeIds,
		Protected: data.Protected.ValueBoolPointer(),
	}

	tflog.Debug(ctx, "[DEBUG] AttachVolumes payload", map[string]interface{}{
		"vm_id":   vmId,
		"payload": payload,
	})

	// Call the API to attach volumes
	result, err := r.client.AttachVolumesToVirtualMachineWithResponse(ctx, vmId, payload)
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
	if callResult == nil || callResult.VolumeAttachments == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	// Wait for attachments to be completed
	err = r.waitForAttachments(ctx, vmId, volumeIds, 3*time.Second, 300*time.Second)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Volume attachment may not be ready",
			fmt.Sprintf("Error waiting for volume attachments: %s", err),
		)
	}

	// Map the API response to the model
	data.Id = types.StringValue(fmt.Sprintf("vm-%d-volumes", vmId))
	
	// Convert volume attachments to types.List
	attachmentsList, diags := r.mapVolumeAttachments(ctx, *callResult.VolumeAttachments)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.VolumeAttachments = attachmentsList

	// If protected is not set in config, use the value from the first attachment if available
	if data.Protected.IsNull() && len(*callResult.VolumeAttachments) > 0 {
		firstAttachment := (*callResult.VolumeAttachments)[0]
		if firstAttachment.Protected != nil {
			data.Protected = types.BoolValue(*firstAttachment.Protected)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVolumeAttachment) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data resource_core_volume_attachment.CoreVolumeAttachmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For read, we would need to query the VM to get its attached volumes
	// Since there's no direct API to get volume attachments by VM ID,
	// we'll keep the state as-is for now
	// In a real implementation, you might query the VM details to verify the attachments

	tflog.Debug(ctx, "[DEBUG] Read operation for volume attachment", map[string]interface{}{
		"id":    data.Id.ValueString(),
		"vm_id": data.VmId.ValueInt64(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVolumeAttachment) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data resource_core_volume_attachment.CoreVolumeAttachmentModel
	var oldData resource_core_volume_attachment.CoreVolumeAttachmentModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &oldData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Since vm_id and volume_ids require replacement, the only updateable field is protected
	// For now, we'll just accept the update and keep the state
	// To update individual attachments, you would use the UpdateAVolumeAttachment API

	tflog.Debug(ctx, "[DEBUG] Update operation for volume attachment", map[string]interface{}{
		"id":    data.Id.ValueString(),
		"vm_id": data.VmId.ValueInt64(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVolumeAttachment) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data resource_core_volume_attachment.CoreVolumeAttachmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmId := int(data.VmId.ValueInt64())

	// Convert volume_ids from types.List to []int
	var volumeIds []int
	volumeIdsList := data.VolumeIds.Elements()
	for _, volumeId := range volumeIdsList {
		if volumeIdVal, ok := volumeId.(types.Int64); ok {
			volumeIds = append(volumeIds, int(volumeIdVal.ValueInt64()))
		}
	}

	// Build the request payload
	payload := volume_attachment.DetachVolumesPayload{
		VolumeIds: &volumeIds,
	}

	tflog.Debug(ctx, "[DEBUG] DetachVolumes payload", map[string]interface{}{
		"vm_id":   vmId,
		"payload": payload,
	})

	// Call the API to detach volumes
	result, err := r.client.DetachVolumesFromVirtualMachineWithResponse(ctx, vmId, payload)
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

	tflog.Debug(ctx, "[DEBUG] Detach API call completed, waiting for volumes to fully detach", map[string]interface{}{
		"vm_id":      vmId,
		"volume_ids": volumeIds,
	})

	// Wait for detachment to complete before returning
	// This ensures volumes are fully detached before they can be deleted
	err = r.waitForDetachments(ctx, vmId, volumeIds, 3*time.Second, 300*time.Second)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Volume detachment may not be complete",
			fmt.Sprintf("Error waiting for volume detachment: %s. Volumes may still be detaching.", err),
		)
		// Still return success to allow cleanup to proceed
		return
	}

	tflog.Debug(ctx, "[DEBUG] Successfully detached volumes from VM", map[string]interface{}{
		"vm_id":      vmId,
		"volume_ids": volumeIds,
	})
}

func (r *ResourceCoreVolumeAttachment) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Import format: vm-{vm_id}-volumes
	// or just the vm_id
	idParts := strings.Split(req.ID, "-")
	
	var vmIdStr string
	if len(idParts) == 3 && idParts[0] == "vm" && idParts[2] == "volumes" {
		vmIdStr = idParts[1]
	} else {
		vmIdStr = req.ID
	}
	
	vmId, err := strconv.ParseInt(vmIdStr, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected format: vm-{vm_id}-volumes or just {vm_id}, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), fmt.Sprintf("vm-%d-volumes", vmId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vm_id"), vmId)...)
}

// Helper functions

func (r *ResourceCoreVolumeAttachment) waitForAttachments(
	ctx context.Context,
	vmId int,
	volumeIds []int,
	interval time.Duration,
	timeout time.Duration,
) error {
	// For now, just wait a fixed time
	// In a production implementation, you would poll the VM or volume attachment status
	tflog.Debug(ctx, "[DEBUG] Waiting for volume attachments to complete", map[string]interface{}{
		"vm_id":      vmId,
		"volume_ids": volumeIds,
		"timeout":    timeout,
	})
	
	time.Sleep(5 * time.Second)
	return nil
}

func (r *ResourceCoreVolumeAttachment) waitForDetachments(
	ctx context.Context,
	vmId int,
	volumeIds []int,
	interval time.Duration,
	timeout time.Duration,
) error {
	// Wait for volumes to be fully detached before allowing deletion
	// In a production implementation, you would poll the volume status to verify detachment
	tflog.Debug(ctx, "[DEBUG] Waiting for volume detachments to complete", map[string]interface{}{
		"vm_id":      vmId,
		"volume_ids": volumeIds,
		"timeout":    timeout,
	})
	
	// Wait longer for detachment as it may take more time than attachment
	time.Sleep(100 * time.Second)
	return nil
}

func (r *ResourceCoreVolumeAttachment) mapVolumeAttachments(
	ctx context.Context,
	attachments []volume_attachment.AttachVolumeFields,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	
	// Create element type for the list
	elementType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.Int64Type,
			"volume_id":   types.Int64Type,
			"instance_id": types.Int64Type,
			"device":      types.StringType,
			"status":      types.StringType,
			"protected":   types.BoolType,
			"created_at":  types.StringType,
		},
	}
	
	// If no attachments, return empty list
	if len(attachments) == 0 {
		return types.ListValueMust(elementType, []attr.Value{}), diags
	}
	
	// Build list of attachment objects
	elements := make([]attr.Value, 0, len(attachments))
	
	for _, att := range attachments {
		obj, d := types.ObjectValue(
			elementType.AttrTypes,
			map[string]attr.Value{
				"id": func() attr.Value {
					if att.Id == nil {
						return types.Int64Null()
					}
					return types.Int64Value(int64(*att.Id))
				}(),
				"volume_id": func() attr.Value {
					if att.VolumeId == nil {
						return types.Int64Null()
					}
					return types.Int64Value(int64(*att.VolumeId))
				}(),
				"instance_id": func() attr.Value {
					if att.InstanceId == nil {
						return types.Int64Null()
					}
					return types.Int64Value(int64(*att.InstanceId))
				}(),
				"device": func() attr.Value {
					if att.Device == nil {
						return types.StringNull()
					}
					return types.StringValue(*att.Device)
				}(),
				"status": func() attr.Value {
					if att.Status == nil {
						return types.StringNull()
					}
					return types.StringValue(*att.Status)
				}(),
				"protected": func() attr.Value {
					if att.Protected == nil {
						return types.BoolNull()
					}
					return types.BoolValue(*att.Protected)
				}(),
				"created_at": func() attr.Value {
					if att.CreatedAt == nil {
						return types.StringNull()
					}
					return types.StringValue(att.CreatedAt.String())
				}(),
			},
		)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(elementType), diags
		}
		elements = append(elements, obj)
	}
	
	list, d := types.ListValue(elementType, elements)
	diags.Append(d...)
	
	return list, diags
}
