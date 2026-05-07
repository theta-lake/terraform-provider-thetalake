package usergroup

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
	extId := "ext-456"

	plan := &userGroupPlanModel{
		Name:        types.StringValue("Test UG"),
		Description: types.StringValue(desc),
		ExternalId:  types.StringValue(extId),
		UserIds: types.SetValueMust(types.Int64Type, []attr.Value{
			types.Int64Value(201),
			types.Int64Value(202),
		}),
	}

	apiModel := toApiModel(plan)

	assert.Equal(t, "Test UG", apiModel.Name)
	assert.NotNil(t, apiModel.Description)
	assert.Equal(t, desc, *apiModel.Description)
	assert.NotNil(t, apiModel.ExternalId)
	assert.Equal(t, extId, *apiModel.ExternalId)
	assert.ElementsMatch(t, []int64{201, 202}, apiModel.UserIds)
}

func TestToApiModel_NullOptionalFields(t *testing.T) {
	plan := &userGroupPlanModel{
		Name:        types.StringValue("Minimal UG"),
		Description: types.StringNull(),
		ExternalId:  types.StringNull(),
		UserIds:     types.SetNull(types.Int64Type),
	}

	apiModel := toApiModel(plan)

	assert.Equal(t, "Minimal UG", apiModel.Name)
	assert.Nil(t, apiModel.Description)
	assert.Nil(t, apiModel.ExternalId)
	assert.Nil(t, apiModel.UserIds)
}

func TestToApiModel_UnknownOptionalFields(t *testing.T) {
	plan := &userGroupPlanModel{
		Name:        types.StringValue("Unknown UG"),
		Description: types.StringUnknown(),
		ExternalId:  types.StringUnknown(),
		UserIds:     types.SetUnknown(types.Int64Type),
	}

	apiModel := toApiModel(plan)

	assert.Equal(t, "Unknown UG", apiModel.Name)
	assert.Nil(t, apiModel.Description)
	assert.Nil(t, apiModel.ExternalId)
	assert.Nil(t, apiModel.UserIds)
}

func TestFromApiModel(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2024-06-01T09:00:00Z")
	desc := "A description"
	extId := "ext-456"

	apiUg := thetalake.UserGroup{
		Id:          75,
		Name:        "Test UG",
		Description: &desc,
		ExternalId:  &extId,
		CreatedAt:   &ts,
		UpdatedAt:   &ts,
		UserIds:     []int64{201, 202},
	}

	state := fromApiModel(apiUg)

	assert.Equal(t, int64(75), state.Id.ValueInt64())
	assert.Equal(t, "Test UG", state.Name.ValueString())
	assert.Equal(t, desc, state.Description.ValueString())
	assert.Equal(t, extId, state.ExternalId.ValueString())
	assert.Equal(t, ts.Format(time.RFC3339), state.CreatedAt.ValueString())
	assert.Equal(t, ts.Format(time.RFC3339), state.UpdatedAt.ValueString())

	var ids []int64
	state.UserIds.ElementsAs(context.Background(), &ids, false)
	assert.ElementsMatch(t, []int64{201, 202}, ids)
}

func TestFromApiModel_NilOptionalFields(t *testing.T) {
	apiUg := thetalake.UserGroup{
		Id:          2,
		Name:        "Minimal UG",
		Description: nil,
		ExternalId:  nil,
		CreatedAt:   nil,
		UpdatedAt:   nil,
		UserIds:     nil,
	}

	state := fromApiModel(apiUg)

	assert.Equal(t, int64(2), state.Id.ValueInt64())
	assert.True(t, state.Description.IsNull())
	assert.True(t, state.ExternalId.IsNull())
	assert.True(t, state.CreatedAt.IsNull())
	assert.True(t, state.UpdatedAt.IsNull())
	assert.Equal(t, 0, len(state.UserIds.Elements()))
}
