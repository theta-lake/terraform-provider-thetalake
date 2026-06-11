package role

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type rolePlanModel struct {
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
	Permissions types.Set    `tfsdk:"permissions"`
}

type roleStateModel struct {
	CreatedAt     timetypes.RFC3339 `tfsdk:"created_at"`
	Default       types.Bool        `tfsdk:"default"`
	Description   types.String      `tfsdk:"description"`
	Id            types.Int64       `tfsdk:"id"`
	IsBuiltIn     types.Bool        `tfsdk:"is_built_in"`
	Name          types.String      `tfsdk:"name"`
	NumberOfUsers types.Int64       `tfsdk:"number_of_users"`
	Permissions   types.Set         `tfsdk:"permissions"`
	UpdatedAt     timetypes.RFC3339 `tfsdk:"updated_at"`
}

func toApiModel(plan *rolePlanModel) thetalake.Role {
	role := thetalake.Role{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	if !plan.Permissions.IsNull() && !plan.Permissions.IsUnknown() {
		var permissions []string
		for _, v := range plan.Permissions.Elements() {
			if s, ok := v.(types.String); ok {
				permissions = append(permissions, s.ValueString())
			}
		}
		role.Permissions = permissions
	}

	return role
}

func fromApiModel(r thetalake.Role) roleStateModel {
	state := roleStateModel{
		Default:       types.BoolValue(r.Default),
		Description:   types.StringValue(r.Description),
		Id:            types.Int64Value(r.Id),
		IsBuiltIn:     types.BoolValue(r.IsBuiltIn),
		Name:          types.StringValue(r.Name),
		NumberOfUsers: types.Int64Value(r.NumberOfUsers),
	}

	if r.CreatedAt != nil {
		state.CreatedAt = timetypes.NewRFC3339TimeValue(*r.CreatedAt)
	} else {
		state.CreatedAt = timetypes.NewRFC3339Null()
	}

	if r.UpdatedAt != nil {
		state.UpdatedAt = timetypes.NewRFC3339TimeValue(*r.UpdatedAt)
	} else {
		state.UpdatedAt = timetypes.NewRFC3339Null()
	}

	permValues := make([]attr.Value, len(r.Permissions))
	for i, p := range r.Permissions {
		permValues[i] = types.StringValue(p)
	}
	state.Permissions = types.SetValueMust(types.StringType, permValues)

	return state
}
