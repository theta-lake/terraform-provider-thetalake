package user

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type RoleModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type CurrentWorkspaceModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type SecurityFilterModel struct {
	SearchId types.Int64  `tfsdk:"search_id"`
	Name     types.String `tfsdk:"name"`
}

type UserResourceModel struct {
	CreatedAt             timetypes.RFC3339     `tfsdk:"created_at"`
	Id                    types.Int64           `tfsdk:"id"`
	Email                 types.String          `tfsdk:"email"`
	Name                  types.String          `tfsdk:"name"`
	CurrentWorkspace      CurrentWorkspaceModel `tfsdk:"current_org_unit"`
	Disabled              types.Bool            `tfsdk:"disabled"`
	DisabledAt            timetypes.RFC3339     `tfsdk:"disabled_at"`
	DefaultUserTimezone   types.String          `tfsdk:"default_user_timezone"`
	ForceSso              types.Bool            `tfsdk:"force_sso"`
	HasMultipleWorkspaces types.Bool            `tfsdk:"has_multiple_workspaces"`
	LastLogin             timetypes.RFC3339     `tfsdk:"last_login"`
	OtpEnabled            types.Bool            `tfsdk:"otp_enabled"`
	OtpEnabledAt          timetypes.RFC3339     `tfsdk:"otp_enabled_at"`
	Password              types.String          `tfsdk:"password"`
	PasswordConfirmation  types.String          `tfsdk:"password_confirmation"`
	PasswordChangedAt     timetypes.RFC3339     `tfsdk:"password_changed_at"`
	QueuePaused           types.Bool            `tfsdk:"queue_paused"`
	RoleId                types.Int64           `tfsdk:"role_id"`
	SearchId              types.Int64           `tfsdk:"search_id"`
	Role                  RoleModel             `tfsdk:"role"`
	SecurityFilter        SecurityFilterModel   `tfsdk:"security_filter"`
	UpdatedAt             timetypes.RFC3339     `tfsdk:"updated_at"`
}

func (userModel *UserResourceModel) ToApiModel() thetalake.User {
	newUser := thetalake.User{
		Email:                userModel.Email.ValueString(),
		Name:                 userModel.Name.ValueString(),
		Password:             userModel.Password.ValueString(),
		PasswordConfirmation: userModel.Password.ValueString(),
		RoleId:               userModel.RoleId.ValueInt64(),
	}

	if userModel.SearchId.ValueInt64() != 0 {
		newUser.SearchId = userModel.SearchId.ValueInt64()
	}

	return newUser
}

func FromApiModel(user thetalake.User) UserResourceModel {
	userModel := UserResourceModel{
		CreatedAt:             timetypes.NewRFC3339TimeValue(user.CreatedAt),
		Disabled:              types.BoolValue(user.Disabled),
		Email:                 types.StringValue(user.Email),
		ForceSso:              types.BoolValue(user.ForceSso),
		HasMultipleWorkspaces: types.BoolValue(user.HasMultipleWorkspaces),
		Id:                    types.Int64Value(user.Id),
		Name:                  types.StringValue(user.Name),
		OtpEnabled:            types.BoolValue(user.OtpEnabled),
		QueuePaused:           types.BoolValue(user.QueuePaused),
		RoleId:                types.Int64Value(user.RoleId),
	}

	if user.SecurityFilter != nil {
		userModel.SecurityFilter = SecurityFilterModel{
			SearchId: types.Int64Value(user.SecurityFilter.SearchId),
			Name:     types.StringValue(user.SecurityFilter.Name),
		}
	}

	userModel.CurrentWorkspace = CurrentWorkspaceModel{
		Id:   types.Int64Value(user.CurrentWorkspace.Id),
		Name: types.StringValue(user.CurrentWorkspace.Name),
	}

	userModel.Role = RoleModel{
		Id:   types.Int64Value(user.Role.Id),
		Name: types.StringValue(user.Role.Name),
	}

	if user.DefaultUserTimezone != nil {
		userModel.DefaultUserTimezone = types.StringValue(*user.DefaultUserTimezone)
	} else {
		userModel.DefaultUserTimezone = types.StringNull()
	}

	if user.LastLogin == nil {
		userModel.LastLogin = timetypes.NewRFC3339Null()
	} else {
		userModel.LastLogin = timetypes.NewRFC3339TimeValue(*user.LastLogin)
	}

	if user.DisabledAt == nil {
		userModel.DisabledAt = timetypes.NewRFC3339Null()
	} else {
		userModel.DisabledAt = timetypes.NewRFC3339TimeValue(*user.DisabledAt)
	}

	if user.OtpEnabledAt == nil {
		userModel.OtpEnabledAt = timetypes.NewRFC3339Null()
	} else {
		userModel.OtpEnabledAt = timetypes.NewRFC3339TimeValue(*user.OtpEnabledAt)
	}

	if user.PasswordChangedAt == nil {
		userModel.PasswordChangedAt = timetypes.NewRFC3339Null()
	} else {
		userModel.PasswordChangedAt = timetypes.NewRFC3339TimeValue(*user.PasswordChangedAt)
	}

	if user.UpdatedAt == nil {
		userModel.UpdatedAt = timetypes.NewRFC3339Null()
	} else {
		userModel.UpdatedAt = timetypes.NewRFC3339TimeValue(*user.UpdatedAt)
	}

	return userModel
}
