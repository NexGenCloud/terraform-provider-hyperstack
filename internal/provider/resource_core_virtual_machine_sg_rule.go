package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/virtual_machine"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/resource_core_virtual_machine_sg_rule"
)

var _ resource.Resource = &ResourceCoreVirtualMachineSgRule{}
var _ resource.ResourceWithImportState = &ResourceCoreVirtualMachineSgRule{}

func NewResourceCoreVirtualMachineSgRule() resource.Resource {
	return &ResourceCoreVirtualMachineSgRule{}
}

type ResourceCoreVirtualMachineSgRule struct {
	hyperstack *client.HyperstackClient
	client     *virtual_machine.ClientWithResponses
}

func (r *ResourceCoreVirtualMachineSgRule) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_virtual_machine_sg_rule"
}

func (r *ResourceCoreVirtualMachineSgRule) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_core_virtual_machine_sg_rule.CoreVirtualMachineSgRuleResourceSchema(ctx)
}

func (r *ResourceCoreVirtualMachineSgRule) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceCoreVirtualMachineSgRule) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data resource_core_virtual_machine_sg_rule.CoreVirtualMachineSgRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vmId := int(data.VirtualMachineId.ValueInt64())

	result, err := r.client.AddSecurityRuleWithResponse(ctx, vmId, func() virtual_machine.AddSecurityRuleJSONRequestBody {
		return virtual_machine.AddSecurityRuleJSONRequestBody{
			Direction:      data.Direction.String(),
			Ethertype:      data.Ethertype.String(),
			Protocol:       data.Protocol.String(),
			RemoteIpPrefix: data.RemoteIpPrefix.String(),
			PortRangeMin: func() *int {
				if data.PortRangeMin.IsNull() {
					return nil
				}

				val := int(data.PortRangeMin.ValueInt64())
				return &val
			}(),
			PortRangeMax: func() *int {
				if data.PortRangeMax.IsNull() {
					return nil
				}

				val := int(data.PortRangeMax.ValueInt64())
				return &val
			}(),
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

	callResult := result.JSON200.SecurityRule
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	data = r.ApiToModel(ctx, &resp.Diagnostics, callResult, vmId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVirtualMachineSgRule) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data resource_core_virtual_machine_sg_rule.CoreVirtualMachineSgRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sgId := int(data.Id.ValueInt64())
	vmId := int(data.VirtualMachineId.ValueInt64())
	result, err := r.client.GetAnInstanceDetailsWithResponse(ctx, vmId)
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

	sg := func() *virtual_machine.SecurityGroupRuleFields {
		for _, sgRow := range *result.JSON200.Instance.SecurityRules {
			if *sgRow.Id == sgId {
				return &virtual_machine.SecurityGroupRuleFields{
					CreatedAt:      sgRow.CreatedAt,
					Direction:      sgRow.Direction,
					Ethertype:      sgRow.Ethertype,
					Id:             sgRow.Id,
					PortRangeMin:   sgRow.PortRangeMin,
					PortRangeMax:   sgRow.PortRangeMax,
					Protocol:       sgRow.Protocol,
					RemoteIpPrefix: sgRow.RemoteIpPrefix,
					Status:         sgRow.Status,
				}
			}
		}

		return nil
	}()

	if sg == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data = r.ApiToModel(ctx, &resp.Diagnostics, sg, vmId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreVirtualMachineSgRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported for VM SG rules resources",
		"",
	)
	return
}

func (r *ResourceCoreVirtualMachineSgRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_core_virtual_machine_sg_rule.CoreVirtualMachineSgRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sgId := int(data.Id.ValueInt64())
	vmId := int(data.VirtualMachineId.ValueInt64())

	result, err := r.client.DeleteASecurityRuleWithResponse(ctx, vmId, sgId)
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

func (r *ResourceCoreVirtualMachineSgRule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ResourceCoreVirtualMachineSgRule) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *virtual_machine.SecurityGroupRuleFields,
	vmId int,
) resource_core_virtual_machine_sg_rule.CoreVirtualMachineSgRuleModel {
	return resource_core_virtual_machine_sg_rule.CoreVirtualMachineSgRuleModel{
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
		Direction: func() types.String {
			if response.Direction == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Direction)
		}(),
		Ethertype: func() types.String {
			if response.Ethertype == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Ethertype)
		}(),
		Id: func() types.Int64 {
			if response.Id == nil {
				return types.Int64Null()
			}
			return types.Int64Value(int64(*response.Id))
		}(),
		PortRangeMax: func() types.Int64 {
			if response.PortRangeMax == nil {
				return types.Int64Null()
			}
			return types.Int64Value(int64(*response.PortRangeMax))
		}(),
		PortRangeMin: func() types.Int64 {
			if response.PortRangeMin == nil {
				return types.Int64Null()
			}
			return types.Int64Value(int64(*response.PortRangeMin))
		}(),
		Protocol: func() types.String {
			if response.Protocol == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Protocol)
		}(),
		RemoteIpPrefix: func() types.String {
			if response.RemoteIpPrefix == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.RemoteIpPrefix)
		}(),
		Status: func() types.String {
			if response.Status == nil {
				return types.StringNull()
			}
			return types.StringValue(*response.Status)
		}(),
		VirtualMachineId: func() types.Int64 {
			return types.Int64Value(int64(vmId))
		}(),
	}
}
