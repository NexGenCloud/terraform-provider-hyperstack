package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/virtual_machine"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/resource_core_virtual_machine"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"math/big"
	"time"
)

var _ resource.Resource = &ResourceCoreVirtualMachine{}
var _ resource.ResourceWithImportState = &ResourceCoreVirtualMachine{}

func NewResourceCoreVirtualMachine() resource.Resource {
	return &ResourceCoreVirtualMachine{}
}

type ResourceCoreVirtualMachine struct {
	hyperstack *client.HyperstackClient
	client     *virtual_machine.ClientWithResponses
}

func (r *ResourceCoreVirtualMachine) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_virtual_machine"
}

func (r *ResourceCoreVirtualMachine) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_core_virtual_machine.CoreVirtualMachineResourceSchema(ctx)
}

func (r *ResourceCoreVirtualMachine) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	r.client, err = virtual_machine.NewClientWithResponses(
		r.hyperstack.ApiServer,
		virtual_machine.WithRequestEditorFn(r.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (r *ResourceCoreVirtualMachine) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var dataOld resource_core_virtual_machine.CoreVirtualMachineModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &dataOld)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateVirtualMachineWithResponse(ctx, func() virtual_machine.CreateVirtualMachineJSONRequestBody {
		return virtual_machine.CreateVirtualMachineJSONRequestBody{
			Name:                 dataOld.Name.ValueString(),
			EnvironmentName:      dataOld.EnvironmentName.ValueString(),
			ImageName:            dataOld.ImageName.ValueStringPointer(),
			VolumeName:           dataOld.VolumeName.ValueStringPointer(),
			CreateBootableVolume: dataOld.CreateBootableVolume.ValueBoolPointer(),
			FlavorName:           dataOld.FlavorName.ValueString(),
			// TODO: disable setting here
			Flavor: func() *virtual_machine.FlavorObjectFields {
				if dataOld.Flavor.IsNull() {
					return nil
				}
				return &virtual_machine.FlavorObjectFields{
					Cpu: func() *int {
						if dataOld.Flavor.Cpu.IsNull() {
							return nil
						}

						val := int(dataOld.Flavor.Cpu.ValueInt64())
						return &val
					}(),
					Disk: func() *int {
						if dataOld.Flavor.Disk.IsNull() {
							return nil
						}

						val := int(dataOld.Flavor.Disk.ValueInt64())
						return &val
					}(),
					Gpu: dataOld.Flavor.Gpu.ValueStringPointer(),
					GpuCount: func() *int {
						if dataOld.Flavor.GpuCount.IsNull() {
							return nil
						}

						val := int(dataOld.Flavor.GpuCount.ValueInt64())
						return &val
					}(),
					Ram: func() *float32 {
						if dataOld.Flavor.Ram.IsNull() {
							return nil
						}

						val, _ := dataOld.Flavor.Ram.ValueBigFloat().Float32()
						return &val
					}(),
				}
			}(),
			KeyName:          dataOld.KeyName.ValueString(),
			UserData:         dataOld.UserData.ValueStringPointer(),
			CallbackUrl:      dataOld.CallbackUrl.ValueStringPointer(),
			AssignFloatingIp: dataOld.AssignFloatingIp.ValueBoolPointer(),
			// TODO: disable setting here
			Profile: func() *virtual_machine.ProfileObjectFields {
				if dataOld.Profile.IsNull() {
					return nil
				}
				prof := dataOld.Profile.Elements()[0].(resource_core_virtual_machine.ProfileValue)
				return &virtual_machine.ProfileObjectFields{
					Name:        prof.Name.ValueString(),
					Description: prof.Description.ValueStringPointer(),
				}
			}(),
			Count: 1,
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

	callResult := result.JSON200.Instances
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No dataOld in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	if len(*callResult) != 1 {
		resp.Diagnostics.AddError(
			"Wrong instance count",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	list := *callResult
	instanceModel := &list[0]

	id := *instanceModel.Id
	err = r.WaitForResult(
		ctx,
		3*time.Second,
		300*time.Second,
		func(ctx context.Context) (bool, error) {
			result, err := r.client.RetrieveVirtualMachineDetailsWithResponse(ctx, id)
			if err != nil {
				return false, err
			}

			if result.JSON200 == nil {
				return false, fmt.Errorf("Wrong API result: %s", result.StatusCode())
			}

			status := *result.JSON200.Instance.Status
			if status == "CREATING" || status == "BUILD" {
				return false, nil
			}

			instanceModel = result.JSON200.Instance
			return true, nil
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
	if *instanceModel.Status == "ERROR" {
		resp.Diagnostics.AddError(
			"Failed creating instance: status %s",
			*instanceModel.Status,
		)
		return
	}

	data := r.ApiToModel(ctx, &resp.Diagnostics, instanceModel)
	r.MergeData(&data, &dataOld)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVirtualMachine) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var dataOld resource_core_virtual_machine.CoreVirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataOld)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.RetrieveVirtualMachineDetailsWithResponse(ctx, int(dataOld.Id.ValueInt64()))
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

	callResult := result.JSON200.Instance
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No dataOld in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	data := r.ApiToModel(ctx, &resp.Diagnostics, callResult)
	r.MergeData(&data, &dataOld)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// TODO: deduplicate
func (r *ResourceCoreVirtualMachine) WaitForResult(
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

func (r *ResourceCoreVirtualMachine) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var dataOld resource_core_virtual_machine.CoreVirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataOld)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data resource_core_virtual_machine.CoreVirtualMachineModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	labels := make([]string, len(dataOld.Labels.Elements()))
	dataOld.Labels.ElementsAs(ctx, labels, false)

	// TODO: finish implementing
	resp.Diagnostics.AddError(
		"Update not supported for VM resources",
		fmt.Sprintf("%+v", labels),
	)
	r.MergeData(&data, &dataOld)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVirtualMachine) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_core_virtual_machine.CoreVirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt64())

	result, err := r.client.DeleteVirtualMachineWithResponse(ctx, id)
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

	err = r.WaitForResult(
		ctx,
		3*time.Second,
		120*time.Second,
		func(ctx context.Context) (bool, error) {
			result, err := r.client.RetrieveVirtualMachineDetailsWithResponse(ctx, id)
			if err != nil {
				return false, err
			}

			return result.JSON404 != nil, nil
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

func (r *ResourceCoreVirtualMachine) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ResourceCoreVirtualMachine) MergeData(
	data *resource_core_virtual_machine.CoreVirtualMachineModel,
	dataOld *resource_core_virtual_machine.CoreVirtualMachineModel,
) {
	// Assign all values that are available only during creation stage
	data.AssignFloatingIp = dataOld.AssignFloatingIp
	data.CreateBootableVolume = dataOld.CreateBootableVolume
	data.EnvironmentName = dataOld.EnvironmentName
	data.FlavorName = dataOld.FlavorName
	data.ImageName = dataOld.ImageName
	data.KeyName = dataOld.KeyName
	data.Profile = dataOld.Profile

	// TODO: not implemented as there is no field in output
	if !dataOld.CallbackUrl.IsUnknown() {
		data.CallbackUrl = dataOld.CallbackUrl
	}
	if !dataOld.UserData.IsUnknown() {
		data.UserData = dataOld.UserData
	}

	// TODO: verify
	//if !dataOld.FloatingIp.IsUnknown() {
	//	data.FloatingIp = dataOld.FloatingIp
	//}
	//if !dataOld.VolumeName.IsUnknown() {
	//	data.VolumeName = dataOld.VolumeName
	//}
}

func (r *ResourceCoreVirtualMachine) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *virtual_machine.InstanceFields,
) resource_core_virtual_machine.CoreVirtualMachineModel {
	return resource_core_virtual_machine.CoreVirtualMachineModel{
		Id: func() types.Int64 {
			if response.Id == nil {
				return types.Int64Null()
			}
			return types.Int64Value(int64(*response.Id))
		}(),
		Name: types.StringPointerValue(response.Name),
		Status: func() types.String {
			if response.Status == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Status)
		}(),
		PowerState: func() types.String {
			if response.PowerState == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.PowerState)
		}(),
		VmState: func() types.String {
			if response.VmState == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.VmState)
		}(),
		FixedIp:    types.StringPointerValue(response.FixedIp),
		FloatingIp: types.StringPointerValue(response.FloatingIp),
		FloatingIpStatus: func() types.String {
			if response.FloatingIpStatus == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.FloatingIpStatus)
		}(),
		Keypair:           r.MapKeypair(ctx, diags, *response.Keypair),
		Environment:       r.MapEnvironment(ctx, diags, *response.Environment),
		Image:             r.MapImage(ctx, diags, *response.Image),
		Labels:            r.MapLabels(ctx, diags, *response.Labels),
		Flavor:            r.MapFlavor(ctx, diags, *response.Flavor),
		VolumeAttachments: r.MapVolumeAttachments(ctx, diags, *response.VolumeAttachments),
		SecurityRules:     r.MapSecurityRules(ctx, diags, *response.SecurityRules, int64(*response.Id)),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
		Locked: types.BoolPointerValue(response.Locked),
		Os:     types.StringPointerValue(response.Os),
		// Doesn't make sense for reads (only used during creation)
		AssignFloatingIp:     types.BoolNull(),
		CreateBootableVolume: types.BoolNull(),
		EnvironmentName:      types.StringNull(),
		FlavorName:           types.StringNull(),
		ImageName:            types.StringNull(),
		KeyName:              types.StringNull(),
		Profile:              types.ListNull(resource_core_virtual_machine.ProfileType{}),
		VolumeName:           types.StringNull(),
		// TODO: not implemented as there is no field in output
		CallbackUrl: types.StringUnknown(),
		UserData:    types.StringUnknown(),
	}
}

