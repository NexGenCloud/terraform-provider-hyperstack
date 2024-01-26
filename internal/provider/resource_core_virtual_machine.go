package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/virtual_machine"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/resource_core_virtual_machine"
	"math/big"
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
	var data resource_core_virtual_machine.CoreVirtualMachineModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateInstancesWithResponse(ctx, func() virtual_machine.CreateInstancesJSONRequestBody {
		return virtual_machine.CreateInstancesJSONRequestBody{
			Name:                 data.Name.ValueString(),
			EnvironmentName:      data.EnvironmentName.ValueString(),
			ImageName:            data.ImageName.ValueStringPointer(),
			VolumeName:           data.VolumeName.ValueStringPointer(),
			CreateBootableVolume: data.CreateBootableVolume.ValueBoolPointer(),
			FlavorName:           data.FlavorName.ValueStringPointer(),
			Flavor: func() *virtual_machine.FlavorObjectFields {
				if data.Flavor.IsNull() {
					return nil
				}
				return &virtual_machine.FlavorObjectFields{
					Cpu: func() *int {
						if data.Flavor.Cpu.IsNull() {
							return nil
						}

						val := int(data.Flavor.Cpu.ValueInt64())
						return &val
					}(),
					Disk: func() *int {
						if data.Flavor.Disk.IsNull() {
							return nil
						}

						val := int(data.Flavor.Disk.ValueInt64())
						return &val
					}(),
					Gpu: data.Flavor.Gpu.ValueStringPointer(),
					GpuCount: func() *int {
						if data.Flavor.GpuCount.IsNull() {
							return nil
						}

						val := int(data.Flavor.GpuCount.ValueInt64())
						return &val
					}(),
					Ram: func() *float32 {
						if data.Flavor.Ram.IsNull() {
							return nil
						}

						val, _ := data.Flavor.Ram.ValueBigFloat().Float32()
						return &val
					}(),
				}
			}(),
			KeyName:          data.KeyName.ValueString(),
			UserData:         data.UserData.ValueStringPointer(),
			CallbackUrl:      data.CallbackUrl.ValueStringPointer(),
			AssignFloatingIp: data.AssignFloatingIp.ValueBoolPointer(),
			Profile: func() *virtual_machine.ProfileObjectFields {
				if data.Profile.IsNull() {
					return nil
				}
				return &virtual_machine.ProfileObjectFields{
					Name:        data.Profile.Name.ValueString(),
					Description: data.Profile.Description.ValueStringPointer(),
				}
			}(),
			Count: 0,
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
			"No data in API result",
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
	data = r.ApiToModel(ctx, &resp.Diagnostics, &list[0])
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVirtualMachine) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data resource_core_virtual_machine.CoreVirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GetAnInstanceDetailsWithResponse(ctx, int(data.Id.ValueInt64()))
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
			"No data in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	data = r.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVirtualMachine) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported for VM resources",
		"",
	)
	return
}

func (r *ResourceCoreVirtualMachine) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_core_virtual_machine.CoreVirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt64())

	result, err := r.client.DeleteAnInstanceWithResponse(ctx, id)
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

func (r *ResourceCoreVirtualMachine) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ResourceCoreVirtualMachine) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *virtual_machine.InstanceAdminFields,
) resource_core_virtual_machine.CoreVirtualMachineModel {
	return resource_core_virtual_machine.CoreVirtualMachineModel{
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
		FixedIp: func() types.String {
			if response.FixedIp == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.FixedIp)
		}(),
		FloatingIp: func() types.String {
			if response.FloatingIp == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.FloatingIp)
		}(),
		FloatingIpStatus: func() types.String {
			if response.FloatingIpStatus == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.FloatingIpStatus)
		}(),
		OpenstackId: func() types.String {
			if response.OpenstackId == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.OpenstackId)
		}(),
		Environment:       r.MapEnvironment(ctx, diags, *response.Environment),
		Image:             r.MapImage(ctx, diags, *response.Image),
		Flavor:            r.MapFlavor(ctx, diags, *response.Flavor),
		VolumeAttachments: r.MapVolumeAttachments(ctx, diags, *response.VolumeAttachments),
		SecurityRules:     r.MapSecurityRules(ctx, diags, *response.SecurityRules),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),

		AssignFloatingIp:     types.Bool{},
		CallbackUrl:          types.String{},
		CreateBootableVolume: types.Bool{},
		EnvironmentName:      types.String{},
		FlavorName:           types.String{},
		ImageName:            types.String{},
		KeyName:              types.String{},
		Keypair:              resource_core_virtual_machine.KeypairValue{},
		Profile:              resource_core_virtual_machine.ProfileValue{},
		UserData:             types.String{},
		VolumeName:           types.String{},
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
) attr.Value {
	model, diagnostic := resource_core_virtual_machine.NewKeypairValue(
		resource_core_virtual_machine.KeypairValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"name": types.StringValue(*data.Name),
		},
	)
	diags.Append(diagnostic...)

	result, diagnostic := model.ToObjectValue(ctx)
	diags.Append(diagnostic...)

	return result
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
						"device": types.StringValue(*row.Device),
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
			"id":          types.Int64Value(int64(*data.Id)),
			"name":        types.StringValue(*data.Name),
			"description": types.StringValue(*data.Description),
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
) types.List {
	model, diagnostic := types.ListValue(
		resource_core_virtual_machine.SecurityRulesValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := resource_core_virtual_machine.NewSecurityRulesValue(
					resource_core_virtual_machine.SecurityRulesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id":        types.Int64Value(int64(*row.Id)),
						"direction": types.StringValue(*row.Direction),
						"protocol":  types.StringValue(*row.Protocol),
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
