package role

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type roleTestRoute struct {
	method  string
	path    string
	handler http.HandlerFunc
}

func newRoleTestClient(t *testing.T, routes ...roleTestRoute) *thetalake.Client {
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

func roleSchemaForTest(t *testing.T) resourceschema.Schema {
	t.Helper()
	r := &roleResource{}
	resp := &frameworkresource.SchemaResponse{}
	r.Schema(context.Background(), frameworkresource.SchemaRequest{}, resp)
	return resp.Schema
}

func rolePlanForTest(t *testing.T, model roleStateModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: roleSchemaForTest(t)}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to seed plan: %v", diags)
	}
	return plan
}

func TestRoleCreateInvalidPermissionDiagnostic(t *testing.T) {
	createCalled := false
	r := &roleResource{client: newRoleTestClient(t,
		roleTestRoute{
			method: http.MethodGet,
			path:   "/roles/permissions",
			handler: func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"permissions":["cases:read","cases:create"]}`))
			},
		},
		roleTestRoute{
			method: http.MethodPost,
			path:   "/roles",
			handler: func(w http.ResponseWriter, req *http.Request) {
				createCalled = true
				w.WriteHeader(http.StatusCreated)
			},
		},
	)}
	resp := &frameworkresource.CreateResponse{State: tfsdk.State{Schema: roleSchemaForTest(t)}}

	r.Create(context.Background(), frameworkresource.CreateRequest{Plan: rolePlanForTest(t, roleStateModel{
		CreatedAt:     timetypes.NewRFC3339Null(),
		Default:       types.BoolNull(),
		Description:   types.StringValue("Role used to test invalid permissions"),
		Id:            types.Int64Null(),
		IsBuiltIn:     types.BoolNull(),
		Name:          types.StringValue("Invalid Permission Role"),
		NumberOfUsers: types.Int64Null(),
		Permissions: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("cases:read"),
			types.StringValue("cases:not-real"),
		}),
		UpdatedAt: timetypes.NewRFC3339Null(),
	})}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected invalid permissions to add an error diagnostic")
	}
	if resp.Diagnostics.ErrorsCount() != 1 {
		t.Fatalf("expected 1 error diagnostic, got %d", resp.Diagnostics.ErrorsCount())
	}
	if createCalled {
		t.Fatal("expected create request to be blocked before POST /roles")
	}

	errDiag := resp.Diagnostics.Errors()[0]
	if errDiag.Summary() != "Invalid role permissions" {
		t.Fatalf("expected invalid permission summary, got %q", errDiag.Summary())
	}
	if !strings.Contains(errDiag.Detail(), "cases:not-real") {
		t.Fatalf("expected invalid permission detail to mention the bad permission, got %q", errDiag.Detail())
	}
}
