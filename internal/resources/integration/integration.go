package integration

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type integrationResource struct {
	client *thetalake.Client
}

func NewIntegrationResource() resource.Resource {
	return &integrationResource{}
}

func (r *integrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_integration", req.ProviderTypeName)
}

func (r *integrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*thetalake.Client)
}

func (r *integrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := integrationPlanModel{}

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("generic_journaling"), &plan.GenericJournaling)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("google_workspace_email"), &plan.GoogleWorkspaceEmail)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("paused"), &plan.Paused)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("theta_lake_api"), &plan.ThetaLakeApi)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create Integration", "Create failed to read Integration resource plan data")
		return
	}

	apiModel, diags := toApiModel(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create Integration", "Create failed to map Integration plan data")
		return
	}

	createdIntegration, err := r.client.CreateIntegration(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Integration", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	// Unlike Update, Create has no prior state to fall back on for the id: if the
	// create response omits it, there is no correct id to track this resource by (or
	// to query GetIntegrationConfiguration with), so surface that explicitly instead
	// of silently querying /integrations/0/configuration.
	if createdIntegration.Id == 0 {
		resp.Diagnostics.AddError(
			"Failed to create Integration",
			"Create response did not include an integration id; the resource cannot be reliably tracked.",
		)
		return
	}

	options, err := r.resolveServiceParams(ctx, createdIntegration.Id, createdIntegration.ServiceParams, apiModel.Type)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Integration", fmt.Sprintf("Create failed to get Integration configuration: %s", err.Error()))
		return
	}

	// Use apiModel.Type (known authoritatively from the plan via toApiModel) rather than
	// re-deriving the type from the response's integration_type/integration_type_id, which
	// may not resolve through thetalake.IntegrationTypeSlug the same way for every server.
	state := fromApiModel(createdIntegration, options, apiModel.Type)

	// paused is sent in the create request, so the plan's value is authoritative; carry
	// it over from the plan (as Update does for its own reasons below) rather than
	// trusting a create response that might omit service_paused, which would otherwise
	// silently decode to false and produce an inconsistent-result error when paused was
	// requested true.
	state.Paused = plan.Paused

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// resolveServiceParams returns serviceParams if the create/update response included
// them, falling back to a GetIntegrationConfiguration call when it didn't. This avoids
// ever having to fall back to the plan's (possibly-unknown, in a real Terraform plan)
// type block to populate state. theta_lake_api has no configuration options, so the
// round-trip is skipped for that type; fromApiModel populates its (always-empty) block
// regardless of whether options is nil.
//
// id must be the caller's authoritative integration id (e.g. prior state's id in
// Update), not necessarily the create/update response's own id field: that field is
// itself only Computed with UseStateForUnknown, so an omitted response id decodes to
// 0, and querying /integrations/0/configuration would fail even though the update
// otherwise succeeded.
func (r *integrationResource) resolveServiceParams(ctx context.Context, id int64, serviceParams *thetalake.IntegrationOptions, typeSlug string) (*thetalake.IntegrationOptions, error) {
	if serviceParams != nil || typeSlug == thetalake.IntegrationTypeThetaLakeApi {
		return serviceParams, nil
	}

	configuration, err := r.client.GetIntegrationConfiguration(ctx, id)
	if err != nil {
		return nil, err
	}

	return &configuration.Options, nil
}

func (r *integrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := integrationStateModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	apiIntegration, err := r.client.GetIntegrationById(ctx, state.Id.ValueInt64())
	if err != nil {
		if errors.Is(err, thetalake.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read Integration", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	typeSlug := thetalake.IntegrationTypeSlug(apiIntegration.IntegrationTypeId, apiIntegration.IntegrationType)
	if typeSlug == "" {
		resp.Diagnostics.AddError(
			"Unsupported Integration Type",
			fmt.Sprintf(
				"Integration %d has integration_type %q (integration_type_id %d), which is not one of the types this provider supports.",
				state.Id.ValueInt64(), apiIntegration.IntegrationType, apiIntegration.IntegrationTypeId,
			),
		)
		return
	}

	// theta_lake_api has no configuration options, so skip the round-trip entirely for
	// that type; fromApiModel doesn't use options for it regardless.
	var options *thetalake.IntegrationOptions
	if typeSlug != thetalake.IntegrationTypeThetaLakeApi {
		configuration, err := r.client.GetIntegrationConfiguration(ctx, state.Id.ValueInt64())
		if err != nil {
			if errors.Is(err, thetalake.ErrNotFound) {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("Failed to read Integration", fmt.Sprintf("Read failed to get Integration configuration: %s", err.Error()))
			return
		}
		options = &configuration.Options
	}

	updatedState := fromApiModel(apiIntegration, options, typeSlug)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *integrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan integrationPlanModel
	var state integrationStateModel

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("generic_journaling"), &plan.GenericJournaling)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("google_workspace_email"), &plan.GoogleWorkspaceEmail)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("paused"), &plan.Paused)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("theta_lake_api"), &plan.ThetaLakeApi)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update Integration", "Update failed to read Integration resource plan data")
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
		resp.Diagnostics.AddError("Failed to update Integration", "Update failed to map Integration plan data")
		return
	}
	apiModel.Id = state.Id.ValueInt64()

	updatedIntegration, err := r.client.UpdateIntegration(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Integration", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	// Use the prior state's id (already known authoritative, and what apiModel.Id was
	// set to for the PUT request above), not updatedIntegration.Id: id has
	// UseStateForUnknown, so a PUT response that omits it decodes to 0, and querying
	// /integrations/0/configuration would fail even though the update itself succeeded.
	options, err := r.resolveServiceParams(ctx, state.Id.ValueInt64(), updatedIntegration.ServiceParams, apiModel.Type)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Integration", fmt.Sprintf("Update failed to get Integration configuration: %s", err.Error()))
		return
	}

	// Use apiModel.Type (known authoritatively from the plan) rather than re-deriving
	// the type from the response, same as Create.
	updatedState := fromApiModel(updatedIntegration, options, apiModel.Type)

	// created_at, id, integration_group, integration_type, and integration_type_id are
	// all immutable after creation and Computed with UseStateForUnknown, so the plan
	// already carries a known prior value for each; take them from prior state rather
	// than the PUT response in case the response omits any of them.
	preserveImmutableComputedFieldsFromState(&updatedState, &state)

	// UpdateIntegrationRequest has no paused field, so this update can't have changed
	// the paused state: carry the prior state's value rather than the PUT response's
	// service_paused. Persisting that value means a subsequent pause/start failure
	// records the still-accurate paused state instead of the requested one, while
	// still matching the plan whenever plan and prior state already agree (in which
	// case the pause/start branch below is skipped).
	updatedState.Paused = state.Paused

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Paused.ValueBool() != state.Paused.ValueBool() {
		if plan.Paused.ValueBool() {
			err = r.client.PauseIntegration(ctx, state.Id.ValueInt64())
		} else {
			err = r.client.StartIntegration(ctx, state.Id.ValueInt64())
		}
		if err != nil {
			resp.Diagnostics.AddError("Failed to update Integration", fmt.Sprintf("Update failed to set paused state: %s", err.Error()))
			return
		}

		updatedState.Paused = plan.Paused
		resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
	}
}

func (r *integrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state integrationStateModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	err := r.client.DeleteIntegration(ctx, state.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Integration", fmt.Sprintf("Delete failed with error: %s", err.Error()))
	}
}

func (r *integrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric Integration ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