func (r *ResourceCoreVirtualMachine) MapEnvironment(
	ctx context.Context,
	diags *diag.Diagnostics,
	data virtual_machine.InstanceEnvironmentFields,
) resource_core_virtual_machine.EnvironmentValue {
	model, diagnostic := resource_core_virtual_machine.NewEnvironmentValue(
		resource_core_virtual_machine.EnvironmentValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":     types.Int64Value(int64(*data.Id)),
			"name":   types.StringValue(*data.Name),
			"org_id": types.Int64Value(int64(*data.OrgId)),
			"region": types.StringValue(*data.Region),
		},
	)
	diags.Append(diagnostic...)

	return model
}

func (r *ResourceCoreVirtualMachine) MapImage(
	ctx context.Context,
	diags *diag.Diagnostics,
	data virtual_machine.InstanceImageFields,
) resource_core_virtual_machine.ImageValue {
	model, diagnostic := resource_core_virtual_machine.NewImageValue(
		resource_core_virtual_machine.ImageValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"name": types.StringValue(*data.Name),
		},
	)
	diags.Append(diagnostic...)

	return model
}

func (r *ResourceCoreVirtualMachine) MapFlavor(
	ctx context.Context,
	diags *diag.Diagnostics,
	data virtual_machine.InstanceFlavorFields,
) resource_core_virtual_machine.FlavorValue {
	model, diagnostic := resource_core_virtual_machine.NewFlavorValue(
		resource_core_virtual_machine.FlavorValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":        types.Int64Value(int64(*data.Id)),
			"ephemeral": types.Int64Value(int64(*data.Ephemeral)),
			"name":      types.StringValue(*data.Name),
			"cpu":       types.Int64Value(int64(*data.Cpu)),
			"ram":       types.NumberValue(big.NewFloat(float64(*data.Ram))),
			"disk":      types.Int64Value(int64(*data.Disk)),
			"gpu":       types.StringValue(*data.Gpu),
			"gpu_count": types.Int64Value(int64(*data.GpuCount)),
		},
	)
	diags.Append(diagnostic...)

	return model
}

