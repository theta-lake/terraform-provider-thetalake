package swrvrule

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

var swrvRuleInputSourceAttrTypes = map[string]attr.Type{
	"id":   types.Int64Type,
	"type": types.StringType,
}

var swrvRuleInputSourceObjectType = types.ObjectType{AttrTypes: swrvRuleInputSourceAttrTypes}

type swrvRuleInputSourceModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type swrvRulePlanModel struct {
	Description        types.String `tfsdk:"description"`
	InputSources       types.List   `tfsdk:"input_sources"`
	Name               types.String `tfsdk:"name"`
	PolicyId           types.Int64  `tfsdk:"policy_id"`
	Priority           types.Int64  `tfsdk:"priority"`
	RetentionLibraryId types.Int64  `tfsdk:"retention_library_id"`
	SupervisionSpaceId types.Int64  `tfsdk:"supervision_space_id"`
	WorkflowId         types.Int64  `tfsdk:"workflow_id"`
}

type swrvRuleStateModel struct {
	DefaultRule          types.Bool   `tfsdk:"default_rule"`
	Description          types.String `tfsdk:"description"`
	Id                   types.Int64  `tfsdk:"id"`
	InputSources         types.List   `tfsdk:"input_sources"`
	IsBuiltIn            types.Bool   `tfsdk:"is_built_in"`
	Name                 types.String `tfsdk:"name"`
	PolicyId             types.Int64  `tfsdk:"policy_id"`
	PolicyName           types.String `tfsdk:"policy_name"`
	Priority             types.Int64  `tfsdk:"priority"`
	RetentionLibraryId   types.Int64  `tfsdk:"retention_library_id"`
	RetentionLibraryName types.String `tfsdk:"retention_library_name"`
	SearchId             types.Int64  `tfsdk:"search_id"`
	SearchName           types.String `tfsdk:"search_name"`
	SupervisionSpaceId   types.Int64  `tfsdk:"supervision_space_id"`
	SupervisionSpaceName types.String `tfsdk:"supervision_space_name"`
	WorkflowId           types.Int64  `tfsdk:"workflow_id"`
	WorkflowName         types.String `tfsdk:"workflow_name"`
}

func toApiModel(ctx context.Context, plan *swrvRulePlanModel) (thetalake.SwrvRule, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	rule := thetalake.SwrvRule{
		Name:               plan.Name.ValueString(),
		PolicyId:           plan.PolicyId.ValueInt64(),
		RetentionLibraryId: plan.RetentionLibraryId.ValueInt64(),
		WorkflowId:         plan.WorkflowId.ValueInt64(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		rule.Description = &description
	}

	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		priority := plan.Priority.ValueInt64()
		rule.Priority = &priority
	}

	if !plan.SupervisionSpaceId.IsNull() && !plan.SupervisionSpaceId.IsUnknown() {
		supervisionSpaceID := plan.SupervisionSpaceId.ValueInt64()
		rule.SupervisionSpaceId = &supervisionSpaceID
	}

	var inputSources []swrvRuleInputSourceModel
	diagnostics.Append(plan.InputSources.ElementsAs(ctx, &inputSources, false)...)
	if diagnostics.HasError() {
		return thetalake.SwrvRule{}, diagnostics
	}

	rule.InputSources = make([]thetalake.SwrvRuleInputSource, 0, len(inputSources))
	for _, inputSource := range inputSources {
		apiInputSource := thetalake.SwrvRuleInputSource{
			Type: inputSource.Type.ValueString(),
		}

		if !inputSource.Id.IsNull() && !inputSource.Id.IsUnknown() {
			apiInputSource.Id = inputSource.Id.ValueInt64()
		}

		rule.InputSources = append(rule.InputSources, apiInputSource)
	}

	return rule, diagnostics
}

func fromApiModel(rule thetalake.SwrvRule) swrvRuleStateModel {
	state := swrvRuleStateModel{
		DefaultRule: types.BoolValue(rule.DefaultRule),
		Id:          types.Int64Value(rule.Id),
		InputSources: inputSourcesFromAPI(
			rule.InputSource,
		),
		IsBuiltIn: types.BoolValue(rule.IsBuiltIn),
		Name:      types.StringValue(rule.Name),
	}

	if rule.Description == nil {
		state.Description = types.StringNull()
	} else {
		state.Description = types.StringValue(*rule.Description)
	}

	if rule.Policy == nil {
		state.PolicyId = types.Int64Null()
		state.PolicyName = types.StringNull()
	} else {
		state.PolicyId = types.Int64Value(rule.Policy.Id)
		state.PolicyName = types.StringValue(rule.Policy.Name)
	}

	if rule.Priority == nil {
		state.Priority = types.Int64Null()
	} else {
		state.Priority = types.Int64Value(*rule.Priority)
	}

	if rule.RetentionLibrary == nil {
		state.RetentionLibraryId = types.Int64Null()
		state.RetentionLibraryName = types.StringNull()
	} else {
		state.RetentionLibraryId = types.Int64Value(rule.RetentionLibrary.Id)
		state.RetentionLibraryName = types.StringValue(rule.RetentionLibrary.Name)
	}

	if rule.Search == nil {
		state.SearchId = types.Int64Null()
		state.SearchName = types.StringNull()
	} else {
		state.SearchId = types.Int64Value(rule.Search.Id)
		state.SearchName = types.StringValue(rule.Search.Name)
	}

	if rule.SupervisionSpace == nil {
		state.SupervisionSpaceId = types.Int64Null()
		state.SupervisionSpaceName = types.StringNull()
	} else {
		state.SupervisionSpaceId = types.Int64Value(rule.SupervisionSpace.Id)
		state.SupervisionSpaceName = types.StringValue(rule.SupervisionSpace.Name)
	}

	if rule.Workflow == nil {
		state.WorkflowId = types.Int64Null()
		state.WorkflowName = types.StringNull()
	} else {
		state.WorkflowId = types.Int64Value(rule.Workflow.Id)
		state.WorkflowName = types.StringValue(rule.Workflow.Name)
	}

	return state
}

func inputSourcesFromAPI(inputSources []thetalake.SwrvInputSource) types.List {
	values := make([]attr.Value, 0, len(inputSources))
	for _, inputSource := range inputSources {
		attributes := map[string]attr.Value{
			"id":   types.Int64Null(),
			"type": types.StringNull(),
		}

		switch {
		case inputSource.Integration != nil:
			attributes["id"] = types.Int64Value(inputSource.Integration.Id)
			attributes["type"] = types.StringValue(inputSource.Integration.Type)
		case inputSource.Name != "":
			mappedType, ok := swrvInputSourceTypeFromName(inputSource.Name)
			if !ok {
				continue
			}
			attributes["type"] = types.StringValue(mappedType)
		default:
			continue
		}

		values = append(values, types.ObjectValueMust(swrvRuleInputSourceAttrTypes, attributes))
	}

	if len(values) == 0 {
		return types.ListValueMust(swrvRuleInputSourceObjectType, []attr.Value{})
	}

	return types.ListValueMust(swrvRuleInputSourceObjectType, values)
}

func swrvInputSourceTypeFromName(name string) (string, bool) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	switch normalized {
	case "all integration uploads":
		return "all_integration_uploads", true
	case "all submission portal uploads":
		return "all_submission_portal_uploads", true
	case "all uploads":
		return "all_uploads", true
	case "all user uploads":
		return "all_user_uploads", true
	default:
		return "", false
	}
}

func inputSourcesHaveValues(inputSources types.List) bool {
	if inputSources.IsNull() || inputSources.IsUnknown() {
		return false
	}

	return len(inputSources.Elements()) > 0
}
