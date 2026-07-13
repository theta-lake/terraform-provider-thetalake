package legalcase

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestToApiModel(t *testing.T) {
	openDate, _ := time.Parse(time.RFC3339, "2024-01-15T10:00:00Z")

	plan := &casePlanModel{
		Name:        types.StringValue("Test Case"),
		Number:      types.StringValue("CASE-001"),
		Description: types.StringValue("A description"),
		Visibility:  types.StringValue("PRIVATE"),
		OpenDate:    timetypes.NewRFC3339TimeValue(openDate),
		ManagerIds: types.SetValueMust(types.Int64Type, []attr.Value{
			types.Int64Value(101),
			types.Int64Value(102),
		}),
	}

	apiModel := toApiModel(plan)

	assert.Equal(t, "Test Case", apiModel.Name)
	assert.Equal(t, "CASE-001", apiModel.Number)
	assert.Equal(t, "A description", apiModel.Description)
	assert.Equal(t, "PRIVATE", apiModel.Visibility)
	assert.True(t, openDate.Equal(apiModel.OpenDate))
	assert.ElementsMatch(t, []int64{101, 102}, apiModel.ManagerIds)
}

func TestToApiModel_NullOptionalFields(t *testing.T) {
	plan := &casePlanModel{
		Name:        types.StringValue("Minimal Case"),
		Number:      types.StringValue("CASE-002"),
		Visibility:  types.StringValue("PUBLIC"),
		Description: types.StringNull(),
		OpenDate:    timetypes.NewRFC3339Null(),
		ManagerIds:  types.SetNull(types.Int64Type),
	}

	apiModel := toApiModel(plan)

	assert.Equal(t, "Minimal Case", apiModel.Name)
	assert.Equal(t, "", apiModel.Description)
	assert.True(t, apiModel.OpenDate.IsZero())
	assert.Nil(t, apiModel.ManagerIds)
}

func TestFromApiModel(t *testing.T) {
	openDate, _ := time.Parse(time.RFC3339, "2024-01-15T10:00:00Z")
	closeDate, _ := time.Parse(time.RFC3339, "2024-02-01T10:00:00Z")
	createdAt, _ := time.Parse(time.RFC3339, "2024-01-15T10:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2024-02-01T10:00:00Z")

	apiCase := thetalake.Case{
		Id:           628,
		Name:         "Test Case",
		Number:       "CASE-001",
		Description:  "A description",
		Status:       "CLOSED",
		Visibility:   "PRIVATE",
		RecordsCount: 3,
		OpenDate:     openDate,
		CloseDate:    &closeDate,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		ManagerIds:   []int64{101, 102},
	}

	state := fromApiModel(apiCase)

	assert.Equal(t, int64(628), state.Id.ValueInt64())
	assert.Equal(t, "Test Case", state.Name.ValueString())
	assert.Equal(t, "CASE-001", state.Number.ValueString())
	assert.Equal(t, "A description", state.Description.ValueString())
	assert.Equal(t, "CLOSED", state.Status.ValueString())
	assert.Equal(t, "PRIVATE", state.Visibility.ValueString())
	assert.Equal(t, int64(3), state.RecordsCount.ValueInt64())
	assert.Equal(t, openDate.Format(time.RFC3339), state.OpenDate.ValueString())
	assert.Equal(t, closeDate.Format(time.RFC3339), state.CloseDate.ValueString())
	assert.Equal(t, createdAt.Format(time.RFC3339), state.CreatedAt.ValueString())
	assert.Equal(t, updatedAt.Format(time.RFC3339), state.UpdatedAt.ValueString())

	var ids []int64
	state.ManagerIds.ElementsAs(context.Background(), &ids, false)
	assert.ElementsMatch(t, []int64{101, 102}, ids)
}

func TestFromApiModel_NilOptionalFields(t *testing.T) {
	apiCase := thetalake.Case{
		Id:         1,
		Name:       "Minimal Case",
		CloseDate:  nil,
		ManagerIds: nil,
	}

	state := fromApiModel(apiCase)

	assert.Equal(t, int64(1), state.Id.ValueInt64())
	assert.True(t, state.CloseDate.IsNull())
	assert.True(t, state.CreatedAt.IsNull())
	assert.Equal(t, 0, len(state.ManagerIds.Elements()))
}