func (r *ResourceCoreVirtualMachine) MapKeypair(
	ctx context.Context,
	diags *diag.Diagnostics,
	data virtual_machine.InstanceKeypairFields,
) resource_core_virtual_machine.KeypairValue {
	model, diagnostic := resource_core_virtual_machine.NewKeypairValue(
		resource_core_virtual_machine.KeypairValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"name": types.StringValue(*data.Name),
		},
	)
	diags.Append(diagnostic...)

	return model
}

func (r *ResourceCoreVirtualMachine) MapVolumeAttachments(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []virtual_machine.VolumeAttachmentFields,
) types.List {
	model, diagnostic := types.ListValue(
		resource_core_virtual_machine.VolumeAttachmentsValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := resource_core_virtual_machine.NewVolumeAttachmentsValue(
					resource_core_virtual_machine.VolumeAttachmentsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"status": types.StringValue(*row.Status),
						"device": func() attr.Value {
							if row.Device == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.Device)
						}(),
						"created_at": func() attr.Value {
							if row.CreatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.CreatedAt.String())
						}(),
						"volume": r.MapVolume(ctx, diags, *row.Volume),
					},
				)
				diags.Append(diagnostic...)
				roles = append(roles, model)
			}
			return roles
		}(),
	)
	diags.Append(diagnostic...)
	return model
}

