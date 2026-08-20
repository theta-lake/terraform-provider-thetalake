package integration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

var genericJournalingAttrTypes = map[string]attr.Type{
	"download_o365_onedrive_links": types.BoolType,
	"download_salesforce_doclinks": types.BoolType,
	"index_headers":                types.StringType,
	"sender_spf_override":          types.StringType,
	"undeliverable_disabled":       types.BoolType,
	"undeliverable_email_address":  types.StringType,
	"undeliverable_email_password": types.StringType,
	"undeliverable_email_port":     types.Int64Type,
	"undeliverable_email_server":   types.StringType,
	"undeliverable_email_user":     types.StringType,
}

var googleWorkspaceEmailAttrTypes = map[string]attr.Type{
	"download_o365_onedrive_links": types.BoolType,
	"index_headers":                types.StringType,
	"sender_spf_override":          types.StringType,
	"undeliverable_email_address":  types.StringType,
}

// thetaLakeApiAttrTypes is empty: the Theta Lake API integration type has no
// configuration options in the spec.
var thetaLakeApiAttrTypes = map[string]attr.Type{}

type genericJournalingModel struct {
	DownloadO365OnedriveLinks  types.Bool   `tfsdk:"download_o365_onedrive_links"`
	DownloadSalesforceDoclinks types.Bool   `tfsdk:"download_salesforce_doclinks"`
	IndexHeaders               types.String `tfsdk:"index_headers"`
	SenderSpfOverride          types.String `tfsdk:"sender_spf_override"`
	UndeliverableDisabled      types.Bool   `tfsdk:"undeliverable_disabled"`
	UndeliverableEmailAddress  types.String `tfsdk:"undeliverable_email_address"`
	UndeliverableEmailPassword types.String `tfsdk:"undeliverable_email_password"`
	UndeliverableEmailPort     types.Int64  `tfsdk:"undeliverable_email_port"`
	UndeliverableEmailServer   types.String `tfsdk:"undeliverable_email_server"`
	UndeliverableEmailUser     types.String `tfsdk:"undeliverable_email_user"`
}

type googleWorkspaceEmailModel struct {
	DownloadO365OnedriveLinks types.Bool   `tfsdk:"download_o365_onedrive_links"`
	IndexHeaders              types.String `tfsdk:"index_headers"`
	SenderSpfOverride         types.String `tfsdk:"sender_spf_override"`
	UndeliverableEmailAddress types.String `tfsdk:"undeliverable_email_address"`
}

type integrationPlanModel struct {
	GenericJournaling    types.Object `tfsdk:"generic_journaling"`
	GoogleWorkspaceEmail types.Object `tfsdk:"google_workspace_email"`
	Name                 types.String `tfsdk:"name"`
	Paused               types.Bool   `tfsdk:"paused"`
	ThetaLakeApi         types.Object `tfsdk:"theta_lake_api"`
}

type integrationStateModel struct {
	CreatedAt            timetypes.RFC3339 `tfsdk:"created_at"`
	GenericJournaling    types.Object      `tfsdk:"generic_journaling"`
	GoogleWorkspaceEmail types.Object      `tfsdk:"google_workspace_email"`
	Id                   types.Int64       `tfsdk:"id"`
	IntegrationGroup     types.String      `tfsdk:"integration_group"`
	IntegrationType      types.String      `tfsdk:"integration_type"`
	IntegrationTypeId    types.Int64       `tfsdk:"integration_type_id"`
	Name                 types.String      `tfsdk:"name"`
	Paused               types.Bool        `tfsdk:"paused"`
	ThetaLakeApi         types.Object      `tfsdk:"theta_lake_api"`
	Type                 types.String      `tfsdk:"type"`
}

