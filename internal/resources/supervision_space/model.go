package supervisionspace

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type supervisionSpacePlanModel struct {
	AllParticipants                   types.Bool   `tfsdk:"all_participants"`
	AllUsers                          types.Bool   `tfsdk:"all_users"`
	Description                       types.String `tfsdk:"description"`
	DirectoryGroupIds                 types.List   `tfsdk:"directory_group_ids"`
	ExternalId                        types.String `tfsdk:"external_id"`
	HardEnforce                       types.Bool   `tfsdk:"hard_enforce"`
	ID                                types.Int64  `tfsdk:"id"`
	IntegrationIds                    types.List   `tfsdk:"integration_ids"`
	MediaTypes                        types.List   `tfsdk:"media_types"`
	Name                              types.String `tfsdk:"name"`
	RetentionLibraryIds               types.List   `tfsdk:"retention_library_ids"`
	RequestedSupervisionSpacePriority types.Int64  `tfsdk:"requested_supervision_space_priority"`
	UserGroupIds                      types.List   `tfsdk:"user_group_ids"`
	UserIds                           types.List   `tfsdk:"user_ids"`
}

type supervisionSpaceStateModel struct {
	AllParticipants                   types.Bool        `tfsdk:"all_participants"`
	AllUsers                          types.Bool        `tfsdk:"all_users"`
	CanDelete                         types.Bool        `tfsdk:"can_delete"`
	CanEnableAllParticipants          types.Bool        `tfsdk:"can_enable_all_participants"`
	CompiledUserList                  types.List        `tfsdk:"compiled_user_list"`
	CreatedAt                         timetypes.RFC3339 `tfsdk:"created_at"`
	Description                       types.String      `tfsdk:"description"`
	DirectoryGroupIds                 types.List        `tfsdk:"directory_group_ids"`
	Disabled                          types.Bool        `tfsdk:"disabled"`
	ExternalId                        types.String      `tfsdk:"external_id"`
	HardEnforce                       types.Bool        `tfsdk:"hard_enforce"`
	Id                                types.Int64       `tfsdk:"id"`
	IntegrationIds                    types.List        `tfsdk:"integration_ids"`
	MediaTypes                        types.List        `tfsdk:"media_types"`
	Name                              types.String      `tfsdk:"name"`
	RetentionLibraryIds               types.List        `tfsdk:"retention_library_ids"`
	RequestedSupervisionSpacePriority types.Int64       `tfsdk:"requested_supervision_space_priority"`
	AssignedSupervisionSpacePriority  types.Int64       `tfsdk:"assigned_supervision_space_priority"`
	UpdatedAt                         timetypes.RFC3339 `tfsdk:"updated_at"`
	UserGroupIds                      types.List        `tfsdk:"user_group_ids"`
	UserIds                           types.List        `tfsdk:"user_ids"`
}

