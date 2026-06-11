package swrvrule

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	providerschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type swrvRuleTestRoute struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func newSwrvRuleTestClient(t *testing.T, routes ...swrvRuleTestRoute) *thetalake.Client {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"test-access-token","token_type":"Bearer","expires_in":3600}`))
	})

	for _, route := range routes {
		rt := route
		mux.HandleFunc("/api/v1"+rt.path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != rt.method {
				t.Fatalf("expected %s %s, got %s %s", rt.method, rt.path, r.Method, r.URL.Path)
			}
			if got := r.Header.Get("Authorization"); got != "Bearer test-access-token" {
				t.Fatalf("expected bearer token header, got %q", got)
			}
			rt.handler(w, r)
		})
	}

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return thetalake.NewClient(server.URL, "test-client-id", "test-client-secret")
}

func swrvRuleSchema(t *testing.T) providerschema.Schema {
	t.Helper()
	r := &swrvRuleResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	return resp.Schema
}

func swrvRulePlan(t *testing.T, model swrvRuleStateModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: swrvRuleSchema(t)}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to seed plan: %v", diags)
	}
	return plan
}

func swrvRuleState(t *testing.T, model swrvRuleStateModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: swrvRuleSchema(t)}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to seed state: %v", diags)
	}
	return state
}

func swrvRuleInputSourcesList(values ...map[string]attr.Value) types.List {
	objects := make([]attr.Value, 0, len(values))
	for _, value := range values {
		objects = append(objects, types.ObjectValueMust(swrvRuleInputSourceAttrTypes, value))
	}
	return types.ListValueMust(swrvRuleInputSourceObjectType, objects)
}

func TestNewSwrvRuleResource(t *testing.T) {
	if _, ok := NewSwrvRuleResource().(*swrvRuleResource); !ok {
		t.Fatal("expected NewSwrvRuleResource to return *swrvRuleResource")
	}
}

func TestSwrvRuleMetadata(t *testing.T) {
	r := &swrvRuleResource{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "thetalake"}, resp)

	if resp.TypeName != "thetalake_swrv_rule" {
		t.Fatalf("expected type name thetalake_swrv_rule, got %q", resp.TypeName)
	}
}

func TestSwrvRuleConfigure(t *testing.T) {
	r := &swrvRuleResource{}
	resp := &resource.ConfigureResponse{}
	client := &thetalake.Client{}

	r.Configure(context.Background(), resource.ConfigureRequest{}, resp)
	if r.client != nil {
		t.Fatal("expected client to remain nil when provider data is absent")
	}

	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: client}, resp)
	if r.client != client {
		t.Fatal("expected client to be assigned from provider data")
	}
}

func TestSwrvRuleSchema(t *testing.T) {
	r := &swrvRuleResource{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Fatal("expected schema description to be populated")
	}

	inputSourcesAttr, ok := resp.Schema.Attributes["input_sources"].(providerschema.ListNestedAttribute)
	if !ok {
		t.Fatal("expected input_sources to be a list nested attribute")
	}
	if !inputSourcesAttr.Required {
		t.Fatal("expected input_sources to be required")
	}
	if len(inputSourcesAttr.Validators) != 1 {
		t.Fatalf("expected input_sources to have 1 validator, got %d", len(inputSourcesAttr.Validators))
	}
	if _, ok := inputSourcesAttr.Validators[0].(inputSourcesValidator); !ok {
		t.Fatalf("expected input_sources validator to be %T, got %T", inputSourcesValidator{}, inputSourcesAttr.Validators[0])
	}

	priorityAttr, ok := resp.Schema.Attributes["priority"].(providerschema.Int64Attribute)
	if !ok {
		t.Fatal("expected priority to be an int64 attribute")
	}
	if !priorityAttr.Optional || !priorityAttr.Computed {
		t.Fatal("expected priority to be optional and computed")
	}

	workflowAttr, ok := resp.Schema.Attributes["workflow_id"].(providerschema.Int64Attribute)
	if !ok {
		t.Fatal("expected workflow_id to be an int64 attribute")
	}
	if !workflowAttr.Required {
		t.Fatal("expected workflow_id to be required")
	}
}

func TestInputSourcesValidatorRequiresIDForIntegration(t *testing.T) {
	listValue := types.ListValueMust(swrvRuleInputSourceObjectType, []attr.Value{
		types.ObjectValueMust(swrvRuleInputSourceAttrTypes, map[string]attr.Value{
			"id":   types.Int64Null(),
			"type": types.StringValue("integration"),
		}),
	})

	request := validator.ListRequest{
		Path:        path.Root("input_sources"),
		ConfigValue: listValue,
	}
	response := &validator.ListResponse{}

	inputSourcesValidator{}.ValidateList(context.Background(), request, response)

	if !response.Diagnostics.HasError() {
		t.Fatal("expected validator to reject integration input source without id")
	}
	if got := len(response.Diagnostics); got != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", got)
	}
}

func TestInputSourcesValidatorAllowsNonIntegrationWithoutID(t *testing.T) {
	listValue := types.ListValueMust(swrvRuleInputSourceObjectType, []attr.Value{
		types.ObjectValueMust(swrvRuleInputSourceAttrTypes, map[string]attr.Value{
			"id":   types.Int64Null(),
			"type": types.StringValue("all_uploads"),
		}),
	})

	request := validator.ListRequest{
		Path:        path.Root("input_sources"),
		ConfigValue: listValue,
	}
	response := &validator.ListResponse{}

	inputSourcesValidator{}.ValidateList(context.Background(), request, response)

	if response.Diagnostics.HasError() {
		t.Fatalf("expected validator to allow non-integration input source without id, got diagnostics: %v", response.Diagnostics)
	}
}

func TestSwrvRuleImportStateInvalidID(t *testing.T) {
	r := &swrvRuleResource{}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "not-a-number"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected invalid import ID to add an error diagnostic")
	}
}

func TestSwrvRuleCreate(t *testing.T) {
	r := &swrvRuleResource{client: newSwrvRuleTestClient(t, swrvRuleTestRoute{
		method: http.MethodPost,
		path:   "/workflows/swrv_rules",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"swrv_rule":{"default_rule":false,"description":"My SWRV rule","id":2337,"input_source":[{"name":"All Uploads","integration":null}],"is_built_in":false,"name":"swrv-example","policy":{"id":147,"name":"All Detections Active"},"priority":4,"retention_library":{"id":1,"name":"Test Bucket"},"search":null,"supervision_space":null,"workflow":{"id":14536,"name":"My workflow"}}}`))
		},
	})}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: swrvRuleSchema(t)}}

	r.Create(context.Background(), resource.CreateRequest{Plan: swrvRulePlan(t, swrvRuleStateModel{
		DefaultRule:          types.BoolNull(),
		Description:          types.StringValue("My SWRV rule"),
		Id:                   types.Int64Null(),
		InputSources:         swrvRuleInputSourcesList(map[string]attr.Value{"id": types.Int64Null(), "type": types.StringValue("all_uploads")}),
		IsBuiltIn:            types.BoolNull(),
		Name:                 types.StringValue("swrv-example"),
		PolicyId:             types.Int64Value(147),
		PolicyName:           types.StringNull(),
		Priority:             types.Int64Value(4),
		RetentionLibraryId:   types.Int64Value(1),
		RetentionLibraryName: types.StringNull(),
		SearchId:             types.Int64Null(),
		SearchName:           types.StringNull(),
		SupervisionSpaceId:   types.Int64Null(),
		SupervisionSpaceName: types.StringNull(),
		WorkflowId:           types.Int64Value(14536),
		WorkflowName:         types.StringNull(),
	})}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected create to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state swrvRuleStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read create state: %v", diags)
	}
	if state.Id.ValueInt64() != 2337 {
		t.Fatalf("expected created SWRV rule id 2337, got %d", state.Id.ValueInt64())
	}
	if got := len(state.InputSources.Elements()); got != 1 {
		t.Fatalf("expected input_sources to be preserved from plan, got %d entries", got)
	}
}

