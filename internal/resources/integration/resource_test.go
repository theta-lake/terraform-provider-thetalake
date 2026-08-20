package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	providerschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type integrationTestRoute struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func newIntegrationTestClient(t *testing.T, routes ...integrationTestRoute) *thetalake.Client {
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
			rt.handler(w, r)
		})
	}

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return thetalake.NewClient(server.URL, "test-client-id", "test-client-secret")
}

func integrationSchemaForTest(t *testing.T) providerschema.Schema {
	t.Helper()
	r := &integrationResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	return resp.Schema
}

func integrationPlanForTest(t *testing.T, model integrationStateModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: integrationSchemaForTest(t)}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to seed plan: %v", diags)
	}
	return plan
}

func integrationStateForTest(t *testing.T, model integrationStateModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: integrationSchemaForTest(t)}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to seed state: %v", diags)
	}
	return state
}

// thetaLakeApiOnlyModel returns a minimal integrationStateModel selecting the
// theta_lake_api type, suitable for seeding a plan/state in tests that don't care
// about the other computed fields.
func thetaLakeApiOnlyModel(id int64, name string, paused bool) integrationStateModel {
	return integrationStateModel{
		CreatedAt:            timetypes.NewRFC3339Null(),
		GenericJournaling:    types.ObjectNull(genericJournalingAttrTypes),
		GoogleWorkspaceEmail: types.ObjectNull(googleWorkspaceEmailAttrTypes),
		Id:                   types.Int64Value(id),
		IntegrationGroup:     types.StringNull(),
		IntegrationType:      types.StringNull(),
		IntegrationTypeId:    types.Int64Null(),
		Name:                 types.StringValue(name),
		Paused:               types.BoolValue(paused),
		ThetaLakeApi:         types.ObjectValueMust(thetaLakeApiAttrTypes, map[string]attr.Value{}),
		Type:                 types.StringNull(),
	}
}

func TestNewIntegrationResource(t *testing.T) {
	if _, ok := NewIntegrationResource().(*integrationResource); !ok {
		t.Fatal("expected NewIntegrationResource to return *integrationResource")
	}
}

func TestIntegrationMetadata(t *testing.T) {
	r := &integrationResource{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "thetalake"}, resp)

	if resp.TypeName != "thetalake_integration" {
		t.Fatalf("expected type name thetalake_integration, got %q", resp.TypeName)
	}
}

func TestIntegrationConfigure(t *testing.T) {
	r := &integrationResource{}
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

func TestIntegrationSchema(t *testing.T) {
	schema := integrationSchemaForTest(t)

	if schema.MarkdownDescription == "" {
		t.Fatal("expected schema description to be populated")
	}

	nameAttr, ok := schema.Attributes["name"].(providerschema.StringAttribute)
	if !ok || !nameAttr.Required {
		t.Fatal("expected name to be a required string attribute")
	}

	pausedAttr, ok := schema.Attributes["paused"].(providerschema.BoolAttribute)
	if !ok || !pausedAttr.Optional || !pausedAttr.Computed {
		t.Fatal("expected paused to be optional and computed")
	}

	typeAttr, ok := schema.Attributes["type"].(providerschema.StringAttribute)
	if !ok || !typeAttr.Computed {
		t.Fatal("expected type to be computed")
	}

	genericJournalingAttr, ok := schema.Attributes["generic_journaling"].(providerschema.SingleNestedAttribute)
	if !ok || !genericJournalingAttr.Optional {
		t.Fatal("expected generic_journaling to be an optional single nested attribute")
	}
	if len(genericJournalingAttr.PlanModifiers) == 0 {
		t.Fatal("expected generic_journaling to force replacement on type change")
	}
	passwordAttr, ok := genericJournalingAttr.Attributes["undeliverable_email_password"].(providerschema.StringAttribute)
	if !ok || !passwordAttr.Sensitive {
		t.Fatal("expected undeliverable_email_password to be sensitive")
	}
	if !passwordAttr.Optional || !passwordAttr.Computed {
		t.Fatal("expected undeliverable_email_password to be optional and computed, like its mailbox companion attributes")
	}
	if len(passwordAttr.PlanModifiers) != 0 {
		t.Fatal("expected undeliverable_email_password to have no UseStateForUnknown plan modifier, so an out-of-band change is always surfaced as drift on refresh")
	}

	thetaLakeApiAttr, ok := schema.Attributes["theta_lake_api"].(providerschema.SingleNestedAttribute)
	if !ok || !thetaLakeApiAttr.Optional {
		t.Fatal("expected theta_lake_api to be an optional single nested attribute")
	}
	if len(thetaLakeApiAttr.Attributes) != 0 {
		t.Fatal("expected theta_lake_api to have no nested attributes")
	}
}

func TestIntegrationImportStateInvalidID(t *testing.T) {
	r := &integrationResource{}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "not-a-number"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected invalid import ID to add an error diagnostic")
	}
}

