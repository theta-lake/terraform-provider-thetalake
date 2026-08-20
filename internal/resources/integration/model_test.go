package integration

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestToApiModel_GenericJournaling(t *testing.T) {
	plan := integrationPlanModel{
		Name:                 types.StringValue("Custom Generic Journaling Integration"),
		Paused:               types.BoolValue(false),
		GoogleWorkspaceEmail: types.ObjectNull(googleWorkspaceEmailAttrTypes),
		ThetaLakeApi:         types.ObjectNull(thetaLakeApiAttrTypes),
		GenericJournaling: types.ObjectValueMust(genericJournalingAttrTypes, map[string]attr.Value{
			"download_o365_onedrive_links": types.BoolValue(true),
			"download_salesforce_doclinks": types.BoolValue(false),
			"index_headers":                types.StringValue("X-Header-Score,X-Routed-Via"),
			"sender_spf_override":          types.StringNull(),
			"undeliverable_disabled":       types.BoolValue(true),
			"undeliverable_email_address":  types.StringValue("undeliverable@example.com"),
			"undeliverable_email_password": types.StringValue("SecretPassword"),
			"undeliverable_email_port":     types.Int64Value(993),
			"undeliverable_email_server":   types.StringValue("email.realbank.com"),
			"undeliverable_email_user":     types.StringValue("Undeliverable User"),
		}),
	}

	apiModel, diags := toApiModel(context.Background(), &plan)
	if diags.HasError() {
		t.Fatalf("expected no diagnostics, got: %v", diags)
	}

	if apiModel.Type != thetalake.IntegrationTypeGenericJournaling {
		t.Fatalf("expected type %q, got %q", thetalake.IntegrationTypeGenericJournaling, apiModel.Type)
	}
	if apiModel.Options == nil {
		t.Fatal("expected options to be populated")
	}
	if apiModel.Options.IndexHeaders == nil || *apiModel.Options.IndexHeaders != "X-Header-Score,X-Routed-Via" {
		t.Fatalf("expected index_headers to round-trip, got %v", apiModel.Options.IndexHeaders)
	}
	if apiModel.Options.UndeliverableEmailPassword == nil || *apiModel.Options.UndeliverableEmailPassword != "SecretPassword" {
		t.Fatalf("expected undeliverable_email_password to round-trip, got %v", apiModel.Options.UndeliverableEmailPassword)
	}
	if apiModel.Options.SenderSpfOverride != nil {
		t.Fatalf("expected sender_spf_override to remain unset, got %v", *apiModel.Options.SenderSpfOverride)
	}
	if apiModel.Paused == nil || *apiModel.Paused != false {
		t.Fatalf("expected paused to be false, got %v", apiModel.Paused)
	}
}

func TestToApiModel_GoogleWorkspaceEmail(t *testing.T) {
	plan := integrationPlanModel{
		Name:              types.StringValue("Custom Google Workspace Email Integration"),
		Paused:            types.BoolValue(false),
		GenericJournaling: types.ObjectNull(genericJournalingAttrTypes),
		ThetaLakeApi:      types.ObjectNull(thetaLakeApiAttrTypes),
		GoogleWorkspaceEmail: types.ObjectValueMust(googleWorkspaceEmailAttrTypes, map[string]attr.Value{
			"download_o365_onedrive_links": types.BoolValue(true),
			"index_headers":                types.StringValue("X-Header-Score,X-Routed-Via"),
			"sender_spf_override":          types.StringValue("v=spf1 ip4:127.0.0.1/32 -all"),
			"undeliverable_email_address":  types.StringValue("undeliverable@example.com"),
		}),
	}

	apiModel, diags := toApiModel(context.Background(), &plan)
	if diags.HasError() {
		t.Fatalf("expected no diagnostics, got: %v", diags)
	}

	if apiModel.Type != thetalake.IntegrationTypeGoogleWorkspaceEmail {
		t.Fatalf("expected type %q, got %q", thetalake.IntegrationTypeGoogleWorkspaceEmail, apiModel.Type)
	}
	if apiModel.Options == nil || apiModel.Options.UndeliverableEmailAddress == nil {
		t.Fatal("expected undeliverable_email_address to be populated")
	}
	// Fields not present in this variant's schema must never leak in.
	if apiModel.Options.UndeliverableEmailPassword != nil {
		t.Fatal("expected undeliverable_email_password to be absent from google_workspace_email options")
	}
}

