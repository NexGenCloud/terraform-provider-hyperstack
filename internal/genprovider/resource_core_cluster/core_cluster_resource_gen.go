// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_core_cluster

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CoreClusterResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_address": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"enable_public_ip": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"environment_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"image_name": schema.StringAttribute{
				Required: true,
			},
			"keypair_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kube_config": schema.StringAttribute{
				Computed: true,
			},
			"kubernetes_version": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"master_flavor_name": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"node_addresses": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"node_count": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"node_flavor": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"cpu": schema.Int64Attribute{
						Computed: true,
					},
					"disk": schema.Int64Attribute{
						Computed: true,
					},
					"ephemeral": schema.Int64Attribute{
						Computed: true,
					},
					"gpu": schema.StringAttribute{
						Computed: true,
					},
					"gpu_count": schema.Int64Attribute{
						Computed: true,
					},
					"id": schema.Int64Attribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
					"ram": schema.NumberAttribute{
						Computed: true,
					},
				},
				CustomType: NodeFlavorType{
					ObjectType: types.ObjectType{
						AttrTypes: NodeFlavorValue{}.AttributeTypes(ctx),
					},
				},
				Computed: true,
			},
			"node_flavor_name": schema.StringAttribute{
				Required: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"status_reason": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

type CoreClusterModel struct {
	ApiAddress        types.String    `tfsdk:"api_address"`
	CreatedAt         types.String    `tfsdk:"created_at"`
	EnablePublicIp    types.Bool      `tfsdk:"enable_public_ip"`
	EnvironmentName   types.String    `tfsdk:"environment_name"`
	Id                types.Int64     `tfsdk:"id"`
	ImageName         types.String    `tfsdk:"image_name"`
	KeypairName       types.String    `tfsdk:"keypair_name"`
	KubeConfig        types.String    `tfsdk:"kube_config"`
	KubernetesVersion types.String    `tfsdk:"kubernetes_version"`
	MasterFlavorName  types.String    `tfsdk:"master_flavor_name"`
	Name              types.String    `tfsdk:"name"`
	NodeAddresses     types.List      `tfsdk:"node_addresses"`
	NodeCount         types.Int64     `tfsdk:"node_count"`
	NodeFlavor        NodeFlavorValue `tfsdk:"node_flavor"`
	NodeFlavorName    types.String    `tfsdk:"node_flavor_name"`
	Status            types.String    `tfsdk:"status"`
	StatusReason      types.String    `tfsdk:"status_reason"`
}

var _ basetypes.ObjectTypable = NodeFlavorType{}

type NodeFlavorType struct {
	basetypes.ObjectType
}

func (t NodeFlavorType) Equal(o attr.Type) bool {
	other, ok := o.(NodeFlavorType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t NodeFlavorType) String() string {
	return "NodeFlavorType"
}

func (t NodeFlavorType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	cpuAttribute, ok := attributes["cpu"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`cpu is missing from object`)

		return nil, diags
	}

	cpuVal, ok := cpuAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`cpu expected to be basetypes.Int64Value, was: %T`, cpuAttribute))
	}

	diskAttribute, ok := attributes["disk"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk is missing from object`)

		return nil, diags
	}

	diskVal, ok := diskAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk expected to be basetypes.Int64Value, was: %T`, diskAttribute))
	}

	ephemeralAttribute, ok := attributes["ephemeral"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ephemeral is missing from object`)

		return nil, diags
	}

	ephemeralVal, ok := ephemeralAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ephemeral expected to be basetypes.Int64Value, was: %T`, ephemeralAttribute))
	}

	gpuAttribute, ok := attributes["gpu"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`gpu is missing from object`)

		return nil, diags
	}

	gpuVal, ok := gpuAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`gpu expected to be basetypes.StringValue, was: %T`, gpuAttribute))
	}

	gpuCountAttribute, ok := attributes["gpu_count"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`gpu_count is missing from object`)

		return nil, diags
	}

	gpuCountVal, ok := gpuCountAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`gpu_count expected to be basetypes.Int64Value, was: %T`, gpuCountAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return nil, diags
	}

	idVal, ok := idAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.Int64Value, was: %T`, idAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return nil, diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	ramAttribute, ok := attributes["ram"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ram is missing from object`)

		return nil, diags
	}

	ramVal, ok := ramAttribute.(basetypes.NumberValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ram expected to be basetypes.NumberValue, was: %T`, ramAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return NodeFlavorValue{
		Cpu:       cpuVal,
		Disk:      diskVal,
		Ephemeral: ephemeralVal,
		Gpu:       gpuVal,
		GpuCount:  gpuCountVal,
		Id:        idVal,
		Name:      nameVal,
		Ram:       ramVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewNodeFlavorValueNull() NodeFlavorValue {
	return NodeFlavorValue{
		state: attr.ValueStateNull,
	}
}

func NewNodeFlavorValueUnknown() NodeFlavorValue {
	return NodeFlavorValue{
		state: attr.ValueStateUnknown,
	}
}

func NewNodeFlavorValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (NodeFlavorValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing NodeFlavorValue Attribute Value",
				"While creating a NodeFlavorValue value, a missing attribute value was detected. "+
					"A NodeFlavorValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("NodeFlavorValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid NodeFlavorValue Attribute Type",
				"While creating a NodeFlavorValue value, an invalid attribute value was detected. "+
					"A NodeFlavorValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("NodeFlavorValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("NodeFlavorValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra NodeFlavorValue Attribute Value",
				"While creating a NodeFlavorValue value, an extra attribute value was detected. "+
					"A NodeFlavorValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra NodeFlavorValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewNodeFlavorValueUnknown(), diags
	}

	cpuAttribute, ok := attributes["cpu"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`cpu is missing from object`)

		return NewNodeFlavorValueUnknown(), diags
	}

	cpuVal, ok := cpuAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`cpu expected to be basetypes.Int64Value, was: %T`, cpuAttribute))
	}

	diskAttribute, ok := attributes["disk"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disk is missing from object`)

		return NewNodeFlavorValueUnknown(), diags
	}

	diskVal, ok := diskAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disk expected to be basetypes.Int64Value, was: %T`, diskAttribute))
	}

	ephemeralAttribute, ok := attributes["ephemeral"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ephemeral is missing from object`)

		return NewNodeFlavorValueUnknown(), diags
	}

	ephemeralVal, ok := ephemeralAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ephemeral expected to be basetypes.Int64Value, was: %T`, ephemeralAttribute))
	}

	gpuAttribute, ok := attributes["gpu"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`gpu is missing from object`)

		return NewNodeFlavorValueUnknown(), diags
	}

	gpuVal, ok := gpuAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`gpu expected to be basetypes.StringValue, was: %T`, gpuAttribute))
	}

	gpuCountAttribute, ok := attributes["gpu_count"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`gpu_count is missing from object`)

		return NewNodeFlavorValueUnknown(), diags
	}

	gpuCountVal, ok := gpuCountAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`gpu_count expected to be basetypes.Int64Value, was: %T`, gpuCountAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewNodeFlavorValueUnknown(), diags
	}

	idVal, ok := idAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.Int64Value, was: %T`, idAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewNodeFlavorValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	ramAttribute, ok := attributes["ram"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ram is missing from object`)

		return NewNodeFlavorValueUnknown(), diags
	}

	ramVal, ok := ramAttribute.(basetypes.NumberValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ram expected to be basetypes.NumberValue, was: %T`, ramAttribute))
	}

	if diags.HasError() {
		return NewNodeFlavorValueUnknown(), diags
	}

	return NodeFlavorValue{
		Cpu:       cpuVal,
		Disk:      diskVal,
		Ephemeral: ephemeralVal,
		Gpu:       gpuVal,
		GpuCount:  gpuCountVal,
		Id:        idVal,
		Name:      nameVal,
		Ram:       ramVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewNodeFlavorValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) NodeFlavorValue {
	object, diags := NewNodeFlavorValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewNodeFlavorValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t NodeFlavorType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewNodeFlavorValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewNodeFlavorValueUnknown(), nil
	}

	if in.IsNull() {
		return NewNodeFlavorValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewNodeFlavorValueMust(NodeFlavorValue{}.AttributeTypes(ctx), attributes), nil
}