// toApiModel builds a thetalake.Integration from whichever of the three type blocks is
// set in the plan. The ExactlyOneOf schema validator guarantees exactly one is set by
// the time this runs against real plan data.
func toApiModel(ctx context.Context, plan *integrationPlanModel) (thetalake.Integration, diag.Diagnostics) {
	var diagnostics diag.Diagnostics

	apiModel := thetalake.Integration{
		Name: plan.Name.ValueString(),
	}

	paused := plan.Paused.ValueBool()
	apiModel.Paused = &paused

	switch {
	case !plan.GenericJournaling.IsNull() && !plan.GenericJournaling.IsUnknown():
		var options genericJournalingModel
		diagnostics.Append(plan.GenericJournaling.As(ctx, &options, basetypes.ObjectAsOptions{})...)
		if diagnostics.HasError() {
			return thetalake.Integration{}, diagnostics
		}

		apiModel.Type = thetalake.IntegrationTypeGenericJournaling
		apiModel.Options = genericJournalingToApiOptions(options)
	case !plan.GoogleWorkspaceEmail.IsNull() && !plan.GoogleWorkspaceEmail.IsUnknown():
		var options googleWorkspaceEmailModel
		diagnostics.Append(plan.GoogleWorkspaceEmail.As(ctx, &options, basetypes.ObjectAsOptions{})...)
		if diagnostics.HasError() {
			return thetalake.Integration{}, diagnostics
		}

		apiModel.Type = thetalake.IntegrationTypeGoogleWorkspaceEmail
		apiModel.Options = googleWorkspaceEmailToApiOptions(options)
	case !plan.ThetaLakeApi.IsNull() && !plan.ThetaLakeApi.IsUnknown():
		apiModel.Type = thetalake.IntegrationTypeThetaLakeApi
		apiModel.Options = &thetalake.IntegrationOptions{}
	default:
		diagnostics.AddError(
			"Invalid Integration Configuration",
			"Exactly one of generic_journaling, google_workspace_email, or theta_lake_api must be set",
		)
	}

	return apiModel, diagnostics
}

func genericJournalingToApiOptions(plan genericJournalingModel) *thetalake.IntegrationOptions {
	options := &thetalake.IntegrationOptions{}

	if !plan.DownloadO365OnedriveLinks.IsNull() && !plan.DownloadO365OnedriveLinks.IsUnknown() {
		value := plan.DownloadO365OnedriveLinks.ValueBool()
		options.DownloadO365OnedriveLinks = &value
	}

	if !plan.DownloadSalesforceDoclinks.IsNull() && !plan.DownloadSalesforceDoclinks.IsUnknown() {
		value := plan.DownloadSalesforceDoclinks.ValueBool()
		options.DownloadSalesforceDoclinks = &value
	}

	if !plan.IndexHeaders.IsNull() && !plan.IndexHeaders.IsUnknown() {
		value := plan.IndexHeaders.ValueString()
		options.IndexHeaders = &value
	}

	if !plan.SenderSpfOverride.IsNull() && !plan.SenderSpfOverride.IsUnknown() {
		value := plan.SenderSpfOverride.ValueString()
		options.SenderSpfOverride = &value
	}

	if !plan.UndeliverableDisabled.IsNull() && !plan.UndeliverableDisabled.IsUnknown() {
		value := plan.UndeliverableDisabled.ValueBool()
		options.UndeliverableDisabled = &value
	}

	if !plan.UndeliverableEmailAddress.IsNull() && !plan.UndeliverableEmailAddress.IsUnknown() {
		value := plan.UndeliverableEmailAddress.ValueString()
		options.UndeliverableEmailAddress = &value
	}

	if !plan.UndeliverableEmailPassword.IsNull() && !plan.UndeliverableEmailPassword.IsUnknown() {
		value := plan.UndeliverableEmailPassword.ValueString()
		options.UndeliverableEmailPassword = &value
	}

	if !plan.UndeliverableEmailPort.IsNull() && !plan.UndeliverableEmailPort.IsUnknown() {
		value := plan.UndeliverableEmailPort.ValueInt64()
		options.UndeliverableEmailPort = &value
	}

	if !plan.UndeliverableEmailServer.IsNull() && !plan.UndeliverableEmailServer.IsUnknown() {
		value := plan.UndeliverableEmailServer.ValueString()
		options.UndeliverableEmailServer = &value
	}

	if !plan.UndeliverableEmailUser.IsNull() && !plan.UndeliverableEmailUser.IsUnknown() {
		value := plan.UndeliverableEmailUser.ValueString()
		options.UndeliverableEmailUser = &value
	}

	return options
}