func TestToApiModel_ThetaLakeApi(t *testing.T) {
	plan := integrationPlanModel{
		Name:                 types.StringValue("Custom Theta Lake API Integration"),
		Paused:               types.BoolValue(false),
		GenericJournaling:    types.ObjectNull(genericJournalingAttrTypes),
		GoogleWorkspaceEmail: types.ObjectNull(googleWorkspaceEmailAttrTypes),
		ThetaLakeApi:         types.ObjectValueMust(thetaLakeApiAttrTypes, map[string]attr.Value{}),
	}

	apiModel, diags := toApiModel(context.Background(), &plan)
	if diags.HasError() {
		t.Fatalf("expected no diagnostics, got: %v", diags)
	}

	if apiModel.Type != thetalake.IntegrationTypeThetaLakeApi {
		t.Fatalf("expected type %q, got %q", thetalake.IntegrationTypeThetaLakeApi, apiModel.Type)
	}
	if apiModel.Options == nil {
		t.Fatal("expected an empty options object, got nil")
	}
	if apiModel.Options.IndexHeaders != nil {
		t.Fatal("expected theta_lake_api options to have no fields set")
	}
}

func TestToApiModel_NoBlockSet(t *testing.T) {
	plan := integrationPlanModel{
		Name:                 types.StringValue("Missing Type"),
		Paused:               types.BoolValue(false),
		GenericJournaling:    types.ObjectNull(genericJournalingAttrTypes),
		GoogleWorkspaceEmail: types.ObjectNull(googleWorkspaceEmailAttrTypes),
		ThetaLakeApi:         types.ObjectNull(thetaLakeApiAttrTypes),
	}

	_, diags := toApiModel(context.Background(), &plan)
	if !diags.HasError() {
		t.Fatal("expected an error diagnostic when no type block is set")
	}
}

func TestFromApiModel_GenericJournaling(t *testing.T) {
	status := "Queued"
	apiModel := thetalake.Integration{
		Id:                302,
		IntegrationGroup:  "Collaboration Platform",
		IntegrationType:   "Generic Journaling",
		IntegrationTypeId: 41,
		Name:              "Custom Generic Journaling Integration",
		ServicePaused:     false,
		Status:            &status,
	}
	options := &thetalake.IntegrationOptions{
		UndeliverableEmailServer:   new("email.realbank.com"),
		UndeliverableEmailPort:     new(int64(993)),
		UndeliverableEmailPassword: new("SecretPassword"),
	}

	state := fromApiModel(apiModel, options, thetalake.IntegrationTypeGenericJournaling)

	if state.Type.ValueString() != thetalake.IntegrationTypeGenericJournaling {
		t.Fatalf("expected derived type %q, got %q", thetalake.IntegrationTypeGenericJournaling, state.Type.ValueString())
	}
	if state.GenericJournaling.IsNull() {
		t.Fatal("expected generic_journaling to be populated")
	}
	if !state.GoogleWorkspaceEmail.IsNull() || !state.ThetaLakeApi.IsNull() {
		t.Fatal("expected the other two type blocks to remain null")
	}

	var decoded genericJournalingModel
	if diags := state.GenericJournaling.As(context.Background(), &decoded, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("failed to decode generic_journaling: %v", diags)
	}
	if decoded.UndeliverableEmailServer.ValueString() != "email.realbank.com" {
		t.Fatalf("expected undeliverable_email_server to round-trip, got %q", decoded.UndeliverableEmailServer.ValueString())
	}
	// The API does return the password, and fromApiModel must surface it so that Read
	// can report an out-of-band change as drift.
	if decoded.UndeliverableEmailPassword.ValueString() != "SecretPassword" {
		t.Fatalf("expected undeliverable_email_password from the API, got %q", decoded.UndeliverableEmailPassword.ValueString())
	}
}

