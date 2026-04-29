package usergroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type userGroupPlanModel struct {
	Description types.String `tfsdk:"description"`
	ExternalId  types.String `tfsdk:"external_id"`
	Name        types.String `tfsdk:"name"`
	UserIds     types.List   `tfsdk:"user_ids"`
}

type userGroupStateModel struct {
	CreatedAt   timetypes.RFC3339 `tfsdk:"created_at"`
	Description types.String      `tfsdk:"description"`
	ExternalId  types.String      `tfsdk:"external_id"`
	Id          types.Int64       `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	UpdatedAt   timetypes.RFC3339 `tfsdk:"updated_at"`
	UserIds     types.List        `tfsdk:"user_ids"`
}

func toApiModel(plan *userGroupPlanModel) thetalake.UserGroup {
	userGroup := thetalake.UserGroup{
		Description: plan.Description.ValueString(),
		Name:        plan.Name.ValueString(),
	}

	if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
		externalId := plan.ExternalId.ValueString()
		userGroup.ExternalId = &externalId
	}

	var userIds []int64
	for _, v := range plan.UserIds.Elements() {
		if id, ok := v.(types.Int64); ok {
			userIds = append(userIds, id.ValueInt64())
		}
	}
	userGroup.UserIds = userIds

	return userGroup
}

func fromApiModel(userGroup thetalake.UserGroup) userGroupStateModel {
	state := userGroupStateModel{
		Description: types.StringValue(userGroup.Description),
		Id:          types.Int64Value(userGroup.Id),
		Name:        types.StringValue(userGroup.Name),
	}

	if userGroup.ExternalId != nil {
		state.ExternalId = types.StringValue(*userGroup.ExternalId)
	} else {
		state.ExternalId = types.StringNull()
	}

	if userGroup.CreatedAt != nil {
		state.CreatedAt = timetypes.NewRFC3339TimeValue(*userGroup.CreatedAt)
	} else {
		state.CreatedAt = timetypes.NewRFC3339Null()
	}

	if userGroup.UpdatedAt != nil {
		state.UpdatedAt = timetypes.NewRFC3339TimeValue(*userGroup.UpdatedAt)
	} else {
		state.UpdatedAt = timetypes.NewRFC3339Null()
	}

	var userIdValues []attr.Value
	for _, id := range userGroup.UserIds {
		userIdValues = append(userIdValues, types.Int64Value(id))
	}
	state.UserIds = types.ListValueMust(types.Int64Type, userIdValues)

	return state
}