func (r *ResourceCoreVirtualMachine) MapVolume(
	ctx context.Context,
	diags *diag.Diagnostics,
	data virtual_machine.VolumeFieldsForInstance,
) attr.Value {
	model, diagnostic := resource_core_virtual_machine.NewVolumeValue(
		resource_core_virtual_machine.VolumeValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"id":   types.Int64Value(int64(*data.Id)),
			"name": types.StringValue(*data.Name),
			"description": func() attr.Value {
				if data.Description == nil {
					return types.StringNull()
				}
				return types.StringValue(*data.Description)
			}(),
			"volume_type": types.StringValue(*data.VolumeType),
			"size":        types.Int64Value(int64(*data.Size)),
		},
	)
	diags.Append(diagnostic...)

	result, diagnostic := model.ToObjectValue(ctx)
	diags.Append(diagnostic...)

	return result
}

func (r *ResourceCoreVirtualMachine) MapSecurityRules(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []virtual_machine.SecurityRulesFieldsForInstance,
	vmId int64,
) types.List {
	model, diagnostic := types.ListValue(
		resource_core_virtual_machine.SecurityRulesValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := resource_core_virtual_machine.NewSecurityRulesValue(
					resource_core_virtual_machine.SecurityRulesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":                 types.Int64Value(int64(*row.Id)),
						"virtual_machine_id": types.Int64Value(vmId),
						"direction":          types.StringValue(*row.Direction),
						"protocol":           types.StringValue(*row.Protocol),
						"port_range_min": func() attr.Value {
							if row.PortRangeMin == nil {
								return types.Int64Null()
							}
							return types.Int64Value(int64(*row.PortRangeMin))
						}(),
						"port_range_max": func() attr.Value {
							if row.PortRangeMax == nil {
								return types.Int64Null()
							}
							return types.Int64Value(int64(*row.PortRangeMax))
						}(),
						"ethertype":        types.StringValue(*row.Ethertype),
						"remote_ip_prefix": types.StringValue(*row.RemoteIpPrefix),
						"status":           types.StringValue(*row.Status),
						"created_at": func() attr.Value {
							if row.CreatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.CreatedAt.String())
						}(),
					},
				)
				diags.Append(diagnostic...)
				roles = append(roles, model)
			}
			return roles
		}(),
	)
	diags.Append(diagnostic...)
	return model
}

func (r *ResourceCoreVirtualMachine) MapLabels(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []string,
) types.List {
	model, diagnostic := types.ListValue(
		types.StringType,
		func() []attr.Value {
			labels := make([]attr.Value, 0)
			for _, row := range data {
				model := types.StringValue(row)
				labels = append(labels, model)
			}
			return labels
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
