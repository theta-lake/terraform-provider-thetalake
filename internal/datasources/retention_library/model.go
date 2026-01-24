package retentionlibrary

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

type RetentionLibraryDataSourceModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func FromApiModel(retentionLibrary thetalake.RetentionLibrary) RetentionLibraryDataSourceModel {
	retentionLibraryModel := RetentionLibraryDataSourceModel{
		Id:   types.Int64Value(retentionLibrary.Id),
		Name: types.StringValue(retentionLibrary.Name),
	}

	return retentionLibraryModel
}
