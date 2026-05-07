package directorygroup

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestToApiModel(t *testing.T) {
	desc := "A description"
	extId := "ext-123"

	plan := &directoryGroupPlanModel{
		Name:        types.StringValue("Test DG"),
		Description: types.StringValue(desc),
		ExternalId:  types.StringValue(extId),
		IdentityIds: types.SetValueMust(types.Int64Type, []attr.Value{
			types.Int64Value(101),
			types.Int64Value(102),
		}),
	}

	apiModel := toApiModel(plan)

	assert.Equal(t, "Test DG", apiModel.Name)
	assert.NotNil(t, apiModel.Description)
	assert.Equal(t, desc, *apiModel.Description)
	assert.NotNil(t, apiModel.ExternalId)
	assert.Equal(t, extId, *apiModel.ExternalId)
	assert.ElementsMatch(t, []int64{101, 102}, apiModel.IdentityIds)
}

func TestToApiModel_NullOptionalFields(t *testing.T) {
	plan := &directoryGroupPlanModel{
		Name:        types.StringValue("Minimal DG"),
		Description: types.StringNull(),
		ExternalId:  types.StringNull(),
		IdentityIds: types.SetNull(types.Int64Type),
	}

	apiModel := toApiModel(plan)

	assert.Equal(t, "Minimal DG", apiModel.Name)
	assert.Nil(t, apiModel.Description)
	assert.Nil(t, apiModel.ExternalId)
	assert.Nil(t, apiModel.IdentityIds)
}

func TestToApiModel_UnknownOptionalFields(t *testing.T) {
	plan := &directoryGroupPlanModel{
		Name:        types.StringValue("Unknown DG"),
		Description: types.StringUnknown(),
		ExternalId:  types.StringUnknown(),
		IdentityIds: types.SetUnknown(types.Int64Type),
	}

	apiModel := toApiModel(plan)

	assert.Equal(t, "Unknown DG", apiModel.Name)
	assert.Nil(t, apiModel.Description)
	assert.Nil(t, apiModel.ExternalId)
	assert.Nil(t, apiModel.IdentityIds)
}

func TestFromApiModel(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	desc := "A description"
	extId := "ext-123"

	apiDg := thetalake.DirectoryGroup{
		Id:          1996,
		Name:        "Test DG",
		Description: &desc,
		ExternalId:  &extId,
		CreatedAt:   &ts,
		UpdatedAt:   &ts,
		IdentityIds: []int64{101, 102},
	}

	state := fromApiModel(apiDg)

	assert.Equal(t, int64(1996), state.Id.ValueInt64())
	assert.Equal(t, "Test DG", state.Name.ValueString())
	assert.Equal(t, desc, state.Description.ValueString())
	assert.Equal(t, extId, state.ExternalId.ValueString())
	assert.Equal(t, ts.Format(time.RFC3339), state.CreatedAt.ValueString())
	assert.Equal(t, ts.Format(time.RFC3339), state.UpdatedAt.ValueString())

	var ids []int64
	state.IdentityIds.ElementsAs(context.Background(), &ids, false)
	assert.ElementsMatch(t, []int64{101, 102}, ids)
}

func TestFromApiModel_NilOptionalFields(t *testing.T) {
	apiDg := thetalake.DirectoryGroup{
		Id:          1,
		Name:        "Minimal DG",
		Description: nil,
		ExternalId:  nil,
		CreatedAt:   nil,
		UpdatedAt:   nil,
		IdentityIds: nil,
	}

	state := fromApiModel(apiDg)

	assert.Equal(t, int64(1), state.Id.ValueInt64())
	assert.True(t, state.Description.IsNull())
	assert.True(t, state.ExternalId.IsNull())
	assert.True(t, state.CreatedAt.IsNull())
	assert.True(t, state.UpdatedAt.IsNull())
	assert.Equal(t, 0, len(state.IdentityIds.Elements()))
}