func TestIntegrationCreate(t *testing.T) {
	r := &integrationResource{client: newIntegrationTestClient(t, integrationTestRoute{
		method: http.MethodPost,
		path:   "/integrations",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"integration":{"id":305,"integration_group":"Collaboration Platform","integration_type":"Theta Lake API","integration_type_id":80,"name":"Custom Theta Lake API Integration","service_paused":false,"status":"Queued","service_params":{}}}`))
		},
	})}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	r.Create(context.Background(), resource.CreateRequest{
		Plan: integrationPlanForTest(t, thetaLakeApiOnlyModel(0, "Custom Theta Lake API Integration", false)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected create to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state integrationStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read create state: %v", diags)
	}
	if state.Id.ValueInt64() != 305 {
		t.Fatalf("expected created integration id 305, got %d", state.Id.ValueInt64())
	}
	if state.Type.ValueString() != thetalake.IntegrationTypeThetaLakeApi {
		t.Fatalf("expected derived type %q, got %q", thetalake.IntegrationTypeThetaLakeApi, state.Type.ValueString())
	}
}

func TestIntegrationCreateMissingResponseIdSurfacesError(t *testing.T) {
	// Unlike Update, Create has no prior state to fall back on for the id: if the
	// create response omits it, there is no correct id to track the resource by, so
	// this must surface as an explicit error diagnostic rather than silently querying
	// GetIntegrationConfiguration for /integrations/0/configuration.
	r := &integrationResource{client: newIntegrationTestClient(t, integrationTestRoute{
		method: http.MethodPost,
		path:   "/integrations",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"integration":{"integration_group":"Collaboration Platform","integration_type":"Generic Journaling","integration_type_id":41,"name":"Custom Generic Journaling Integration","service_paused":false,"status":"Queued"}}`))
		},
	})}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	planOptions := types.ObjectValueMust(genericJournalingAttrTypes, map[string]attr.Value{
		"download_o365_onedrive_links": types.BoolValue(false),
		"download_salesforce_doclinks": types.BoolValue(false),
		"index_headers":                types.StringNull(),
		"sender_spf_override":          types.StringNull(),
		"undeliverable_disabled":       types.BoolValue(false),
		"undeliverable_email_address":  types.StringNull(),
		"undeliverable_email_password": types.StringValue("SecretPassword"),
		"undeliverable_email_port":     types.Int64Null(),
		"undeliverable_email_server":   types.StringValue("email.realbank.com"),
		"undeliverable_email_user":     types.StringNull(),
	})

	r.Create(context.Background(), resource.CreateRequest{
		Plan: integrationPlanForTest(t, genericJournalingOnlyModel(0, "Custom Generic Journaling Integration", false, planOptions)),
	}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected create to fail with an error diagnostic when the response omits the integration id")
	}
}

// genericJournalingOnlyModel returns a minimal integrationStateModel selecting the
// generic_journaling type, suitable for seeding a plan/state in tests that don't care
// about the other computed fields.
func genericJournalingOnlyModel(id int64, name string, paused bool, options types.Object) integrationStateModel {
	return integrationStateModel{
		CreatedAt:            timetypes.NewRFC3339Null(),
		GenericJournaling:    options,
		GoogleWorkspaceEmail: types.ObjectNull(googleWorkspaceEmailAttrTypes),
		Id:                   types.Int64Value(id),
		IntegrationGroup:     types.StringNull(),
		IntegrationType:      types.StringNull(),
		IntegrationTypeId:    types.Int64Null(),
		Name:                 types.StringValue(name),
		Paused:               types.BoolValue(paused),
		ThetaLakeApi:         types.ObjectNull(thetaLakeApiAttrTypes),
		Type:                 types.StringNull(),
	}
}

