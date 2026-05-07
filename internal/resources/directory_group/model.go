package directorygroup

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type directoryGroupPlanModel struct {
	Description types.String `tfsdk:"description"`
	ExternalId  types.String `tfsdk:"external_id"`
	IdentityIds types.Set    `tfsdk:"identity_ids"`
	Name        types.String `tfsdk:"name"`
}

type directoryGroupStateModel struct {
	CreatedAt   timetypes.RFC3339 `tfsdk:"created_at"`
	Description types.String      `tfsdk:"description"`
	ExternalId  types.String      `tfsdk:"external_id"`
	Id          types.Int64       `tfsdk:"id"`
	IdentityIds types.Set         `tfsdk:"identity_ids"`
	Name        types.String      `tfsdk:"name"`
	UpdatedAt   timetypes.RFC3339 `tfsdk:"updated_at"`
}

func toApiModel(plan *directoryGroupPlanModel) thetalake.DirectoryGroup {
	dg := thetalake.DirectoryGroup{
		Name: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		desc := plan.Description.ValueString()
		dg.Description = &desc
	}

	if !plan.ExternalId.IsNull() && !plan.ExternalId.IsUnknown() {
		extId := plan.ExternalId.ValueString()
		dg.ExternalId = &extId
	}

	if !plan.IdentityIds.IsNull() && !plan.IdentityIds.IsUnknown() {
		var identityIds []int64
		for _, v := range plan.IdentityIds.Elements() {
			if id, ok := v.(types.Int64); ok {
				identityIds = append(identityIds, id.ValueInt64())
			}
		}
		dg.IdentityIds = identityIds
	}

	return dg
}

func fromApiModel(dg thetalake.DirectoryGroup) directoryGroupStateModel {
	state := directoryGroupStateModel{
		Id:   types.Int64Value(dg.Id),
		Name: types.StringValue(dg.Name),
	}

	if dg.Description != nil {
		state.Description = types.StringValue(*dg.Description)
	} else {
		state.Description = types.StringNull()
	}

	if dg.ExternalId != nil {
		state.ExternalId = types.StringValue(*dg.ExternalId)
	} else {
		state.ExternalId = types.StringNull()
	}

	if dg.CreatedAt != nil {
		state.CreatedAt = timetypes.NewRFC3339TimeValue(*dg.CreatedAt)
	} else {
		state.CreatedAt = timetypes.NewRFC3339Null()
	}

	if dg.UpdatedAt != nil {
		state.UpdatedAt = timetypes.NewRFC3339TimeValue(*dg.UpdatedAt)
	} else {
		state.UpdatedAt = timetypes.NewRFC3339Null()
	}

	var identityIdValues []attr.Value
	for _, id := range dg.IdentityIds {
		identityIdValues = append(identityIdValues, types.Int64Value(id))
	}
	state.IdentityIds = types.SetValueMust(types.Int64Type, identityIdValues)

	return state
}
