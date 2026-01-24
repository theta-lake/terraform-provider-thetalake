package usergroup

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type UserGroupDataSourceModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func FromApiModel(userGroup thetalake.UserGroup) UserGroupDataSourceModel {
	userGroupModel := UserGroupDataSourceModel{
		Id:   types.Int64Value(userGroup.Id),
		Name: types.StringValue(userGroup.Name),
	}

	return userGroupModel
}