func TestIntegrationCreateGenericJournalingMissingServiceParamsFallsBackToConfigurationEndpoint(t *testing.T) {
	// The create response omits service_params entirely for a generic_journaling
	// integration. Create must fall back to a GetIntegrationConfiguration call to
	// authoritatively populate generic_journaling, rather than copying the plan's
	// block into state: in a real Terraform plan, any Optional+Computed member the
	// user didn't configure (index_headers, sender_spf_override,
	// undeliverable_email_address/port/server/user) is unknown, not null, and an
	// unknown value must never end up in the applied state.
	r := &integrationResource{client: newIntegrationTestClient(t,
		integrationTestRoute{
			method: http.MethodPost,
			path:   "/integrations",
			handler: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"integration":{"id":302,"integration_group":"Collaboration Platform","integration_type":"Generic Journaling","integration_type_id":41,"name":"Custom Generic Journaling Integration","service_paused":false,"status":"Queued"}}`))
			},
		},
		integrationTestRoute{
			method: http.MethodGet,
			path:   "/integrations/302/configuration",
			handler: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"configuration":{"integration_type":"Generic Journaling","integration_type_id":41,"options":{"download_o365_onedrive_links":false,"download_salesforce_doclinks":false,"index_headers":"X-Header-Score","sender_spf_override":"v=spf1 -all","undeliverable_disabled":false,"undeliverable_email_address":"undeliverable@example.com","undeliverable_email_password":"SecretPassword","undeliverable_email_port":993,"undeliverable_email_server":"email.realbank.com","undeliverable_email_user":"Undeliverable User"}}}`))
			},
		},
	)}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	// A real Terraform plan leaves unconfigured Optional+Computed members unknown,
	// not null.
	planOptions := types.ObjectValueMust(genericJournalingAttrTypes, map[string]attr.Value{
		"download_o365_onedrive_links": types.BoolValue(false),
		"download_salesforce_doclinks": types.BoolValue(false),
		"index_headers":                types.StringUnknown(),
		"sender_spf_override":          types.StringUnknown(),
		"undeliverable_disabled":       types.BoolValue(false),
		"undeliverable_email_address":  types.StringUnknown(),
		"undeliverable_email_password": types.StringValue("SecretPassword"),
		"undeliverable_email_port":     types.Int64Unknown(),
		"undeliverable_email_server":   types.StringValue("email.realbank.com"),
		"undeliverable_email_user":     types.StringUnknown(),
	})

	r.Create(context.Background(), resource.CreateRequest{
		Plan: integrationPlanForTest(t, genericJournalingOnlyModel(0, "Custom Generic Journaling Integration", false, planOptions)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected create to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	if !resp.State.Raw.IsFullyKnown() {
		t.Fatal("expected no attribute in the resulting state to be unknown")
	}

	var state integrationStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read create state: %v", diags)
	}
	if state.GenericJournaling.IsNull() || state.GenericJournaling.IsUnknown() {
		t.Fatal("expected generic_journaling to be populated from the configuration endpoint fallback")
	}

	var decoded genericJournalingModel
	if diags := state.GenericJournaling.As(context.Background(), &decoded, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("failed to decode generic_journaling: %v", diags)
	}
	if decoded.UndeliverableEmailServer.ValueString() != "email.realbank.com" {
		t.Fatalf("expected undeliverable_email_server from the configuration endpoint, got %q", decoded.UndeliverableEmailServer.ValueString())
	}
	if decoded.UndeliverableEmailPassword.ValueString() != "SecretPassword" {
		t.Fatalf("expected undeliverable_email_password from the configuration endpoint, got %q", decoded.UndeliverableEmailPassword.ValueString())
	}
	if decoded.IndexHeaders.ValueString() != "X-Header-Score" {
		t.Fatalf("expected index_headers from the configuration endpoint, got %q", decoded.IndexHeaders.ValueString())
	}
}

func TestIntegrationUpdateNoPauseChangePersistsPlanPausedNotResponse(t *testing.T) {
	// When plan.Paused == state.Paused == true, the pause/start branch is skipped, so
	// the persisted value must be the prior state's paused (which equals the plan's),
	// not whatever the PUT response happens to report for service_paused.
	r := &integrationResource{client: newIntegrationTestClient(t, integrationTestRoute{
		method: http.MethodPut,
		path:   "/integrations/305",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"integration":{"id":305,"integration_group":"Collaboration Platform","integration_type":"Theta Lake API","integration_type_id":80,"name":"Updated Theta Lake API Integration","service_paused":false,"status":"Queued","service_params":{}}}`))
		},
	})}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	// No pause/start route is registered; if either were called the test client
	// would fail the test with an unexpected request.
	r.Update(context.Background(), resource.UpdateRequest{
		Plan:  integrationPlanForTest(t, thetaLakeApiOnlyModel(305, "Updated Theta Lake API Integration", true)),
		State: integrationStateForTest(t, thetaLakeApiOnlyModel(305, "Custom Theta Lake API Integration", true)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected update to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state integrationStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags)
	}
	if !state.Paused.ValueBool() {
		t.Fatal("expected paused to remain true, taken from the plan, even though the update response reported service_paused: false")
	}
}

func TestIntegrationUpdateIdSurvivesResponseMissingId(t *testing.T) {
	// The PUT response for this test omits the integration id (defaults to 0 when
	// decoded). id has UseStateForUnknown, so state must keep the prior state's id
	// rather than taking the (missing) value from the response.
	r := &integrationResource{client: newIntegrationTestClient(t, integrationTestRoute{
		method: http.MethodPut,
		path:   "/integrations/305",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"integration":{"integration_group":"Collaboration Platform","integration_type":"Theta Lake API","integration_type_id":80,"name":"Updated Theta Lake API Integration","service_paused":false,"status":"Queued","service_params":{}}}`))
		},
	})}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	r.Update(context.Background(), resource.UpdateRequest{
		Plan:  integrationPlanForTest(t, thetaLakeApiOnlyModel(305, "Updated Theta Lake API Integration", false)),
		State: integrationStateForTest(t, thetaLakeApiOnlyModel(305, "Custom Theta Lake API Integration", false)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected update to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state integrationStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags)
	}
	if state.Id.ValueInt64() != 305 {
		t.Fatalf("expected id 305 to be carried over from prior state, got %d", state.Id.ValueInt64())
	}
}

