package customlexicon

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

// dateOnlyFormat is the format the Theta Lake API expects/returns for
// start_date/end_date on write (YYYY-MM-DD), and the format we truncate the
// RFC3339 read-back values to, to avoid perpetual diffs.
const dateOnlyFormat = "2006-01-02"

// formatDatePtr renders an API date pointer in dateOnlyFormat, preserving
// nil. UpdateCustomLexiconRequest.StartDate/EndDate have no `omitempty`, so
// every update call must carry an explicit value (or explicit nil) for them —
// leaving them unset would send JSON null and clear the lexicon's existing
// dates.
func formatDatePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format(dateOnlyFormat)
	return &formatted
}

// customLexiconModel covers every attribute in the schema, so it can be used
// as the target for both `Plan.Get` and `State.Get` — those calls require the
// struct to define a field for every schema attribute, including the
// computed-only ones.
type customLexiconModel struct {
	AcceptsInput               types.Bool        `tfsdk:"accepts_input"`
	AttachmentsEnabled         types.Bool        `tfsdk:"attachments_enabled"`
	BoilerplateEnabled         types.Bool        `tfsdk:"boilerplate_enabled"`
	ChatroomNameAnalyzed       types.Bool        `tfsdk:"chatroom_name_analyzed"`
	CommunicationDirection     types.Set         `tfsdk:"communication_direction"`
	CountProximityByCharacters types.Bool        `tfsdk:"count_proximity_by_characters"`
	CreatedAt                  timetypes.RFC3339 `tfsdk:"created_at"`
	Description                types.String      `tfsdk:"description"`
	Disabled                   types.Bool        `tfsdk:"disabled"`
	DisabledAt                 timetypes.RFC3339 `tfsdk:"disabled_at"`
	EmailSmartBody             types.Bool        `tfsdk:"email_smart_body"`
	EmailSubjectAnalyzed       types.Bool        `tfsdk:"email_subject_analyzed"`
	EndDate                    types.String      `tfsdk:"end_date"`
	FilenameAnalyzed           types.Bool        `tfsdk:"filename_analyzed"`
	Id                         types.Int64       `tfsdk:"id"`
	MaxParticipants            types.Int64       `tfsdk:"max_participants"`
	MinNumRulesWithHits        types.Int64       `tfsdk:"min_num_rules_with_hits"`
	Name                       types.String      `tfsdk:"name"`
	PolicyIds                  types.Set         `tfsdk:"policy_ids"`
	RiskType                   types.String      `tfsdk:"risk_type"`
	RuleScope                  types.Set         `tfsdk:"rule_scope"`
	Rules                      types.Set         `tfsdk:"rules"`
	StartDate                  types.String      `tfsdk:"start_date"`
	UpdatedAt                  timetypes.RFC3339 `tfsdk:"updated_at"`
}

func stringSetToSlice(set types.Set) []string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	values := make([]string, 0, len(set.Elements()))
	for _, v := range set.Elements() {
		if s, ok := v.(types.String); ok {
			values = append(values, s.ValueString())
		}
	}
	return values
}

func int64SetToSlice(set types.Set) []int64 {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	values := make([]int64, 0, len(set.Elements()))
	for _, v := range set.Elements() {
		if i, ok := v.(types.Int64); ok {
			values = append(values, i.ValueInt64())
		}
	}
	return values
}

func stringSliceToSet(values []string) types.Set {
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.StringValue(v))
	}
	return types.SetValueMust(types.StringType, elements)
}

func int64SliceToSet(values []int64) types.Set {
	elements := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elements = append(elements, types.Int64Value(v))
	}
	return types.SetValueMust(types.Int64Type, elements)
}

// toCreateRequest maps the plan to the create request body. Create-only
// fields are only sent when explicitly configured (non-null, non-unknown);
// omitted optional fields let the API apply its own defaults.
func toCreateRequest(plan *customLexiconModel) thetalake.CreateCustomLexiconRequest {
	request := thetalake.CreateCustomLexiconRequest{
		Description: plan.Description.ValueString(),
		Name:        plan.Name.ValueString(),
		RiskType:    plan.RiskType.ValueString(),
		Rules:       stringSetToSlice(plan.Rules),
	}

	if !plan.AttachmentsEnabled.IsNull() && !plan.AttachmentsEnabled.IsUnknown() {
		v := plan.AttachmentsEnabled.ValueBool()
		request.AttachmentsEnabled = &v
	}
	if !plan.BoilerplateEnabled.IsNull() && !plan.BoilerplateEnabled.IsUnknown() {
		v := plan.BoilerplateEnabled.ValueBool()
		request.BoilerplateEnabled = &v
	}
	if !plan.ChatroomNameAnalyzed.IsNull() && !plan.ChatroomNameAnalyzed.IsUnknown() {
		v := plan.ChatroomNameAnalyzed.ValueBool()
		request.ChatroomNameAnalyzed = &v
	}
	if !plan.CommunicationDirection.IsNull() && !plan.CommunicationDirection.IsUnknown() {
		request.CommunicationDirection = stringSetToSlice(plan.CommunicationDirection)
	}
	if !plan.CountProximityByCharacters.IsNull() && !plan.CountProximityByCharacters.IsUnknown() {
		v := plan.CountProximityByCharacters.ValueBool()
		request.CountProximityByCharacters = &v
	}
	if !plan.EmailSmartBody.IsNull() && !plan.EmailSmartBody.IsUnknown() {
		v := plan.EmailSmartBody.ValueBool()
		request.EmailSmartBody = &v
	}
	if !plan.EmailSubjectAnalyzed.IsNull() && !plan.EmailSubjectAnalyzed.IsUnknown() {
		v := plan.EmailSubjectAnalyzed.ValueBool()
		request.EmailSubjectAnalyzed = &v
	}
	if !plan.EndDate.IsNull() && !plan.EndDate.IsUnknown() {
		v := plan.EndDate.ValueString()
		request.EndDate = &v
	}
	if !plan.FilenameAnalyzed.IsNull() && !plan.FilenameAnalyzed.IsUnknown() {
		v := plan.FilenameAnalyzed.ValueBool()
		request.FilenameAnalyzed = &v
	}
	if !plan.MaxParticipants.IsNull() && !plan.MaxParticipants.IsUnknown() {
		v := plan.MaxParticipants.ValueInt64()
		request.MaxParticipants = &v
	}
	if !plan.MinNumRulesWithHits.IsNull() && !plan.MinNumRulesWithHits.IsUnknown() {
		v := plan.MinNumRulesWithHits.ValueInt64()
		request.MinNumRulesWithHits = &v
	}
	if !plan.PolicyIds.IsNull() && !plan.PolicyIds.IsUnknown() {
		request.Policies = int64SetToSlice(plan.PolicyIds)
	}
	if !plan.RuleScope.IsNull() && !plan.RuleScope.IsUnknown() {
		request.RuleScope = stringSetToSlice(plan.RuleScope)
	}
	if !plan.StartDate.IsNull() && !plan.StartDate.IsUnknown() {
		v := plan.StartDate.ValueString()
		request.StartDate = &v
	}

	return request
}

