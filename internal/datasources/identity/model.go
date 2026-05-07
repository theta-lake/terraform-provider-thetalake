package identity

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type IdentityDataSourceModel struct {
	Email      types.String `tfsdk:"email"`
	ExternalId types.String `tfsdk:"external_id"`
	Id         types.Int64  `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
}

func FromApiModel(identity thetalake.Identity) IdentityDataSourceModel {
	model := IdentityDataSourceModel{
		Id:   types.Int64Value(identity.Id),
		Name: types.StringValue(identity.Name),
	}

	if identity.Email != nil {
		model.Email = types.StringValue(*identity.Email)
	} else {
		model.Email = types.StringNull()
	}

	if identity.ExternalId != nil {
		model.ExternalId = types.StringValue(*identity.ExternalId)
	} else {
		model.ExternalId = types.StringNull()
	}

	return model
}
