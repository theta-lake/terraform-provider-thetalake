package supervisionspace

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
	plan := &supervisionSpacePlanModel{
		AllParticipants:                   types.BoolValue(true),
		AllUsers:                          types.BoolValue(false),
		Description:                       types.StringValue("Test Supervision Space"),
		ExternalId:                        types.StringValue("ext-123"),
		HardEnforce:                       types.BoolValue(true),
		Name:                              types.StringValue("Space Name"),
		RequestedSupervisionSpacePriority: types.Int64Value(42),
	}

	// Populate list fields
	plan.DirectoryGroupIds = types.ListValueMust(types.Int64Type, []attr.Value{
		types.Int64Value(1001),
		types.Int64Value(1002),
	})
	plan.IntegrationIds = types.ListValueMust(types.Int64Type, []attr.Value{
		types.Int64Value(2001),
	})
	plan.MediaTypes = types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("chat"),
		types.StringValue("email"),
	})
	plan.RetentionLibraryIds = types.ListValueMust(types.Int64Type, []attr.Value{
		types.Int64Value(3001),
	})
	plan.UserGroupIds = types.ListValueMust(types.Int64Type, []attr.Value{
		types.Int64Value(4001),
	})
	plan.UserIds = types.ListValueMust(types.Int64Type, []attr.Value{
		types.Int64Value(5001),
		types.Int64Value(5002),
	})

	apiModel := toApiModel(plan)

	// Scalar fields
	assert.Equal(t, true, apiModel.AllParticipants)
	assert.Equal(t, false, apiModel.AllUsers)
	assert.Equal(t, plan.Description.ValueString(), apiModel.Description)
	assert.Equal(t, plan.ExternalId.ValueString(), apiModel.ExternalId)
	assert.Equal(t, plan.HardEnforce.ValueBool(), apiModel.HardEnforce)
	assert.Equal(t, plan.Name.ValueString(), apiModel.Name)
	assert.Equal(t, int(plan.RequestedSupervisionSpacePriority.ValueInt64()), apiModel.SupervisionSpacePriority)

	// DirectoryGroupIds
	assert.Equal(t, []int64{1001, 1002}, apiModel.DirectoryGroupIds)
	// IntegrationIds
	assert.Equal(t, []int64{2001}, apiModel.IntegrationIds)
	// RetentionLibraryIds
	assert.Equal(t, []int64{3001}, apiModel.RetentionLibraryIds)
	// UserGroupIds
	assert.Equal(t, []int64{4001}, apiModel.UserGroupIds)
	// UserIds
	assert.Equal(t, []int64{5001, 5002}, apiModel.UserIds)

	// MediaTypes -> MediaTypeIds mapping
	expectedMediaTypeIds := thetalake.MediaTypesNamesToIds([]string{"chat", "email"})
	assert.Equal(t, expectedMediaTypeIds, apiModel.MediaTypeIds)
}

func TestFromApiModel(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")

	apiSpace := thetalake.SupervisionSpace{
		AllParticipants:          true,
		AllUsers:                 false,
		CanDelete:                true,
		CanEnableAllParticipants: false,
		CreatedAt:                ts,
		Description:              "Test Supervision Space",
		Disabled:                 false,
		ExternalId:               "ext-123",
		HardEnforce:              true,
		Id:                       1234,
		Name:                     "Space Name",
		SupervisionSpacePriority: 7,
		UpdatedAt:                ts,
		DirectoryGroups: []thetalake.DirectoryGroup{
			{Id: 1001, Name: "DG1"},
			{Id: 1002, Name: "DG2"},
		},
		Integrations: []thetalake.Integration{
			{Id: 2001, Name: "Int1"},
		},
		MediaTypes: []thetalake.MediaType{
			{Id: 3, Name: "Chat"},
			{Id: 5, Name: "EMAIL"},
		},
		RetentionLibraries: []thetalake.RetentionLibrary{
			{Id: 3001, Name: "RL1"},
		},
		UserGroups: []thetalake.UserGroup{
			{Id: 4001, Name: "UG1"},
		},
		UserIds: []int64{5001, 5002},
	}

	state := fromApiModel(apiSpace)

	// Scalar fields
	assert.Equal(t, true, state.AllParticipants.ValueBool())
	assert.Equal(t, false, state.AllUsers.ValueBool())
	assert.Equal(t, true, state.CanDelete.ValueBool())
	assert.Equal(t, false, state.CanEnableAllParticipants.ValueBool())
	assert.Equal(t, ts.Format(time.RFC3339), state.CreatedAt.ValueString())
	assert.Equal(t, "Test Supervision Space", state.Description.ValueString())
	assert.Equal(t, false, state.Disabled.ValueBool())
	assert.Equal(t, "ext-123", state.ExternalId.ValueString())
	assert.Equal(t, true, state.HardEnforce.ValueBool())
	assert.Equal(t, int64(1234), state.Id.ValueInt64())
	assert.Equal(t, "Space Name", state.Name.ValueString())
	assert.Equal(t, int64(7), state.AssignedSupervisionSpacePriority.ValueInt64())
	assert.Equal(t, ts.Format(time.RFC3339), state.UpdatedAt.ValueString())

	// DirectoryGroupIds
	var dirIds []int64
	state.DirectoryGroupIds.ElementsAs(context.Background(), &dirIds, false)
	assert.Equal(t, []int64{1001, 1002}, dirIds)

	// IntegrationIds
	var integrationIds []int64
	state.IntegrationIds.ElementsAs(context.Background(), &integrationIds, false)
	assert.Equal(t, []int64{2001}, integrationIds)

	// MediaTypes (names, lowercased)
	var mediaTypes []string
	state.MediaTypes.ElementsAs(context.Background(), &mediaTypes, false)
	assert.Equal(t, []string{"chat", "email"}, mediaTypes)

	// RetentionLibraryIds
	var retentionIds []int64
	state.RetentionLibraryIds.ElementsAs(context.Background(), &retentionIds, false)
	assert.Equal(t, []int64{3001}, retentionIds)

	// UserGroupIds
	var userGroupIds []int64
	state.UserGroupIds.ElementsAs(context.Background(), &userGroupIds, false)
	assert.Equal(t, []int64{4001}, userGroupIds)

	// UserIds
	var userIds []int64
	state.UserIds.ElementsAs(context.Background(), &userIds, false)
	assert.Equal(t, []int64{5001, 5002}, userIds)
}
