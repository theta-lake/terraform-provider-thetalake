package supervisionspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type supervisionSpaceResource struct {
	client *thetalake.Client
}

func NewSupervisionSpaceResource() resource.Resource {
	return &supervisionSpaceResource{}
}

func (r *supervisionSpaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_supervision_space", req.ProviderTypeName)
}

func (r *supervisionSpaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*thetalake.Client)
}

func (r *supervisionSpaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// plan := &supervisionSpaceModel{}

	// // Read Terraform plan data
	// resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	// if resp.Diagnostics.HasError() {
	// 	resp.Diagnostics.AddError("Internal Error", "Failed to read plan data")
	// 	return
	// }

	// apiDirectoryGroupIds := []int{}
	// resp.Diagnostics.Append(
	// 	plan.DirectoryGroupIds.ElementsAs(ctx, &apiDirectoryGroupIds, false)...,
	// )
	// if resp.Diagnostics.HasError() {
	// 	resp.Diagnostics.AddError("Internal Error", "Failed to iterate over directory_group_ids")
	// 	return
	// }

	// apiIntegrationIds := []int{}
	// resp.Diagnostics.Append(
	// 	plan.IntegrationIds.ElementsAs(ctx, &apiIntegrationIds, false)...,
	// )
	// if resp.Diagnostics.HasError() {
	// 	resp.Diagnostics.AddError("Internal Error", "Failed to iterate over integration_ids")
	// 	return
	// }

	// apiMediaTypeIds := []int{}
	// resp.Diagnostics.Append(
	// 	plan.MediaTypeIds.ElementsAs(ctx, &apiMediaTypeIds, false)...,
	// )
	// if resp.Diagnostics.HasError() {
	// 	resp.Diagnostics.AddError("Internal Error", "Failed to iterate over media_type_ids")
	// 	return
	// }

	// apiRetentionLibraryIds := []int{}
	// resp.Diagnostics.Append(
	// 	plan.RetentionLibraryIds.ElementsAs(ctx, &apiRetentionLibraryIds, false)...,
	// )
	// if resp.Diagnostics.HasError() {
	// 	resp.Diagnostics.AddError("Internal Error", "Failed to iterate over retention_library_ids")
	// 	return
	// }

	// createSupervisionSpaceReq := &apiCreateSupervisionSpaceRequest{
	// 	AllParticipants:          plan.AllParticipants.ValueBool(),
	// 	AllUsers:                 plan.AllParticipants.ValueBool(),
	// 	Description:              plan.Description.ValueString(),
	// 	DirectoryGroupIds:        apiDirectoryGroupIds,
	// 	ExternalId:               plan.ExternalId.ValueString(),
	// 	HardEnforce:              plan.HardEnforce.ValueBool(),
	// 	IntegrationIds:           apiIntegrationIds,
	// 	MediaTypeIds:             apiMediaTypeIds,
	// 	Name:                     plan.Name.ValueString(),
	// 	RetentionLibraryIds:      apiRetentionLibraryIds,
	// 	SupervisionSpacePriority: int(plan.SupervisionSpacePriority.ValueInt64()),
	// }

	// body, err := json.Marshal(createSupervisionSpaceReq)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Bad Request", fmt.Sprintf("Failed to marshal request: %s", err.Error()))
	// 	return
	// }

	// endpoint := fmt.Sprintf("%s/supervision_spaces", r.apiClient.endpoint)
	// token := fmt.Sprintf("Bearer %s", r.apiClient.token)

	// supervisionSpaceReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	// if err != nil {
	// 	resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Failed to create request: %s", err.Error()))
	// 	return
	// }
	// supervisionSpaceReq.Header.Set("Authorization", token)
	// supervisionSpaceReq.Header.Set("Content-Type", "application/json")

	// supervisionSpaceResp, err := r.apiClient.client.Do(supervisionSpaceReq)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d, Error: %s", supervisionSpaceResp.StatusCode, err.Error()))
	// 	return
	// }

	// if supervisionSpaceResp.StatusCode != http.StatusOK {
	// 	resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d", supervisionSpaceResp.StatusCode))
	// 	return
	// }

	// createSupervisionSpaceResp := &apiCreateSupervisionSpaceResponse{}

	// if err := json.NewDecoder(supervisionSpaceResp.Body).Decode(createSupervisionSpaceResp); err != nil {
	// 	resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Failed to decode response: %s", err.Error()))
	// 	return
	// }

	// state := &supervisionSpaceModel{
	// 	AllParticipants:          types.BoolValue(createSupervisionSpaceResp.SupervisionSpace.AllParticipants),
	// 	AllUsers:                 types.BoolValue(createSupervisionSpaceResp.SupervisionSpace.AllUsers),
	// 	Description:              types.StringValue(createSupervisionSpaceResp.SupervisionSpace.Description),
	// 	DirectoryGroupIds:        intSlicetoInt64List(GetSupervisionSpaceDirectoryGroupIds(createSupervisionSpaceResp.SupervisionSpace.DirectoryGroups)),
	// 	ExternalId:               types.StringValue(createSupervisionSpaceResp.SupervisionSpace.ExternalId),
	// 	HardEnforce:              types.BoolValue(createSupervisionSpaceResp.SupervisionSpace.HardEnforce),
	// 	ID:                       types.Int64Value(int64(createSupervisionSpaceResp.SupervisionSpace.ID)),
	// 	IntegrationIds:           intSlicetoInt64List(GetSupervisionSpaceIntegrationIds(createSupervisionSpaceResp.SupervisionSpace.Integrations)),
	// 	MediaTypeIds:             intSlicetoInt64List(GetSupervisionSpaceMediaTypeIds(createSupervisionSpaceResp.SupervisionSpace.MediaTypes)),
	// 	Name:                     types.StringValue(createSupervisionSpaceResp.SupervisionSpace.Name),
	// 	RetentionLibraryIds:      intSlicetoInt64List(GetSupervisionSpaceRetentionLibraryIds(createSupervisionSpaceResp.SupervisionSpace.RetentionLibraries)),
	// 	SupervisionSpacePriority: types.Int64Value(int64(createSupervisionSpaceResp.SupervisionSpace.SupervisionSpacePriority)),
	// }

	// // Save data into Terraform state
	// resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

}

