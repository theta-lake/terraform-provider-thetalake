package directorygroup

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type DirectoryGroupDataSourceModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func FromApiModel(directoryGroup thetalake.DirectoryGroup) DirectoryGroupDataSourceModel {
	directoryGroupModel := DirectoryGroupDataSourceModel{
		Id:   types.Int64Value(directoryGroup.Id),
		Name: types.StringValue(directoryGroup.Name),
	}

	return directoryGroupModel
}