func TestSwrvRuleRead(t *testing.T) {
	r := &swrvRuleResource{client: newSwrvRuleTestClient(t, swrvRuleTestRoute{
		method: http.MethodGet,
		path:   "/workflows/swrv_rules/2337",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"swrv_rule":{"default_rule":false,"description":"Updated SWRV rule","id":2337,"input_source":[{"name":"All Uploads","integration":null}],"is_built_in":false,"name":"swrv-example","policy":{"id":147,"name":"All Detections Active"},"priority":0,"retention_library":{"id":1,"name":"Test Bucket"},"search":{"id":739139,"name":"All org users"},"supervision_space":{"id":10420,"name":"Executive Team"},"workflow":{"id":14536,"name":"My workflow"}}}`))
		},
	})}
	currentState := swrvRuleState(t, swrvRuleStateModel{
		Description:        types.StringValue("My SWRV rule"),
		Id:                 types.Int64Value(2337),
		InputSources:       swrvRuleInputSourcesList(map[string]attr.Value{"id": types.Int64Null(), "type": types.StringValue("all_uploads")}),
		Name:               types.StringValue("swrv-example"),
		PolicyId:           types.Int64Value(147),
		Priority:           types.Int64Value(4),
		RetentionLibraryId: types.Int64Value(1),
		WorkflowId:         types.Int64Value(14536),
	})
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: swrvRuleSchema(t)}}

	r.Read(context.Background(), resource.ReadRequest{State: currentState}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected read to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state swrvRuleStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read refreshed state: %v", diags)
	}
	if state.Priority.ValueInt64() != 0 {
		t.Fatalf("expected priority 0, got %d", state.Priority.ValueInt64())
	}
	if state.SearchId.ValueInt64() != 739139 {
		t.Fatalf("expected search_id 739139, got %d", state.SearchId.ValueInt64())
	}
	if state.SupervisionSpaceId.ValueInt64() != 10420 {
		t.Fatalf("expected supervision_space_id 10420, got %d", state.SupervisionSpaceId.ValueInt64())
	}
	if got := len(state.InputSources.Elements()); got != 1 {
		t.Fatalf("expected input_sources to remain populated, got %d entries", got)
	}
}

func TestSwrvRuleUpdate(t *testing.T) {
	r := &swrvRuleResource{client: newSwrvRuleTestClient(t, swrvRuleTestRoute{
		method: http.MethodPut,
		path:   "/workflows/swrv_rules/2337",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"swrv_rule":{"default_rule":false,"description":"Updated SWRV rule","id":2337,"input_source":[{"name":"All User Uploads","integration":null}],"is_built_in":false,"name":"swrv-example-updated","policy":{"id":147,"name":"All Detections Active"},"priority":0,"retention_library":{"id":1,"name":"Test Bucket"},"search":null,"supervision_space":null,"workflow":{"id":14536,"name":"My workflow"}}}`))
		},
	})}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: swrvRuleSchema(t)}}

	r.Update(context.Background(), resource.UpdateRequest{
		Plan: swrvRulePlan(t, swrvRuleStateModel{
			Description:          types.StringValue("Updated SWRV rule"),
			Id:                   types.Int64Value(2337),
			InputSources:         swrvRuleInputSourcesList(map[string]attr.Value{"id": types.Int64Null(), "type": types.StringValue("all_user_uploads")}),
			Name:                 types.StringValue("swrv-example-updated"),
			PolicyId:             types.Int64Value(147),
			Priority:             types.Int64Value(0),
			RetentionLibraryId:   types.Int64Value(1),
			WorkflowId:           types.Int64Value(14536),
			DefaultRule:          types.BoolNull(),
			IsBuiltIn:            types.BoolNull(),
			PolicyName:           types.StringNull(),
			RetentionLibraryName: types.StringNull(),
			SearchId:             types.Int64Null(),
			SearchName:           types.StringNull(),
			SupervisionSpaceId:   types.Int64Null(),
			SupervisionSpaceName: types.StringNull(),
			WorkflowName:         types.StringNull(),
		}),
		State: swrvRuleState(t, swrvRuleStateModel{Id: types.Int64Value(2337), InputSources: types.ListNull(swrvRuleInputSourceObjectType)}),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected update to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state swrvRuleStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags)
	}
	if state.Name.ValueString() != "swrv-example-updated" {
		t.Fatalf("expected updated name, got %q", state.Name.ValueString())
	}
	if state.Priority.ValueInt64() != 0 {
		t.Fatalf("expected updated priority 0, got %d", state.Priority.ValueInt64())
	}
	if got := len(state.InputSources.Elements()); got != 1 {
		t.Fatalf("expected updated input_sources to be preserved, got %d entries", got)
	}
}

func TestSwrvRuleDelete(t *testing.T) {
	r := &swrvRuleResource{client: newSwrvRuleTestClient(t, swrvRuleTestRoute{
		method: http.MethodDelete,
		path:   "/workflows/swrv_rules/2337",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"Workflow deleted successfully"}`))
		},
	})}
	resp := &resource.DeleteResponse{}

	r.Delete(context.Background(), resource.DeleteRequest{State: swrvRuleState(t, swrvRuleStateModel{Id: types.Int64Value(2337), InputSources: types.ListNull(swrvRuleInputSourceObjectType)})}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected delete to succeed, got diagnostics: %v", resp.Diagnostics)
	}
}
