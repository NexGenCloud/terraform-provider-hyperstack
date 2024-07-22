package provider

import (
	"context"
	"fmt"
	"github.com/NexGenCloud/hyperstack-sdk-go/lib/clusters"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/client"
	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/genprovider/resource_core_cluster"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"time"
)

var _ resource.Resource = &ResourceCoreCluster{}
var _ resource.ResourceWithImportState = &ResourceCoreCluster{}

func NewResourceCoreCluster() resource.Resource {
	return &ResourceCoreCluster{}
}

type ResourceCoreCluster struct {
	hyperstack *client.HyperstackClient
	client     *clusters.ClientWithResponses
}

func (r *ResourceCoreCluster) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_core_cluster"
}

func (r *ResourceCoreCluster) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_core_cluster.CoreClusterResourceSchema(ctx)
}

func (r *ResourceCoreCluster) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.hyperstack, _ = req.ProviderData.(*client.HyperstackClient)

	var err error
	r.client, err = clusters.NewClientWithResponses(
		r.hyperstack.ApiServer,
		clusters.WithRequestEditorFn(r.hyperstack.GetAddHeadersFn()),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected error",
			fmt.Sprintf("%s", err),
		)
		return
	}
}

func (r *ResourceCoreCluster) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var dataOld resource_core_cluster.CoreClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &dataOld)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.CreateClusterWithResponse(ctx, func() clusters.CreateClusterJSONRequestBody {
		return clusters.CreateClusterJSONRequestBody{
			EnablePublicIp:    dataOld.EnablePublicIp.ValueBoolPointer(),
			EnvironmentName:   dataOld.EnvironmentName.ValueString(),
			ImageName:         dataOld.ImageName.ValueString(),
			KeypairName:       dataOld.KeypairName.ValueString(),
			KubernetesVersion: dataOld.KubernetesVersion.ValueString(),
			MasterFlavorName:  dataOld.MasterFlavorName.ValueString(),
			Name:              dataOld.Name.ValueString(),
			NodeCount:         int(dataOld.NodeCount.ValueInt64()),
			NodeFlavorName:    dataOld.NodeFlavorName.ValueString(),
		}
	}())
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if result.JSON201 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	callResult := result.JSON201.Cluster
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No dataOld in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	clusterModel := callResult

	id := *clusterModel.Id
	err = r.WaitForResult(
		ctx,
		3*time.Second,
		900*time.Second,
		func(ctx context.Context) (bool, error) {
			result, err := r.client.GettingClusterDetailWithResponse(ctx, id)
			if err != nil {
				return false, err
			}

			if result.JSON200 == nil {
				return false, fmt.Errorf("Wrong API result: %s", result.StatusCode())
			}

			status := *result.JSON200.Cluster.Status
			if status == "CREATING" {
				return false, nil
			}

			clusterModel = result.JSON200.Cluster
			return true, nil
		},
	)

	// TODO: doesn't save resource info in state
	if err != nil {
		resp.Diagnostics.AddError(
			"Waiting for state change error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	// TODO: doesn't save resource info in state
	if *clusterModel.Status == "ERROR" {
		resp.Diagnostics.AddError(
			"Failed creating instance: status %s",
			*clusterModel.Status,
		)
		return
	}

	data := r.ApiToModel(ctx, &resp.Diagnostics, callResult)
	r.MergeData(&data, &dataOld)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreCluster) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var dataOld resource_core_cluster.CoreClusterModel
	resp.Diagnostics.Append(req.State.Get(ctx, &dataOld)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.GettingClusterDetailWithResponse(ctx, int(dataOld.Id.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if result.JSON404 != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if result.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	callResult := result.JSON200.Cluster
	if callResult == nil {
		resp.Diagnostics.AddWarning(
			"No dataOld in API result",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	data := r.ApiToModel(ctx, &resp.Diagnostics, callResult)
	r.MergeData(&data, &dataOld)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// WaitForResult TODO: deduplicate
func (r *ResourceCoreCluster) WaitForResult(
	ctx context.Context,
	pollInterval,
	timeout time.Duration,
	checker func(ctx context.Context) (bool, error),
) error {
	timeoutTimer := time.NewTimer(timeout)
	pollTicker := time.NewTicker(pollInterval)

	defer timeoutTimer.Stop()
	defer pollTicker.Stop()

	for {
		select {
		case <-timeoutTimer.C:
			return fmt.Errorf(
				"Timeout %s reached waiting for resource state change",
				timeout.String(),
			)
		case <-pollTicker.C:
			status, err := checker(ctx)
			if err != nil {
				return fmt.Errorf(
					"Error calling checker while waiting: %e",
					err,
				)
			}
			if status {
				return nil // Resource is created
			}
			// Continue waiting
		}
	}
}

func (r *ResourceCoreCluster) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var dataOld resource_core_cluster.CoreClusterModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &dataOld)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data resource_core_cluster.CoreClusterModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.MergeData(&data, &dataOld)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceCoreCluster) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_core_cluster.CoreClusterModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := int(data.Id.ValueInt64())

	result, err := r.client.DeleteAClusterWithResponse(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"API request error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	if result.JSON200 == nil {
		resp.Diagnostics.AddError(
			"Wrong API json response",
			fmt.Sprintf("Status: %s, body: %s", result.Status(), string(result.Body)),
		)
		return
	}

	err = r.WaitForResult(
		ctx,
		3*time.Second,
		120*time.Second,
		func(ctx context.Context) (bool, error) {
			result, err := r.client.GettingClusterDetailWithResponse(ctx, id)
			if err != nil {
				return false, err
			}

			return result.JSON404 != nil, nil
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Waiting for state change error",
			fmt.Sprintf("%s", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *ResourceCoreCluster) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseFloat(req.ID, 64)

	if err != nil {
		resp.Diagnostics.AddError(
			"Wrong ID specified",
			req.ID,
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *ResourceCoreCluster) MergeData(
	data *resource_core_cluster.CoreClusterModel,
	dataOld *resource_core_cluster.CoreClusterModel,
) {
	// Assign all values that are available only during creation stage
	data.ImageName = dataOld.ImageName
	data.NodeFlavorName = dataOld.NodeFlavorName
	data.MasterFlavorName = dataOld.MasterFlavorName
}

func (r *ResourceCoreCluster) ApiToModel(
	ctx context.Context,
	diags *diag.Diagnostics,
	response *clusters.ClusterFields,
) resource_core_cluster.CoreClusterModel {
	return resource_core_cluster.CoreClusterModel{
		ApiAddress: types.StringPointerValue(response.ApiAddress),
		CreatedAt: func() types.String {
			if response.CreatedAt == nil {
				return types.StringNull()
			}
			return types.StringValue(response.CreatedAt.String())
		}(),
		EnablePublicIp:  types.BoolPointerValue(response.EnablePublicIp),
		EnvironmentName: types.StringPointerValue(response.EnvironmentName),
		Id: func() types.Int64 {
			if response.Id == nil {
				return types.Int64Null()
			}
			return types.Int64Value(int64(*response.Id))
		}(),
		ImageName:         types.StringNull(),
		KeypairName:       types.StringPointerValue(response.KeypairName),
		KubeConfig:        types.StringPointerValue(response.KubeConfig),
		KubernetesVersion: types.StringPointerValue(response.KubernetesVersion),
		MasterFlavorName:  types.StringNull(),
		Name:              types.StringPointerValue(response.Name),
		NodeAddresses:     r.MapNodeAddresses(ctx, diags, *response.NodeAddresses),
		NodeCount: func() types.Int64 {
			if response.NodeCount == nil {
				return types.Int64Null()
			}
			return types.Int64Value(int64(*response.NodeCount))
		}(),
		NodeFlavor:     resource_core_cluster.NodeFlavorValue{},
		NodeFlavorName: types.StringNull(),
		Status:         types.StringPointerValue(response.Status),
		StatusReason:   types.StringPointerValue(response.StatusReason),
	}
}

func (r *ResourceCoreCluster) MapNodeAddresses(
	ctx context.Context,
	diags *diag.Diagnostics,
	data []string,
) types.List {
	model, diagnostic := types.ListValue(
		types.StringType,
		func() []attr.Value {
			labels := make([]attr.Value, 0)
			for _, row := range data {
				model := types.StringValue(row)
				labels = append(labels, model)
			}
			return labels
		}(),
	)
	diags.Append(diagnostic...)
	return model
}