func (r *supervisionSpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *supervisionSpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *supervisionSpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// state := &supervisionSpaceModel{}

	// // Read Terraform state data
	// resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	// if resp.Diagnostics.HasError() {
	// 	resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
	// 	return
	// }

	// if state.ID.IsNull() || state.ID.IsUnknown() {
	// 	resp.Diagnostics.AddError("Internal Error", "Failed to delete the supervision_space: the id is missing in state")
	// 	return
	// }

	// endpoint := fmt.Sprintf("%s/supervision_spaces/%v", r.apiClient.endpoint, state.ID.ValueInt64())
	// token := fmt.Sprintf("Bearer %s", r.apiClient.token)

	// supervisionSpaceReq, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Failed to create request: %s", err.Error()))
	// 	return
	// }
	// supervisionSpaceReq.Header.Set("Authorization", token)
	// supervisionSpaceReq.Header.Set("Content-Type", "application/json")

	// supervisionSpaceResp, err := r.apiClient.client.Do(supervisionSpaceReq)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d, Error: %s", supervisionSpaceResp.StatusCode, err.Error()))
	// 	return
	// }

	// if supervisionSpaceResp.StatusCode != http.StatusOK {
	// 	resp.Diagnostics.AddError("Internal Error", fmt.Sprintf("Status: %d", supervisionSpaceResp.StatusCode))
	// 	return
	// }

	// resp.State.RemoveResource(ctx)
}

// func GetSupervisionSpaceDirectoryGroupIds(groups []apiSupervisionSpaceDirectoryGroup) []int {
// 	ids := make([]int, len(groups))
// 	for i, group := range groups {
// 		ids[i] = group.ID
// 	}
// 	return ids
// }

// func GetSupervisionSpaceIntegrationIds(groups []apiSupervisionSpaceIntegration) []int {
// 	ids := make([]int, len(groups))
// 	for i, group := range groups {
// 		ids[i] = group.ID
// 	}
// 	return ids
// }

// func GetSupervisionSpaceMediaTypeIds(groups []apiSupervisionSpaceMediaType) []int {
// 	ids := make([]int, len(groups))
// 	for i, group := range groups {
// 		ids[i] = group.ID
// 	}
// 	return ids
// }

// func GetSupervisionSpaceRetentionLibraryIds(groups []apiSupervisionSpaceRetentionLibrary) []int {
// 	ids := make([]int, len(groups))
// 	for i, group := range groups {
// 		ids[i] = group.ID
// 	}
// 	return ids
// }
