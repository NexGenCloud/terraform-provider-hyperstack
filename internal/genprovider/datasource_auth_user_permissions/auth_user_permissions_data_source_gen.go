// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_auth_user_permissions

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

func AuthUserPermissionsDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"auth_user_permissions": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"permission": schema.StringAttribute{
							Computed: true,
						},
						"resource": schema.StringAttribute{
							Computed: true,
						},
					},
					CustomType: AuthUserPermissionsType{
						ObjectType: types.ObjectType{
							AttrTypes: AuthUserPermissionsValue{}.AttributeTypes(ctx),
						},
					},
				},
				Computed: true,
			},
			"id": schema.Int64Attribute{
				Required: true,
			},
		},
	}
}

type AuthUserPermissionsModel struct {
	AuthUserPermissions types.Set   `tfsdk:"auth_user_permissions"`
	Id                  types.Int64 `tfsdk:"id"`
}

var _ basetypes.ObjectTypable = AuthUserPermissionsType{}

type AuthUserPermissionsType struct {
	basetypes.ObjectType
}

func (t AuthUserPermissionsType) Equal(o attr.Type) bool {
	other, ok := o.(AuthUserPermissionsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t AuthUserPermissionsType) String() string {
	return "AuthUserPermissionsType"
}

func (t AuthUserPermissionsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

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

	permissionAttribute, ok := attributes["permission"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`permission is missing from object`)

		return nil, diags
	}

	permissionVal, ok := permissionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`permission expected to be basetypes.StringValue, was: %T`, permissionAttribute))
	}

	resourceAttribute, ok := attributes["resource"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`resource is missing from object`)

		return nil, diags
	}

	resourceVal, ok := resourceAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`resource expected to be basetypes.StringValue, was: %T`, resourceAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return AuthUserPermissionsValue{
		Id:         idVal,
		Permission: permissionVal,
		Resource:   resourceVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewAuthUserPermissionsValueNull() AuthUserPermissionsValue {
	return AuthUserPermissionsValue{
		state: attr.ValueStateNull,
	}
}

func NewAuthUserPermissionsValueUnknown() AuthUserPermissionsValue {
	return AuthUserPermissionsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewAuthUserPermissionsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (AuthUserPermissionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing AuthUserPermissionsValue Attribute Value",
				"While creating a AuthUserPermissionsValue value, a missing attribute value was detected. "+
					"A AuthUserPermissionsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AuthUserPermissionsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid AuthUserPermissionsValue Attribute Type",
				"While creating a AuthUserPermissionsValue value, an invalid attribute value was detected. "+
					"A AuthUserPermissionsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AuthUserPermissionsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("AuthUserPermissionsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra AuthUserPermissionsValue Attribute Value",
				"While creating a AuthUserPermissionsValue value, an extra attribute value was detected. "+
					"A AuthUserPermissionsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra AuthUserPermissionsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewAuthUserPermissionsValueUnknown(), diags
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewAuthUserPermissionsValueUnknown(), diags
	}

	idVal, ok := idAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be basetypes.Int64Value, was: %T`, idAttribute))
	}

	permissionAttribute, ok := attributes["permission"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`permission is missing from object`)

		return NewAuthUserPermissionsValueUnknown(), diags
	}

	permissionVal, ok := permissionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`permission expected to be basetypes.StringValue, was: %T`, permissionAttribute))
	}

	resourceAttribute, ok := attributes["resource"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`resource is missing from object`)

		return NewAuthUserPermissionsValueUnknown(), diags
	}

	resourceVal, ok := resourceAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`resource expected to be basetypes.StringValue, was: %T`, resourceAttribute))
	}

	if diags.HasError() {
		return NewAuthUserPermissionsValueUnknown(), diags
	}

	return AuthUserPermissionsValue{
		Id:         idVal,
		Permission: permissionVal,
		Resource:   resourceVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewAuthUserPermissionsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) AuthUserPermissionsValue {
	object, diags := NewAuthUserPermissionsValue(attributeTypes, attributes)

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

		panic("NewAuthUserPermissionsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t AuthUserPermissionsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewAuthUserPermissionsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewAuthUserPermissionsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewAuthUserPermissionsValueNull(), nil
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

	return NewAuthUserPermissionsValueMust(AuthUserPermissionsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t AuthUserPermissionsType) ValueType(ctx context.Context) attr.Value {
	return AuthUserPermissionsValue{}
}

var _ basetypes.ObjectValuable = AuthUserPermissionsValue{}

type AuthUserPermissionsValue struct {
	Id         basetypes.Int64Value  `tfsdk:"id"`
	Permission basetypes.StringValue `tfsdk:"permission"`
	Resource   basetypes.StringValue `tfsdk:"resource"`
	state      attr.ValueState
}

func (v AuthUserPermissionsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["id"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["permission"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["resource"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.Id.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["id"] = val

		val, err = v.Permission.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["permission"] = val

		val, err = v.Resource.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["resource"] = val

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

func (v AuthUserPermissionsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v AuthUserPermissionsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v AuthUserPermissionsValue) String() string {
	return "AuthUserPermissionsValue"
}

func (v AuthUserPermissionsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"id":         basetypes.Int64Type{},
		"permission": basetypes.StringType{},
		"resource":   basetypes.StringType{},
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
			"id":         v.Id,
			"permission": v.Permission,
			"resource":   v.Resource,
		})

	return objVal, diags
}

func (v AuthUserPermissionsValue) Equal(o attr.Value) bool {
	other, ok := o.(AuthUserPermissionsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.Permission.Equal(other.Permission) {
		return false
	}

	if !v.Resource.Equal(other.Resource) {
		return false
	}

	return true
}

func (v AuthUserPermissionsValue) Type(ctx context.Context) attr.Type {
	return AuthUserPermissionsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v AuthUserPermissionsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":         basetypes.Int64Type{},
		"permission": basetypes.StringType{},
		"resource":   basetypes.StringType{},
	}
}