func TestIntegrationUpdateConfigurationFallbackUsesPriorStateIdNotResponseId(t *testing.T) {
	// The PUT response for this test omits both id and service_params for a
	// generic_journaling integration (a type that has configuration options).
	// resolveServiceParams must query GetIntegrationConfiguration using the prior
	// state's id (305), not the response's missing id (which would decode to 0 and
	// target /integrations/0/configuration, failing even though the update itself
	// succeeded).
	configurationRequested := false
	r := &integrationResource{client: newIntegrationTestClient(t,
		integrationTestRoute{
			method: http.MethodPut,
			path:   "/integrations/305",
			handler: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"integration":{"name":"Updated Generic Journaling Integration","integration_type":"Generic Journaling","integration_type_id":41,"service_paused":false,"status":"Queued"}}`))
			},
		},
		integrationTestRoute{
			method: http.MethodGet,
			path:   "/integrations/305/configuration",
			handler: func(w http.ResponseWriter, req *http.Request) {
				configurationRequested = true
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"configuration":{"integration_type":"Generic Journaling","integration_type_id":41,"options":{"undeliverable_email_server":"email.realbank.com","undeliverable_email_password":"SecretPassword"}}}`))
			},
		},
	)}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	planOptions := types.ObjectValueMust(genericJournalingAttrTypes, map[string]attr.Value{
		"download_o365_onedrive_links": types.BoolValue(false),
		"download_salesforce_doclinks": types.BoolValue(false),
		"index_headers":                types.StringNull(),
		"sender_spf_override":          types.StringNull(),
		"undeliverable_disabled":       types.BoolValue(false),
		"undeliverable_email_address":  types.StringNull(),
		"undeliverable_email_password": types.StringValue("SecretPassword"),
		"undeliverable_email_port":     types.Int64Null(),
		"undeliverable_email_server":   types.StringValue("email.realbank.com"),
		"undeliverable_email_user":     types.StringNull(),
	})

	r.Update(context.Background(), resource.UpdateRequest{
		Plan:  integrationPlanForTest(t, genericJournalingOnlyModel(305, "Updated Generic Journaling Integration", false, planOptions)),
		State: integrationStateForTest(t, genericJournalingOnlyModel(305, "Custom Generic Journaling Integration", false, planOptions)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected update to succeed, got diagnostics: %v", resp.Diagnostics)
	}
	if !configurationRequested {
		t.Fatal("expected the configuration endpoint to be queried as a fallback")
	}

	var state integrationStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags)
	}
	if state.Id.ValueInt64() != 305 {
		t.Fatalf("expected id 305 to be carried over from prior state, got %d", state.Id.ValueInt64())
	}
	if state.GenericJournaling.IsNull() || state.GenericJournaling.IsUnknown() {
		t.Fatal("expected generic_journaling to be populated from the configuration endpoint fallback")
	}

	var decoded genericJournalingModel
	if diags := state.GenericJournaling.As(context.Background(), &decoded, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("failed to decode generic_journaling: %v", diags)
	}
	if decoded.UndeliverableEmailServer.ValueString() != "email.realbank.com" {
		t.Fatalf("expected undeliverable_email_server from the configuration endpoint, got %q", decoded.UndeliverableEmailServer.ValueString())
	}
}

