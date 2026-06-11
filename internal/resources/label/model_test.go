package label

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestToApiModel(t *testing.T) {
	plan := &labelPlanModel{
		BackgroundColor: types.StringValue("#FFC906"),
		Hidden:          types.BoolValue(true),
		LongName:        types.StringValue("Label description"),
		ShortName:       types.StringValue("Label"),
	}

	apiModel := toApiModel(plan)

	if apiModel.BackgroundColor != "#FFC906" {
		t.Fatalf("expected background color to be propagated, got %q", apiModel.BackgroundColor)
	}
	if !apiModel.Hidden {
		t.Fatal("expected hidden to be propagated")
	}
	if apiModel.LongName != "Label description" {
		t.Fatalf("expected long_name to be propagated, got %q", apiModel.LongName)
	}
	if apiModel.ShortName != "Label" {
		t.Fatalf("expected short_name to be propagated, got %q", apiModel.ShortName)
	}
}

func TestFromApiModelWithUpdatedAt(t *testing.T) {
	createdAt := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := createdAt.Add(2 * time.Hour)

	state := fromApiModel(thetalake.Label{
		BackgroundColor:   "#FFC906",
		CreatedAt:         createdAt,
		Hidden:            true,
		Id:                5,
		LongName:          "Label description",
		OrgUnitId:         108,
		ShortName:         "Label",
		TaggedDatumsCount: 7,
		UpdatedAt:         &updatedAt,
		UserId:            422,
	})

	if state.BackgroundColor.ValueString() != "#FFC906" {
		t.Fatalf("expected background color to map into state, got %q", state.BackgroundColor.ValueString())
	}
	if !state.Hidden.ValueBool() {
		t.Fatal("expected hidden to map into state")
	}
	if state.Id.ValueInt64() != 5 {
		t.Fatalf("expected id 5, got %d", state.Id.ValueInt64())
	}
	if state.ShortName.ValueString() != "Label" {
		t.Fatalf("expected short_name to map into state, got %q", state.ShortName.ValueString())
	}
	if state.UpdatedAt.IsNull() {
		t.Fatal("expected updated_at to be set")
	}
	if got, diags := state.CreatedAt.ValueRFC3339Time(); diags.HasError() || !got.Equal(createdAt) {
		t.Fatalf("expected created_at %s, got %s", createdAt, got)
	}
	if got, diags := state.UpdatedAt.ValueRFC3339Time(); diags.HasError() || !got.Equal(updatedAt) {
		t.Fatalf("expected updated_at %s, got %s", updatedAt, got)
	}
}

func TestFromApiModelWithoutUpdatedAt(t *testing.T) {
	state := fromApiModel(thetalake.Label{
		BackgroundColor: "#FFC906",
		CreatedAt:       time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC),
		Id:              5,
		LongName:        "Label description",
		ShortName:       "Label",
	})

	if !state.UpdatedAt.IsNull() {
		t.Fatal("expected updated_at to be null when API omits it")
	}
	if state.CreatedAt.Equal(timetypes.NewRFC3339Null()) {
		t.Fatal("expected created_at to be populated")
	}
}