func googleWorkspaceEmailToApiOptions(plan googleWorkspaceEmailModel) *thetalake.IntegrationOptions {
	options := &thetalake.IntegrationOptions{}

	if !plan.DownloadO365OnedriveLinks.IsNull() && !plan.DownloadO365OnedriveLinks.IsUnknown() {
		value := plan.DownloadO365OnedriveLinks.ValueBool()
		options.DownloadO365OnedriveLinks = &value
	}

	if !plan.IndexHeaders.IsNull() && !plan.IndexHeaders.IsUnknown() {
		value := plan.IndexHeaders.ValueString()
		options.IndexHeaders = &value
	}

	if !plan.SenderSpfOverride.IsNull() && !plan.SenderSpfOverride.IsUnknown() {
		value := plan.SenderSpfOverride.ValueString()
		options.SenderSpfOverride = &value
	}

	if !plan.UndeliverableEmailAddress.IsNull() && !plan.UndeliverableEmailAddress.IsUnknown() {
		value := plan.UndeliverableEmailAddress.ValueString()
		options.UndeliverableEmailAddress = &value
	}

	return options
}

// fromApiModel maps an API integration (plus, when available, its type-specific
// options) into Terraform state. options should come from the create/update
// response's ServiceParams, or from a GetIntegrationConfiguration call (Create/Update
// fall back to that call when the response omits service_params; Read always uses it
// for types that have configuration options); pass nil only when typeSlug is
// thetalake.IntegrationTypeThetaLakeApi, which has no options to populate.
//
// typeSlug identifies which of the three type blocks to populate. Callers that
// already know the type authoritatively (Create/Update, from the plan-derived
// request) should pass that instead of re-deriving it from the response, since the
// response's integration_type/integration_type_id aren't guaranteed to round-trip
// through thetalake.IntegrationTypeSlug the same request. Read, which only has an
// ID to start from, must derive typeSlug from the response itself.
func fromApiModel(apiModel thetalake.Integration, options *thetalake.IntegrationOptions, typeSlug string) integrationStateModel {
	state := integrationStateModel{
		GenericJournaling:    types.ObjectNull(genericJournalingAttrTypes),
		GoogleWorkspaceEmail: types.ObjectNull(googleWorkspaceEmailAttrTypes),
		Id:                   types.Int64Value(apiModel.Id),
		IntegrationGroup:     types.StringValue(apiModel.IntegrationGroup),
		IntegrationType:      types.StringValue(apiModel.IntegrationType),
		IntegrationTypeId:    types.Int64Value(apiModel.IntegrationTypeId),
		Name:                 types.StringValue(apiModel.Name),
		Paused:               types.BoolValue(apiModel.ServicePaused),
		ThetaLakeApi:         types.ObjectNull(thetaLakeApiAttrTypes),
	}

	if apiModel.CreatedAt != nil {
		state.CreatedAt = timetypes.NewRFC3339TimeValue(*apiModel.CreatedAt)
	} else {
		state.CreatedAt = timetypes.NewRFC3339Null()
	}

	state.Type = types.StringValue(typeSlug)

	// theta_lake_api has no configuration options, so the API's create/update responses
	// may omit service_params entirely for this type. Set the (always-empty) block
	// regardless of whether options is nil, since it never depends on options data.
	if typeSlug == thetalake.IntegrationTypeThetaLakeApi {
		state.ThetaLakeApi = types.ObjectValueMust(thetaLakeApiAttrTypes, map[string]attr.Value{})
	}

	if options == nil {
		return state
	}

	switch typeSlug {
	case thetalake.IntegrationTypeGenericJournaling:
		state.GenericJournaling = genericJournalingFromApiOptions(*options)
	case thetalake.IntegrationTypeGoogleWorkspaceEmail:
		state.GoogleWorkspaceEmail = googleWorkspaceEmailFromApiOptions(*options)
	}

	return state
}

