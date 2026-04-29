package user

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type UserDataSourceModel struct {
	Email types.String `tfsdk:"email"`
	Id    types.Int64  `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
}

func FromApiModel(u thetalake.User) UserDataSourceModel {
	return UserDataSourceModel{
		Email: types.StringValue(u.Email),
		Id:    types.Int64Value(u.Id),
		Name:  types.StringValue(u.Name),
	}
}
