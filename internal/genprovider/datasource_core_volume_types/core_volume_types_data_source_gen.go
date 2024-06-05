// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_core_volume_types

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func CoreVolumeTypesDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"core_volume_types": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

type CoreVolumeTypesModel struct {
	CoreVolumeTypes types.Set `tfsdk:"core_volume_types"`
}