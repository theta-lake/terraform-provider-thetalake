package supervisionspace

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
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
	plan := supervisionSpacePlanModel{}

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("all_participants"), &plan.AllParticipants)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("all_users"), &plan.AllUsers)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("directory_group_ids"), &plan.DirectoryGroupIds)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("external_id"), &plan.ExternalId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("hard_enforce"), &plan.HardEnforce)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("integration_ids"), &plan.IntegrationIds)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("media_types"), &plan.MediaTypes)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("retention_library_ids"), &plan.RetentionLibraryIds)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("requested_supervision_space_priority"), &plan.RequestedSupervisionSpacePriority)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("user_group_ids"), &plan.UserGroupIds)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("user_ids"), &plan.UserIds)...)

	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create Supervision Space", "Create failed to read Supervision Space resource plan data")
		return
	}

	// Map plan model to API model
	apiModel := toApiModel(&plan)

	space, err := r.client.CreateSupervisionSpace(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Supervision Space", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	state := fromApiModel(space)

	state.RequestedSupervisionSpacePriority = plan.RequestedSupervisionSpacePriority

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *supervisionSpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := supervisionSpaceStateModel{}

	// Read Terraform state data
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	space, err := r.client.GetSupervisionSpaceById(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read Supervision Space", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	updatedState := fromApiModel(space)

	// Preserve the requested priority from existing state so Terraform
	// continues to see the configured value, while the assigned
	// priority reflects what the API actually uses.
	updatedState.RequestedSupervisionSpacePriority = state.RequestedSupervisionSpacePriority

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *supervisionSpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan supervisionSpacePlanModel
	var state supervisionSpaceStateModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("all_participants"), &plan.AllParticipants)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("all_users"), &plan.AllUsers)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("directory_group_ids"), &plan.DirectoryGroupIds)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("external_id"), &plan.ExternalId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("hard_enforce"), &plan.HardEnforce)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("integration_ids"), &plan.IntegrationIds)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("media_types"), &plan.MediaTypes)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("retention_library_ids"), &plan.RetentionLibraryIds)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("requested_supervision_space_priority"), &plan.RequestedSupervisionSpacePriority)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("user_group_ids"), &plan.UserGroupIds)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("user_ids"), &plan.UserIds)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update Supervision Space", "Update failed to read Supervision Space resource plan data")
		return
	}

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	// Map plan model to API model
	apiModel := toApiModel(&plan)
	apiModel.Id = state.Id.ValueInt64()

	updatedSpace, err := r.client.UpdateSupervisionSpace(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Supervision Space", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	// Map API model to state model
	updatedState := fromApiModel(updatedSpace)

	updatedState.RequestedSupervisionSpacePriority = plan.RequestedSupervisionSpacePriority

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *supervisionSpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state supervisionSpaceStateModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	err := r.client.DeleteSupervisionSpace(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Supervision Space", fmt.Sprintf("Delete failed with error: %s", err.Error()))
		return
	}
}

// ImportState allows existing users to be brought under Terraform
// management by specifying their ID. The ID from the import command is
// mapped directly to the "id" attribute, after which Read will
// populate the rest of the state.
func (r *supervisionSpaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the string ID provided by Terraform into an int64 so it
	// matches the Int64 "id" attribute in the schema.
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric user ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
