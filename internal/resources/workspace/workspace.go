package workspace

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type workspaceResource struct {
	client *thetalake.Client
}

func NewWorkspaceResource() resource.Resource {
	return &workspaceResource{}
}

func (r *workspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_workspace", req.ProviderTypeName)
}

func (r *workspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*thetalake.Client)
}

func (r *workspaceResource) Create(ctx context.Context, _ resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Cannot create workspace",
		"Workspaces cannot be created via the Theta Lake API. Use `terraform import thetalake_workspace.<name> <id>` to manage an existing workspace.",
	)
}

func (r *workspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := workspaceStateModel{}

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Internal Error", "Failed to read state data")
		return
	}

	workspace, err := r.client.GetWorkspaceById(ctx, state.Id.ValueInt64())
	if err != nil {
		if errors.Is(err, thetalake.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read Workspace", fmt.Sprintf("Read failed with error: %s", err.Error()))
		return
	}

	updatedState := fromApiModel(workspace)

	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
}

func (r *workspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan  workspacePlanModel
		state workspaceStateModel
	)

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("allow_anonymous_via_shared_links"), &plan.AllowAnonymousViaSharedLinks)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("analysis_supervision_space_ids"), &plan.AnalysisSupervisionSpaceIds)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("audit_log_retention_period"), &plan.AuditLogRetentionPeriod)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("case_management_manager_assignment"), &plan.CaseManagementManagerAssignment)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("default_transcription_language"), &plan.DefaultTranscriptionLanguage)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("default_workspace_timezone"), &plan.DefaultWorkspaceTimezone)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("delete_on_expiration"), &plan.DeleteOnExpiration)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("hide_attachments_from_search"), &plan.HideAttachmentsFromSearch)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("preferred_languages"), &plan.PreferredLanguages)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("reauthenticate_on_network_change"), &plan.ReauthenticateOnNetworkChange)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("shared_links_expiration_period"), &plan.SharedLinksExpirationPeriod)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("show_system_messages_in_chat"), &plan.ShowSystemMessagesInChat)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("use_name_matcher"), &plan.UseNameMatcher)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("use_owner_only_space_matcher"), &plan.UseOwnerOnlySpaceMatcher)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("user_ids"), &plan.UserIds)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Failed to update Workspace", "Update failed to read Workspace resource plan data")
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
		resp.Diagnostics.AddError("Failed to update Workspace", "Update failed to map Workspace plan data")
		return
	}

	updatedWorkspace, err := r.client.UpdateWorkspace(ctx, apiModel)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update Workspace", fmt.Sprintf("Update failed with error: %s", err.Error()))
		return
	}

	// Persist the settings update now so it isn't lost if a later membership
	// change fails partway through.
	updatedState := fromApiModel(updatedWorkspace)
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.UserIds.IsNull() && !plan.UserIds.IsUnknown() {
		var planUserIds []int64
		resp.Diagnostics.Append(plan.UserIds.ElementsAs(ctx, &planUserIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		var stateUserIds []int64
		resp.Diagnostics.Append(state.UserIds.ElementsAs(ctx, &stateUserIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		stateSet := make(map[int64]struct{}, len(stateUserIds))
		for _, id := range stateUserIds {
			stateSet[id] = struct{}{}
		}

		planSet := make(map[int64]struct{}, len(planUserIds))
		for _, id := range planUserIds {
			planSet[id] = struct{}{}
		}

		var (
			membershipErrDetail  string
			membershipErrSummary string
			usersChanged         bool
		)

		for _, id := range planUserIds {
			if _, exists := stateSet[id]; !exists {
				if err := r.client.AddUserToWorkspace(ctx, id); err != nil {
					membershipErrSummary = "Failed to add user to workspace"
					membershipErrDetail = fmt.Sprintf("Failed to add user %d: %s", id, err.Error())
					break
				}
				usersChanged = true
			}
		}

		if membershipErrSummary == "" {
			for _, id := range stateUserIds {
				if _, exists := planSet[id]; !exists {
					if err := r.client.RemoveUserFromWorkspace(ctx, id); err != nil {
						membershipErrSummary = "Failed to remove user from workspace"
						membershipErrDetail = fmt.Sprintf("Failed to remove user %d: %s", id, err.Error())
						break
					}
					usersChanged = true
				}
			}
		}

		// Re-sync state with whatever membership changes actually landed on
		// the server, even if the loop above stopped early on an error.
		if usersChanged {
			refreshedWorkspace, refreshErr := r.client.GetWorkspaceById(ctx, state.Id.ValueInt64())
			if refreshErr != nil {
				resp.Diagnostics.AddError("Failed to read Workspace after user changes", fmt.Sprintf("Read failed with error: %s", refreshErr.Error()))
				return
			}

			updatedState = fromApiModel(refreshedWorkspace)
			resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if membershipErrSummary != "" {
			resp.Diagnostics.AddError(membershipErrSummary, membershipErrDetail)
			return
		}
	}
}

func (r *workspaceResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Workspaces cannot be deleted via the Theta Lake API.
	// Removing the resource from Terraform state without making an API call.
}

// ImportState allows existing workspaces to be brought under Terraform
// management by specifying their ID. The ID from the import command is
// mapped directly to the "id" attribute, after which Read will
// populate the rest of the state.
func (r *workspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected numeric workspace ID, got %q: %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