func TestIntegrationReadNotFoundRemovesResource(t *testing.T) {
	r := &integrationResource{client: newIntegrationTestClient(t, integrationTestRoute{
		method: http.MethodGet,
		path:   "/integrations/305",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"not found"}`))
		},
	})}
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	r.Read(context.Background(), resource.ReadRequest{
		State: integrationStateForTest(t, thetaLakeApiOnlyModel(305, "Custom Theta Lake API Integration", false)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected not-found read to remove state without diagnostics, got: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatal("expected state to be removed on not found")
	}
}

func TestIntegrationUpdate(t *testing.T) {
	pauseCalled := false
	r := &integrationResource{client: newIntegrationTestClient(t,
		integrationTestRoute{
			method: http.MethodPut,
			path:   "/integrations/305",
			handler: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"integration":{"id":305,"integration_group":"Collaboration Platform","integration_type":"Theta Lake API","integration_type_id":80,"name":"Updated Theta Lake API Integration","service_paused":false,"status":"Queued","service_params":{}}}`))
			},
		},
		integrationTestRoute{
			method: http.MethodPut,
			path:   "/integrations/305/pause",
			handler: func(w http.ResponseWriter, req *http.Request) {
				pauseCalled = true
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status":"The integration has been paused"}`))
			},
		},
	)}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	r.Update(context.Background(), resource.UpdateRequest{
		Plan:  integrationPlanForTest(t, thetaLakeApiOnlyModel(305, "Updated Theta Lake API Integration", true)),
		State: integrationStateForTest(t, thetaLakeApiOnlyModel(305, "Custom Theta Lake API Integration", false)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected update to succeed, got diagnostics: %v", resp.Diagnostics)
	}
	if !pauseCalled {
		t.Fatal("expected pause endpoint to be called when paused changes from false to true")
	}

	var state integrationStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags)
	}
	if !state.Paused.ValueBool() {
		t.Fatal("expected paused to be true after update")
	}
	if state.Name.ValueString() != "Updated Theta Lake API Integration" {
		t.Fatalf("expected updated name, got %q", state.Name.ValueString())
	}
}

