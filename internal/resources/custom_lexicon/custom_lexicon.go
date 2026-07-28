package customlexicon

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type customLexiconResource struct {
	client *thetalake.Client
}

func NewCustomLexiconResource() resource.Resource {
	return &customLexiconResource{}
}

func (r *customLexiconResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_custom_lexicon", req.ProviderTypeName)
}

func (r *customLexiconResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*thetalake.Client)
}

// ValidateConfig rejects disabled=true combined with a non-empty policy_ids:
// disabling a lexicon removes all of its policy associations, so configuring
// both is contradictory.
func (r *customLexiconResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config customLexiconModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Disabled.IsUnknown() || config.PolicyIds.IsUnknown() {
		return
	}

	if config.Disabled.ValueBool() && !config.PolicyIds.IsNull() && len(config.PolicyIds.Elements()) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("policy_ids"),
			"Invalid Custom Lexicon Configuration",
			"policy_ids must be empty when disabled = true: disabling a custom lexicon removes all policies associated with it.",
		)
	}

	// If rule_scope is explicitly set, it must include "email" when enabling email analysis.
	if !config.RuleScope.IsNull() && !config.RuleScope.IsUnknown() {
		emailSmart := !config.EmailSmartBody.IsNull() && !config.EmailSmartBody.IsUnknown() && config.EmailSmartBody.ValueBool()
		emailSubject := !config.EmailSubjectAnalyzed.IsNull() && !config.EmailSubjectAnalyzed.IsUnknown() && config.EmailSubjectAnalyzed.ValueBool()
		if emailSmart || emailSubject {
			hasEmail := slices.Contains(stringSetToSlice(config.RuleScope), "email")
			if !hasEmail {
				resp.Diagnostics.AddAttributeError(
					path.Root("rule_scope"),
					"Invalid Custom Lexicon Configuration",
					"rule_scope must include \"email\" when email_smart_body or email_subject_analyzed is true.",
				)
			}
		}
	}
}

func (r *customLexiconResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customLexiconModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to create Custom Lexicon", "Create failed to read Custom Lexicon resource plan data")
		return
	}

	createRequest := toCreateRequest(&plan)

	lexicon, err := r.client.CreateCustomLexicon(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Custom Lexicon", fmt.Sprintf("Create failed with error: %s", err.Error()))
		return
	}

	// The create endpoint cannot set disabled=true; if the plan requests a
	// disabled lexicon, follow up with an explicit disable call. StartDate and
	// EndDate must be carried forward explicitly.
	if plan.Disabled.ValueBool() {
		disabled := true
		lexicon, err = r.client.UpdateCustomLexicon(ctx, lexicon.Id, thetalake.UpdateCustomLexiconRequest{
			Disabled:  &disabled,
			StartDate: formatDatePtr(lexicon.StartDate),
			EndDate:   formatDatePtr(lexicon.EndDate),
		})
		if err != nil {
			resp.Diagnostics.AddError("Failed to create Custom Lexicon", fmt.Sprintf("Create succeeded but failed to disable the new lexicon: %s", err.Error()))
			return
		}
	}

	state := fromApiModel(lexicon)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *customLexiconResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state customLexiconModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	lexicon, err := r.client.GetCustomLexiconById(ctx, state.Id.ValueInt64())
	if err != nil {
		if errors.Is(err, thetalake.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read Custom Lexicon", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(lexicon)

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *customLexiconResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan customLexiconModel
	var state customLexiconModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update Custom Lexicon", "Update failed to read Custom Lexicon resource plan data")
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	updateRequest := toUpdateRequest(&plan, &state)

	lexicon, err := r.client.UpdateCustomLexicon(ctx, state.Id.ValueInt64(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Custom Lexicon", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(lexicon)

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *customLexiconResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customLexiconModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	// The API has no delete endpoint; teardown is done by disabling the
	// lexicon. An already-disabled lexicon cannot be disabled again, so skip
	// the call in that case.
	if state.Disabled.ValueBool() {
		return
	}

	// StartDate and EndDate must be carried forward explicitly: those fields
	// have no `omitempty`, so leaving them unset would send JSON null and
	// clear the lexicon's existing dates.
	disabled := true
	_, err := r.client.UpdateCustomLexicon(ctx, state.Id.ValueInt64(), thetalake.UpdateCustomLexiconRequest{
		Disabled:  &disabled,
		StartDate: state.StartDate.ValueStringPointer(),
		EndDate:   state.EndDate.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Custom Lexicon", fmt.Sprintf("Delete (disable) failed with error: %s", err.Error()))
		return
	}
}

// ImportState allows existing custom lexicons to be brought under Terraform
// management by specifying their ID.
func (r *customLexiconResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric custom lexicon ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
