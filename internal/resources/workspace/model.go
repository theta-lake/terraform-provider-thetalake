package workspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

var supervisionSpaceAttrTypes = map[string]attr.Type{
	"id":   types.Int64Type,
	"name": types.StringType,
}

var supervisionSpaceObjectType = types.ObjectType{AttrTypes: supervisionSpaceAttrTypes}

var userAttrTypes = map[string]attr.Type{
	"email": types.StringType,
	"id":    types.Int64Type,
	"name":  types.StringType,
}

var userObjectType = types.ObjectType{AttrTypes: userAttrTypes}

type workspacePlanModel struct {
	AllowAnonymousViaSharedLinks    types.Bool   `tfsdk:"allow_anonymous_via_shared_links"`
	AnalysisSupervisionSpaceIds     types.Set    `tfsdk:"analysis_supervision_space_ids"`
	AuditLogRetentionPeriod         types.Int64  `tfsdk:"audit_log_retention_period"`
	CaseManagementManagerAssignment types.Bool   `tfsdk:"case_management_manager_assignment"`
	DefaultTranscriptionLanguage    types.String `tfsdk:"default_transcription_language"`
	DefaultWorkspaceTimezone        types.String `tfsdk:"default_workspace_timezone"`
	DeleteOnExpiration              types.Bool   `tfsdk:"delete_on_expiration"`
	HideAttachmentsFromSearch       types.Bool   `tfsdk:"hide_attachments_from_search"`
	PreferredLanguages              types.List   `tfsdk:"preferred_languages"`
	ReauthenticateOnNetworkChange   types.Bool   `tfsdk:"reauthenticate_on_network_change"`
	SharedLinksExpirationPeriod     types.Int64  `tfsdk:"shared_links_expiration_period"`
	ShowSystemMessagesInChat        types.Bool   `tfsdk:"show_system_messages_in_chat"`
	UseNameMatcher                  types.Bool   `tfsdk:"use_name_matcher"`
	UseOwnerOnlySpaceMatcher        types.Bool   `tfsdk:"use_owner_only_space_matcher"`
	UserIds                         types.Set    `tfsdk:"user_ids"`
}

type workspaceStateModel struct {
	AllowAnonymousViaSharedLinks    types.Bool        `tfsdk:"allow_anonymous_via_shared_links"`
	AnalysisSupervisionSpaceIds     types.Set         `tfsdk:"analysis_supervision_space_ids"`
	AnalysisSupervisionSpaces       types.List        `tfsdk:"analysis_supervision_spaces"`
	AuditLogRetentionPeriod         types.Int64       `tfsdk:"audit_log_retention_period"`
	CaseManagementManagerAssignment types.Bool        `tfsdk:"case_management_manager_assignment"`
	CreatedAt                       timetypes.RFC3339 `tfsdk:"created_at"`
	DefaultTranscriptionLanguage    types.String      `tfsdk:"default_transcription_language"`
	DefaultWorkspaceTimezone        types.String      `tfsdk:"default_workspace_timezone"`
	DeleteOnExpiration              types.Bool        `tfsdk:"delete_on_expiration"`
	Description                     types.String      `tfsdk:"description"`
	Disabled                        types.Bool        `tfsdk:"disabled"`
	DisabledAt                      timetypes.RFC3339 `tfsdk:"disabled_at"`
	HideAttachmentsFromSearch       types.Bool        `tfsdk:"hide_attachments_from_search"`
	Id                              types.Int64       `tfsdk:"id"`
	Name                            types.String      `tfsdk:"name"`
	PreferredLanguages              types.List        `tfsdk:"preferred_languages"`
	ReauthenticateOnNetworkChange   types.Bool        `tfsdk:"reauthenticate_on_network_change"`
	SharedLinksExpirationPeriod     types.Int64       `tfsdk:"shared_links_expiration_period"`
	ShowSystemMessagesInChat        types.Bool        `tfsdk:"show_system_messages_in_chat"`
	UpdatedAt                       timetypes.RFC3339 `tfsdk:"updated_at"`
	UseNameMatcher                  types.Bool        `tfsdk:"use_name_matcher"`
	UseOwnerOnlySpaceMatcher        types.Bool        `tfsdk:"use_owner_only_space_matcher"`
	UserIds                         types.Set         `tfsdk:"user_ids"`
	Users                           types.List        `tfsdk:"users"`
}

