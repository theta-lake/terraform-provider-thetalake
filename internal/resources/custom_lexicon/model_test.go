package customlexicon

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestToCreateRequest(t *testing.T) {
	plan := &customLexiconModel{
		AttachmentsEnabled: types.BoolValue(true),
		Description:        types.StringValue("My Lexicon description"),
		Disabled:           types.BoolValue(false),
		EndDate:            types.StringValue("2025-06-17"),
		MaxParticipants:    types.Int64Value(10),
		Name:               types.StringValue("My Lexicon"),
		PolicyIds: types.SetValueMust(types.Int64Type, []attr.Value{
			types.Int64Value(1),
			types.Int64Value(2),
		}),
		RiskType: types.StringValue("risk"),
		RuleScope: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("chat"),
			types.StringValue("email"),
		}),
		Rules: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("word1"),
			types.StringValue("word2"),
		}),
		StartDate: types.StringValue("2021-06-16"),
		// Unset create-only fields left as zero-value (null).
		BoilerplateEnabled:         types.BoolNull(),
		ChatroomNameAnalyzed:       types.BoolNull(),
		CommunicationDirection:     types.SetNull(types.StringType),
		CountProximityByCharacters: types.BoolNull(),
		EmailSmartBody:             types.BoolNull(),
		EmailSubjectAnalyzed:       types.BoolNull(),
		FilenameAnalyzed:           types.BoolNull(),
		MinNumRulesWithHits:        types.Int64Null(),
	}

	request := toCreateRequest(plan)

	assert.Equal(t, "My Lexicon description", request.Description)
	assert.Equal(t, "My Lexicon", request.Name)
	assert.Equal(t, "risk", request.RiskType)
	assert.ElementsMatch(t, []string{"word1", "word2"}, request.Rules)
	assert.ElementsMatch(t, []int64{1, 2}, request.Policies)
	assert.ElementsMatch(t, []string{"chat", "email"}, request.RuleScope)
	if assert.NotNil(t, request.AttachmentsEnabled) {
		assert.True(t, *request.AttachmentsEnabled)
	}
	assert.Nil(t, request.BoilerplateEnabled)
	if assert.NotNil(t, request.MaxParticipants) {
		assert.Equal(t, int64(10), *request.MaxParticipants)
	}
	assert.Nil(t, request.MinNumRulesWithHits)
	if assert.NotNil(t, request.StartDate) {
		assert.Equal(t, "2021-06-16", *request.StartDate)
	}
	if assert.NotNil(t, request.EndDate) {
		assert.Equal(t, "2025-06-17", *request.EndDate)
	}
}

func TestToUpdateRequest_DisabledUnchanged(t *testing.T) {
	plan := &customLexiconModel{
		Description: types.StringValue("Updated description"),
		Disabled:    types.BoolValue(false),
		Name:        types.StringValue("Updated Name"),
		PolicyIds: types.SetValueMust(types.Int64Type, []attr.Value{
			types.Int64Value(5),
		}),
	}
	state := &customLexiconModel{
		Disabled: types.BoolValue(false),
	}

	request := toUpdateRequest(plan, state)

	assert.Equal(t, "Updated Name", *request.Name)
	assert.Equal(t, "Updated description", *request.Description)
	assert.Nil(t, request.Disabled)
	if assert.NotNil(t, request.PolicyIds) {
		assert.Equal(t, []int64{5}, *request.PolicyIds)
	}
}

func TestToUpdateRequest_DisabledChanged(t *testing.T) {
	plan := &customLexiconModel{
		Description: types.StringValue("Updated description"),
		Disabled:    types.BoolValue(true),
		Name:        types.StringValue("Updated Name"),
		PolicyIds:   types.SetNull(types.Int64Type),
	}
	state := &customLexiconModel{
		Disabled: types.BoolValue(false),
	}

	request := toUpdateRequest(plan, state)

	if assert.NotNil(t, request.Disabled) {
		assert.True(t, *request.Disabled)
	}
	assert.Nil(t, request.PolicyIds)
}

func TestToUpdateRequest_EmptyPolicyIdsClears(t *testing.T) {
	plan := &customLexiconModel{
		Description: types.StringValue("d"),
		Disabled:    types.BoolValue(false),
		Name:        types.StringValue("n"),
		PolicyIds:   types.SetValueMust(types.Int64Type, []attr.Value{}),
	}
	state := &customLexiconModel{
		Disabled: types.BoolValue(false),
	}

	request := toUpdateRequest(plan, state)

	if assert.NotNil(t, request.PolicyIds) {
		assert.Equal(t, []int64{}, *request.PolicyIds)
	}
}

func TestFromApiModel(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2021-06-16T01:37:04.262Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2022-10-12T02:29:49.146Z")
	startDate, _ := time.Parse(time.RFC3339, "2021-06-16T00:00:00.000Z")
	endDate, _ := time.Parse(time.RFC3339, "2025-06-17T00:00:00.000Z")
	maxParticipants := int64(10)
	minNumRulesWithHits := int64(5)

	lexicon := thetalake.CustomLexicon{
		AcceptsInput:           true,
		AttachmentsEnabled:     true,
		CommunicationDirection: []string{"inbound", "outbound"},
		CreatedAt:              createdAt,
		Description:            "My Lexicon description",
		DisabledAt:             nil,
		EndDate:                &endDate,
		Id:                     1235,
		MaxParticipants:        &maxParticipants,
		MinNumRulesWithHits:    &minNumRulesWithHits,
		Name:                   "My Lexicon",
		PolicyIds:              []int64{1, 2, 3},
		RiskType:               "risk",
		Rules: map[string]string{
			"a": "word1",
			"b": "word2",
			"c": "word3",
		},
		RuleScope: []string{"chat", "email", "doc"},
		StartDate: &startDate,
		UpdatedAt: updatedAt,
	}

	state := fromApiModel(lexicon)

	assert.Equal(t, int64(1235), state.Id.ValueInt64())
	assert.Equal(t, "My Lexicon", state.Name.ValueString())
	assert.Equal(t, "risk", state.RiskType.ValueString())
	assert.False(t, state.Disabled.ValueBool())
	assert.True(t, state.DisabledAt.IsNull())
	assert.Equal(t, "2021-06-16", state.StartDate.ValueString())
	assert.Equal(t, "2025-06-17", state.EndDate.ValueString())
	assert.Equal(t, int64(10), state.MaxParticipants.ValueInt64())
	assert.Equal(t, int64(5), state.MinNumRulesWithHits.ValueInt64())
	assert.Equal(t, 3, len(state.Rules.Elements()))
	assert.Equal(t, 3, len(state.PolicyIds.Elements()))
	assert.Equal(t, 2, len(state.CommunicationDirection.Elements()))
}

func TestFromApiModel_Disabled(t *testing.T) {
	disabledAt := time.Now().UTC()

	lexicon := thetalake.CustomLexicon{
		Id:         1235,
		Name:       "My Lexicon",
		DisabledAt: &disabledAt,
		Rules:      map[string]string{},
		PolicyIds:  nil,
	}

	state := fromApiModel(lexicon)

	assert.True(t, state.Disabled.ValueBool())
	assert.False(t, state.DisabledAt.IsNull())
	assert.True(t, state.MaxParticipants.IsNull())
	assert.True(t, state.MinNumRulesWithHits.IsNull())
	assert.True(t, state.StartDate.IsNull())
	assert.True(t, state.EndDate.IsNull())
}
