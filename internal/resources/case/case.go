package legalcase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type caseResource struct {
	client *thetalake.Client
}

func NewCaseResource() resource.Resource {
	return &caseResource{}
}

func (r *caseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_case", req.ProviderTypeName)
}

func (r *caseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*thetalake.Client)
}

func (r *caseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan casePlanModel

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("close_date"), &plan.CloseDate)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("number"), &plan.Number)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("open_date"), &plan.OpenDate)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("visibility"), &plan.Visibility)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("manager_ids"), &plan.ManagerIds)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create Case", "Create failed to read Case resource plan data")
		return
	}

	apiModel := toApiModel(&plan)

	c, err := r.client.CreateCase(ctx, apiModel)
	if c.Id != 0 {
		// The case was created on the server even though err may be set (e.g. a
		// manager failed to attach), so persist it now to avoid leaving an
		// orphaned, untracked case behind.
		state := fromApiModel(c)
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	}
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Case", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	c, err = r.reconcileCloseState(ctx, c, plan.CloseDate)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Case", fmt.Sprintf("Create failed to set case close state: %s", err.Error()))
		return
	}

	state := fromApiModel(c)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *caseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state caseStateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, err := r.client.GetCaseById(ctx, state.Id.ValueInt64())
	if err != nil {
		if errors.Is(err, thetalake.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read Case", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(c)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *caseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan casePlanModel
	var state caseStateModel

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("close_date"), &plan.CloseDate)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("number"), &plan.Number)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("open_date"), &plan.OpenDate)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("visibility"), &plan.Visibility)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("manager_ids"), &plan.ManagerIds)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update Case", "Update failed to read Case resource plan data")
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiModel := toApiModel(&plan)
	apiModel.Id = state.Id.ValueInt64()

	updatedCase, err := r.client.UpdateCase(ctx, apiModel)
	if updatedCase.Id != 0 {
		// The update was applied on the server even though err may be set (e.g.
		// a manager failed to attach/detach), so persist it now rather than
		// leaving the state stale relative to the backend.
		updatedState := fromApiModel(updatedCase)
		resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
	}
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Case", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	updatedCase, err = r.reconcileCloseState(ctx, updatedCase, plan.CloseDate)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Case", fmt.Sprintf("Update failed to set case close state: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(updatedCase)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

// reconcileCloseState closes or reopens a case so its close state matches
// desiredCloseDate. The close and reopen endpoints don't return case
// managers (unlike get-by-id), so the managers from c are carried over onto
// the result. An unknown desiredCloseDate (e.g. it references an attribute
// of another resource not yet known at plan time) is left untouched.
func (r *caseResource) reconcileCloseState(ctx context.Context, c thetalake.Case, desiredCloseDate timetypes.RFC3339) (thetalake.Case, error) {
	if desiredCloseDate.IsUnknown() {
		return c, nil
	}

	if desiredCloseDate.IsNull() {
		if c.Status != "CLOSED" {
			return c, nil
		}

		reopened, err := r.client.ReopenCase(ctx, c.Id)
		if err != nil {
			return c, err
		}
		reopened.ManagerIds = c.ManagerIds
		return reopened, nil
	}

	desired, diags := desiredCloseDate.ValueRFC3339Time()
	if diags.HasError() {
		return c, nil
	}

	if c.CloseDate != nil && c.CloseDate.Equal(desired) {
		return c, nil
	}

	closed, err := r.client.CloseCase(ctx, c.Id, desired)
	if err != nil {
		return c, err
	}
	closed.ManagerIds = c.ManagerIds
	return closed, nil
}

func (r *caseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state caseStateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCase(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Case", fmt.Sprintf("Delete failed with error: %s", err.Error()))
	}
}

func (r *caseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric case ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
