// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_core_volume

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CoreVolumeResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bootable": schema.BoolAttribute{
				Computed: true,
			},
			"callback_url": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed: true,
					},
				},
				CustomType: EnvironmentType{
					ObjectType: types.ObjectType{
						AttrTypes: EnvironmentValue{}.AttributeTypes(ctx),
					},
				},
				Computed: true,
			},
			"environment_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"image_id": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
			"volume_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

type CoreVolumeModel struct {
	Bootable        types.Bool       `tfsdk:"bootable"`
	CallbackUrl     types.String     `tfsdk:"callback_url"`
	CreatedAt       types.String     `tfsdk:"created_at"`
	Description     types.String     `tfsdk:"description"`
	Environment     EnvironmentValue `tfsdk:"environment"`
	EnvironmentName types.String     `tfsdk:"environment_name"`
	Id              types.Int64      `tfsdk:"id"`
	ImageId         types.Int64      `tfsdk:"image_id"`
	Name            types.String     `tfsdk:"name"`
	Size            types.Int64      `tfsdk:"size"`
	Status          types.String     `tfsdk:"status"`
	UpdatedAt       types.String     `tfsdk:"updated_at"`
	VolumeType      types.String     `tfsdk:"volume_type"`
}

var _ basetypes.ObjectTypable = EnvironmentType{}

type EnvironmentType struct {
	basetypes.ObjectType
}

func (t EnvironmentType) Equal(o attr.Type) bool {
	other, ok := o.(EnvironmentType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t EnvironmentType) String() string {
	return "EnvironmentType"
}

func (t EnvironmentType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

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

	if diags.HasError() {
		return nil, diags
	}

	return EnvironmentValue{
		Name:  nameVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewEnvironmentValueNull() EnvironmentValue {
	return EnvironmentValue{
		state: attr.ValueStateNull,
	}
}

func NewEnvironmentValueUnknown() EnvironmentValue {
	return EnvironmentValue{
		state: attr.ValueStateUnknown,
	}
}

func NewEnvironmentValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (EnvironmentValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing EnvironmentValue Attribute Value",
				"While creating a EnvironmentValue value, a missing attribute value was detected. "+
					"A EnvironmentValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("EnvironmentValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid EnvironmentValue Attribute Type",
				"While creating a EnvironmentValue value, an invalid attribute value was detected. "+
					"A EnvironmentValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("EnvironmentValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("EnvironmentValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra EnvironmentValue Attribute Value",
				"While creating a EnvironmentValue value, an extra attribute value was detected. "+
					"A EnvironmentValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra EnvironmentValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewEnvironmentValueUnknown(), diags
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewEnvironmentValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	if diags.HasError() {
		return NewEnvironmentValueUnknown(), diags
	}

	return EnvironmentValue{
		Name:  nameVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewEnvironmentValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) EnvironmentValue {
	object, diags := NewEnvironmentValue(attributeTypes, attributes)

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

		panic("NewEnvironmentValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t EnvironmentType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewEnvironmentValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewEnvironmentValueUnknown(), nil
	}

	if in.IsNull() {
		return NewEnvironmentValueNull(), nil
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

	return NewEnvironmentValueMust(EnvironmentValue{}.AttributeTypes(ctx), attributes), nil
}

func (t EnvironmentType) ValueType(ctx context.Context) attr.Value {
	return EnvironmentValue{}
}

var _ basetypes.ObjectValuable = EnvironmentValue{}

type EnvironmentValue struct {
	Name  basetypes.StringValue `tfsdk:"name"`
	state attr.ValueState
}

func (v EnvironmentValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 1)

	var val tftypes.Value
	var err error

	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 1)

		val, err = v.Name.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

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

func (v EnvironmentValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v EnvironmentValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v EnvironmentValue) String() string {
	return "EnvironmentValue"
}

func (v EnvironmentValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"name": basetypes.StringType{},
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
			"name": v.Name,
		})

	return objVal, diags
}

func (v EnvironmentValue) Equal(o attr.Value) bool {
	other, ok := o.(EnvironmentValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	return true
}

func (v EnvironmentValue) Type(ctx context.Context) attr.Type {
	return EnvironmentType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v EnvironmentValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"name": basetypes.StringType{},
	}
}
