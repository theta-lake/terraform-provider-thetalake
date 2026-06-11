package label

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

type labelTestRoute struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func newLabelTestClient(t *testing.T, routes ...labelTestRoute) *thetalake.Client {
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

func labelSchema(t *testing.T) providerschema.Schema {
	t.Helper()
	r := &labelResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	return resp.Schema
}

func labelPlan(t *testing.T, model labelStateModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: labelSchema(t)}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to seed plan: %v", diags)
	}
	return plan
}

func labelState(t *testing.T, model labelStateModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: labelSchema(t)}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to seed state: %v", diags)
	}
	return state
}

func TestNewLabelResource(t *testing.T) {
	if _, ok := NewLabelResource().(*labelResource); !ok {
		t.Fatal("expected NewLabelResource to return *labelResource")
	}
}

func TestLabelMetadata(t *testing.T) {
	r := &labelResource{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "thetalake"}, resp)

	if resp.TypeName != "thetalake_label" {
		t.Fatalf("expected type name thetalake_label, got %q", resp.TypeName)
	}
}

func TestLabelConfigure(t *testing.T) {
	r := &labelResource{}
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

func TestLabelSchema(t *testing.T) {
	r := &labelResource{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Fatal("expected schema description to be populated")
	}

	backgroundColorAttr, ok := resp.Schema.Attributes["background_color"].(providerschema.StringAttribute)
	if !ok {
		t.Fatal("expected background_color to be a string attribute")
	}
	if !backgroundColorAttr.Required {
		t.Fatal("expected background_color to be required")
	}

	hiddenAttr, ok := resp.Schema.Attributes["hidden"].(providerschema.BoolAttribute)
	if !ok {
		t.Fatal("expected hidden to be a bool attribute")
	}
	if !hiddenAttr.Optional || !hiddenAttr.Computed {
		t.Fatal("expected hidden to be optional and computed")
	}

	createdAtAttr, ok := resp.Schema.Attributes["created_at"].(providerschema.StringAttribute)
	if !ok {
		t.Fatal("expected created_at to be a string attribute")
	}
	if !createdAtAttr.Computed {
		t.Fatal("expected created_at to be computed")
	}
}

func TestLabelImportStateInvalidID(t *testing.T) {
	r := &labelResource{}
	resp := &resource.ImportStateResponse{}

	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "not-a-number"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected invalid import ID to add an error diagnostic")
	}
}

func TestLabelCreate(t *testing.T) {
	r := &labelResource{client: newLabelTestClient(t, labelTestRoute{
		method: http.MethodPost,
		path:   "/labels",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"label":{"background_color":"#FFC906","created_at":"2024-01-02T03:04:05Z","hidden":false,"id":5,"long_name":"Label description","org_unit_id":108,"short_name":"Label","tagged_datums_count":7,"updated_at":"2024-01-02T05:04:05Z","user_id":422}}`))
		},
	})}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: labelSchema(t)}}

	r.Create(context.Background(), resource.CreateRequest{Plan: labelPlan(t, labelStateModel{
		BackgroundColor:   types.StringValue("#FFC906"),
		CreatedAt:         timetypes.NewRFC3339Null(),
		Hidden:            types.BoolValue(false),
		Id:                types.Int64Null(),
		LongName:          types.StringValue("Label description"),
		OrgUnitId:         types.Int64Null(),
		ShortName:         types.StringValue("Label"),
		TaggedDatumsCount: types.Int64Null(),
		UpdatedAt:         timetypes.NewRFC3339Null(),
		UserId:            types.Int64Null(),
	})}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected create to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state labelStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read create state: %v", diags)
	}
	if state.Id.ValueInt64() != 5 || state.ShortName.ValueString() != "Label" {
		t.Fatalf("unexpected created state: %#v", state)
	}
}

func TestLabelRead(t *testing.T) {
	r := &labelResource{client: newLabelTestClient(t, labelTestRoute{
		method: http.MethodGet,
		path:   "/labels/5",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"label":{"background_color":"#FF0000","created_at":"2024-01-02T03:04:05Z","hidden":true,"id":5,"long_name":"Updated label description","org_unit_id":108,"short_name":"Updated Label","tagged_datums_count":9,"updated_at":"2024-01-02T06:04:05Z","user_id":422}}`))
		},
	})}
	currentState := labelState(t, labelStateModel{
		BackgroundColor: types.StringValue("#FFC906"),
		CreatedAt:       timetypes.NewRFC3339TimeValue(mustParseTime(t, "2024-01-02T03:04:05Z")),
		Hidden:          types.BoolValue(false),
		Id:              types.Int64Value(5),
		LongName:        types.StringValue("Label description"),
		OrgUnitId:       types.Int64Value(108),
		ShortName:       types.StringValue("Label"),
		UserId:          types.Int64Value(422),
	})
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: labelSchema(t)}}

	r.Read(context.Background(), resource.ReadRequest{State: currentState}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected read to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state labelStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read refreshed state: %v", diags)
	}
	if state.ShortName.ValueString() != "Updated Label" || !state.Hidden.ValueBool() {
		t.Fatalf("unexpected refreshed state: %#v", state)
	}
}

func TestLabelUpdate(t *testing.T) {
	r := &labelResource{client: newLabelTestClient(t, labelTestRoute{
		method: http.MethodPut,
		path:   "/labels/5",
		handler: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"label":{"background_color":"#FF0000","created_at":"2024-01-02T03:04:05Z","hidden":true,"id":5,"long_name":"Updated label description","org_unit_id":108,"short_name":"Updated Label","tagged_datums_count":9,"updated_at":"2024-01-02T06:04:05Z","user_id":422}}`))
		},
	})}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: labelSchema(t)}}

	r.Update(context.Background(), resource.UpdateRequest{
		Plan: labelPlan(t, labelStateModel{
			BackgroundColor:   types.StringValue("#FF0000"),
			CreatedAt:         timetypes.NewRFC3339Null(),
			Hidden:            types.BoolValue(true),
			Id:                types.Int64Value(5),
			LongName:          types.StringValue("Updated label description"),
			OrgUnitId:         types.Int64Null(),
			ShortName:         types.StringValue("Updated Label"),
			TaggedDatumsCount: types.Int64Null(),
			UpdatedAt:         timetypes.NewRFC3339Null(),
			UserId:            types.Int64Null(),
		}),
		State: labelState(t, labelStateModel{Id: types.Int64Value(5)}),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected update to succeed, got diagnostics: %v", resp.Diagnostics)
	}

	var state labelStateModel
	if diags := resp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags)
	}
	if state.ShortName.ValueString() != "Updated Label" || state.BackgroundColor.ValueString() != "#FF0000" {
		t.Fatalf("unexpected updated state: %#v", state)
	}
}

func TestLabelDelete(t *testing.T) {
	deleteCalled := false
	r := &labelResource{client: newLabelTestClient(t, labelTestRoute{
		method: http.MethodDelete,
		path:   "/labels/5",
		handler: func(w http.ResponseWriter, req *http.Request) {
			deleteCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		},
	})}
	resp := &resource.DeleteResponse{}

	r.Delete(context.Background(), resource.DeleteRequest{State: labelState(t, labelStateModel{Id: types.Int64Value(5)})}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected delete to succeed, got diagnostics: %v", resp.Diagnostics)
	}
	if !deleteCalled {
		t.Fatal("expected delete endpoint to be called")
	}
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("failed to parse time %q: %v", value, err)
	}
	return parsed
}
