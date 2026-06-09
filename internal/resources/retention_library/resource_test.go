package retentionlibrary

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	providerschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("failed to parse time %q: %v", value, err)
	}
	return parsed
}

type retentionLibraryTestRoute struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func newRetentionLibraryTestClient(t *testing.T, routes ...retentionLibraryTestRoute) *thetalake.Client {
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

func retentionLibrarySchemaForTest(t *testing.T) providerschema.Schema {
	t.Helper()
	r := &retentionLibraryResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	return resp.Schema
}

func retentionLibraryPlanForTest(t *testing.T, model retentionLibraryStateModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: retentionLibrarySchemaForTest(t)}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to seed plan: %v", diags)
	}
	return plan
}

func retentionLibraryStateForTest(t *testing.T, model retentionLibraryStateModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: retentionLibrarySchemaForTest(t)}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to seed state: %v", diags)
	}
	return state
}

func TestNewRetentionLibraryResource(t *testing.T) {
	if _, ok := NewRetentionLibraryResource().(*retentionLibraryResource); !ok {
		t.Fatal("expected NewRetentionLibraryResource to return *retentionLibraryResource")
	}
}

func TestRetentionLibraryMetadata(t *testing.T) {
	r := &retentionLibraryResource{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "thetalake"}, resp)

	if resp.TypeName != "thetalake_retention_library" {
		t.Fatalf("expected type name thetalake_retention_library, got %q", resp.TypeName)
	}
}

func TestRetentionLibraryConfigure(t *testing.T) {
	r := &retentionLibraryResource{}
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

func TestRetentionLibrarySchema(t *testing.T) {
	r := &retentionLibraryResource{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Fatal("expected schema description to be populated")
	}

	nameAttr, ok := resp.Schema.Attributes["name"].(providerschema.StringAttribute)
	if !ok {
		t.Fatal("expected name to be a string attribute")
	}
	if !nameAttr.Required {
		t.Fatal("expected name to be required")
	}

	descriptionAttr, ok := resp.Schema.Attributes["description"].(providerschema.StringAttribute)
	if !ok {
		t.Fatal("expected description to be a string attribute")
	}
	if !descriptionAttr.Optional || !descriptionAttr.Computed {
		t.Fatal("expected description to be optional and computed")
	}
	if len(descriptionAttr.PlanModifiers) == 0 {
		t.Fatal("expected description to preserve state for unknown values")
	}

	retainInReviewAttr, ok := resp.Schema.Attributes["retain_in_review"].(providerschema.BoolAttribute)
	if !ok {
		t.Fatal("expected retain_in_review to be a bool attribute")
	}
	if !retainInReviewAttr.Optional || !retainInReviewAttr.Computed {
		t.Fatal("expected retain_in_review to be optional and computed")
	}

	storageAccountAttr, ok := resp.Schema.Attributes["storage_account_id"].(providerschema.Int64Attribute)
	if !ok {
		t.Fatal("expected storage_account_id to be an int64 attribute")
	}
	if !storageAccountAttr.Required {
		t.Fatal("expected storage_account_id to be required")
	}
}

func TestRetentionLibraryImportStateInvalidID(t *testing.T) {
	r := &retentionLibraryResource{}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "not-a-number"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected invalid import ID to add an error diagnostic")
	}
}

