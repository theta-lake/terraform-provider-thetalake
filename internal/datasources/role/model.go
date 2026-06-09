package role

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type RoleDataSourceModel struct {
	CreatedAt     timetypes.RFC3339 `tfsdk:"created_at"`
	Default       types.Bool        `tfsdk:"default"`
	Description   types.String      `tfsdk:"description"`
	Id            types.Int64       `tfsdk:"id"`
	IsBuiltIn     types.Bool        `tfsdk:"is_built_in"`
	Name          types.String      `tfsdk:"name"`
	NumberOfUsers types.Int64       `tfsdk:"number_of_users"`
	UpdatedAt     timetypes.RFC3339 `tfsdk:"updated_at"`
}

func FromApiModel(role thetalake.Role) RoleDataSourceModel {
	roleModel := RoleDataSourceModel{
		Default:       types.BoolValue(role.Default),
		Description:   types.StringValue(role.Description),
		Id:            types.Int64Value(role.Id),
		IsBuiltIn:     types.BoolValue(role.IsBuiltIn),
		Name:          types.StringValue(role.Name),
		NumberOfUsers: types.Int64Value(role.NumberOfUsers),
	}

	if role.CreatedAt != nil {
		roleModel.CreatedAt = timetypes.NewRFC3339TimeValue(*role.CreatedAt)
	} else {
		roleModel.CreatedAt = timetypes.NewRFC3339Null()
	}

	if role.UpdatedAt != nil {
		roleModel.UpdatedAt = timetypes.NewRFC3339TimeValue(*role.UpdatedAt)
	} else {
		roleModel.UpdatedAt = timetypes.NewRFC3339Null()
	}

	return roleModel
}
