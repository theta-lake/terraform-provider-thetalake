package label

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type labelPlanModel struct {
	BackgroundColor types.String `tfsdk:"background_color"`
	Hidden          types.Bool   `tfsdk:"hidden"`
	LongName        types.String `tfsdk:"long_name"`
	ShortName       types.String `tfsdk:"short_name"`
}

type labelStateModel struct {
	BackgroundColor   types.String      `tfsdk:"background_color"`
	CreatedAt         timetypes.RFC3339 `tfsdk:"created_at"`
	Hidden            types.Bool        `tfsdk:"hidden"`
	Id                types.Int64       `tfsdk:"id"`
	LongName          types.String      `tfsdk:"long_name"`
	OrgUnitId         types.Int64       `tfsdk:"org_unit_id"`
	ShortName         types.String      `tfsdk:"short_name"`
	TaggedDatumsCount types.Int64       `tfsdk:"tagged_datums_count"`
	UpdatedAt         timetypes.RFC3339 `tfsdk:"updated_at"`
	UserId            types.Int64       `tfsdk:"user_id"`
}

func toApiModel(plan *labelPlanModel) thetalake.Label {
	return thetalake.Label{
		BackgroundColor: plan.BackgroundColor.ValueString(),
		Hidden:          plan.Hidden.ValueBool(),
		LongName:        plan.LongName.ValueString(),
		ShortName:       plan.ShortName.ValueString(),
	}
}

func fromApiModel(label thetalake.Label) labelStateModel {
	state := labelStateModel{
		BackgroundColor:   types.StringValue(label.BackgroundColor),
		CreatedAt:         timetypes.NewRFC3339TimeValue(label.CreatedAt),
		Hidden:            types.BoolValue(label.Hidden),
		Id:                types.Int64Value(label.Id),
		LongName:          types.StringValue(label.LongName),
		OrgUnitId:         types.Int64Value(label.OrgUnitId),
		ShortName:         types.StringValue(label.ShortName),
		TaggedDatumsCount: types.Int64Value(label.TaggedDatumsCount),
		UserId:            types.Int64Value(label.UserId),
	}

	if label.UpdatedAt == nil {
		state.UpdatedAt = timetypes.NewRFC3339Null()
	} else {
		state.UpdatedAt = timetypes.NewRFC3339TimeValue(*label.UpdatedAt)
	}

	return state
}
