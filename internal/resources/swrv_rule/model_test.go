package swrvrule

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestSwrvRuleToApiModel(t *testing.T) {
	plan := &swrvRulePlanModel{
		Description: types.StringValue("My SWRV rule"),
		InputSources: types.ListValueMust(swrvRuleInputSourceObjectType, []attr.Value{
			types.ObjectValueMust(swrvRuleInputSourceAttrTypes, map[string]attr.Value{
				"id":   types.Int64Null(),
				"type": types.StringValue("all_uploads"),
			}),
			types.ObjectValueMust(swrvRuleInputSourceAttrTypes, map[string]attr.Value{
				"id":   types.Int64Value(2345),
				"type": types.StringValue("integration"),
			}),
		}),
		Name:               types.StringValue("swrv-example"),
		PolicyId:           types.Int64Value(147),
		Priority:           types.Int64Value(4),
		RetentionLibraryId: types.Int64Value(1),
		SupervisionSpaceId: types.Int64Value(10420),
		WorkflowId:         types.Int64Value(14536),
	}

	apiModel, diags := toApiModel(context.Background(), plan)
	if diags.HasError() {
		t.Fatalf("expected toApiModel to succeed, got diagnostics: %v", diags)
	}
	if apiModel.Description == nil || *apiModel.Description != "My SWRV rule" {
		t.Fatal("expected description to map into API model")
	}
	if len(apiModel.InputSources) != 2 {
		t.Fatalf("expected 2 input sources, got %d", len(apiModel.InputSources))
	}
	if apiModel.InputSources[0].Type != "all_uploads" {
		t.Fatalf("expected first input source type all_uploads, got %q", apiModel.InputSources[0].Type)
	}
	if apiModel.InputSources[1].Id != 2345 {
		t.Fatalf("expected integration input source id 2345, got %d", apiModel.InputSources[1].Id)
	}
	if apiModel.Priority == nil || *apiModel.Priority != 4 {
		t.Fatal("expected priority to map into API model")
	}
	if apiModel.SupervisionSpaceId == nil || *apiModel.SupervisionSpaceId != 10420 {
		t.Fatal("expected supervision_space_id to map into API model")
	}
}

func TestSwrvRuleFromApiModel(t *testing.T) {
	description := "My SWRV rule"
	priority := int64(4)

	state := fromApiModel(thetalake.SwrvRule{
		DefaultRule: false,
		Description: &description,
		Id:          2337,
		InputSource: []thetalake.SwrvInputSource{
			{Name: "All Uploads"},
			{Integration: &thetalake.SwrvInputSourceIntegration{Id: 2345, Type: "integration", Name: "Zoom Chat"}},
		},
		IsBuiltIn: false,
		Name:      "swrv-example",
		Policy: &thetalake.SwrvPolicy{
			Id:   147,
			Name: "All Detections Active",
		},
		Priority: &priority,
		RetentionLibrary: &thetalake.SwrvRetentionLibrary{
			Id:   1,
			Name: "Test Bucket",
		},
		Search: nil,
		SupervisionSpace: &thetalake.SwrvSupervisionSpace{
			Id:   10420,
			Name: "Executive Team",
		},
		Workflow: &thetalake.SwrvWorkflow{
			Id:   14536,
			Name: "My workflow",
		},
	})

	if state.Id.ValueInt64() != 2337 {
		t.Fatalf("expected id 2337, got %d", state.Id.ValueInt64())
	}
	if state.Description.IsNull() || state.Description.ValueString() != "My SWRV rule" {
		t.Fatal("expected description to map into state")
	}
	if state.PolicyId.ValueInt64() != 147 || state.PolicyName.ValueString() != "All Detections Active" {
		t.Fatalf("expected policy details to map into state, got id=%d name=%q", state.PolicyId.ValueInt64(), state.PolicyName.ValueString())
	}
	if state.Priority.IsNull() || state.Priority.ValueInt64() != 4 {
		t.Fatal("expected priority to map into state")
	}
	if state.SupervisionSpaceId.IsNull() || state.SupervisionSpaceId.ValueInt64() != 10420 {
		t.Fatal("expected supervision space id to map into state")
	}
	if got := len(state.InputSources.Elements()); got != 2 {
		t.Fatalf("expected 2 input sources, got %d", got)
	}
	if !state.SearchId.IsNull() || !state.SearchName.IsNull() {
		t.Fatal("expected search fields to be null when API omits search")
	}
}
