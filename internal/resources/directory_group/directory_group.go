package directorygroup

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type directoryGroupResource struct {
	client *thetalake.Client
}

func NewDirectoryGroupResource() resource.Resource {
	return &directoryGroupResource{}
}

func (r *directoryGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_group", req.ProviderTypeName)
}

func (r *directoryGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*thetalake.Client)
}

func (r *directoryGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan directoryGroupPlanModel

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("external_id"), &plan.ExternalId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("identity_ids"), &plan.IdentityIds)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create Directory Group", "Create failed to read Directory Group resource plan data")
		return
	}

	apiModel := toApiModel(&plan)

	dg, err := r.client.CreateDirectoryGroup(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Directory Group", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	state := fromApiModel(dg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *directoryGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state directoryGroupStateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dg, err := r.client.GetDirectoryGroupById(ctx, state.Id.ValueInt64())
	if err != nil {
		if errors.Is(err, thetalake.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read Directory Group", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(dg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *directoryGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan directoryGroupPlanModel
	var state directoryGroupStateModel

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("external_id"), &plan.ExternalId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("identity_ids"), &plan.IdentityIds)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update Directory Group", "Update failed to read Directory Group resource plan data")
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiModel := toApiModel(&plan)
	apiModel.Id = state.Id.ValueInt64()

	updatedDg, err := r.client.UpdateDirectoryGroup(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Directory Group", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(updatedDg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *directoryGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state directoryGroupStateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDirectoryGroup(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Directory Group", fmt.Sprintf("Delete failed with error: %s", err.Error()))
	}
}

func (r *directoryGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric directory group ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