func TestIntegrationUpdateNoPauseChangeDoesNotCallPauseOrStart(t *testing.T) {
	r := &integrationResource{client: newIntegrationTestClient(t, integrationTestRoute{
		method: http.MethodPut,
		path:   "/integrations/305",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"integration":{"id":305,"integration_group":"Collaboration Platform","integration_type":"Theta Lake API","integration_type_id":80,"name":"Updated Theta Lake API Integration","service_paused":false,"status":"Queued","service_params":{}}}`))
		},
	})}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	// No pause/start route is registered; if either were called the test client
	// would fail the test with an unexpected request.
	r.Update(context.Background(), resource.UpdateRequest{
		Plan:  integrationPlanForTest(t, thetaLakeApiOnlyModel(305, "Updated Theta Lake API Integration", false)),
		State: integrationStateForTest(t, thetaLakeApiOnlyModel(305, "Custom Theta Lake API Integration", false)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected update to succeed, got diagnostics: %v", resp.Diagnostics)
	}
}

func TestIntegrationDelete(t *testing.T) {
	deleteCalled := false
	r := &integrationResource{client: newIntegrationTestClient(t, integrationTestRoute{
		method: http.MethodDelete,
		path:   "/integrations/305",
		handler: func(w http.ResponseWriter, req *http.Request) {
			deleteCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"The integration has been removed"}`))
		},
	})}
	resp := &resource.DeleteResponse{}

	r.Delete(context.Background(), resource.DeleteRequest{
		State: integrationStateForTest(t, thetaLakeApiOnlyModel(305, "Custom Theta Lake API Integration", false)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected delete to succeed, got diagnostics: %v", resp.Diagnostics)
	}
	if !deleteCalled {
		t.Fatal("expected delete endpoint to be called")
	}
}

func TestIntegrationReadConfigurationNotFoundRemovesResource(t *testing.T) {
	// A delete landing between the GetIntegrationById and GetIntegrationConfiguration
	// calls (or any integration type whose configuration endpoint 404s) must be
	// treated the same as GetIntegrationById returning not found: remove the resource
	// from state without an error diagnostic.
	r := &integrationResource{client: newIntegrationTestClient(t,
		integrationTestRoute{
			method: http.MethodGet,
			path:   "/integrations/302",
			handler: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"integration":{"id":302,"integration_group":"Collaboration Platform","integration_type":"Generic Journaling","integration_type_id":41,"name":"Custom Generic Journaling Integration","service_paused":false,"status":"Queued"}}`))
			},
		},
		integrationTestRoute{
			method: http.MethodGet,
			path:   "/integrations/302/configuration",
			handler: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"message":"not found"}`))
			},
		},
	)}
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	r.Read(context.Background(), resource.ReadRequest{
		State: integrationStateForTest(t, genericJournalingOnlyModel(302, "Custom Generic Journaling Integration", false,
			types.ObjectNull(genericJournalingAttrTypes))),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected not-found configuration read to remove state without diagnostics, got: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatal("expected state to be removed when the configuration endpoint reports not found")
	}
}

