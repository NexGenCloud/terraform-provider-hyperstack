// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_core_environments

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

func CoreEnvironmentsDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"core_environments": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
					},
					CustomType: CoreEnvironmentsType{
						ObjectType: types.ObjectType{
							AttrTypes: CoreEnvironmentsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Computed: true,
			},
			"page": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Page Number",
				MarkdownDescription: "Page Number",
			},
			"page_size": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Data Per Page",
				MarkdownDescription: "Data Per Page",
			},
			"search": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Search By ID or Name or Region",
				MarkdownDescription: "Search By ID or Name or Region",
			},
		},
	}
}

type CoreEnvironmentsModel struct {
	CoreEnvironments types.Set    `tfsdk:"core_environments"`
	Page             types.String `tfsdk:"page"`
	PageSize         types.String `tfsdk:"page_size"`
	Search           types.String `tfsdk:"search"`
}

var _ basetypes.ObjectTypable = CoreEnvironmentsType{}

type CoreEnvironmentsType struct {
	basetypes.ObjectType
}

func (t CoreEnvironmentsType) Equal(o attr.Type) bool {
	other, ok := o.(CoreEnvironmentsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CoreEnvironmentsType) String() string {
	return "CoreEnvironmentsType"
}

func (t CoreEnvironmentsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return nil, diags
	}

	createdAtVal, ok := createdAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be basetypes.StringValue, was: %T`, createdAtAttribute))
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

	regionAttribute, ok := attributes["region"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region is missing from object`)

		return nil, diags
	}

	regionVal, ok := regionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be basetypes.StringValue, was: %T`, regionAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CoreEnvironmentsValue{
		CreatedAt: createdAtVal,
		Id:        idVal,
		Name:      nameVal,
		Region:    regionVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewCoreEnvironmentsValueNull() CoreEnvironmentsValue {
	return CoreEnvironmentsValue{
		state: attr.ValueStateNull,
	}
}

func NewCoreEnvironmentsValueUnknown() CoreEnvironmentsValue {
	return CoreEnvironmentsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCoreEnvironmentsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CoreEnvironmentsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CoreEnvironmentsValue Attribute Value",
				"While creating a CoreEnvironmentsValue value, a missing attribute value was detected. "+
					"A CoreEnvironmentsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CoreEnvironmentsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CoreEnvironmentsValue Attribute Type",
				"While creating a CoreEnvironmentsValue value, an invalid attribute value was detected. "+
					"A CoreEnvironmentsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CoreEnvironmentsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CoreEnvironmentsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CoreEnvironmentsValue Attribute Value",
				"While creating a CoreEnvironmentsValue value, an extra attribute value was detected. "+
					"A CoreEnvironmentsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CoreEnvironmentsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCoreEnvironmentsValueUnknown(), diags
	}

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return NewCoreEnvironmentsValueUnknown(), diags
	}

	createdAtVal, ok := createdAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be basetypes.StringValue, was: %T`, createdAtAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewCoreEnvironmentsValueUnknown(), diags
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

		return NewCoreEnvironmentsValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	regionAttribute, ok := attributes["region"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region is missing from object`)

		return NewCoreEnvironmentsValueUnknown(), diags
	}

	regionVal, ok := regionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be basetypes.StringValue, was: %T`, regionAttribute))
	}

	if diags.HasError() {
		return NewCoreEnvironmentsValueUnknown(), diags
	}

	return CoreEnvironmentsValue{
		CreatedAt: createdAtVal,
		Id:        idVal,
		Name:      nameVal,
		Region:    regionVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewCoreEnvironmentsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CoreEnvironmentsValue {
	object, diags := NewCoreEnvironmentsValue(attributeTypes, attributes)

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

		panic("NewCoreEnvironmentsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CoreEnvironmentsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCoreEnvironmentsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCoreEnvironmentsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCoreEnvironmentsValueNull(), nil
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

	return NewCoreEnvironmentsValueMust(CoreEnvironmentsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CoreEnvironmentsType) ValueType(ctx context.Context) attr.Value {
	return CoreEnvironmentsValue{}
}

var _ basetypes.ObjectValuable = CoreEnvironmentsValue{}

type CoreEnvironmentsValue struct {
	CreatedAt basetypes.StringValue `tfsdk:"created_at"`
	Id        basetypes.Int64Value  `tfsdk:"id"`
	Name      basetypes.StringValue `tfsdk:"name"`
	Region    basetypes.StringValue `tfsdk:"region"`
	state     attr.ValueState
}

func (v CoreEnvironmentsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 4)

	var val tftypes.Value
	var err error

	attrTypes["created_at"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["id"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["region"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 4)

		val, err = v.CreatedAt.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["created_at"] = val

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

		val, err = v.Region.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["region"] = val

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

func (v CoreEnvironmentsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CoreEnvironmentsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CoreEnvironmentsValue) String() string {
	return "CoreEnvironmentsValue"
}

func (v CoreEnvironmentsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"created_at": basetypes.StringType{},
		"id":         basetypes.Int64Type{},
		"name":       basetypes.StringType{},
		"region":     basetypes.StringType{},
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
			"created_at": v.CreatedAt,
			"id":         v.Id,
			"name":       v.Name,
			"region":     v.Region,
		})

	return objVal, diags
}

func (v CoreEnvironmentsValue) Equal(o attr.Value) bool {
	other, ok := o.(CoreEnvironmentsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.CreatedAt.Equal(other.CreatedAt) {
		return false
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	if !v.Region.Equal(other.Region) {
		return false
	}

	return true
}

func (v CoreEnvironmentsValue) Type(ctx context.Context) attr.Type {
	return CoreEnvironmentsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CoreEnvironmentsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"created_at": basetypes.StringType{},
		"id":         basetypes.Int64Type{},
		"name":       basetypes.StringType{},
		"region":     basetypes.StringType{},
	}
}
