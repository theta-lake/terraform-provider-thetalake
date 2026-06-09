package retentionlibrary

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type retentionLibraryPlanModel struct {
	Description                types.String `tfsdk:"description"`
	ExternalId                 types.String `tfsdk:"external_id"`
	Name                       types.String `tfsdk:"name"`
	RetainInReview             types.Bool   `tfsdk:"retain_in_review"`
	RetentionPeriodDays        types.Int64  `tfsdk:"retention_period_days"`
	RetentionPeriodEnabled     types.Bool   `tfsdk:"retention_period_enabled"`
	SecCompliantStorageEnabled types.Bool   `tfsdk:"sec_compliant_storage_enabled"`
	StorageAccountId           types.Int64  `tfsdk:"storage_account_id"`
}

type retentionLibraryStateModel struct {
	CreatedAt                    timetypes.RFC3339 `tfsdk:"created_at"`
	DatumCount                   types.Int64       `tfsdk:"datum_count"`
	DatumSize                    types.Int64       `tfsdk:"datum_size"`
	DeleteOnExpiration           types.Bool        `tfsdk:"delete_on_expiration"`
	Description                  types.String      `tfsdk:"description"`
	DisplayName                  types.String      `tfsdk:"display_name"`
	ExternalId                   types.String      `tfsdk:"external_id"`
	Id                           types.Int64       `tfsdk:"id"`
	LegalHoldCount               types.Int64       `tfsdk:"legal_hold_count"`
	Name                         types.String      `tfsdk:"name"`
	RetainInReview               types.Bool        `tfsdk:"retain_in_review"`
	RetentionPeriodDays          types.Int64       `tfsdk:"retention_period_days"`
	RetentionPeriodEnabled       types.Bool        `tfsdk:"retention_period_enabled"`
	RetentionSummaryText         types.String      `tfsdk:"retention_summary_text"`
	SecCompliantStorageConfirmed types.Bool        `tfsdk:"sec_compliant_storage_confirmed"`
	SecCompliantStorageEnabled   types.Bool        `tfsdk:"sec_compliant_storage_enabled"`
	StorageAccountId             types.Int64       `tfsdk:"storage_account_id"`
	SwrvRuleCount                types.Int64       `tfsdk:"swrv_rule_count"`
	UpdatedAt                    timetypes.RFC3339 `tfsdk:"updated_at"`
}

func toApiModel(plan *retentionLibraryPlanModel) thetalake.RetentionLibrary {
	library := thetalake.RetentionLibrary{
		Description:                plan.Description.ValueString(),
		Name:                       plan.Name.ValueString(),
		RetainInReview:             plan.RetainInReview.ValueBool(),
		RetentionPeriodDays:        plan.RetentionPeriodDays.ValueInt64(),
		RetentionPeriodEnabled:     plan.RetentionPeriodEnabled.ValueBool(),
		SecCompliantStorageEnabled: plan.SecCompliantStorageEnabled.ValueBool(),
		StorageAccountId:           plan.StorageAccountId.ValueInt64(),
	}

	if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
		externalId := plan.ExternalId.ValueString()
		library.ExternalId = &externalId
	}

	return library
}

func fromApiModel(library thetalake.RetentionLibrary) retentionLibraryStateModel {
	state := retentionLibraryStateModel{
		CreatedAt:                    timetypes.NewRFC3339TimeValue(library.CreatedAt),
		DatumCount:                   types.Int64Value(library.DatumCount),
		DatumSize:                    types.Int64Value(library.DatumSize),
		DeleteOnExpiration:           types.BoolValue(library.DeleteOnExpiration),
		Description:                  types.StringValue(library.Description),
		DisplayName:                  types.StringValue(library.DisplayName),
		Id:                           types.Int64Value(library.Id),
		LegalHoldCount:               types.Int64Value(library.LegalHoldCount),
		Name:                         types.StringValue(library.Name),
		RetainInReview:               types.BoolValue(library.RetainInReview),
		RetentionPeriodDays:          types.Int64Value(library.RetentionPeriodDays),
		RetentionPeriodEnabled:       types.BoolValue(library.RetentionPeriodEnabled),
		RetentionSummaryText:         types.StringValue(library.RetentionSummaryText),
		SecCompliantStorageConfirmed: types.BoolValue(library.SecCompliantStorageConfirmed),
		SecCompliantStorageEnabled:   types.BoolValue(library.SecCompliantStorageEnabled),
		StorageAccountId:             types.Int64Value(library.StorageAccountId),
		SwrvRuleCount:                types.Int64Value(library.SwrvRuleCount),
		UpdatedAt:                    timetypes.NewRFC3339TimeValue(library.UpdatedAt),
	}

	if library.ExternalId == nil {
		state.ExternalId = types.StringNull()
	} else {
		state.ExternalId = types.StringValue(*library.ExternalId)
	}

	return state
}
