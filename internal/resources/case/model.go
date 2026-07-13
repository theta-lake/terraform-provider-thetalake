package legalcase

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type casePlanModel struct {
	CloseDate   timetypes.RFC3339 `tfsdk:"close_date"`
	Description types.String      `tfsdk:"description"`
	ManagerIds  types.Set         `tfsdk:"manager_ids"`
	Name        types.String      `tfsdk:"name"`
	Number      types.String      `tfsdk:"number"`
	OpenDate    timetypes.RFC3339 `tfsdk:"open_date"`
	Visibility  types.String      `tfsdk:"visibility"`
}

type caseStateModel struct {
	CloseDate    timetypes.RFC3339 `tfsdk:"close_date"`
	CreatedAt    timetypes.RFC3339 `tfsdk:"created_at"`
	Description  types.String      `tfsdk:"description"`
	Id           types.Int64       `tfsdk:"id"`
	ManagerIds   types.Set         `tfsdk:"manager_ids"`
	Name         types.String      `tfsdk:"name"`
	Number       types.String      `tfsdk:"number"`
	OpenDate     timetypes.RFC3339 `tfsdk:"open_date"`
	RecordsCount types.Int64       `tfsdk:"records_count"`
	Status       types.String      `tfsdk:"status"`
	UpdatedAt    timetypes.RFC3339 `tfsdk:"updated_at"`
	Visibility   types.String      `tfsdk:"visibility"`
}

func toApiModel(plan *casePlanModel) thetalake.Case {
	c := thetalake.Case{
		Name:       plan.Name.ValueString(),
		Number:     plan.Number.ValueString(),
		Visibility: plan.Visibility.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		c.Description = plan.Description.ValueString()
	}

	if !plan.OpenDate.IsNull() && !plan.OpenDate.IsUnknown() {
		openDate, diags := plan.OpenDate.ValueRFC3339Time()
		if !diags.HasError() {
			c.OpenDate = openDate
		}
	}

	if !plan.ManagerIds.IsNull() && !plan.ManagerIds.IsUnknown() {
		var managerIds []int64
		for _, v := range plan.ManagerIds.Elements() {
			if id, ok := v.(types.Int64); ok {
				managerIds = append(managerIds, id.ValueInt64())
			}
		}
		c.ManagerIds = managerIds
	}

	return c
}

func fromApiModel(c thetalake.Case) caseStateModel {
	state := caseStateModel{
		Description:  types.StringValue(c.Description),
		Id:           types.Int64Value(c.Id),
		Name:         types.StringValue(c.Name),
		Number:       types.StringValue(c.Number),
		OpenDate:     timetypes.NewRFC3339TimeValue(c.OpenDate),
		RecordsCount: types.Int64Value(c.RecordsCount),
		Status:       types.StringValue(c.Status),
		UpdatedAt:    timetypes.NewRFC3339TimeValue(c.UpdatedAt),
		Visibility:   types.StringValue(c.Visibility),
	}

	if c.CloseDate != nil {
		state.CloseDate = timetypes.NewRFC3339TimeValue(*c.CloseDate)
	} else {
		state.CloseDate = timetypes.NewRFC3339Null()
	}

	if !c.CreatedAt.IsZero() {
		state.CreatedAt = timetypes.NewRFC3339TimeValue(c.CreatedAt)
	} else {
		state.CreatedAt = timetypes.NewRFC3339Null()
	}

	var managerIdValues []attr.Value
	for _, id := range c.ManagerIds {
		managerIdValues = append(managerIdValues, types.Int64Value(id))
	}
	state.ManagerIds = types.SetValueMust(types.Int64Type, managerIdValues)

	return state
}