func TestFromApiModel_ThetaLakeApi(t *testing.T) {
	apiModel := thetalake.Integration{
		Id:                305,
		IntegrationType:   "Theta Lake API",
		IntegrationTypeId: 80,
		Name:              "Custom Theta Lake API Integration",
	}

	state := fromApiModel(apiModel, &thetalake.IntegrationOptions{}, thetalake.IntegrationTypeThetaLakeApi)

	if state.Type.ValueString() != thetalake.IntegrationTypeThetaLakeApi {
		t.Fatalf("expected derived type %q, got %q", thetalake.IntegrationTypeThetaLakeApi, state.Type.ValueString())
	}
	if state.ThetaLakeApi.IsNull() {
		t.Fatal("expected theta_lake_api to be populated with an empty object")
	}
}

func TestFromApiModel_ThetaLakeApi_NilOptions(t *testing.T) {
	// theta_lake_api has no configuration options, so the API's create/update responses
	// may omit service_params entirely for this type (options is nil here). The
	// theta_lake_api block must still be populated with an empty object so it matches
	// the plan's cty.EmptyObjectVal instead of producing a null/non-null consistency error.
	apiModel := thetalake.Integration{
		Id:                305,
		IntegrationType:   "Theta Lake API",
		IntegrationTypeId: 80,
		Name:              "Custom Theta Lake API Integration",
	}

	state := fromApiModel(apiModel, nil, thetalake.IntegrationTypeThetaLakeApi)

	if state.Type.ValueString() != thetalake.IntegrationTypeThetaLakeApi {
		t.Fatalf("expected derived type %q, got %q", thetalake.IntegrationTypeThetaLakeApi, state.Type.ValueString())
	}
	if state.ThetaLakeApi.IsNull() {
		t.Fatal("expected theta_lake_api to be populated with an empty object even when options is nil")
	}
}

func TestFromApiModel_UnresolvableType(t *testing.T) {
	apiModel := thetalake.Integration{
		Id:                999,
		IntegrationType:   "Zoom",
		IntegrationTypeId: 1,
		Name:              "Unsupported",
	}

	typeSlug := thetalake.IntegrationTypeSlug(apiModel.IntegrationTypeId, apiModel.IntegrationType)
	state := fromApiModel(apiModel, nil, typeSlug)

	if state.Type.ValueString() != "" {
		t.Fatalf("expected empty type slug for an unrecognized integration type, got %q", state.Type.ValueString())
	}
	if !state.GenericJournaling.IsNull() || !state.GoogleWorkspaceEmail.IsNull() || !state.ThetaLakeApi.IsNull() {
		t.Fatal("expected all type blocks to remain null when the type can't be resolved")
	}
}

func TestPreserveImmutableComputedFieldsFromState(t *testing.T) {
	priorState := &integrationStateModel{
		CreatedAt:         timetypes.NewRFC3339Null(),
		Id:                types.Int64Value(305),
		IntegrationGroup:  types.StringValue("Collaboration Platform"),
		IntegrationType:   types.StringValue("Theta Lake API"),
		IntegrationTypeId: types.Int64Value(80),
	}

	// Simulates a PUT response that omitted created_at/integration_group/
	// integration_type/integration_type_id/id: fromApiModel would have decoded these
	// to their zero values.
	updated := &integrationStateModel{
		CreatedAt:         timetypes.NewRFC3339Null(),
		Id:                types.Int64Value(0),
		IntegrationGroup:  types.StringValue(""),
		IntegrationType:   types.StringValue(""),
		IntegrationTypeId: types.Int64Value(0),
	}

	preserveImmutableComputedFieldsFromState(updated, priorState)

	if updated.Id.ValueInt64() != 305 {
		t.Fatalf("expected id to be carried over from prior state, got %d", updated.Id.ValueInt64())
	}
	if updated.IntegrationGroup.ValueString() != "Collaboration Platform" {
		t.Fatalf("expected integration_group to be carried over from prior state, got %q", updated.IntegrationGroup.ValueString())
	}
	if updated.IntegrationType.ValueString() != "Theta Lake API" {
		t.Fatalf("expected integration_type to be carried over from prior state, got %q", updated.IntegrationType.ValueString())
	}
	if updated.IntegrationTypeId.ValueInt64() != 80 {
		t.Fatalf("expected integration_type_id to be carried over from prior state, got %d", updated.IntegrationTypeId.ValueInt64())
	}
}
