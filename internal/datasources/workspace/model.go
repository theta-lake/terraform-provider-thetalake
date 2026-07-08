package workspace

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type WorkspaceDataSourceModel struct {
	Description types.String `tfsdk:"description"`
	Id          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
}

func FromApiModel(w thetalake.Workspace) WorkspaceDataSourceModel {
	return WorkspaceDataSourceModel{
		Description: types.StringValue(w.Description),
		Id:          types.Int64Value(w.Id),
		Name:        types.StringValue(w.Name),
	}
}