func genericJournalingFromApiOptions(options thetalake.IntegrationOptions) types.Object {
	values := map[string]attr.Value{
		"download_o365_onedrive_links": boolValueOrFalse(options.DownloadO365OnedriveLinks),
		"download_salesforce_doclinks": boolValueOrFalse(options.DownloadSalesforceDoclinks),
		"index_headers":                stringValueOrNull(options.IndexHeaders),
		"sender_spf_override":          stringValueOrNull(options.SenderSpfOverride),
		"undeliverable_disabled":       boolValueOrFalse(options.UndeliverableDisabled),
		"undeliverable_email_address":  stringValueOrNull(options.UndeliverableEmailAddress),
		"undeliverable_email_password": nonEmptyStringValueOrNull(options.UndeliverableEmailPassword),
		"undeliverable_email_port":     int64ValueOrNull(options.UndeliverableEmailPort),
		"undeliverable_email_server":   stringValueOrNull(options.UndeliverableEmailServer),
		"undeliverable_email_user":     stringValueOrNull(options.UndeliverableEmailUser),
	}

	return types.ObjectValueMust(genericJournalingAttrTypes, values)
}

func googleWorkspaceEmailFromApiOptions(options thetalake.IntegrationOptions) types.Object {
	values := map[string]attr.Value{
		"download_o365_onedrive_links": boolValueOrFalse(options.DownloadO365OnedriveLinks),
		"index_headers":                stringValueOrNull(options.IndexHeaders),
		"sender_spf_override":          stringValueOrNull(options.SenderSpfOverride),
		"undeliverable_email_address":  stringValueOrNull(options.UndeliverableEmailAddress),
	}

	return types.ObjectValueMust(googleWorkspaceEmailAttrTypes, values)
}

func boolValueOrFalse(value *bool) types.Bool {
	if value == nil {
		return types.BoolValue(false)
	}

	return types.BoolValue(*value)
}

func stringValueOrNull(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}

	return types.StringValue(*value)
}

// nonEmptyStringValueOrNull is like stringValueOrNull, but also treats an empty
// string as absent. The API returns undeliverable_email_password's configured value
// when one is set, and "" (rather than omitting the field) when none is, so mapping
// that "" through stringValueOrNull would surface "" where the configuration has no
// password at all and show as permanent drift on every refresh.
func nonEmptyStringValueOrNull(value *string) types.String {
	if value == nil || *value == "" {
		return types.StringNull()
	}

	return types.StringValue(*value)
}

func int64ValueOrNull(value *int64) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}

	return types.Int64Value(*value)
}

// preserveImmutableComputedFieldsFromState carries created_at, id, integration_group,
// integration_type, and integration_type_id over from prior state into updated. All
// five are Computed with a UseStateForUnknown plan modifier and are immutable after
// creation, so the plan already carries a known prior value for each; Update must not
// let a PUT response that omits one of them (decoding to its zero value) clobber that
// known value and trigger a "provider produced inconsistent result after apply" error.
func preserveImmutableComputedFieldsFromState(updated *integrationStateModel, priorState *integrationStateModel) {
	updated.CreatedAt = priorState.CreatedAt
	updated.Id = priorState.Id
	updated.IntegrationGroup = priorState.IntegrationGroup
	updated.IntegrationType = priorState.IntegrationType
	updated.IntegrationTypeId = priorState.IntegrationTypeId
}
