// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_core_regions

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func CoreRegionsDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"core_regions": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
					CustomType: CoreRegionsType{
						ObjectType: types.ObjectType{
							AttrTypes: CoreRegionsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Computed: true,
			},
		},
	}
}

type CoreRegionsModel struct {
	CoreRegions types.Set `tfsdk:"core_regions"`
}

var _ basetypes.ObjectTypable = CoreRegionsType{}

type CoreRegionsType struct {
	basetypes.ObjectType
}

func (t CoreRegionsType) Equal(o attr.Type) bool {
	other, ok := o.(CoreRegionsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CoreRegionsType) String() string {
	return "CoreRegionsType"
}

func (t CoreRegionsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return nil, diags
	}

	descriptionVal, ok := descriptionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`description expected to be basetypes.StringValue, was: %T`, descriptionAttribute))
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

	if diags.HasError() {
		return nil, diags
	}

	return CoreRegionsValue{
		Description: descriptionVal,
		Id:          idVal,
		Name:        nameVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewCoreRegionsValueNull() CoreRegionsValue {
	return CoreRegionsValue{
		state: attr.ValueStateNull,
	}
}

func NewCoreRegionsValueUnknown() CoreRegionsValue {
	return CoreRegionsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCoreRegionsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CoreRegionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CoreRegionsValue Attribute Value",
				"While creating a CoreRegionsValue value, a missing attribute value was detected. "+
					"A CoreRegionsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CoreRegionsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CoreRegionsValue Attribute Type",
				"While creating a CoreRegionsValue value, an invalid attribute value was detected. "+
					"A CoreRegionsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CoreRegionsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CoreRegionsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CoreRegionsValue Attribute Value",
				"While creating a CoreRegionsValue value, an extra attribute value was detected. "+
					"A CoreRegionsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CoreRegionsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCoreRegionsValueUnknown(), diags
	}

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return NewCoreRegionsValueUnknown(), diags
	}

	descriptionVal, ok := descriptionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`description expected to be basetypes.StringValue, was: %T`, descriptionAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewCoreRegionsValueUnknown(), diags
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

		return NewCoreRegionsValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	if diags.HasError() {
		return NewCoreRegionsValueUnknown(), diags
	}

	return CoreRegionsValue{
		Description: descriptionVal,
		Id:          idVal,
		Name:        nameVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewCoreRegionsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CoreRegionsValue {
	object, diags := NewCoreRegionsValue(attributeTypes, attributes)

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

		panic("NewCoreRegionsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CoreRegionsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCoreRegionsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCoreRegionsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCoreRegionsValueNull(), nil
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

	return NewCoreRegionsValueMust(CoreRegionsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CoreRegionsType) ValueType(ctx context.Context) attr.Value {
	return CoreRegionsValue{}
}

var _ basetypes.ObjectValuable = CoreRegionsValue{}

type CoreRegionsValue struct {
	Description basetypes.StringValue `tfsdk:"description"`
	Id          basetypes.Int64Value  `tfsdk:"id"`
	Name        basetypes.StringValue `tfsdk:"name"`
	state       attr.ValueState
}

func (v CoreRegionsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["description"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["id"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.Description.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["description"] = val

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

func (v CoreRegionsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CoreRegionsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CoreRegionsValue) String() string {
	return "CoreRegionsValue"
}

func (v CoreRegionsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"description": basetypes.StringType{},
			"id":          basetypes.Int64Type{},
			"name":        basetypes.StringType{},
		},
		map[string]attr.Value{
			"description": v.Description,
			"id":          v.Id,
			"name":        v.Name,
		})

	return objVal, diags
}

func (v CoreRegionsValue) Equal(o attr.Value) bool {
	other, ok := o.(CoreRegionsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Description.Equal(other.Description) {
		return false
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	return true
}

func (v CoreRegionsValue) Type(ctx context.Context) attr.Type {
	return CoreRegionsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CoreRegionsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"description": basetypes.StringType{},
		"id":          basetypes.Int64Type{},
		"name":        basetypes.StringType{},
	}
}