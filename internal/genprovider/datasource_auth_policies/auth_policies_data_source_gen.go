// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_auth_policies

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

func AuthPoliciesDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"auth_policies": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"permissions": schema.ListNestedAttribute{
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
								CustomType: PermissionsType{
									ObjectType: types.ObjectType{
										AttrTypes: PermissionsValue{}.AttributeTypes(ctx),
									},
								},
							},
							Computed: true,
						},
					},
					CustomType: AuthPoliciesType{
						ObjectType: types.ObjectType{
							AttrTypes: AuthPoliciesValue{}.AttributeTypes(ctx),
						},
					},
				},
				Computed: true,
			},
		},
	}
}

type AuthPoliciesModel struct {
	AuthPolicies types.Set `tfsdk:"auth_policies"`
}

var _ basetypes.ObjectTypable = AuthPoliciesType{}

type AuthPoliciesType struct {
	basetypes.ObjectType
}

func (t AuthPoliciesType) Equal(o attr.Type) bool {
	other, ok := o.(AuthPoliciesType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t AuthPoliciesType) String() string {
	return "AuthPoliciesType"
}

func (t AuthPoliciesType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

	permissionsAttribute, ok := attributes["permissions"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`permissions is missing from object`)

		return nil, diags
	}

	permissionsVal, ok := permissionsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`permissions expected to be basetypes.ListValue, was: %T`, permissionsAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return AuthPoliciesValue{
		CreatedAt:   createdAtVal,
		Description: descriptionVal,
		Id:          idVal,
		Name:        nameVal,
		Permissions: permissionsVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewAuthPoliciesValueNull() AuthPoliciesValue {
	return AuthPoliciesValue{
		state: attr.ValueStateNull,
	}
}

func NewAuthPoliciesValueUnknown() AuthPoliciesValue {
	return AuthPoliciesValue{
		state: attr.ValueStateUnknown,
	}
}

func NewAuthPoliciesValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (AuthPoliciesValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing AuthPoliciesValue Attribute Value",
				"While creating a AuthPoliciesValue value, a missing attribute value was detected. "+
					"A AuthPoliciesValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AuthPoliciesValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid AuthPoliciesValue Attribute Type",
				"While creating a AuthPoliciesValue value, an invalid attribute value was detected. "+
					"A AuthPoliciesValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("AuthPoliciesValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("AuthPoliciesValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra AuthPoliciesValue Attribute Value",
				"While creating a AuthPoliciesValue value, an extra attribute value was detected. "+
					"A AuthPoliciesValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra AuthPoliciesValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewAuthPoliciesValueUnknown(), diags
	}

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return NewAuthPoliciesValueUnknown(), diags
	}

	createdAtVal, ok := createdAtAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be basetypes.StringValue, was: %T`, createdAtAttribute))
	}

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return NewAuthPoliciesValueUnknown(), diags
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

		return NewAuthPoliciesValueUnknown(), diags
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

		return NewAuthPoliciesValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	permissionsAttribute, ok := attributes["permissions"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`permissions is missing from object`)

		return NewAuthPoliciesValueUnknown(), diags
	}

	permissionsVal, ok := permissionsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`permissions expected to be basetypes.ListValue, was: %T`, permissionsAttribute))
	}

	if diags.HasError() {
		return NewAuthPoliciesValueUnknown(), diags
	}

	return AuthPoliciesValue{
		CreatedAt:   createdAtVal,
		Description: descriptionVal,
		Id:          idVal,
		Name:        nameVal,
		Permissions: permissionsVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewAuthPoliciesValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) AuthPoliciesValue {
	object, diags := NewAuthPoliciesValue(attributeTypes, attributes)

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

		panic("NewAuthPoliciesValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t AuthPoliciesType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewAuthPoliciesValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewAuthPoliciesValueUnknown(), nil
	}

	if in.IsNull() {
		return NewAuthPoliciesValueNull(), nil
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

	return NewAuthPoliciesValueMust(AuthPoliciesValue{}.AttributeTypes(ctx), attributes), nil
}

func (t AuthPoliciesType) ValueType(ctx context.Context) attr.Value {
	return AuthPoliciesValue{}
}

var _ basetypes.ObjectValuable = AuthPoliciesValue{}

type AuthPoliciesValue struct {
	CreatedAt   basetypes.StringValue `tfsdk:"created_at"`
	Description basetypes.StringValue `tfsdk:"description"`
	Id          basetypes.Int64Value  `tfsdk:"id"`
	Name        basetypes.StringValue `tfsdk:"name"`
	Permissions basetypes.ListValue   `tfsdk:"permissions"`
	state       attr.ValueState
}

func (v AuthPoliciesValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 5)

	var val tftypes.Value
	var err error

	attrTypes["created_at"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["description"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["id"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["permissions"] = basetypes.ListType{
		ElemType: PermissionsValue{}.Type(ctx),
	}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 5)

		val, err = v.CreatedAt.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["created_at"] = val

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

		val, err = v.Permissions.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["permissions"] = val

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

func (v AuthPoliciesValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v AuthPoliciesValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v AuthPoliciesValue) String() string {
	return "AuthPoliciesValue"
}

func (v AuthPoliciesValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	permissions := types.ListValueMust(
		PermissionsType{
			basetypes.ObjectType{
				AttrTypes: PermissionsValue{}.AttributeTypes(ctx),
			},
		},
		v.Permissions.Elements(),
	)

	if v.Permissions.IsNull() {
		permissions = types.ListNull(
			PermissionsType{
				basetypes.ObjectType{
					AttrTypes: PermissionsValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	if v.Permissions.IsUnknown() {
		permissions = types.ListUnknown(
			PermissionsType{
				basetypes.ObjectType{
					AttrTypes: PermissionsValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	attributeTypes := map[string]attr.Type{
		"created_at":  basetypes.StringType{},
		"description": basetypes.StringType{},
		"id":          basetypes.Int64Type{},
		"name":        basetypes.StringType{},
		"permissions": basetypes.ListType{
			ElemType: PermissionsValue{}.Type(ctx),
		},
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
			"created_at":  v.CreatedAt,
			"description": v.Description,
			"id":          v.Id,
			"name":        v.Name,
			"permissions": permissions,
		})

	return objVal, diags
}

func (v AuthPoliciesValue) Equal(o attr.Value) bool {
	other, ok := o.(AuthPoliciesValue)

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

	if !v.Description.Equal(other.Description) {
		return false
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	if !v.Permissions.Equal(other.Permissions) {
		return false
	}

	return true
}

func (v AuthPoliciesValue) Type(ctx context.Context) attr.Type {
	return AuthPoliciesType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v AuthPoliciesValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"created_at":  basetypes.StringType{},
		"description": basetypes.StringType{},
		"id":          basetypes.Int64Type{},
		"name":        basetypes.StringType{},
		"permissions": basetypes.ListType{
			ElemType: PermissionsValue{}.Type(ctx),
		},
	}
}

var _ basetypes.ObjectTypable = PermissionsType{}

type PermissionsType struct {
	basetypes.ObjectType
}

func (t PermissionsType) Equal(o attr.Type) bool {
	other, ok := o.(PermissionsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t PermissionsType) String() string {
	return "PermissionsType"
}

func (t PermissionsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

	return PermissionsValue{
		Id:         idVal,
		Permission: permissionVal,
		Resource:   resourceVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewPermissionsValueNull() PermissionsValue {
	return PermissionsValue{
		state: attr.ValueStateNull,
	}
}

func NewPermissionsValueUnknown() PermissionsValue {
	return PermissionsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewPermissionsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (PermissionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing PermissionsValue Attribute Value",
				"While creating a PermissionsValue value, a missing attribute value was detected. "+
					"A PermissionsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PermissionsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid PermissionsValue Attribute Type",
				"While creating a PermissionsValue value, an invalid attribute value was detected. "+
					"A PermissionsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PermissionsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("PermissionsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra PermissionsValue Attribute Value",
				"While creating a PermissionsValue value, an extra attribute value was detected. "+
					"A PermissionsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra PermissionsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewPermissionsValueUnknown(), diags
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewPermissionsValueUnknown(), diags
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

		return NewPermissionsValueUnknown(), diags
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

		return NewPermissionsValueUnknown(), diags
	}

	resourceVal, ok := resourceAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`resource expected to be basetypes.StringValue, was: %T`, resourceAttribute))
	}

	if diags.HasError() {
		return NewPermissionsValueUnknown(), diags
	}

	return PermissionsValue{
		Id:         idVal,
		Permission: permissionVal,
		Resource:   resourceVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewPermissionsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) PermissionsValue {
	object, diags := NewPermissionsValue(attributeTypes, attributes)

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

		panic("NewPermissionsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t PermissionsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewPermissionsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewPermissionsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewPermissionsValueNull(), nil
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

	return NewPermissionsValueMust(PermissionsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t PermissionsType) ValueType(ctx context.Context) attr.Value {
	return PermissionsValue{}
}

var _ basetypes.ObjectValuable = PermissionsValue{}

type PermissionsValue struct {
	Id         basetypes.Int64Value  `tfsdk:"id"`
	Permission basetypes.StringValue `tfsdk:"permission"`
	Resource   basetypes.StringValue `tfsdk:"resource"`
	state      attr.ValueState
}

func (v PermissionsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
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

func (v PermissionsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v PermissionsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v PermissionsValue) String() string {
	return "PermissionsValue"
}

func (v PermissionsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (v PermissionsValue) Equal(o attr.Value) bool {
	other, ok := o.(PermissionsValue)

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

func (v PermissionsValue) Type(ctx context.Context) attr.Type {
	return PermissionsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v PermissionsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":         basetypes.Int64Type{},
		"permission": basetypes.StringType{},
		"resource":   basetypes.StringType{},
	}
}