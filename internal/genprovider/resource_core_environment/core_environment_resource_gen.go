// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package resource_core_environment

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CoreEnvironmentResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the environment being created.",
				MarkdownDescription: "The name of the environment being created.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"region": schema.StringAttribute{
				Required:            true,
				Description:         "The geographic location of the data center where the environment is being created. To learn more about regions, [**click here**](https://infrahub-doc.nexgencloud.com/docs/features/regions).",
				MarkdownDescription: "The geographic location of the data center where the environment is being created. To learn more about regions, [**click here**](https://infrahub-doc.nexgencloud.com/docs/features/regions).",
			},
		},
	}
}

type CoreEnvironmentModel struct {
	CreatedAt types.String `tfsdk:"created_at"`
	Id        types.Int64  `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Region    types.String `tfsdk:"region"`
}
