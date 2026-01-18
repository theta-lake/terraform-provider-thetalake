package user

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type userPlanModel struct {
	Email            types.String `tfsdk:"email"`
	Name             types.String `tfsdk:"name"`
	Password         types.String `tfsdk:"password"`
	RoleId           types.Int64  `tfsdk:"role_id"`
	SecurityFilterId types.Int64  `tfsdk:"security_filter_id"`
}

type userStateModel struct {
	CreatedAt          timetypes.RFC3339 `tfsdk:"created_at"`
	Id                 types.Int64       `tfsdk:"id"`
	Email              types.String      `tfsdk:"email"`
	Name               types.String      `tfsdk:"name"`
	Disabled           types.Bool        `tfsdk:"disabled"`
	OtpEnabled         types.Bool        `tfsdk:"otp_enabled"`
	Password           types.String      `tfsdk:"password"`
	QueuePaused        types.Bool        `tfsdk:"queue_paused"`
	Role               types.String      `tfsdk:"role"`
	RoleId             types.Int64       `tfsdk:"role_id"`
	SecurityFilterId   types.Int64       `tfsdk:"security_filter_id"`
	SecurityFilterName types.String      `tfsdk:"security_filter_name"`
	UpdatedAt          timetypes.RFC3339 `tfsdk:"updated_at"`
}

func toApiModel(userModel *userPlanModel) thetalake.User {
	newUser := thetalake.User{
		Email:                userModel.Email.ValueString(),
		Name:                 userModel.Name.ValueString(),
		Password:             userModel.Password.ValueString(),
		PasswordConfirmation: userModel.Password.ValueString(),
		RoleId:               userModel.RoleId.ValueInt64(),
	}

	if userModel.SecurityFilterId.ValueInt64() != 0 {
		newUser.SearchId = userModel.SecurityFilterId.ValueInt64()
	}

	return newUser
}

func fromApiModel(user thetalake.User) userStateModel {
	userModel := userStateModel{
		CreatedAt:   timetypes.NewRFC3339TimeValue(user.CreatedAt),
		Disabled:    types.BoolValue(user.Disabled),
		Email:       types.StringValue(user.Email),
		Id:          types.Int64Value(user.Id),
		Name:        types.StringValue(user.Name),
		OtpEnabled:  types.BoolValue(user.OtpEnabled),
		QueuePaused: types.BoolValue(user.QueuePaused),
		Role:        types.StringValue(user.Role.Name),
		RoleId:      types.Int64Value(user.Role.Id),
	}

	if user.SecurityFilter != nil && user.SecurityFilter.SearchId != 0 {
		userModel.SecurityFilterName = types.StringValue(user.SecurityFilter.Name)
		userModel.SecurityFilterId = types.Int64Value(user.SecurityFilter.SearchId)
	} else {
		userModel.SecurityFilterName = types.StringNull()
		userModel.SecurityFilterId = types.Int64Null()
	}

	if user.UpdatedAt == nil {
		userModel.UpdatedAt = timetypes.NewRFC3339Null()
	} else {
		userModel.UpdatedAt = timetypes.NewRFC3339TimeValue(*user.UpdatedAt)
	}

	return userModel
}
