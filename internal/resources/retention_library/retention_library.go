package retentionlibrary

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type retentionLibraryResource struct {
	client *thetalake.Client
}

func NewRetentionLibraryResource() resource.Resource {
	return &retentionLibraryResource{}
}

func (r *retentionLibraryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_retention_library", req.ProviderTypeName)
}

func (r *retentionLibraryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*thetalake.Client)
}

func (r *retentionLibraryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := retentionLibraryPlanModel{}

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("external_id"), &plan.ExternalId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("retain_in_review"), &plan.RetainInReview)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("retention_period_days"), &plan.RetentionPeriodDays)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("retention_period_enabled"), &plan.RetentionPeriodEnabled)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("sec_compliant_storage_enabled"), &plan.SecCompliantStorageEnabled)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("storage_account_id"), &plan.StorageAccountId)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create Retention Library", "Create failed to read Retention Library resource plan data")
		return
	}

	apiModel := toApiModel(&plan)

	library, err := r.client.CreateRetentionLibrary(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Retention Library", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	state := fromApiModel(library)
	// retain_in_review is not returned by the API; preserve the plan value
	state.RetainInReview = plan.RetainInReview

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *retentionLibraryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := retentionLibraryStateModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	library, err := r.client.GetRetentionLibraryById(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read Retention Library", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(library)
	// retain_in_review is not returned by the API; preserve the existing state value
	updatedState.RetainInReview = state.RetainInReview

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *retentionLibraryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan retentionLibraryPlanModel
	var state retentionLibraryStateModel

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("external_id"), &plan.ExternalId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("retain_in_review"), &plan.RetainInReview)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("retention_period_days"), &plan.RetentionPeriodDays)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("retention_period_enabled"), &plan.RetentionPeriodEnabled)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("sec_compliant_storage_enabled"), &plan.SecCompliantStorageEnabled)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("storage_account_id"), &plan.StorageAccountId)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update Retention Library", "Update failed to read Retention Library resource plan data")
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	apiModel := toApiModel(&plan)
	apiModel.Id = state.Id.ValueInt64()

	library, err := r.client.UpdateRetentionLibrary(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Retention Library", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(library)
	// retain_in_review is not returned by the API; preserve the plan value
	updatedState.RetainInReview = plan.RetainInReview

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *retentionLibraryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state retentionLibraryStateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	err := r.client.DeleteRetentionLibrary(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Retention Library", fmt.Sprintf("Delete failed with error: %s", err.Error()))
		return
	}
}

// ImportState allows existing retention libraries to be brought under Terraform
// management by specifying their ID.
func (r *retentionLibraryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric retention library ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