func (t NodeFlavorType) ValueType(ctx context.Context) attr.Value {
	return NodeFlavorValue{}
}

var _ basetypes.ObjectValuable = NodeFlavorValue{}

type NodeFlavorValue struct {
	Cpu       basetypes.Int64Value  `tfsdk:"cpu"`
	Disk      basetypes.Int64Value  `tfsdk:"disk"`
	Ephemeral basetypes.Int64Value  `tfsdk:"ephemeral"`
	Gpu       basetypes.StringValue `tfsdk:"gpu"`
	GpuCount  basetypes.Int64Value  `tfsdk:"gpu_count"`
	Id        basetypes.Int64Value  `tfsdk:"id"`
	Name      basetypes.StringValue `tfsdk:"name"`
	Ram       basetypes.NumberValue `tfsdk:"ram"`
	state     attr.ValueState
}

func (v NodeFlavorValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 8)

	var val tftypes.Value
	var err error

	attrTypes["cpu"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["disk"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["ephemeral"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["gpu"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["gpu_count"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["id"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["ram"] = basetypes.NumberType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 8)

		val, err = v.Cpu.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["cpu"] = val

		val, err = v.Disk.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["disk"] = val

		val, err = v.Ephemeral.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["ephemeral"] = val

		val, err = v.Gpu.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["gpu"] = val

		val, err = v.GpuCount.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["gpu_count"] = val

		val, err = v.Id.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["id"] = val

		val, err = v.Name.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

		val, err = v.Ram.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["ram"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v NodeFlavorValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v NodeFlavorValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v NodeFlavorValue) String() string {
	return "NodeFlavorValue"
}

func (v NodeFlavorValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"cpu":       basetypes.Int64Type{},
		"disk":      basetypes.Int64Type{},
		"ephemeral": basetypes.Int64Type{},
		"gpu":       basetypes.StringType{},
		"gpu_count": basetypes.Int64Type{},
		"id":        basetypes.Int64Type{},
		"name":      basetypes.StringType{},
		"ram":       basetypes.NumberType{},
	}

	if v.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	}

	if v.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	}

	objVal, diags := types.ObjectValue(
		attributeTypes,
		map[string]attr.Value{
			"cpu":       v.Cpu,
			"disk":      v.Disk,
			"ephemeral": v.Ephemeral,
			"gpu":       v.Gpu,
			"gpu_count": v.GpuCount,
			"id":        v.Id,
			"name":      v.Name,
			"ram":       v.Ram,
		})

	return objVal, diags
}

func (v NodeFlavorValue) Equal(o attr.Value) bool {
	other, ok := o.(NodeFlavorValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Cpu.Equal(other.Cpu) {
		return false
	}

	if !v.Disk.Equal(other.Disk) {
		return false
	}

	if !v.Ephemeral.Equal(other.Ephemeral) {
		return false
	}

	if !v.Gpu.Equal(other.Gpu) {
		return false
	}

	if !v.GpuCount.Equal(other.GpuCount) {
		return false
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	if !v.Ram.Equal(other.Ram) {
		return false
	}

	return true
}

func (v NodeFlavorValue) Type(ctx context.Context) attr.Type {
	return NodeFlavorType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v NodeFlavorValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"cpu":       basetypes.Int64Type{},
		"disk":      basetypes.Int64Type{},
		"ephemeral": basetypes.Int64Type{},
		"gpu":       basetypes.StringType{},
		"gpu_count": basetypes.Int64Type{},
		"id":        basetypes.Int64Type{},
		"name":      basetypes.StringType{},
		"ram":       basetypes.NumberType{},
	}
}