func TestRetentionLibraryCreate(t *testing.T) {
	externalID := "external-123"
	r := &retentionLibraryResource{client: newRetentionLibraryTestClient(t, retentionLibraryTestRoute{
		method: http.MethodPost,
		path:   "/retention_libraries",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"retention_library":{"created_at":"2024-02-03T04:05:06Z","datum_count":12,"datum_size":34,"delete_on_expiration":true,"description":"Retention library description","display_name":"Retention Library (us-east-1)","external_id":"external-123","id":477,"legal_hold_count":5,"name":"Retention Library","retain_in_review":false,"retention_period_days":30,"retention_period_enabled":true,"retention_summary_text":"Retained for 30 days","sec_compliant_storage_confirmed":true,"sec_compliant_storage_enabled":false,"storage_account_id":7,"swrv_rule_count":2,"updated_at":"2024-02-03T07:05:06Z"}}`))
		},
	})}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: retentionLibrarySchemaForTest(t)}}

	r.Create(context.Background(), resource.CreateRequest{Plan: retentionLibraryPlanForTest(t, retentionLibraryStateModel{
		CreatedAt:                    timetypes.NewRFC3339Null(),
		DatumCount:                   types.Int64Null(),
		DatumSize:                    types.Int64Null(),
		DeleteOnExpiration:           types.BoolNull(),
		Description:                  types.StringValue("Retention library description"),
		DisplayName:                  types.StringNull(),
		ExternalId:                   types.StringValue(externalID),
		Id:                           types.Int64Null(),
		LegalHoldCount:               types.Int64Null(),
		Name:                         types.StringValue("Retention Library"),
		RetainInReview:               types.BoolValue(true),
		RetentionPeriodDays:          types.Int64Value(30),
		RetentionPeriodEnabled:       types.BoolValue(true),
		RetentionSummaryText:         types.StringNull(),
		SecCompliantStorageConfirmed: types.BoolNull(),
		SecCompliantStorageEnabled:   types.BoolValue(false),
		StorageAccountId:             types.Int64Value(7),
		SwrvRuleCount:                types.Int64Null(),
		UpdatedAt:                    timetypes.NewRFC3339Null(),
	})}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected create to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state retentionLibraryStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read create state: %v", diags)
	}
	if state.Id.ValueInt64() != 477 {
		t.Fatalf("expected created retention library id 477, got %d", state.Id.ValueInt64())
	}
	if !state.RetainInReview.ValueBool() {
		t.Fatal("expected retain_in_review to be preserved from plan")
	}
}

func TestRetentionLibraryReadPreservesRetainInReview(t *testing.T) {
	r := &retentionLibraryResource{client: newRetentionLibraryTestClient(t, retentionLibraryTestRoute{
		method: http.MethodGet,
		path:   "/retention_libraries/477",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"retention_library":{"created_at":"2024-02-03T04:05:06Z","datum_count":12,"datum_size":34,"delete_on_expiration":true,"description":"Updated retention library description","display_name":"Retention Library (us-east-1)","external_id":"external-123","id":477,"legal_hold_count":5,"name":"Retention Library","retain_in_review":false,"retention_period_days":90,"retention_period_enabled":true,"retention_summary_text":"Retained for 90 days","sec_compliant_storage_confirmed":true,"sec_compliant_storage_enabled":false,"storage_account_id":7,"swrv_rule_count":2,"updated_at":"2024-02-03T07:05:06Z"}}`))
		},
	})}
	currentState := retentionLibraryStateForTest(t, retentionLibraryStateModel{
		CreatedAt:              timetypes.NewRFC3339TimeValue(mustParseTime(t, "2024-02-03T04:05:06Z")),
		Description:            types.StringValue("Retention library description"),
		ExternalId:             types.StringValue("external-123"),
		Id:                     types.Int64Value(477),
		Name:                   types.StringValue("Retention Library"),
		RetainInReview:         types.BoolValue(true),
		RetentionPeriodDays:    types.Int64Value(30),
		RetentionPeriodEnabled: types.BoolValue(true),
		StorageAccountId:       types.Int64Value(7),
		UpdatedAt:              timetypes.NewRFC3339TimeValue(mustParseTime(t, "2024-02-03T07:05:06Z")),
	})
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: retentionLibrarySchemaForTest(t)}}

	r.Read(context.Background(), resource.ReadRequest{State: currentState}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected read to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state retentionLibraryStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read refreshed state: %v", diags)
	}
	if !state.RetainInReview.ValueBool() {
		t.Fatal("expected retain_in_review to be preserved from prior state")
	}
	if state.RetentionPeriodDays.ValueInt64() != 90 {
		t.Fatalf("expected refreshed retention_period_days 90, got %d", state.RetentionPeriodDays.ValueInt64())
	}
}

