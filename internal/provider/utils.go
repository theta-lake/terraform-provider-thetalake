package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func intSlicetoInt64List(s []int) basetypes.ListValue {
	values := make([]attr.Value, len(s))
	for i, value := range s {
		values[i] = types.Int64Value(int64(value))
	}
	return types.ListValueMust(
		types.Int64Type,
		values,
	)
}
