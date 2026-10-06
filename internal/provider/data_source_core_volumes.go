package provider

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/NexGenCloud/hyperstack-sdk-go/lib/volume"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/datasource_core_volumes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DataSourceCoreVolumes{}

func NewDataSourceCoreVolumes() datasource.DataSource {
	return &DataSourceCoreVolumes{}
}

type DataSourceCoreVolumes struct {
	hyperstack *client.HyperstackClient
	client     *volume.ClientWithResponses
}

func (d *DataSourceCoreVolumes) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_volumes"
}

func (d *DataSourceCoreVolumes) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_volumes.CoreVolumesDataSourceSchema(ctx)
}

func (d *DataSourceCoreVolumes) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = volume.NewClientWithResponses(
		d.hyperstack.ApiServer,
		volume.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreVolumes) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_volumes.CoreVolumesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.ListVolumesWithResponse(ctx, func() *volume.ListVolumesParams {
		return &volume.ListVolumesParams{
			Page:     nil,
			PageSize: nil,
			Search:   nil,
		}
	}())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if result.JSON200 == nil {
		bodyBytes, _ := ioutil.ReadAll(result.HTTPResponse.Body)
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	if result.JSON200.Volumes == nil {
		// API returns 'volumes', SDK struct field is 'Volumes'
		resp.Diagnostics.AddWarning(
			"No user data",
			"",
		)
		return
	}

	data = d.ApiToModel(ctx, &resp.Diagnostics, result.JSON200.Volumes)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceCoreVolumes) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]volume.VolumesFields,
) datasource_core_volumes.CoreVolumesModel {
	return datasource_core_volumes.CoreVolumesModel{
		CoreVolumes: func() types.Set {
			return d.MapVolumes(ctx, diags, *response)
		}(),
	}
}

func (d *DataSourceCoreVolumes) MapVolumes(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []volume.VolumesFields,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_core_volumes.CoreVolumesValue{}.Type(ctx),
		func() []attr.Value {
			volumes := make([]attr.Value, 0)
			for _, row := range data {
				// The schema declares name, region and features; a cut-down object
				// type does not match it.
				environment, diagnostic := types.ObjectValue(
					datasource_core_volumes.EnvironmentValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"name":   types.StringPointerValue(row.Environment.Name),
						"region": types.StringPointerValue(row.Environment.Region),
						// Free-form object in the spec; nothing to carry over.
						"features": types.ObjectNull(datasource_core_volumes.FeaturesValue{}.AttributeTypes(ctx)),
					},
				)
				diags.Append(diagnostic...)
				model, diagnostic := datasource_core_volumes.NewCoreVolumesValue(
					datasource_core_volumes.CoreVolumesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id": func() attr.Value {
							if row.Id == nil {
								return types.Int64Null()
							}
							return types.Int64Value(int64(*row.Id))
						}(),
						"name":        types.StringPointerValue(row.Name),
						"environment": environment,
						"description": types.StringPointerValue(row.Description),
						"volume_type": types.StringPointerValue(row.VolumeType),
						"size": func() attr.Value {
							if row.Size == nil {
								return types.Int64Null()
							}
							return types.Int64Value(int64(*row.Size))
						}(),
						"status":   types.StringPointerValue(row.Status),
						"bootable": types.BoolPointerValue(row.Bootable),
						"image_id": func() attr.Value {
							// A volume created without an image has no image_id.
							if row.ImageId == nil {
								return types.Int64Null()
							}
							return types.Int64Value(int64(*row.ImageId))
						}(),
						"callback_url": types.StringPointerValue(row.CallbackUrl),
						"attachments": func() attr.Value {
							if row.Attachments == nil {
								return types.ListNull(datasource_core_volumes.AttachmentsValue{}.Type(ctx))
							}
							return d.MapAttachments(ctx, diags, *row.Attachments)
						}(),
						"created_at": func() attr.Value {
							if row.CreatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.CreatedAt.String())
						}(),
						"updated_at": func() attr.Value {
							if row.UpdatedAt == nil {
								return types.StringNull()
							}
							return types.StringValue(row.UpdatedAt.String())
						}(),
					},
				)
				diags.Append(diagnostic...)
				volumes = append(volumes, model)
			}
			return volumes
		}(),
	)
	diags.Append(diagnostic...)
	return model
}

func (d *DataSourceCoreVolumes) MapAttachments(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []volume.AttachmentsFieldsForVolume,
) attr.Value {
	result, diagnostic := types.ListValue(
		datasource_core_volumes.AttachmentsValue{}.Type(ctx),
		func() []attr.Value {
			attachments := make([]attr.Value, 0, len(data))
			for _, row := range data {
				model, diagnostic := datasource_core_volumes.NewAttachmentsValue(
					datasource_core_volumes.AttachmentsValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"id": func() attr.Value {
							if row.Id == nil {
								return types.Int64Null()
							}
							return types.Int64Value(int64(*row.Id))
						}(),
						"instance_id": func() attr.Value {
							if row.InstanceId == nil {
								return types.Int64Null()
							}
							return types.Int64Value(int64(*row.InstanceId))
						}(),
						"device":    types.StringPointerValue(row.Device),
						"status":    types.StringPointerValue(row.Status),
						"protected": types.BoolPointerValue(row.Protected),
					},
				)
				diags.Append(diagnostic...)
				attachments = append(attachments, model)
			}
			return attachments
		}(),
	)
	diags.Append(diagnostic...)

	return result
}
