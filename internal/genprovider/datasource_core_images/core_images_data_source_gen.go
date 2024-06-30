// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_core_images

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

func CoreImagesDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"core_images": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"images": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"display_size": schema.StringAttribute{
										Computed: true,
									},
									"id": schema.Int64Attribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
									"region_name": schema.StringAttribute{
										Computed: true,
									},
									"size": schema.Int64Attribute{
										Computed: true,
									},
									"type": schema.StringAttribute{
										Computed: true,
									},
									"version": schema.StringAttribute{
										Computed: true,
									},
								},
								CustomType: ImagesType{
									ObjectType: types.ObjectType{
										AttrTypes: ImagesValue{}.AttributeTypes(ctx),
									},
								},
							},
							Computed: true,
						},
						"logo": schema.StringAttribute{
							Computed: true,
						},
						"region_name": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
					},
					CustomType: CoreImagesType{
						ObjectType: types.ObjectType{
							AttrTypes: CoreImagesValue{}.AttributeTypes(ctx),
						},
					},
				},
				Computed: true,
			},
			"region": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Region Name",
				MarkdownDescription: "Region Name",
			},
		},
	}
}

type CoreImagesModel struct {
	CoreImages types.Set    `tfsdk:"core_images"`
	Region     types.String `tfsdk:"region"`
}

var _ basetypes.ObjectTypable = CoreImagesType{}

type CoreImagesType struct {
	basetypes.ObjectType
}

