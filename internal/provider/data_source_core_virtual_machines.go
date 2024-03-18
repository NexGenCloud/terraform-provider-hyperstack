package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/virtual_machine"
	"github.com/NexGenCloud/hyperstack-terraform-provider/internal/client"
	"github.com/NexGenCloud/hyperstack-terraform-provider/internal/genprovider/datasource_core_virtual_machines"
	"math/big"
)

var _ datasource.DataSource = &DataSourceCoreVirtualMachines{}

func NewDataSourceCoreVirtualMachines() datasource.DataSource {
	return &DataSourceCoreVirtualMachines{}
}

type DataSourceCoreVirtualMachines struct {
	hyperstack *client.HyperstackClient
	client     *virtual_machine.ClientWithResponses
}

func (d *DataSourceCoreVirtualMachines) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_virtual_machines"
}

func (d *DataSourceCoreVirtualMachines) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_virtual_machines.CoreVirtualMachinesDataSourceSchema(ctx)
}

func (d *DataSourceCoreVirtualMachines) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = virtual_machine.NewClientWithResponses(
		d.hyperstack.ApiServer,
		virtual_machine.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreVirtualMachines) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_virtual_machines.CoreVirtualMachinesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListInstancesWithResponse(ctx)
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

	data = d.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceCoreVirtualMachines) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]virtual_machine.InstanceAdminFields,
) datasource_core_virtual_machines.CoreVirtualMachinesModel {
	return datasource_core_virtual_machines.CoreVirtualMachinesModel{
		CoreVirtualMachines: func() types.Set {
			return d.MapInstances(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceCoreVirtualMachines) MapInstances(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []virtual_machine.InstanceAdminFields,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_core_virtual_machines.CoreVirtualMachinesValue{}.Type(ctx),
		func() []attr.Value {
			list := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_core_virtual_machines.NewCoreVirtualMachinesValue(
					datasource_core_virtual_machines.CoreVirtualMachinesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id": func() attr.Value {
							if row.Id == nil {
								return types.Int64Null()
							}
							return types.Int64Value(int64(*row.Id))
						}(),
						"name": func() attr.Value {
							if row.Name == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.Name)
						}(),
						"status": func() attr.Value {
							if row.Status == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.Status)
						}(),
						"power_state": func() attr.Value {
							if row.PowerState == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.PowerState)
						}(),
						"vm_state": func() attr.Value {
							if row.VmState == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.VmState)
						}(),
						"fixed_ip": func() attr.Value {
							if row.FixedIp == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.FixedIp)
						}(),
						"floating_ip": func() attr.Value {
							if row.FloatingIp == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.FloatingIp)
						}(),
						"floating_ip_status": func() attr.Value {
							if row.FloatingIpStatus == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.FloatingIpStatus)
						}(),
						"openstack_id": func() attr.Value {
							if row.OpenstackId == nil {
								return types.StringNull()
							}
							return types.StringValue(*row.OpenstackId)
						}(),
						"created_at": func() attr.Value {
							if row.CreatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.CreatedAt.String())
						}(),
						"environment":        d.MapEnvironment(ctx, diags, *row.Environment),
						"image":              d.MapImage(ctx, diags, *row.Image),
						"flavor":             d.MapFlavor(ctx, diags, *row.Flavor),
						"keypair":            d.MapKeypair(ctx, diags, *row.Keypair),
						"volume_attachments": d.MapVolumeAttachments(ctx, diags, *row.VolumeAttachments),
						"security_rules":     d.MapSecurityRules(ctx, diags, *row.SecurityRules),
					},
				)
				diags.Append(diagnostic...)
				list = append(list, model)
			}
			return list
		}(),
	)
	diags.Append(diagnostic...)
	return model
}

func (d *DataSourceCoreVirtualMachines) MapEnvironment(
	ctx context.Context,
	diags *diag.Diagnostics,
	data virtual_machine.InstanceEnvironmentFields,
) attr.Value {
	model, diagnostic := datasource_core_virtual_machines.NewEnvironmentValue(
		datasource_core_virtual_machines.EnvironmentValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"name":   types.StringValue(*data.Name),
			"org_id": types.Int64Value(int64(*data.OrgId)),
			"region": types.StringValue(*data.Region),
		},
	)
	diags.Append(diagnostic...)

	result, diagnostic := model.ToObjectValue(ctx)
	diags.Append(diagnostic...)

	return result
}

func (d *DataSourceCoreVirtualMachines) MapImage(
	ctx context.Context,
	diags *diag.Diagnostics,
	data virtual_machine.InstanceImageFields,
) attr.Value {
	model, diagnostic := datasource_core_virtual_machines.NewImageValue(
		datasource_core_virtual_machines.ImageValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"name": types.StringValue(*data.Name),
		},
	)
	diags.Append(diagnostic...)

	result, diagnostic := model.ToObjectValue(ctx)
	diags.Append(diagnostic...)

	return result
}

func (d *DataSourceCoreVirtualMachines) MapFlavor(
	ctx context.Context,
	diags *diag.Diagnostics,
	data virtual_machine.InstanceFlavorFields,
) attr.Value {
	model, diagnostic := datasource_core_virtual_machines.NewFlavorValue(
		datasource_core_virtual_machines.FlavorValue{}.AttributeTypes(ctx),
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

	result, diagnostic := model.ToObjectValue(ctx)
	diags.Append(diagnostic...)

	return result
}

func (d *DataSourceCoreVirtualMachines) MapKeypair(
	ctx context.Context,
	diags *diag.Diagnostics,
	data virtual_machine.InstanceKeypairFields,
) attr.Value {
	model, diagnostic := datasource_core_virtual_machines.NewKeypairValue(
		datasource_core_virtual_machines.KeypairValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"name": types.StringValue(*data.Name),
		},
	)
	diags.Append(diagnostic...)

	result, diagnostic := model.ToObjectValue(ctx)
	diags.Append(diagnostic...)

	return result
}

func (d *DataSourceCoreVirtualMachines) MapVolumeAttachments(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []virtual_machine.VolumeAttachmentFields,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_core_virtual_machines.VolumeAttachmentsValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_core_virtual_machines.NewVolumeAttachmentsValue(
					datasource_core_virtual_machines.VolumeAttachmentsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"status": types.StringValue(*row.Status),
						"device": types.StringValue(*row.Device),
						"created_at": func() attr.Value {
							if row.CreatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.CreatedAt.String())
						}(),
						"volume": d.MapVolume(ctx, diags, *row.Volume),
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

func (d *DataSourceCoreVirtualMachines) MapVolume(
	ctx context.Context,
	diags *diag.Diagnostics,
	data virtual_machine.VolumeFieldsForInstance,
) attr.Value {
	model, diagnostic := datasource_core_virtual_machines.NewVolumeValue(
		datasource_core_virtual_machines.VolumeValue{}.AttributeTypes(ctx),
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

func (d *DataSourceCoreVirtualMachines) MapSecurityRules(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []virtual_machine.SecurityRulesFieldsForInstance,
) types.List {
	model, diagnostic := types.ListValue(
		datasource_core_virtual_machines.SecurityRulesValue{}.Type(ctx),
		func() []attr.Value {
			roles := make([]attr.Value, 0)
			for _, row := range data {
				model, diagnostic := datasource_core_virtual_machines.NewSecurityRulesValue(
					datasource_core_virtual_machines.SecurityRulesValue{}.AttributeTypes(ctx),
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