func TestIntegrationReadThetaLakeApiSkipsConfigurationCall(t *testing.T) {
	// theta_lake_api's configuration response is unused by fromApiModel, so Read
	// must skip the GetIntegrationConfiguration round-trip entirely for that type.
	r := &integrationResource{client: newIntegrationTestClient(t,
		integrationTestRoute{
			method: http.MethodGet,
			path:   "/integrations/305",
			handler: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"integration":{"id":305,"integration_group":"Collaboration Platform","integration_type":"Theta Lake API","integration_type_id":80,"name":"Custom Theta Lake API Integration","service_paused":false,"status":"Queued"}}`))
			},
		},
		integrationTestRoute{
			method: http.MethodGet,
			path:   "/integrations/305/configuration",
			handler: func(w http.ResponseWriter, req *http.Request) {
				t.Fatal("expected the configuration endpoint to not be called for a theta_lake_api integration")
			},
		},
	)}
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	r.Read(context.Background(), resource.ReadRequest{
		State: integrationStateForTest(t, thetaLakeApiOnlyModel(305, "Custom Theta Lake API Integration", false)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected read to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state integrationStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	if state.ThetaLakeApi.IsNull() {
		t.Fatal("expected theta_lake_api to be populated with an empty object")
	}
}

func TestIntegrationCreatePausedTrueSurvivesResponseOmittingServicePaused(t *testing.T) {
	// service_paused is a plain bool, so an omitted value in the create response
	// decodes to false. paused is sent in the create request, so the plan's value
	// (true) is authoritative and must be carried over rather than trusting the
	// response.
	r := &integrationResource{client: newIntegrationTestClient(t, integrationTestRoute{
		method: http.MethodPost,
		path:   "/integrations",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"integration":{"id":305,"integration_group":"Collaboration Platform","integration_type":"Theta Lake API","integration_type_id":80,"name":"Custom Theta Lake API Integration","status":"Queued","service_params":{}}}`))
		},
	})}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	r.Create(context.Background(), resource.CreateRequest{
		Plan: integrationPlanForTest(t, thetaLakeApiOnlyModel(0, "Custom Theta Lake API Integration", true)),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected create to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state integrationStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read create state: %v", diags)
	}
	if !state.Paused.ValueBool() {
		t.Fatal("expected paused to remain true, taken from the plan, even though the create response omitted service_paused")
	}
}

func TestIntegrationUpdatePreservesImmutableFieldsWhenResponseOmitsThem(t *testing.T) {
	// created_at, integration_group, integration_type, and integration_type_id are all
	// Computed with UseStateForUnknown and immutable after creation, same as id (see
	// TestIntegrationUpdateIdSurvivesResponseMissingId). A PUT response that omits any
	// of them must not clobber the prior state's known values with zero values.
	r := &integrationResource{client: newIntegrationTestClient(t, integrationTestRoute{
		method: http.MethodPut,
		path:   "/integrations/305",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"integration":{"name":"Updated Theta Lake API Integration","service_paused":false,"status":"Queued","service_params":{}}}`))
		},
	})}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: integrationSchemaForTest(t)}}

	priorState := integrationStateModel{
		CreatedAt:            timetypes.NewRFC3339ValueMust("2026-07-31T12:00:00Z"),
		GenericJournaling:    types.ObjectNull(genericJournalingAttrTypes),
		GoogleWorkspaceEmail: types.ObjectNull(googleWorkspaceEmailAttrTypes),
		Id:                   types.Int64Value(305),
		IntegrationGroup:     types.StringValue("Collaboration Platform"),
		IntegrationType:      types.StringValue("Theta Lake API"),
		IntegrationTypeId:    types.Int64Value(80),
		Name:                 types.StringValue("Custom Theta Lake API Integration"),
		Paused:               types.BoolValue(false),
		ThetaLakeApi:         types.ObjectValueMust(thetaLakeApiAttrTypes, map[string]attr.Value{}),
		Type:                 types.StringValue(thetalake.IntegrationTypeThetaLakeApi),
	}

	r.Update(context.Background(), resource.UpdateRequest{
		Plan:  integrationPlanForTest(t, thetaLakeApiOnlyModel(305, "Updated Theta Lake API Integration", false)),
		State: integrationStateForTest(t, priorState),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected update to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state integrationStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags)
	}
	if state.Id.ValueInt64() != 305 {
		t.Fatalf("expected id 305 to be carried over from prior state, got %d", state.Id.ValueInt64())
	}
	if !state.CreatedAt.Equal(priorState.CreatedAt) {
		t.Fatalf("expected created_at to be carried over from prior state, got %v", state.CreatedAt)
	}
	if state.IntegrationGroup.ValueString() != "Collaboration Platform" {
		t.Fatalf("expected integration_group to be carried over from prior state, got %q", state.IntegrationGroup.ValueString())
	}
	if state.IntegrationType.ValueString() != "Theta Lake API" {
		t.Fatalf("expected integration_type to be carried over from prior state, got %q", state.IntegrationType.ValueString())
	}
	if state.IntegrationTypeId.ValueInt64() != 80 {
		t.Fatalf("expected integration_type_id to be carried over from prior state, got %d", state.IntegrationTypeId.ValueInt64())
	}
}
