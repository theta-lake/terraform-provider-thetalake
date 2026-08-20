package customlexicon

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

// The create endpoint can return 201 slightly before the new lexicon is
// visible to the get-by-id / list endpoints. These bound how long Create
// waits for that eventual consistency to resolve before giving up.
const (
	createConsistencyMaxAttempts = 12
	createConsistencyDelay       = 5 * time.Second
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

	// Only skip the disabled/policy_ids check when either of those attributes is
	// unknown; the rule_scope/email checks below are independent and still apply.
	if !config.Disabled.IsUnknown() && !config.PolicyIds.IsUnknown() {
		if config.Disabled.ValueBool() && !config.PolicyIds.IsNull() && len(config.PolicyIds.Elements()) > 0 {
			resp.Diagnostics.AddAttributeError(
				path.Root("policy_ids"),
				"Invalid Custom Lexicon Configuration",
				"policy_ids must be empty when disabled = true: disabling a custom lexicon removes all policies associated with it.",
			)
		}
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

	// The create response can reflect stale/default field values from before
	// the write has fully settled, so use the polled, eventually-consistent
	// lexicon (not the create response) as the source of truth going forward.
	createdLexicon := lexicon
	lexicon, err = r.waitForCustomLexiconConsistency(ctx, lexicon.Id)
	if err != nil {
		state := fromApiModel(createdLexicon)
		state.PolicyIds = reconcilePolicyIds(state.PolicyIds, plan.PolicyIds)
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		resp.Diagnostics.AddError("Failed to create Custom Lexicon", fmt.Sprintf("Create succeeded but the new lexicon did not become available for retrieval: %s", err.Error()))
		return
	}

	// The create endpoint cannot set disabled=true; if the plan requests a
	// disabled lexicon, follow up with an explicit disable call. StartDate and
	// EndDate must be carried forward explicitly.
	if plan.Disabled.ValueBool() {
		existingLexicon := lexicon
		disabled := true
		lexicon, err = r.updateCustomLexiconWithRetry(ctx, lexicon.Id, thetalake.UpdateCustomLexiconRequest{
			Disabled:  &disabled,
			StartDate: formatDatePtr(lexicon.StartDate),
			EndDate:   formatDatePtr(lexicon.EndDate),
		})
		if err != nil {
			state := fromApiModel(existingLexicon)
			state.PolicyIds = reconcilePolicyIds(state.PolicyIds, plan.PolicyIds)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			resp.Diagnostics.AddError("Failed to create Custom Lexicon", fmt.Sprintf("Create succeeded but failed to disable the new lexicon: %s", err.Error()))
			return
		}
	}

	state := fromApiModel(lexicon)
	state.PolicyIds = reconcilePolicyIds(state.PolicyIds, plan.PolicyIds)

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
	updatedState.PolicyIds = reconcilePolicyIds(updatedState.PolicyIds, state.PolicyIds)

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

	lexicon, err := r.updateCustomLexiconWithRetry(ctx, state.Id.ValueInt64(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Custom Lexicon", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(lexicon)
	updatedState.PolicyIds = reconcilePolicyIds(updatedState.PolicyIds, plan.PolicyIds)

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
	_, err := r.updateCustomLexiconWithRetry(ctx, state.Id.ValueInt64(), thetalake.UpdateCustomLexiconRequest{
		Disabled:  &disabled,
		StartDate: state.StartDate.ValueStringPointer(),
		EndDate:   state.EndDate.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Custom Lexicon", fmt.Sprintf("Delete (disable) failed with error: %s", err.Error()))
		return
	}
}

// waitForCustomLexiconConsistency polls GetCustomLexiconById until the
// newly created lexicon becomes visible, working around a brief window of
// eventual consistency between the create endpoint (which returns 201 as
// soon as the write is accepted) and the get-by-id / list endpoints.
func (r *customLexiconResource) waitForCustomLexiconConsistency(ctx context.Context, id int64) (thetalake.CustomLexicon, error) {
	var lastErr error

	for attempt := range createConsistencyMaxAttempts {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return thetalake.CustomLexicon{}, ctx.Err()
			case <-time.After(createConsistencyDelay):
			}
		}

		lexicon, err := r.client.GetCustomLexiconById(ctx, id)
		if err == nil {
			return lexicon, nil
		}
		lastErr = err
	}

	return thetalake.CustomLexicon{}, lastErr
}

// updateCustomLexiconWithRetry calls UpdateCustomLexicon, retrying if the API
// returns a thetalake.RetryableError (a 503 with a Retry-After header). This
// happens when a lexicon is updated shortly after creation, before the
// earlier write has fully settled, and uses the same bounded-attempts
// approach as waitForCustomLexiconConsistency, sleeping for the duration the
// API specifies instead of a fixed delay.
func (r *customLexiconResource) updateCustomLexiconWithRetry(ctx context.Context, id int64, request thetalake.UpdateCustomLexiconRequest) (thetalake.CustomLexicon, error) {
	var lastErr error

	for attempt := range createConsistencyMaxAttempts {
		if attempt > 0 {
			retryAfter := createConsistencyDelay
			if retryableErr, ok := errors.AsType[*thetalake.RetryableError](lastErr); ok {
				retryAfter = retryableErr.RetryAfter
			}

			select {
			case <-ctx.Done():
				return thetalake.CustomLexicon{}, ctx.Err()
			case <-time.After(retryAfter):
			}
		}

		lexicon, err := r.client.UpdateCustomLexicon(ctx, id, request)
		if err == nil {
			return lexicon, nil
		}

		if _, ok := errors.AsType[*thetalake.RetryableError](err); !ok {
			return thetalake.CustomLexicon{}, err
		}
		lastErr = err
	}

	return thetalake.CustomLexicon{}, lastErr
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
