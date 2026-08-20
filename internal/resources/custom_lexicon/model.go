package customlexicon

import (
	"regexp"
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

// dateOnlyPattern validates configured start_date/end_date values. Those
// attributes are optional-but-not-computed, so a value the API would normalise
// differently (e.g. "2024-1-1") would surface as an unhelpful "provider
// produced inconsistent result" error rather than a config error.
var dateOnlyPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

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

// datePtr resolves a planned start_date/end_date for an update: the planned
// value normally wins (null included, which clears the date), but an unresolved
// (unknown) plan value falls back to the current state so the existing date is
// carried forward instead of being cleared.
func datePtr(planned types.String, current types.String) *string {
	if planned.IsUnknown() {
		return current.ValueStringPointer()
	}
	return planned.ValueStringPointer()
}

// reconcilePolicyIds preserves the null-vs-empty distinction that the API does
// not make: the lexicon endpoints always return a policy_ids array, so a lexicon
// with no policies reads back as an empty set. policy_ids is
// optional-but-not-computed, so an unconfigured (null) value must stay null,
// otherwise Terraform reports the applied value as differing from the plan.
func reconcilePolicyIds(apiValue types.Set, configured types.Set) types.Set {
	if configured.IsNull() && len(apiValue.Elements()) == 0 {
		return types.SetNull(types.Int64Type)
	}
	return apiValue
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

// toCreateRequest maps the plan to the create request body. Required fields
// are always populated from the plan; the remaining nullable fields
// (end_date, max_participants, min_num_rules_with_hits, policy_ids,
// start_date) are only sent when explicitly configured (non-null,
// non-unknown), letting the API apply its own defaults when omitted.
func toCreateRequest(plan *customLexiconModel) thetalake.CreateCustomLexiconRequest {
	attachmentsEnabled := plan.AttachmentsEnabled.ValueBool()
	boilerplateEnabled := plan.BoilerplateEnabled.ValueBool()
	chatroomNameAnalyzed := plan.ChatroomNameAnalyzed.ValueBool()
	countProximityByCharacters := plan.CountProximityByCharacters.ValueBool()
	emailSmartBody := plan.EmailSmartBody.ValueBool()
	emailSubjectAnalyzed := plan.EmailSubjectAnalyzed.ValueBool()
	filenameAnalyzed := plan.FilenameAnalyzed.ValueBool()

	request := thetalake.CreateCustomLexiconRequest{
		AttachmentsEnabled:         &attachmentsEnabled,
		BoilerplateEnabled:         &boilerplateEnabled,
		ChatroomNameAnalyzed:       &chatroomNameAnalyzed,
		CommunicationDirection:     stringSetToSlice(plan.CommunicationDirection),
		CountProximityByCharacters: &countProximityByCharacters,
		Description:                plan.Description.ValueString(),
		EmailSmartBody:             &emailSmartBody,
		EmailSubjectAnalyzed:       &emailSubjectAnalyzed,
		FilenameAnalyzed:           &filenameAnalyzed,
		Name:                       plan.Name.ValueString(),
		RiskType:                   plan.RiskType.ValueString(),
		RuleScope:                  stringSetToSlice(plan.RuleScope),
		Rules:                      stringSetToSlice(plan.Rules),
	}

	if !plan.EndDate.IsNull() && !plan.EndDate.IsUnknown() {
		v := plan.EndDate.ValueString()
		request.EndDate = &v
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
//
// start_date, end_date and policy_ids are optional-but-not-computed, so a null
// plan value means "not configured" rather than "unchanged" and is sent as an
// explicit clear (JSON null for the dates, an empty array for policy_ids).
func toUpdateRequest(plan *customLexiconModel, state *customLexiconModel) thetalake.UpdateCustomLexiconRequest {
	name := plan.Name.ValueString()
	description := plan.Description.ValueString()

	request := thetalake.UpdateCustomLexiconRequest{
		Name:        &name,
		Description: &description,
		StartDate:   datePtr(plan.StartDate, state.StartDate),
		EndDate:     datePtr(plan.EndDate, state.EndDate),
	}

	// An unknown value can only reach here if the config references a value
	// Terraform could not resolve; leave the field unset so the API keeps the
	// existing associations rather than clearing them.
	if !plan.PolicyIds.IsUnknown() {
		policyIds := int64SetToSlice(plan.PolicyIds)
		if policyIds == nil {
			policyIds = []int64{}
		}
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
