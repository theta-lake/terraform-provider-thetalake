package integration

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type IntegrationDataSourceModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func FromApiModel(integration thetalake.Integration) IntegrationDataSourceModel {
	integrationModel := IntegrationDataSourceModel{
		Id:   types.Int64Value(integration.Id),
		Name: types.StringValue(integration.Name),
	}

	return integrationModel
}
