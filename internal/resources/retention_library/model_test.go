package retentionlibrary

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestToApiModelWithExternalID(t *testing.T) {
	plan := &retentionLibraryPlanModel{
		Description:                types.StringValue("Retention library description"),
		ExternalId:                 types.StringValue("external-123"),
		Name:                       types.StringValue("Retention Library"),
		RetainInReview:             types.BoolValue(true),
		RetentionPeriodDays:        types.Int64Value(30),
		RetentionPeriodEnabled:     types.BoolValue(true),
		SecCompliantStorageEnabled: types.BoolValue(false),
		StorageAccountId:           types.Int64Value(7),
	}

	apiModel := toApiModel(plan)

	if apiModel.Description != "Retention library description" {
		t.Fatalf("expected description to map, got %q", apiModel.Description)
	}
	if apiModel.ExternalId == nil || *apiModel.ExternalId != "external-123" {
		t.Fatal("expected external_id to map into API model")
	}
	if !apiModel.RetainInReview {
		t.Fatal("expected retain_in_review to map into API model")
	}
	if apiModel.RetentionPeriodDays != 30 {
		t.Fatalf("expected retention_period_days 30, got %d", apiModel.RetentionPeriodDays)
	}
	if apiModel.StorageAccountId != 7 {
		t.Fatalf("expected storage_account_id 7, got %d", apiModel.StorageAccountId)
	}
}

func TestToApiModelWithoutExternalID(t *testing.T) {
	plan := &retentionLibraryPlanModel{
		Description:                types.StringValue("Retention library description"),
		ExternalId:                 types.StringNull(),
		Name:                       types.StringValue("Retention Library"),
		RetainInReview:             types.BoolValue(false),
		RetentionPeriodDays:        types.Int64Value(0),
		RetentionPeriodEnabled:     types.BoolValue(false),
		SecCompliantStorageEnabled: types.BoolValue(false),
		StorageAccountId:           types.Int64Value(7),
	}

	apiModel := toApiModel(plan)

	if apiModel.ExternalId != nil {
		t.Fatal("expected external_id to stay nil when plan value is null")
	}
}

func TestFromApiModelWithAndWithoutExternalID(t *testing.T) {
	externalID := "external-123"
	createdAt := time.Date(2024, time.February, 3, 4, 5, 6, 0, time.UTC)
	updatedAt := createdAt.Add(3 * time.Hour)

	stateWithExternalID := fromApiModel(thetalake.RetentionLibrary{
		CreatedAt:                    createdAt,
		DatumCount:                   12,
		DatumSize:                    34,
		DeleteOnExpiration:           true,
		Description:                  "Retention library description",
		DisplayName:                  "Retention Library (us-east-1)",
		ExternalId:                   &externalID,
		Id:                           477,
		LegalHoldCount:               5,
		Name:                         "Retention Library",
		RetainInReview:               true,
		RetentionPeriodDays:          30,
		RetentionPeriodEnabled:       true,
		RetentionSummaryText:         "Retained for 30 days",
		SecCompliantStorageConfirmed: true,
		SecCompliantStorageEnabled:   false,
		StorageAccountId:             7,
		SwrvRuleCount:                2,
		UpdatedAt:                    updatedAt,
	})

	if stateWithExternalID.ExternalId.IsNull() {
		t.Fatal("expected external_id to be populated")
	}
	if stateWithExternalID.ExternalId.ValueString() != externalID {
		t.Fatalf("expected external_id %q, got %q", externalID, stateWithExternalID.ExternalId.ValueString())
	}
	if stateWithExternalID.Id.ValueInt64() != 477 {
		t.Fatalf("expected id 477, got %d", stateWithExternalID.Id.ValueInt64())
	}
	if !stateWithExternalID.DeleteOnExpiration.ValueBool() {
		t.Fatal("expected delete_on_expiration to map into state")
	}
	if got, diags := stateWithExternalID.UpdatedAt.ValueRFC3339Time(); diags.HasError() || !got.Equal(updatedAt) {
		t.Fatalf("expected updated_at %s, got %s", updatedAt, got)
	}

	stateWithoutExternalID := fromApiModel(thetalake.RetentionLibrary{
		CreatedAt: createdAt,
		Id:        477,
		Name:      "Retention Library",
		UpdatedAt: updatedAt,
	})

	if !stateWithoutExternalID.ExternalId.IsNull() {
		t.Fatal("expected external_id to be null when API omits it")
	}
}
