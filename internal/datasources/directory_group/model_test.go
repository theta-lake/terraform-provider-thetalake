package directorygroup

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/theta-lake/terraform-provider-thetalake/internal/client/thetalake"
)

func TestFromApiModel(t *testing.T) {
	dg := thetalake.DirectoryGroup{
		Id:   2065,
		Name: "Test Directory Group",
	}

	model := FromApiModel(dg)

	assert.Equal(t, int64(2065), model.Id.ValueInt64())
	assert.Equal(t, "Test Directory Group", model.Name.ValueString())
}