func fromApiModel(space thetalake.SupervisionSpace) supervisionSpaceStateModel {
	spaceModel := supervisionSpaceStateModel{
		AllParticipants:                  types.BoolValue(space.AllParticipants),
		AllUsers:                         types.BoolValue(space.AllUsers),
		CanDelete:                        types.BoolValue(space.CanDelete),
		CanEnableAllParticipants:         types.BoolValue(space.CanEnableAllParticipants),
		CreatedAt:                        timetypes.NewRFC3339TimeValue(space.CreatedAt),
		Description:                      types.StringValue(space.Description),
		Disabled:                         types.BoolValue(space.Disabled),
		ExternalId:                       types.StringValue(space.ExternalId),
		HardEnforce:                      types.BoolValue(space.HardEnforce),
		Id:                               types.Int64Value(int64(space.ID)),
		Name:                             types.StringValue(space.Name),
		AssignedSupervisionSpacePriority: types.Int64Value(int64(space.SupervisionSpacePriority)),
		UpdatedAt:                        timetypes.NewRFC3339TimeValue(space.UpdatedAt),
	}

	// Populate list-typed attributes
	var compiledUserListValues []attr.Value
	for _, user := range space.CompiledUserList {
		compiledUserListValues = append(compiledUserListValues, types.Int64Value(int64(user.Id)))
	}
	spaceModel.CompiledUserList = types.ListValueMust(types.Int64Type, compiledUserListValues)

	var directoryGroupIdsValues []attr.Value
	for _, group := range space.DirectoryGroups {
		directoryGroupIdsValues = append(directoryGroupIdsValues, types.Int64Value(int64(group.Id)))
	}
	spaceModel.DirectoryGroupIds = types.ListValueMust(types.Int64Type, directoryGroupIdsValues)

	var integrationIdsValues []attr.Value
	for _, integration := range space.Integrations {
		integrationIdsValues = append(integrationIdsValues, types.Int64Value(int64(integration.Id)))
	}
	spaceModel.IntegrationIds = types.ListValueMust(types.Int64Type, integrationIdsValues)

	var mediaTypeIdsValues []attr.Value
	for _, mediaType := range space.MediaTypes {
		mediaTypeIdsValues = append(mediaTypeIdsValues, types.StringValue(strings.ToLower(mediaType.Name)))
	}
	spaceModel.MediaTypes = types.ListValueMust(types.StringType, mediaTypeIdsValues)

	var retentionLibraryIdsValues []attr.Value
	for _, library := range space.RetentionLibraries {
		retentionLibraryIdsValues = append(retentionLibraryIdsValues, types.Int64Value(int64(library.Id)))
	}
	spaceModel.RetentionLibraryIds = types.ListValueMust(types.Int64Type, retentionLibraryIdsValues)

	var userGroupIdsValues []attr.Value
	for _, group := range space.UserGroups {
		userGroupIdsValues = append(userGroupIdsValues, types.Int64Value(int64(group.Id)))
	}
	spaceModel.UserGroupIds = types.ListValueMust(types.Int64Type, userGroupIdsValues)

	var userIdsValues []attr.Value
	for _, user := range space.Users {
		userIdsValues = append(userIdsValues, types.Int64Value(int64(user.Id)))
	}
	spaceModel.UserIds = types.ListValueMust(types.Int64Type, userIdsValues)

	return spaceModel
}

func toApiModel(spaceModel *supervisionSpacePlanModel) thetalake.SupervisionSpace {
	newSpace := thetalake.SupervisionSpace{
		AllParticipants:          spaceModel.AllParticipants.ValueBool(),
		AllUsers:                 spaceModel.AllUsers.ValueBool(),
		Description:              spaceModel.Description.ValueString(),
		ExternalId:               spaceModel.ExternalId.ValueString(),
		HardEnforce:              spaceModel.HardEnforce.ValueBool(),
		Name:                     spaceModel.Name.ValueString(),
		SupervisionSpacePriority: int(spaceModel.RequestedSupervisionSpacePriority.ValueInt64()),
		IntegrationIds:           []int64{},
		MediaTypeIds:             []int64{},
		RetentionLibraryIds:      []int64{},
		UserGroupIds:             []int64{},
		UserIds:                  []int64{},
	}

	// Iterate over the DirectoryGroupIds list and populate DirectoryGroupIds field
	var ids []int64
	spaceModel.DirectoryGroupIds.ElementsAs(context.Background(), &ids, false)
	for _, id := range ids {
		newSpace.DirectoryGroupIds = append(newSpace.DirectoryGroupIds, id)
	}

	// Iterate over the IntegrationIds list and populate IntegrationIds field
	ids = []int64{}
	spaceModel.IntegrationIds.ElementsAs(context.Background(), &ids, false)
	for _, id := range ids {
		newSpace.IntegrationIds = append(newSpace.IntegrationIds, id)
	}

	// Iterate over the MediaTypeIds list and populate MediaTypeIds field
	mediaTypes := []string{}
	spaceModel.MediaTypes.ElementsAs(context.Background(), &mediaTypes, false)
	newSpace.MediaTypeIds = thetalake.MediaTypesNamesToIds(mediaTypes)

	// Iterate over the RetentionLibraryIds list and populate RetentionLibraryIds field
	ids = []int64{}
	spaceModel.RetentionLibraryIds.ElementsAs(context.Background(), &ids, false)
	for _, id := range ids {
		newSpace.RetentionLibraryIds = append(newSpace.RetentionLibraryIds, id)
	}

	// Iterate over the UserGroupIds list and populate UserGroupIds field
	ids = []int64{}
	spaceModel.UserGroupIds.ElementsAs(context.Background(), &ids, false)
	for _, id := range ids {
		newSpace.UserGroupIds = append(newSpace.UserGroupIds, id)
	}

	// Iterate over the UserIds list and populate UserIds field
	ids = []int64{}
	spaceModel.UserIds.ElementsAs(context.Background(), &ids, false)
	for _, id := range ids {
		newSpace.UserIds = append(newSpace.UserIds, id)
	}

	return newSpace
}