func toApiModel(ctx context.Context, plan *workspacePlanModel) (thetalake.Workspace, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	workspace := thetalake.Workspace{
		AllowAnonymousViaSharedLinks:    plan.AllowAnonymousViaSharedLinks.ValueBool(),
		CaseManagementManagerAssignment: plan.CaseManagementManagerAssignment.ValueBool(),
		DefaultTranscriptionLanguage:    plan.DefaultTranscriptionLanguage.ValueString(),
		DefaultWorkspaceTimezone:        plan.DefaultWorkspaceTimezone.ValueString(),
		DeleteOnExpiration:              plan.DeleteOnExpiration.ValueBool(),
		HideAttachmentsFromSearch:       plan.HideAttachmentsFromSearch.ValueBool(),
		ReauthenticateOnNetworkChange:   plan.ReauthenticateOnNetworkChange.ValueBool(),
		ShowSystemMessagesInChat:        plan.ShowSystemMessagesInChat.ValueBool(),
		UseNameMatcher:                  plan.UseNameMatcher.ValueBool(),
		UseOwnerOnlySpaceMatcher:        plan.UseOwnerOnlySpaceMatcher.ValueBool(),
	}

	if !plan.AuditLogRetentionPeriod.IsNull() && !plan.AuditLogRetentionPeriod.IsUnknown() {
		v := plan.AuditLogRetentionPeriod.ValueInt64()
		workspace.AuditLogRetentionPeriod = &v
	}

	var ids []int64
	diagnostics.Append(plan.AnalysisSupervisionSpaceIds.ElementsAs(ctx, &ids, false)...)
	if diagnostics.HasError() {
		return thetalake.Workspace{}, diagnostics
	}
	if ids == nil {
		ids = []int64{}
	}
	workspace.AnalysisSupervisionSpaceIds = ids

	var langs []string
	diagnostics.Append(plan.PreferredLanguages.ElementsAs(ctx, &langs, false)...)
	if diagnostics.HasError() {
		return thetalake.Workspace{}, diagnostics
	}
	if langs == nil {
		langs = []string{}
	}
	workspace.PreferredLanguages = langs

	if !plan.SharedLinksExpirationPeriod.IsNull() && !plan.SharedLinksExpirationPeriod.IsUnknown() {
		v := plan.SharedLinksExpirationPeriod.ValueInt64()
		workspace.SharedLinksExpirationPeriod = &v
	}

	return workspace, diagnostics
}

func fromApiModel(workspace thetalake.Workspace) workspaceStateModel {
	state := workspaceStateModel{
		AllowAnonymousViaSharedLinks:    types.BoolValue(workspace.AllowAnonymousViaSharedLinks),
		CaseManagementManagerAssignment: types.BoolValue(workspace.CaseManagementManagerAssignment),
		CreatedAt:                       timetypes.NewRFC3339TimeValue(workspace.CreatedAt),
		DefaultTranscriptionLanguage:    types.StringValue(workspace.DefaultTranscriptionLanguage),
		DefaultWorkspaceTimezone:        types.StringValue(workspace.DefaultWorkspaceTimezone),
		DeleteOnExpiration:              types.BoolValue(workspace.DeleteOnExpiration),
		Description:                     types.StringValue(workspace.Description),
		Disabled:                        types.BoolValue(workspace.Disabled),
		HideAttachmentsFromSearch:       types.BoolValue(workspace.HideAttachmentsFromSearch),
		Id:                              types.Int64Value(workspace.Id),
		Name:                            types.StringValue(workspace.Name),
		ReauthenticateOnNetworkChange:   types.BoolValue(workspace.ReauthenticateOnNetworkChange),
		ShowSystemMessagesInChat:        types.BoolValue(workspace.ShowSystemMessagesInChat),
		UpdatedAt:                       timetypes.NewRFC3339TimeValue(workspace.UpdatedAt),
		UseNameMatcher:                  types.BoolValue(workspace.UseNameMatcher),
		UseOwnerOnlySpaceMatcher:        types.BoolValue(workspace.UseOwnerOnlySpaceMatcher),
	}

	if workspace.AuditLogRetentionPeriod == nil {
		state.AuditLogRetentionPeriod = types.Int64Null()
	} else {
		state.AuditLogRetentionPeriod = types.Int64Value(*workspace.AuditLogRetentionPeriod)
	}

	if workspace.DisabledAt == nil {
		state.DisabledAt = timetypes.NewRFC3339Null()
	} else {
		state.DisabledAt = timetypes.NewRFC3339TimeValue(*workspace.DisabledAt)
	}

	ssIds := make([]attr.Value, 0, len(workspace.AnalysisSupervisionSpaceIds))
	for _, id := range workspace.AnalysisSupervisionSpaceIds {
		ssIds = append(ssIds, types.Int64Value(id))
	}
	state.AnalysisSupervisionSpaceIds = types.SetValueMust(types.Int64Type, ssIds)

	ssObjects := make([]attr.Value, 0, len(workspace.AnalysisSupervisionSpaces))
	for _, ss := range workspace.AnalysisSupervisionSpaces {
		ssObjects = append(ssObjects, types.ObjectValueMust(supervisionSpaceAttrTypes, map[string]attr.Value{
			"id":   types.Int64Value(ss.Id),
			"name": types.StringValue(ss.Name),
		}))
	}
	state.AnalysisSupervisionSpaces = types.ListValueMust(supervisionSpaceObjectType, ssObjects)

	langValues := make([]attr.Value, 0, len(workspace.PreferredLanguages))
	for _, lang := range workspace.PreferredLanguages {
		langValues = append(langValues, types.StringValue(lang))
	}
	state.PreferredLanguages = types.ListValueMust(types.StringType, langValues)

	if workspace.SharedLinksExpirationPeriod == nil {
		state.SharedLinksExpirationPeriod = types.Int64Null()
	} else {
		state.SharedLinksExpirationPeriod = types.Int64Value(*workspace.SharedLinksExpirationPeriod)
	}

	userIdValues := make([]attr.Value, 0, len(workspace.Users))
	userObjects := make([]attr.Value, 0, len(workspace.Users))
	for _, u := range workspace.Users {
		userIdValues = append(userIdValues, types.Int64Value(u.Id))
		userObjects = append(userObjects, types.ObjectValueMust(userAttrTypes, map[string]attr.Value{
			"email": types.StringValue(u.Email),
			"id":    types.Int64Value(u.Id),
			"name":  types.StringValue(u.Name),
		}))
	}
	state.UserIds = types.SetValueMust(types.Int64Type, userIdValues)
	state.Users = types.ListValueMust(userObjectType, userObjects)

	return state
}
