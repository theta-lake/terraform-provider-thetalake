package swrvrule

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type swrvRuleResource struct {
	client *thetalake.Client
}

func NewSwrvRuleResource() resource.Resource {
	return &swrvRuleResource{}
}

func (r *swrvRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_swrv_rule", req.ProviderTypeName)
}

func (r *swrvRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*thetalake.Client)
}

func (r *swrvRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := swrvRulePlanModel{}

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("input_sources"), &plan.InputSources)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("policy_id"), &plan.PolicyId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("priority"), &plan.Priority)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("retention_library_id"), &plan.RetentionLibraryId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("supervision_space_id"), &plan.SupervisionSpaceId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("workflow_id"), &plan.WorkflowId)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create SWRV Rule", "Create failed to read SWRV rule resource plan data")
		return
	}

	apiModel, diags := toApiModel(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create SWRV Rule", "Create failed to map SWRV rule plan data")
		return
	}

	rule, err := r.client.CreateSwrvRule(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create SWRV Rule", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	state := fromApiModel(rule)
	state.InputSources = plan.InputSources

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *swrvRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := swrvRuleStateModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	rule, err := r.client.GetSwrvRuleById(ctx, state.Id.ValueInt64())
	if err != nil {
		if errors.Is(err, thetalake.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read SWRV Rule", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(rule)
	if !inputSourcesHaveValues(updatedState.InputSources) && inputSourcesHaveValues(state.InputSources) {
		updatedState.InputSources = state.InputSources
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *swrvRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan swrvRulePlanModel
	var state swrvRuleStateModel

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("input_sources"), &plan.InputSources)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("policy_id"), &plan.PolicyId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("priority"), &plan.Priority)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("retention_library_id"), &plan.RetentionLibraryId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("supervision_space_id"), &plan.SupervisionSpaceId)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("workflow_id"), &plan.WorkflowId)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update SWRV Rule", "Update failed to read SWRV rule resource plan data")
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	apiModel, diags := toApiModel(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update SWRV Rule", "Update failed to map SWRV rule plan data")
		return
	}
	apiModel.Id = state.Id.ValueInt64()

	rule, err := r.client.UpdateSwrvRule(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update SWRV Rule", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(rule)
	updatedState.InputSources = plan.InputSources

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *swrvRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state swrvRuleStateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	err := r.client.DeleteSwrvRule(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete SWRV Rule", fmt.Sprintf("Delete failed with error: %s", err.Error()))
	}
}

func (r *swrvRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric SWRV rule ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