func TestRetentionLibraryReadNotFoundRemovesResource(t *testing.T) {
	r := &retentionLibraryResource{client: newRetentionLibraryTestClient(t, retentionLibraryTestRoute{
		method: http.MethodGet,
		path:   "/retention_libraries/477",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"not found"}`))
		},
	})}
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: retentionLibrarySchemaForTest(t)}}

	r.Read(context.Background(), resource.ReadRequest{State: retentionLibraryStateForTest(t, retentionLibraryStateModel{Id: types.Int64Value(477)})}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected not-found read to remove state without diagnostics, got: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatal("expected state to be removed on not found")
	}
}

func TestRetentionLibraryUpdate(t *testing.T) {
	externalID := "external-123"
	r := &retentionLibraryResource{client: newRetentionLibraryTestClient(t, retentionLibraryTestRoute{
		method: http.MethodPut,
		path:   "/retention_libraries/477",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"retention_library":{"created_at":"2024-02-03T04:05:06Z","datum_count":12,"datum_size":34,"delete_on_expiration":true,"description":"Updated retention library description","display_name":"Retention Library (us-east-1)","external_id":"external-123","id":477,"legal_hold_count":5,"name":"Retention Library","retain_in_review":false,"retention_period_days":90,"retention_period_enabled":true,"retention_summary_text":"Retained for 90 days","sec_compliant_storage_confirmed":true,"sec_compliant_storage_enabled":false,"storage_account_id":7,"swrv_rule_count":2,"updated_at":"2024-02-03T07:05:06Z"}}`))
		},
	})}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: retentionLibrarySchemaForTest(t)}}

	r.Update(context.Background(), resource.UpdateRequest{
		Plan: retentionLibraryPlanForTest(t, retentionLibraryStateModel{
			CreatedAt:                    timetypes.NewRFC3339Null(),
			DatumCount:                   types.Int64Null(),
			DatumSize:                    types.Int64Null(),
			DeleteOnExpiration:           types.BoolNull(),
			Description:                  types.StringValue("Updated retention library description"),
			DisplayName:                  types.StringNull(),
			ExternalId:                   types.StringValue(externalID),
			Id:                           types.Int64Value(477),
			LegalHoldCount:               types.Int64Null(),
			Name:                         types.StringValue("Retention Library"),
			RetainInReview:               types.BoolValue(true),
			RetentionPeriodDays:          types.Int64Value(90),
			RetentionPeriodEnabled:       types.BoolValue(true),
			RetentionSummaryText:         types.StringNull(),
			SecCompliantStorageConfirmed: types.BoolNull(),
			SecCompliantStorageEnabled:   types.BoolValue(false),
			StorageAccountId:             types.Int64Value(7),
			SwrvRuleCount:                types.Int64Null(),
			UpdatedAt:                    timetypes.NewRFC3339Null(),
		}),
		State: retentionLibraryStateForTest(t, retentionLibraryStateModel{Id: types.Int64Value(477), RetainInReview: types.BoolValue(false)}),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected update to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state retentionLibraryStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags)
	}
	if !state.RetainInReview.ValueBool() {
		t.Fatal("expected retain_in_review to be preserved from plan during update")
	}
	if state.RetentionPeriodDays.ValueInt64() != 90 {
		t.Fatalf("expected updated retention_period_days 90, got %d", state.RetentionPeriodDays.ValueInt64())
	}
}

func TestRetentionLibraryDelete(t *testing.T) {
	deleteCalled := false
	r := &retentionLibraryResource{client: newRetentionLibraryTestClient(t, retentionLibraryTestRoute{
		method: http.MethodDelete,
		path:   "/retention_libraries/477",
		handler: func(w http.ResponseWriter, req *http.Request) {
			deleteCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		},
	})}
	resp := &resource.DeleteResponse{}

	r.Delete(context.Background(), resource.DeleteRequest{State: retentionLibraryStateForTest(t, retentionLibraryStateModel{Id: types.Int64Value(477)})}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected delete to succeed, got diagnostics: %v", resp.Diagnostics)
	}
	if !deleteCalled {
		t.Fatal("expected delete endpoint to be called")
	}
}