// toUpdateRequest maps the plan (and current state, to detect a disabled
// transition) to the update request body. Disabled is only sent when it
// changes vs. state, since the API rejects disabling an already-disabled
// lexicon.
func toUpdateRequest(plan *customLexiconModel, state *customLexiconModel) thetalake.UpdateCustomLexiconRequest {
	name := plan.Name.ValueString()
	description := plan.Description.ValueString()

	request := thetalake.UpdateCustomLexiconRequest{
		Name:        &name,
		Description: &description,
	}

	if plan.StartDate.IsUnknown() {
		request.StartDate = state.StartDate.ValueStringPointer()
	} else {
		request.StartDate = plan.StartDate.ValueStringPointer()
	}

	if plan.EndDate.IsUnknown() {
		request.EndDate = state.EndDate.ValueStringPointer()
	} else {
		request.EndDate = plan.EndDate.ValueStringPointer()
	}

	if !plan.PolicyIds.IsNull() && !plan.PolicyIds.IsUnknown() {
		policyIds := int64SetToSlice(plan.PolicyIds)
		request.PolicyIds = &policyIds
	}

	if plan.Disabled.ValueBool() != state.Disabled.ValueBool() {
		disabled := plan.Disabled.ValueBool()
		request.Disabled = &disabled
	}

	return request
}

func fromApiModel(lexicon thetalake.CustomLexicon) customLexiconModel {
	state := customLexiconModel{
		AcceptsInput:               types.BoolValue(lexicon.AcceptsInput),
		AttachmentsEnabled:         types.BoolValue(lexicon.AttachmentsEnabled),
		BoilerplateEnabled:         types.BoolValue(lexicon.BoilerplateEnabled),
		ChatroomNameAnalyzed:       types.BoolValue(lexicon.ChatroomNameAnalyzed),
		CommunicationDirection:     stringSliceToSet(lexicon.CommunicationDirection),
		CountProximityByCharacters: types.BoolValue(lexicon.CountProximityByCharacters),
		CreatedAt:                  timetypes.NewRFC3339TimeValue(lexicon.CreatedAt),
		Description:                types.StringValue(lexicon.Description),
		Disabled:                   types.BoolValue(lexicon.DisabledAt != nil),
		EmailSmartBody:             types.BoolValue(lexicon.EmailSmartBody),
		EmailSubjectAnalyzed:       types.BoolValue(lexicon.EmailSubjectAnalyzed),
		FilenameAnalyzed:           types.BoolValue(lexicon.FilenameAnalyzed),
		Id:                         types.Int64Value(lexicon.Id),
		Name:                       types.StringValue(lexicon.Name),
		PolicyIds:                  int64SliceToSet(lexicon.PolicyIds),
		RiskType:                   types.StringValue(lexicon.RiskType),
		RuleScope:                  stringSliceToSet(lexicon.RuleScope),
		UpdatedAt:                  timetypes.NewRFC3339TimeValue(lexicon.UpdatedAt),
	}

	rules := make([]string, 0, len(lexicon.Rules))
	for _, rule := range lexicon.Rules {
		rules = append(rules, rule)
	}
	state.Rules = stringSliceToSet(rules)

	if lexicon.DisabledAt == nil {
		state.DisabledAt = timetypes.NewRFC3339Null()
	} else {
		state.DisabledAt = timetypes.NewRFC3339TimeValue(*lexicon.DisabledAt)
	}

	if lexicon.MaxParticipants == nil {
		state.MaxParticipants = types.Int64Null()
	} else {
		state.MaxParticipants = types.Int64Value(*lexicon.MaxParticipants)
	}

	if lexicon.MinNumRulesWithHits == nil {
		state.MinNumRulesWithHits = types.Int64Null()
	} else {
		state.MinNumRulesWithHits = types.Int64Value(*lexicon.MinNumRulesWithHits)
	}

	state.StartDate = types.StringPointerValue(formatDatePtr(lexicon.StartDate))
	state.EndDate = types.StringPointerValue(formatDatePtr(lexicon.EndDate))

	return state
}
