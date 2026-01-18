package mediatypes

import "github.com/hashicorp/terraform-plugin-framework/types"

var (
	mediaTypeMap = map[string]int64{
		"Video":      1,
		"Audio":      2,
		"Chat":       3,
		"Attachment": 4,
		"Email":      5,
		"Image":      6,
	}
)

type MediaTypeDataSourceModel struct {
	Name types.String `tfsdk:"name"`
	Id   types.Int64  `tfsdk:"id"`
}
