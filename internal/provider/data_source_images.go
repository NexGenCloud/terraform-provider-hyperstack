package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nexgen/hyperstack-sdk-go/lib/image"
	"github.com/nexgen/hyperstack-terraform-provider/internal/client"
	"github.com/nexgen/hyperstack-terraform-provider/internal/genprovider/datasource_core_images"
	"io/ioutil"
)

var _ datasource.DataSource = &DataSourceCoreImages{}

func NewDataSourceCoreImages() datasource.DataSource {
	return &DataSourceCoreImages{}
}

func (d *DataSourceCoreImages) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_images"
}

func (d *DataSourceCoreImages) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_core_images.CoreImagesDataSourceSchema(ctx)
}

type DataSourceCoreImages struct {
	hyperstack *client.HyperstackClient
	client     *image.ClientWithResponses
}

func (d *DataSourceCoreImages) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	d.client, err = image.NewClientWithResponses(
		d.hyperstack.ApiServer,
		image.WithRequestEditorFn(d.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (d *DataSourceCoreImages) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_core_images.CoreImagesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Initialize the parameters as nil
	params := (*image.RetrieveImagesParams)(nil)
	result := (*image.RetrieveImagesResponse)(nil)
	err := error(nil)

	// If data.Region is not nil or empty, construct the parameters
	if !data.Region.IsNull() && data.Region.String() != "" {
		stringRegion := string(data.Region.ValueString())

		params = &image.RetrieveImagesParams{
			Region: &stringRegion,
		}
		result, err = d.client.RetrieveImagesWithResponse(ctx, params)
	} else {
		result, err = d.client.RetrieveImagesWithResponse(ctx, nil)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}
	bodyBytes, _ := ioutil.ReadAll(result.HTTPResponse.Body)
	if result.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("%s", string(bodyBytes)))
		return
	}

	callResult := result.JSON200.Images

	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No data in API result",
			fmt.Sprintf("%s", string(bodyBytes)),
		)
		return
	}

	data = d.ApiToModel(ctx, &resp.Diagnostics, callResult)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSourceCoreImages) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *[]image.ImageGetResponse,
) datasource_core_images.CoreImagesModel {
	return datasource_core_images.CoreImagesModel{
		CoreImages: d.MapImages(ctx, diags, *response),
	}
}

func (d *DataSourceCoreImages) MapImages(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []image.ImageGetResponse,
) types.Set {
	model, diagnostic := types.SetValue(
		datasource_core_images.CoreImagesValue{}.Type(ctx),
		func() []attr.Value {
			coreImages := make([]attr.Value, 0)
			for _, row := range data {
				logo := ""
				if row.Logo != nil {
					logo = *row.Logo
				}
				regionName := ""
				if row.RegionName != nil {
					regionName = *row.RegionName
				}
				imageType := ""
				if row.Type != nil {
					imageType = *row.Type
				}
				imagesList, _ := types.ListValue(
					datasource_core_images.ImagesValue{}.Type(ctx),
					func() []attr.Value {
						images := make([]attr.Value, 0)
						for _, imageItem := range *row.Images {
							displaySize := ""
							if imageItem.DisplaySize != nil {
								displaySize = *imageItem.DisplaySize
							}
							id := int64(0)
							if imageItem.Id != nil {
								id = int64(*imageItem.Id)
							}
							name := ""
							if imageItem.Name != nil {
								name = *imageItem.Name
							}
							regionName := ""
							if imageItem.RegionName != nil {
								regionName = *imageItem.RegionName
							}
							size := int64(0)
							if imageItem.Size != nil {
								size = int64(*imageItem.Size)
							}
							imageType := ""
							if imageItem.Type != nil {
								imageType = *imageItem.Type
							}
							version := ""
							if imageItem.Version != nil {
								version = *imageItem.Version
							}

							modelImage, diagnostic := datasource_core_images.NewImagesValue(
								datasource_core_images.ImagesValue{}.AttributeTypes(ctx),
								map[string]attr.Value{
									"display_size": types.StringValue(displaySize),
									"id":           types.Int64Value(id),
									"name":         types.StringValue(name),
									"region_name":  types.StringValue(regionName),
									"size":         types.Int64Value(size),
									"type":         types.StringValue(imageType),
									"version":      types.StringValue(version),
								},
							)
							images = append(images, modelImage)
							diags.Append(diagnostic...)
						}
						return images
					}())
				modelCoreImage, _ := datasource_core_images.NewCoreImagesValue(
					datasource_core_images.CoreImagesValue{}.AttributeTypes(ctx),
					map[string]attr.Value{
						"images":      imagesList,
						"logo":        types.StringValue(logo),
						"region_name": types.StringValue(regionName),
						"type":        types.StringValue(imageType),
					})
				coreImages = append(coreImages, modelCoreImage)
				//diags.Append(diagnostic...)
			}
			return coreImages
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