func (t CoreImagesType) Equal(o attr.Type) bool {
	other, ok := o.(CoreImagesType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CoreImagesType) String() string {
	return "CoreImagesType"
}

func (t CoreImagesType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	imagesAttribute, ok := attributes["images"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`images is missing from object`)

		return nil, diags
	}

	imagesVal, ok := imagesAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`images expected to be basetypes.ListValue, was: %T`, imagesAttribute))
	}

	logoAttribute, ok := attributes["logo"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`logo is missing from object`)

		return nil, diags
	}

	logoVal, ok := logoAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`logo expected to be basetypes.StringValue, was: %T`, logoAttribute))
	}

	regionNameAttribute, ok := attributes["region_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region_name is missing from object`)

		return nil, diags
	}

	regionNameVal, ok := regionNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region_name expected to be basetypes.StringValue, was: %T`, regionNameAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return nil, diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CoreImagesValue{
		Images:         imagesVal,
		Logo:           logoVal,
		RegionName:     regionNameVal,
		CoreImagesType: typeVal,
		state:          attr.ValueStateKnown,
	}, diags
}

func NewCoreImagesValueNull() CoreImagesValue {
	return CoreImagesValue{
		state: attr.ValueStateNull,
	}
}

func NewCoreImagesValueUnknown() CoreImagesValue {
	return CoreImagesValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCoreImagesValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CoreImagesValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CoreImagesValue Attribute Value",
				"While creating a CoreImagesValue value, a missing attribute value was detected. "+
					"A CoreImagesValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CoreImagesValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CoreImagesValue Attribute Type",
				"While creating a CoreImagesValue value, an invalid attribute value was detected. "+
					"A CoreImagesValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CoreImagesValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CoreImagesValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CoreImagesValue Attribute Value",
				"While creating a CoreImagesValue value, an extra attribute value was detected. "+
					"A CoreImagesValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CoreImagesValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCoreImagesValueUnknown(), diags
	}

	imagesAttribute, ok := attributes["images"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`images is missing from object`)

		return NewCoreImagesValueUnknown(), diags
	}

	imagesVal, ok := imagesAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`images expected to be basetypes.ListValue, was: %T`, imagesAttribute))
	}

	logoAttribute, ok := attributes["logo"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`logo is missing from object`)

		return NewCoreImagesValueUnknown(), diags
	}

	logoVal, ok := logoAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`logo expected to be basetypes.StringValue, was: %T`, logoAttribute))
	}

	regionNameAttribute, ok := attributes["region_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region_name is missing from object`)

		return NewCoreImagesValueUnknown(), diags
	}

	regionNameVal, ok := regionNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region_name expected to be basetypes.StringValue, was: %T`, regionNameAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return NewCoreImagesValueUnknown(), diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return NewCoreImagesValueUnknown(), diags
	}

	return CoreImagesValue{
		Images:         imagesVal,
		Logo:           logoVal,
		RegionName:     regionNameVal,
		CoreImagesType: typeVal,
		state:          attr.ValueStateKnown,
	}, diags
}

func NewCoreImagesValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CoreImagesValue {
	object, diags := NewCoreImagesValue(attributeTypes, attributes)

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

		panic("NewCoreImagesValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CoreImagesType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCoreImagesValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCoreImagesValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCoreImagesValueNull(), nil
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

	return NewCoreImagesValueMust(CoreImagesValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CoreImagesType) ValueType(ctx context.Context) attr.Value {
	return CoreImagesValue{}
}

var _ basetypes.ObjectValuable = CoreImagesValue{}

type CoreImagesValue struct {
	Images         basetypes.ListValue   `tfsdk:"images"`
	Logo           basetypes.StringValue `tfsdk:"logo"`
	RegionName     basetypes.StringValue `tfsdk:"region_name"`
	CoreImagesType basetypes.StringValue `tfsdk:"type"`
	state          attr.ValueState
}

func (v CoreImagesValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 4)

	var val tftypes.Value
	var err error

	attrTypes["images"] = basetypes.ListType{
		ElemType: ImagesValue{}.Type(ctx),
	}.TerraformType(ctx)
	attrTypes["logo"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["region_name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["type"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 4)

		val, err = v.Images.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["images"] = val

		val, err = v.Logo.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["logo"] = val

		val, err = v.RegionName.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["region_name"] = val

		val, err = v.CoreImagesType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["type"] = val

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

func (v CoreImagesValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CoreImagesValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CoreImagesValue) String() string {
	return "CoreImagesValue"
}

func (v CoreImagesValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	images := types.ListValueMust(
		ImagesType{
			basetypes.ObjectType{
				AttrTypes: ImagesValue{}.AttributeTypes(ctx),
			},
		},
		v.Images.Elements(),
	)

	if v.Images.IsNull() {
		images = types.ListNull(
			ImagesType{
				basetypes.ObjectType{
					AttrTypes: ImagesValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	if v.Images.IsUnknown() {
		images = types.ListUnknown(
			ImagesType{
				basetypes.ObjectType{
					AttrTypes: ImagesValue{}.AttributeTypes(ctx),
				},
			},
		)
	}

	attributeTypes := map[string]attr.Type{
		"images": basetypes.ListType{
			ElemType: ImagesValue{}.Type(ctx),
		},
		"logo":        basetypes.StringType{},
		"region_name": basetypes.StringType{},
		"type":        basetypes.StringType{},
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
			"images":      images,
			"logo":        v.Logo,
			"region_name": v.RegionName,
			"type":        v.CoreImagesType,
		})

	return objVal, diags
}

func (v CoreImagesValue) Equal(o attr.Value) bool {
	other, ok := o.(CoreImagesValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Images.Equal(other.Images) {
		return false
	}

	if !v.Logo.Equal(other.Logo) {
		return false
	}

	if !v.RegionName.Equal(other.RegionName) {
		return false
	}

	if !v.CoreImagesType.Equal(other.CoreImagesType) {
		return false
	}

	return true
}

func (v CoreImagesValue) Type(ctx context.Context) attr.Type {
	return CoreImagesType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CoreImagesValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"images": basetypes.ListType{
			ElemType: ImagesValue{}.Type(ctx),
		},
		"logo":        basetypes.StringType{},
		"region_name": basetypes.StringType{},
		"type":        basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = ImagesType{}

type ImagesType struct {
	basetypes.ObjectType
}

func (t ImagesType) Equal(o attr.Type) bool {
	other, ok := o.(ImagesType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ImagesType) String() string {
	return "ImagesType"
}

func (t ImagesType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	displaySizeAttribute, ok := attributes["display_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`display_size is missing from object`)

		return nil, diags
	}

	displaySizeVal, ok := displaySizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`display_size expected to be basetypes.StringValue, was: %T`, displaySizeAttribute))
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

	regionNameAttribute, ok := attributes["region_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region_name is missing from object`)

		return nil, diags
	}

	regionNameVal, ok := regionNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region_name expected to be basetypes.StringValue, was: %T`, regionNameAttribute))
	}

	sizeAttribute, ok := attributes["size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`size is missing from object`)

		return nil, diags
	}

	sizeVal, ok := sizeAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`size expected to be basetypes.Int64Value, was: %T`, sizeAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return nil, diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	versionAttribute, ok := attributes["version"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`version is missing from object`)

		return nil, diags
	}

	versionVal, ok := versionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`version expected to be basetypes.StringValue, was: %T`, versionAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ImagesValue{
		DisplaySize: displaySizeVal,
		Id:          idVal,
		Name:        nameVal,
		RegionName:  regionNameVal,
		Size:        sizeVal,
		ImagesType:  typeVal,
		Version:     versionVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewImagesValueNull() ImagesValue {
	return ImagesValue{
		state: attr.ValueStateNull,
	}
}

func NewImagesValueUnknown() ImagesValue {
	return ImagesValue{
		state: attr.ValueStateUnknown,
	}
}

func NewImagesValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ImagesValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ImagesValue Attribute Value",
				"While creating a ImagesValue value, a missing attribute value was detected. "+
					"A ImagesValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ImagesValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ImagesValue Attribute Type",
				"While creating a ImagesValue value, an invalid attribute value was detected. "+
					"A ImagesValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ImagesValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ImagesValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ImagesValue Attribute Value",
				"While creating a ImagesValue value, an extra attribute value was detected. "+
					"A ImagesValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ImagesValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewImagesValueUnknown(), diags
	}

	displaySizeAttribute, ok := attributes["display_size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`display_size is missing from object`)

		return NewImagesValueUnknown(), diags
	}

	displaySizeVal, ok := displaySizeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`display_size expected to be basetypes.StringValue, was: %T`, displaySizeAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewImagesValueUnknown(), diags
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

		return NewImagesValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be basetypes.StringValue, was: %T`, nameAttribute))
	}

	regionNameAttribute, ok := attributes["region_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region_name is missing from object`)

		return NewImagesValueUnknown(), diags
	}

	regionNameVal, ok := regionNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region_name expected to be basetypes.StringValue, was: %T`, regionNameAttribute))
	}

	sizeAttribute, ok := attributes["size"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`size is missing from object`)

		return NewImagesValueUnknown(), diags
	}

	sizeVal, ok := sizeAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`size expected to be basetypes.Int64Value, was: %T`, sizeAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return NewImagesValueUnknown(), diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	versionAttribute, ok := attributes["version"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`version is missing from object`)

		return NewImagesValueUnknown(), diags
	}

	versionVal, ok := versionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`version expected to be basetypes.StringValue, was: %T`, versionAttribute))
	}

	if diags.HasError() {
		return NewImagesValueUnknown(), diags
	}

	return ImagesValue{
		DisplaySize: displaySizeVal,
		Id:          idVal,
		Name:        nameVal,
		RegionName:  regionNameVal,
		Size:        sizeVal,
		ImagesType:  typeVal,
		Version:     versionVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewImagesValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ImagesValue {
	object, diags := NewImagesValue(attributeTypes, attributes)

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

		panic("NewImagesValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ImagesType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewImagesValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewImagesValueUnknown(), nil
	}

	if in.IsNull() {
		return NewImagesValueNull(), nil
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

	return NewImagesValueMust(ImagesValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ImagesType) ValueType(ctx context.Context) attr.Value {
	return ImagesValue{}
}

var _ basetypes.ObjectValuable = ImagesValue{}

type ImagesValue struct {
	DisplaySize basetypes.StringValue `tfsdk:"display_size"`
	Id          basetypes.Int64Value  `tfsdk:"id"`
	Name        basetypes.StringValue `tfsdk:"name"`
	RegionName  basetypes.StringValue `tfsdk:"region_name"`
	Size        basetypes.Int64Value  `tfsdk:"size"`
	ImagesType  basetypes.StringValue `tfsdk:"type"`
	Version     basetypes.StringValue `tfsdk:"version"`
	state       attr.ValueState
}

func (v ImagesValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 7)

	var val tftypes.Value
	var err error

	attrTypes["display_size"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["id"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["region_name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["size"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["type"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["version"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 7)

		val, err = v.DisplaySize.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["display_size"] = val

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

		val, err = v.RegionName.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["region_name"] = val

		val, err = v.Size.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["size"] = val

		val, err = v.ImagesType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["type"] = val

		val, err = v.Version.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["version"] = val

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

func (v ImagesValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ImagesValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ImagesValue) String() string {
	return "ImagesValue"
}

func (v ImagesValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := map[string]attr.Type{
		"display_size": basetypes.StringType{},
		"id":           basetypes.Int64Type{},
		"name":         basetypes.StringType{},
		"region_name":  basetypes.StringType{},
		"size":         basetypes.Int64Type{},
		"type":         basetypes.StringType{},
		"version":      basetypes.StringType{},
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
			"display_size": v.DisplaySize,
			"id":           v.Id,
			"name":         v.Name,
			"region_name":  v.RegionName,
			"size":         v.Size,
			"type":         v.ImagesType,
			"version":      v.Version,
		})

	return objVal, diags
}

func (v ImagesValue) Equal(o attr.Value) bool {
	other, ok := o.(ImagesValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.DisplaySize.Equal(other.DisplaySize) {
		return false
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	if !v.RegionName.Equal(other.RegionName) {
		return false
	}

	if !v.Size.Equal(other.Size) {
		return false
	}

	if !v.ImagesType.Equal(other.ImagesType) {
		return false
	}

	if !v.Version.Equal(other.Version) {
		return false
	}

	return true
}

func (v ImagesValue) Type(ctx context.Context) attr.Type {
	return ImagesType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ImagesValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"display_size": basetypes.StringType{},
		"id":           basetypes.Int64Type{},
		"name":         basetypes.StringType{},
		"region_name":  basetypes.StringType{},
		"size":         basetypes.Int64Type{},
		"type":         basetypes.StringType{},
		"version":      basetypes.StringType{},
	}
}
